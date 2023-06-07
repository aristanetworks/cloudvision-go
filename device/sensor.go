// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/version"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var (
	// precreate more frequently used gnmi paths
	lastErrorKey = pgnmi.PathFromString("last-error")
	lastSeenKey  = pgnmi.Path("last-seen")

	undefinedPathElem = &gnmi.PathElem{Name: "___undefined___"}
)

func catchPanic(desc string, f func() error) func() error {
	return func() (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				err = fmt.Errorf("fatal error in %s: %v", desc, rerr)
			}
		}()
		return f()
	}
}

// CredentialResolver is the interface used to resolve credentials.
type CredentialResolver interface {
	Resolve(ctx context.Context, ref string) (string, error)
}
type passthroughCredResolver struct {
}

func (p *passthroughCredResolver) Resolve(ctx context.Context, ref string) (string, error) {
	return ref, nil
}

var passthroughResolver CredentialResolver = &passthroughCredResolver{}

// datasourceConfig holds the configs received from the server
type datasourceConfig struct {
	name       string
	typ        string
	enabled    bool
	option     map[string]string
	credential map[string]string
	loglevel   logrus.Level
}

func (c datasourceConfig) clone() *datasourceConfig {
	return &datasourceConfig{
		name:       c.name,
		typ:        c.typ,
		enabled:    c.enabled,
		option:     cloneMap(c.option),
		credential: cloneMap(c.credential),
		loglevel:   c.loglevel,
	}
}

func (c datasourceConfig) String() string {
	return fmt.Sprintf("name: %s, typ: %s, enabled: %t, option: %v, loglevel: %s",
		c.name, c.typ, c.enabled, c.option, c.loglevel)
}

func (c datasourceConfig) equals(other *datasourceConfig) bool {
	return c.name == other.name &&
		c.typ == other.typ &&
		c.enabled == other.enabled &&
		mapEquals(c.option, other.option) &&
		mapEquals(c.credential, other.credential)
}

func cloneMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mapEquals(lh map[string]string, rh map[string]string) bool {
	if len(lh) != len(rh) {
		return false
	}
	for k, v := range lh {
		if v2, ok := rh[k]; !ok {
			return false
		} else if v2 != v {
			return false
		}
	}
	return true
}

type datasource struct {
	log           *logrus.Entry
	sensorID      string
	gnmic         gnmi.GNMIClient
	apiaddr       string
	clientFactory func(gnmi.GNMIClient, *Info) cvclient.CVClient
	grpcConnector GRPCConnector // Connector to get gRPC connection
	standalone    bool

	credResolver CredentialResolver

	// Current running config. Config changes require a datasource restart.
	config            *datasourceConfig
	running           atomic.Bool
	failureRetryTimer *provider.BackoffTimer

	redeployTimer *time.Timer
	cancel        context.CancelFunc
	execGroup     *errgroup.Group

	info     *Info
	grpcc    *grpc.ClientConn
	cvClient cvclient.CVClient

	// Holds datasource/state/sensor[id=sensor]/source[name=datasource]
	statePrefix       *gnmi.Path
	heartbeatInterval time.Duration

	monitor *datasourceMonitor
	logRate float64
}

func (d *datasource) submitDatasourceUpdates(ctx context.Context,
	updates ...*gnmi.Update) error {
	d.log.Debugf("submitting updates: %v", updates)
	_, err := d.gnmic.Set(ctx, &gnmi.SetRequest{
		Prefix: d.statePrefix,
		Update: updates,
	})
	return err
}

func (d *datasource) handleDatasourceError(ctx context.Context, e error) {
	d.log.Error(e)
	err := d.submitDatasourceUpdates(ctx,
		pgnmi.Update(lastErrorKey, agnmi.TypedValue(e.Error())))
	if err != nil {
		d.log.Errorf("Failed to publish error: %v. Reason: %v", e, err)
	}
}

func (d *datasource) scheduleRestart(in time.Duration) {
	// If was already triggered we are ok in triggering it again
	_ = d.redeployTimer.Reset(in)
}

func (d *datasource) deploy(ctx context.Context, cfg *datasourceConfig) error {
	d.log.Tracef("Trying to deploy config: %v. Current: %v", cfg, d.config)

	// We don't need to re-run if we are already running the same config.
	if d.running.Load() && cfg.equals(d.config) {
		d.log.Info("Deploy requested with the same config, ignoring request.")
		return nil
	}

	// Start processing new config
	d.stop()
	d.config = cfg.clone()

	if !cfg.enabled {
		d.failureRetryTimer.Reset() // reset failure backoff
		d.log.Infof("Data source %v disabled", d.config.name)
		if err := d.submitDatasourceUpdates(ctx,
			pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(false)),
			pgnmi.Update(lastErrorKey, agnmi.TypedValue("Data source disabled"))); err != nil {
			d.log.Error("Failed to publish disabled status:", err)
		}
		return nil
	}
	d.log.Info("Starting run")

	ctx, cancel := context.WithCancel(ctx)
	d.execGroup, ctx = errgroup.WithContext(ctx)
	d.cancel = cancel

	// Start running the device in the background
	d.execGroup.Go(func() error {
		runCtx, cancel := context.WithCancel(ctx)
		d.running.Store(true)
		defer func() {
			d.running.Store(false)
			cancel() // make sure to cancel everything if the Run function returns for any reason.
		}()

		err := d.Run(runCtx)

		// If main context was canceled, do not retry.
		select {
		case <-ctx.Done():
			if err != nil {
				// Not necessarily an error, just logging in case the device had
				// something relevant to say.
				d.log.Infof("Data source stopped due to cancel request. Returned: %v", err)
			}
			return ctx.Err()
		default:
		}

		if err == nil || IsBadConfigError(err) {
			if err != nil {
				d.handleDatasourceError(ctx, fmt.Errorf("Data source stopped: %v", err))
			} else {
				errS := d.submitDatasourceUpdates(ctx,
					pgnmi.Update(lastErrorKey, agnmi.TypedValue("Data source stopped")))
				if errS != nil {
					d.log.Errorf("failed to publish Data source stopped event due to: %v", errS)
				}
			}
			// Return without scheduling a re-run.
			// We treat returning no error as an acceptable stop.
			// But still cancel run ctx which will mark the data source as inactive.
			return nil
		}

		// Handle return error by pushing it to datasource last-error state
		backoff := d.failureRetryTimer.Backoff()
		err = fmt.Errorf("Datasource stopped unexpectedly: %w. Retrying in %v", err, backoff)
		d.handleDatasourceError(ctx, err)
		d.scheduleRestart(backoff)
		return err
	})

	return nil
}

// Signal stop and wait for datasource to finish
func (d *datasource) stop() {
	if d.redeployTimer != nil {
		_ = d.redeployTimer.Stop()
	}
	if d.cancel != nil {
		d.cancel()
	}
	if d.execGroup != nil {
		if err := d.execGroup.Wait(); err != nil && !errors.Is(err, context.Canceled) {
			d.log.Errorf("Stop encountered error: %v", err)
		}
	}
}

func (d *datasource) resolveCredentials(ctx context.Context,
	configs map[string]string) (map[string]string, error) {
	var err error
	creds := make(map[string]string, len(configs))
	for key, cred := range configs {
		creds[key], err = d.credResolver.Resolve(ctx, cred)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve credential for %v: %w", key, err)
		}
	}
	return creds, nil
}

// Run executes the datasource
func (d *datasource) Run(ctx context.Context) (err error) {
	// Submit initial status information initially as the next operations can be slow
	if err := d.submitDatasourceUpdates(ctx,
		// last-seen is not sent as to not qualify the datasource as streaming yet
		pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
		pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue(d.config.typ)),
		pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
	); err != nil {
		return err
	}

	creds, err := d.resolveCredentials(ctx, d.config.credential)
	if err != nil {
		return err
	}

	// Prepare device to execute based on datasource config
	info, err := NewDeviceInfo(ctx, &Config{
		Device:  d.config.typ,
		Options: mergeOpts(d.config.option, creds),
	}, d.monitor)
	if err != nil {
		return err
	}
	d.cvClient = d.clientFactory(d.gnmic, info)
	d.info = info
	deviceID := info.ID

	var updates []*gnmi.Update
	if len(deviceID) > 0 {
		updates = append(updates, pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue(deviceID)))
		if err := d.submitDatasourceUpdates(ctx, updates...); err != nil {
			return fmt.Errorf("failed to publish startup status: %w", err)
		}
	}

	// Register the device before starting providers. If we can't reach
	// the device right now, we should return an error rather than
	// considering it added.
	if err := d.cvClient.SendDeviceMetadata(ctx); err != nil {
		return fmt.Errorf("error sending device metadata for device %q (%s): %w",
			deviceID, d.info.Config.Device, err)
	}

	errg, ctx := errgroup.WithContext(ctx)

	// Start providers.
	errg.Go(func() error {
		if err := d.runProviders(ctx); err != nil {
			return fmt.Errorf("error starting providers for device %q (%s): %w",
				deviceID, d.info.Config.Device, err)
		}
		return nil
	})

	// Send metadata updates.
	errg.Go(func() error {
		if err := d.sendPeriodicUpdates(ctx); err != nil {
			return fmt.Errorf("error updating device metadata: %w", err)
		}
		return nil
	})

	// Send minoitor logging updates.
	errg.Go(func() error {
		if err := d.sendMonitorLogging(ctx); err != nil {
			return fmt.Errorf("error updating monitor logging: %w", err)
		}
		return nil
	})

	// Start manager, maybe.
	if manager, ok := d.info.Device.(Manager); ok {
		inv := newDatasourceInventory(d.cvClient)
		errg.Go(catchPanic("Manage", func() error {
			if err := manager.Manage(ctx, inv); err != nil {
				return fmt.Errorf("error in Manage: %w", err)
			}
			return nil
		}))
	}

	return errg.Wait()
}

func (d *datasource) runProviders(ctx context.Context) error {
	providers, err := d.info.Device.Providers()
	if err != nil {
		return err
	}

	errg, ctx := errgroup.WithContext(ctx)

	// We may override this grpc using the grpcConnector if available.
	grpcConn := d.grpcc

	for _, p := range providers {
		p := p // scope p for goroutines that use it

		switch pt := p.(type) {
		case provider.GNMIProvider:
			pt.InitGNMI(d.cvClient.ForProvider(pt))
		case provider.GRPCProvider:
			if d.grpcConnector != nil && grpcConn == d.grpcc {
				// lazy initialize the new connection once for the device
				cc := GRPCConnectorConfig{
					d.info.ID,
					d.standalone,
				}
				d.log.Debugf("Opening connector connection to device")
				conn, err := d.grpcConnector.Connect(ctx, d.grpcc, d.apiaddr, cc)
				if err != nil {
					return fmt.Errorf("gRPC connection to device %v failed: %w", cc.DeviceID, err)
				}
				grpcConn = conn
				// if we setup a new connection, close it when the context is canceled
				errg.Go(func() error {
					<-ctx.Done()
					d.log.Debugf("Closing connector connection to device")
					return conn.Close()
				})
			}
			pt.InitGRPC(grpcConn)
		default:
			return fmt.Errorf("unexpected provider type %T", p)
		}

		// Start the provider.
		errg.Go(catchPanic(fmt.Sprintf("%T.Run", p), func() error {
			if err := p.Run(ctx); err != nil {
				return fmt.Errorf("provider %T exiting with error: %w", p, err)
			}
			return nil
		}))
	}

	return errg.Wait()
}

func (d *datasource) sendPeriodicUpdates(ctx context.Context) error {
	ticker := time.NewTicker(d.heartbeatInterval)
	defer ticker.Stop()

	wasFailing := false // used to only log once when device is unhealthy and back alive
	streamingStart := true

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			alive, err := d.info.Device.Alive(ctx)
			if err == nil && alive {
				ts := agnmi.TypedValue(time.Now().UnixNano())
				updates := []*gnmi.Update{pgnmi.Update(lastSeenKey, ts)}

				if wasFailing {
					d.log.Info("Device is back alive")
					wasFailing = false
					updates = append(updates, pgnmi.Update(lastErrorKey,
						agnmi.TypedValue("Device is back alive")))
				}

				if streamingStart {
					updates = append(updates, pgnmi.Update(pgnmi.Path("streaming-start"), ts))
				}

				if err := d.submitDatasourceUpdates(ctx, updates...); err != nil {
					d.log.Error("Publish status failed:", err)
				} else if err == nil {
					// Clear flag only after first successful set
					streamingStart = false
				}

				if err := d.cvClient.SendHeartbeat(ctx, alive); err != nil {
					// Don't give up if an update fails for some reason.
					d.log.Errorf("Error sending heartbeat: %v", err)
				}
			} else if !wasFailing {
				msg := errors.New("Device not alive")
				if err != nil {
					msg = fmt.Errorf("Device not alive: %w", err)
				}
				d.handleDatasourceError(ctx, msg)
				wasFailing = true
			}
		}
	}
}

func (d *datasource) sendMonitorLogging(ctx context.Context) error {
	logLimter := rate.NewLimiter(
		rate.Limit(d.logRate)*rate.Every(time.Minute), 10)
	isLoggingDense := false
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-d.monitor.logCh:
			if !logLimter.Allow() {
				if isLoggingDense {
					continue
				} else {
					msg = "Logging is too dense. Please check in the console"
					isLoggingDense = true
				}
			} else {
				isLoggingDense = false
			}
			err := d.submitDatasourceUpdates(ctx,
				pgnmi.Update(lastErrorKey, agnmi.TypedValue(msg)))
			if err != nil {
				d.log.Errorf("Error sending monitor log: %v", err)
			}
		}
	}
}

func datasourceFromPath(path *gnmi.Path) string {
	for _, elm := range path.Elem {
		if elm.Name == "source" {
			if name, ok := elm.Key["name"]; ok {
				return name
			}
		}
	}
	return ""
}

type datasourceInventory struct {
	client cvclient.CVClient

	rwlock  sync.RWMutex // protects the fields below
	devices map[string]*Info
}

func newDatasourceInventory(client cvclient.CVClient) *datasourceInventory {
	return &datasourceInventory{
		client:  client,
		devices: map[string]*Info{},
	}
}

func (s *datasourceInventory) getManagedIDs() []string {
	ids := make([]string, 0, len(s.devices))
	for _, v := range s.devices {
		ids = append(ids, v.ID)
	}
	sort.Strings(ids) // keep it consistent
	return ids
}

func (s *datasourceInventory) Add(deviceInfo *Info) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	if _, ok := s.devices[deviceInfo.ID]; !ok {
		s.devices[deviceInfo.ID] = deviceInfo
		s.client.SetManagedDevices(s.getManagedIDs())
	}
	return nil
}

func (s *datasourceInventory) Delete(key string) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	delete(s.devices, key)
	s.client.SetManagedDevices(s.getManagedIDs())
	return nil
}

func (s *datasourceInventory) Get(key string) (*Info, error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	return s.devices[key], nil
}

func (s *datasourceInventory) List() []*Info {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	out := make([]*Info, 0, len(s.devices))
	for _, v := range s.devices {
		out = append(out, v)
	}
	return out
}

func mergeOpts(o, c map[string]string) map[string]string {
	x := make(map[string]string, len(o)+len(c))
	for k, v := range o {
		x[k] = v
	}
	for k, v := range c {
		x[k] = v
	}
	return x
}

// Sensor manages the config for multiple datasources
type Sensor struct {
	id                string
	gnmic             gnmi.GNMIClient
	grpcc             *grpc.ClientConn
	apiaddr           string
	grpcConnector     GRPCConnector // Connector to get gRPC connection
	standalone        bool
	heartbeatInterval time.Duration

	credResolver CredentialResolver

	redeployDatasource chan string
	datasourceConfig   map[string]*datasourceConfig
	datasource         map[string]*datasource
	clientFactory      func(gnmi.GNMIClient, *Info) cvclient.CVClient

	// channel to receive custom configs from.
	configCh chan *Config

	deviceRedeployTimer     time.Duration
	failureRetryBackoffMax  time.Duration
	failureRetryBackoffBase time.Duration

	active        bool
	heartbeatLock sync.Mutex // used to stop synchronize and prevent heartbeats when deleting sensor
	log           *logrus.Entry
	logRate       float64
	statePrefix   *gnmi.Path
}

// SensorOption is used to configure the Sensor
type SensorOption func(m *Sensor)

// WithSensorConnectorAddress sets the connector address
func WithSensorConnectorAddress(addr string) SensorOption {
	return func(s *Sensor) { s.apiaddr = addr }
}

// WithSensorHeartbeatInterval sets the duration between sensor heartbeats
func WithSensorHeartbeatInterval(d time.Duration) SensorOption {
	return func(s *Sensor) { s.heartbeatInterval = d }
}

// WithSensorFailureRetryBackoffBase sets the duration between datasource restarts on failure
func WithSensorFailureRetryBackoffBase(d time.Duration) SensorOption {
	return func(s *Sensor) { s.failureRetryBackoffBase = d }
}

// WithSensorFailureRetryBackoffMax sets the max backoff between datasource restarts due to failures
func WithSensorFailureRetryBackoffMax(d time.Duration) SensorOption {
	return func(s *Sensor) { s.failureRetryBackoffMax = d }
}

// WithSensorGNMIClient sets a gNMI client on the Sensor.
func WithSensorGNMIClient(c gnmi.GNMIClient) SensorOption {
	return func(s *Sensor) { s.gnmic = c }
}

// WithSensorGRPCConn sets a gRPC connection on the Sensor.
func WithSensorGRPCConn(c *grpc.ClientConn) SensorOption {
	return func(s *Sensor) { s.grpcc = c }
}

// WithSensorConnector sets a gRPC connector
func WithSensorConnector(c GRPCConnector) SensorOption {
	return func(s *Sensor) { s.grpcConnector = c }
}

// WithSensorCredentialResolver sets a credential resolver
func WithSensorCredentialResolver(c CredentialResolver) SensorOption {
	return func(s *Sensor) { s.credResolver = c }
}

// WithSensorStandaloneStatus sets the stanalone status
func WithSensorStandaloneStatus(standalone bool) SensorOption {
	return func(s *Sensor) { s.standalone = standalone }
}

// WithSensorConfigChan provides a channel for supplying configs to the sensor
func WithSensorConfigChan(configCh chan *Config) SensorOption {
	return func(s *Sensor) { s.configCh = configCh }
}

// WithSensorClientFactory sets a cvclient factory on the Sensor.
func WithSensorClientFactory(f func(gnmi.GNMIClient,
	*Info) cvclient.CVClient) SensorOption {
	return func(s *Sensor) { s.clientFactory = f }
}

func (s *Sensor) handleConfigUpdate(ctx context.Context,
	resp *gnmi.Notification, postSync bool) error {
	// For each deleted datasource name, cancel that datasource and
	// delete it and its config from our collections.
	dsUpdated := map[string]map[string]bool{}
	// dsUpated records the datasources to be restart
	// key is the name of datasource, value is the map of changed fields.
	for _, p := range resp.Delete {
		fullPath := pgnmi.PathJoin(resp.Prefix, p)
		leafName := fullPath.Elem[len(fullPath.Elem)-1].Name
		if leafName == "id" || leafName == "sensor" {
			// If the sensor itself has been deleted, log but leave
			// the datasource deletion/addition until the relevant
			// leaf nodes are deleted.
			s.log.Infof("Sensor deleted: %v", fullPath)
			for name := range s.datasourceConfig {
				s.removeDatasource(ctx, name)
			}

			// Need to stop heartbeats before deleting state to avoid flap
			s.heartbeatLock.Lock()
			s.active = false
			s.heartbeatLock.Unlock()

			// Delete state for this sensor
			if _, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
				Delete: []*gnmi.Path{s.statePrefix},
			}); err != nil {
				return err
			}
		} else if leafName != "name" {
			name := datasourceFromPath(fullPath)
			dscfg, ok := s.datasourceConfig[name]
			if ok {
				curr := 4 // (0)datasource/(1)config/(2)sensor[id]/(3)source[name]/(4)fields
				elemNext := func() *gnmi.PathElem {
					if curr >= len(fullPath.Elem) {
						return undefinedPathElem
					}
					out := fullPath.Elem[curr]
					curr++
					return out
				}
				elem := elemNext()
				if _, ok := dsUpdated[name]; !ok {
					dsUpdated[name] = make(map[string]bool)
				}
				dsUpdated[name][elem.Name] = true
				switch elem.Name {
				case "credential":
					if k, ok := elem.Key["key"]; ok {
						if elemNext().Name == "value" {
							delete(dscfg.credential, k)
						}
					}
				case "option":
					if k, ok := elem.Key["key"]; ok {
						if elemNext().Name == "value" {
							delete(dscfg.option, k)
						}
					}
				default:
					delete(dsUpdated[name], elem.Name)
				}
			}
		} else {
			if name := datasourceFromPath(fullPath); name != "" {
				s.removeDatasource(ctx, name)
			}
		}
	}

	if len(resp.Update) > 0 && !s.active {
		s.heartbeatLock.Lock()
		s.active = true
		// Sync state, assuming no state needs to be deleted as there is not supposed
		// to be a sensor streaming when we notice a config creation after initial sync.
		if postSync {
			if err := s.syncState(ctx, nil); err != nil {
				s.log.Errorf("Failed to sync state on restart: %v", err)
			}
		}
		s.heartbeatLock.Unlock()

		// Restart all existing configs. This only happens if they are coming from
		// the cofigCh and reach the sensor before the sensor config.
		for name := range s.datasourceConfig {
			ds := s.getDatasource(ctx, name)
			ds.scheduleRestart(s.deviceRedeployTimer)
		}
	}

	// For each updated datasource, update the datasource config but
	// hold off on restarting the datasource.
	for _, upd := range resp.Update {
		fullPath := pgnmi.PathJoin(resp.Prefix, upd.Path)
		leafName := fullPath.Elem[len(fullPath.Elem)-1].Name
		if leafName == "id" {
			s.log.Infof("Sensor added: %v", fullPath)
		}
		name := datasourceFromPath(fullPath)
		if name == "" {
			continue
		}

		dscfg, ok := s.datasourceConfig[name]
		if !ok {
			dscfg = &datasourceConfig{
				name:       name,
				option:     map[string]string{},
				credential: map[string]string{},
				loglevel:   logrus.InfoLevel,
			}
			s.datasourceConfig[name] = dscfg
		}

		curr := 4 // (0)datasource/(1)config/(2)sensor[id]/(3)source[name]/(4)fields
		elemNext := func() *gnmi.PathElem {
			if curr >= len(fullPath.Elem) {
				return undefinedPathElem
			}
			out := fullPath.Elem[curr]
			curr++
			return out
		}

		elem := elemNext()
		if _, ok := dsUpdated[name]; !ok {
			dsUpdated[name] = make(map[string]bool)
		}
		dsUpdated[name][elem.Name] = true
		switch elem.Name {
		case "type":
			dscfg.typ = upd.Val.GetStringVal()
		case "enabled":
			dscfg.enabled = upd.Val.GetBoolVal()
		case "credential":
			if k, ok := elem.Key["key"]; ok {
				if elemNext().Name == "value" {
					dscfg.credential[k] = upd.Val.GetStringVal()
				}
			}
		case "option":
			if k, ok := elem.Key["key"]; ok {
				if elemNext().Name == "value" {
					dscfg.option[k] = upd.Val.GetStringVal()
				}
			}
		case "loglevel":
			level := upd.Val.GetStringVal()
			loglevel, ok := logMapping[level]
			if !ok {
				s.log.Errorf("unknown datasource loglevel for %s: %s,"+
					"setting it to INFO level", name, level)
				loglevel = logrus.InfoLevel
			}
			dscfg.loglevel = loglevel
		default:
			delete(dsUpdated[name], elem.Name)
		}
	}

	// If we are still on sync phase, nothing else to do.
	if !postSync {
		return nil
	}

	// Reset redeploy timer so we have a few moments to aggregate more changes
	for name := range dsUpdated {
		ds := s.getDatasource(ctx, name)
		if _, ok := dsUpdated[name]["loglevel"]; ok {
			ds.monitor.SetLoggerLevel(s.datasourceConfig[name].loglevel)
			if len(dsUpdated[name]) == 1 {
				// If only the log level changes in the datasource config, skip the restart.
				continue
			}
		}
		// We don't care if it already ran, we want to reschedule a run after the config changes.
		ds.scheduleRestart(s.deviceRedeployTimer)

	}

	return nil
}

func (s *Sensor) removeDatasource(ctx context.Context, name string) {
	ds, ok := s.datasource[name]
	if !ok {
		return
	}

	s.log.Debugf("Removing datasource: %s", name)
	ds.stop()
	delete(s.datasource, name)
	delete(s.datasourceConfig, name)
	if _, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
		Delete: []*gnmi.Path{ds.statePrefix},
	}); err != nil {
		s.log.Errorf("Failed to delete state for %s: %v", name, err)
	}
	s.log.Infof("Datasource removed: %s", name)
}

func (s *Sensor) getDatasource(ctx context.Context, name string) *datasource {
	if runtime, ok := s.datasource[name]; ok {
		return runtime
	}

	s.log.Infof("New datasource: %v", name)

	prefix := pgnmi.PathFromString(fmt.Sprintf(
		"datasource/state/sensor[id=%s]/source[name=%s]", s.id, name))
	prefix.Origin = "arista"
	prefix.Target = "cv"

	runtime := &datasource{
		log:           s.log.WithField("datasource", name),
		sensorID:      s.id,
		clientFactory: s.clientFactory,
		gnmic:         s.gnmic,
		apiaddr:       s.apiaddr,
		grpcConnector: s.grpcConnector,
		standalone:    s.standalone,
		credResolver:  s.credResolver,
		config: &datasourceConfig{
			name:     name,
			loglevel: logrus.InfoLevel,
		},
		// Create redeploy timer with 1h and stop it, so it only runs when we reset the timer.
		redeployTimer: time.AfterFunc(time.Hour, func() {
			// Queue it in the main loop to be updated.
			// This guarantees a single goroutine is playing with the datasource
			// configs and runtime fields.
			s.redeployDatasource <- name
		}),
		grpcc:             s.grpcc,
		statePrefix:       prefix,
		heartbeatInterval: s.heartbeatInterval,
		logRate:           s.logRate,
	}
	// Setup monitor for datasource
	runtime.monitor = newDatasourceMonitor(runtime.log, runtime.config.loglevel)
	// Setup fatal error retry backoff with long wait periods to avoid flood.
	runtime.failureRetryTimer = provider.NewBackoffTimer(
		provider.WithBackoffBase(s.failureRetryBackoffBase),
		provider.WithBackoffMax(s.failureRetryBackoffMax))

	_ = runtime.redeployTimer.Stop() // we don't want to run it yet.
	s.datasource[name] = runtime
	return runtime
}

func (s *Sensor) runDatasourceConfig(ctx context.Context, name string) error {
	cfg, ok := s.datasourceConfig[name]
	if !ok {
		return fmt.Errorf("config not found: %v", name)
	}

	runtime := s.getDatasource(ctx, name)
	return runtime.deploy(ctx, cfg)
}

type handleUpdateFn func(context.Context, *gnmi.Notification, bool) error
type handleSyncResponseFn func(context.Context) error

func (s *Sensor) subscribe(ctx context.Context, opts *agnmi.SubscribeOptions,
	handleUpdate handleUpdateFn, handleSyncResponse handleSyncResponseFn) error {
	s.log.Debugf("subscribe: %v", opts)

	respCh := make(chan *gnmi.SubscribeResponse)
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return agnmi.SubscribeErr(ctx, s.gnmic, opts, respCh)
	})
	errg.Go(func() error {
		// Main sensor loop, reading new configs or redeploying datasources after config changes.
		postSync := false
		for {
			select {
			case name := <-s.redeployDatasource:
				if err := s.runDatasourceConfig(ctx, name); err != nil {
					s.log.Errorf("redeploy failed: %v", err)
				}
			case resp, ok := <-respCh:
				if !ok {
					return nil
				}
				s.log.Tracef("Got response: %v", resp)
				switch subResp := resp.Response.(type) {
				case *gnmi.SubscribeResponse_Update:
					if err := handleUpdate(ctx, subResp.Update, postSync); err != nil {
						return err
					}
				case *gnmi.SubscribeResponse_SyncResponse:
					postSync = true
					if err := handleSyncResponse(ctx); err != nil {
						return err
					}
				}
			case cfg, ok := <-s.configCh:
				if !ok {
					s.log.Info("configCh closed, will not take anymore configs from it")
					s.configCh = nil // prevent spin
					continue
				}
				if cfg.IsDeleted() { // special case to know when config is removed
					s.removeDatasource(ctx, cfg.Name)
				} else {
					loglevel, ok := logMapping[cfg.LogLevel]
					if !ok {
						s.log.Errorf("unknown datasource loglevel for %s: %s,"+
							"setting it to INFO level", cfg.Name, cfg.LogLevel)
						loglevel = logrus.InfoLevel
					}
					s.datasourceConfig[cfg.Name] = &datasourceConfig{
						name:     cfg.Name,
						typ:      cfg.Device,
						enabled:  true,
						option:   cfg.Options,
						loglevel: loglevel,
					}
					ds := s.getDatasource(ctx, cfg.Name)
					if s.active {
						ds.scheduleRestart(s.deviceRedeployTimer)
					}
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	return errg.Wait()
}

func (s *Sensor) syncState(ctx context.Context, stateNames map[string]struct{}) error {
	s.log.Debugf("syncState found states: %v", stateNames)

	// Delete state that are not present in configs
	var toDelete []*gnmi.Path
	for name := range stateNames {
		// Only clean up datasources that are not configured
		if _, ok := s.datasourceConfig[name]; !ok {
			gnmiPath := pgnmi.PathFromString(fmt.Sprintf("source[name=%s]", name))
			toDelete = append(toDelete, gnmiPath)
		}
	}

	s.log.Infof("Resyncing state. Deleting: %v", toDelete)
	ts := time.Now().UnixNano()
	_, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
		Prefix: s.statePrefix,
		Delete: toDelete,
		Update: []*gnmi.Update{
			pgnmi.Update(pgnmi.Path("version"), agnmi.TypedValue(version.CollectorVersion)),
			pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(ts)),
			pgnmi.Update(lastSeenKey, agnmi.TypedValue(ts)),
			pgnmi.Update(lastErrorKey, agnmi.TypedValue("Sensor started")),
		},
	})
	return err
}

func (s *Sensor) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.heartbeatLock.Lock()
			if !s.active {
				s.heartbeatLock.Unlock()
				continue
			}
			ts := time.Now().UnixNano()
			_, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
				Prefix: s.statePrefix,
				Update: []*gnmi.Update{
					pgnmi.Update(lastSeenKey, agnmi.TypedValue(ts)),
				},
			})
			if err != nil {
				s.log.Errorf("Failed to publish heartbeat: %v", err)
			}
			s.heartbeatLock.Unlock()
		}
	}
}

// Run executes the sensor to start to manage datasources
func (s *Sensor) Run(ctx context.Context) error {
	s.log.Infof("Running sensor %q, version %v", s.id, version.CollectorVersion)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sensorConfigPath := fmt.Sprintf("/datasource/config/sensor[id=%s]", s.id)
	sensorStatePath := fmt.Sprintf("/datasource/state/sensor[id=%s]", s.id)

	// Compile a list of datasources in the state collection.
	stateDatasourceNames := map[string]struct{}{}
	handleStateUpdate := func(ctx context.Context,
		notif *gnmi.Notification, postSync bool) error {
		for _, upd := range notif.Update {
			fullPath := pgnmi.PathJoin(notif.Prefix, upd.Path)
			leafName := fullPath.Elem[len(fullPath.Elem)-1].Name
			if leafName == "name" {
				dsName := upd.Val.GetStringVal()
				stateDatasourceNames[dsName] = struct{}{}
			}
		}
		return nil
	}
	handleStateSync := func(ctx context.Context) error {
		s.log.Info("state sync complete")
		return nil
	}
	if err := s.subscribe(ctx,
		&agnmi.SubscribeOptions{
			Origin: "arista",
			Target: "cv",
			Paths:  [][]string{agnmi.SplitPath(sensorStatePath)},
			Mode:   "once",
		}, handleStateUpdate, handleStateSync); err != nil {
		return fmt.Errorf("state resync failed: %w", err)
	}

	// Start sensor heartbeat loop. This will wait to do work only when we have a valid config.
	go s.heartbeatLoop(ctx)

	// Subscribe to config forever.
	handleConfigUpdate := func(ctx context.Context,
		notif *gnmi.Notification, postSync bool) error {
		return s.handleConfigUpdate(ctx, notif, postSync)
	}
	handleConfigSync := func(ctx context.Context) error {
		// Update sensor and datasources state based on existing configs
		if !s.active {
			s.log.Info("No datasource config found before sync response, waiting for configs.")
			return nil
		}
		if err := s.syncState(ctx, stateDatasourceNames); err != nil {
			return err
		}
		// Run synced configs
		for name := range s.datasourceConfig {
			if err := s.runDatasourceConfig(ctx, name); err != nil {
				s.log.Errorf("resync run failed %s: %v", name, err)
			}
		}
		s.log.Info("Config sync complete")
		return nil
	}

	err := s.subscribe(ctx,
		&agnmi.SubscribeOptions{
			Origin: "arista",
			Target: "cv",
			Paths:  [][]string{agnmi.SplitPath(sensorConfigPath)},
			Mode:   "stream",
		}, handleConfigUpdate, handleConfigSync)
	if err != nil {
		s.log.Infof("Config subscription returned: error: %v", err)
	}

	s.log.Infof("Terminating %d datasources...", len(s.datasource))
	for _, ds := range s.datasource {
		ds.stop()
	}
	s.log.Infof("All datasources closed")

	return err
}

func (s *Sensor) validateOptions() {
	if s.gnmic == nil {
		s.log.Fatal("gNMI client must be set to start sensor")
	}
	//if s.grpcc == nil {
	//	s.log.Fatal("gRPC client must be set to start sensor")
	//}
	if s.clientFactory == nil {
		s.log.Fatal("factory must be set to start sensor")
	}
}

// NewSensor creates a new Sensor
func NewSensor(name string, logRate float64, opts ...SensorOption) *Sensor {
	prefix := pgnmi.PathFromString(fmt.Sprintf("datasource/state/sensor[id=%s]", name))
	prefix.Origin = "arista"
	prefix.Target = "cv"

	s := &Sensor{
		id:                      name,
		log:                     logrus.WithField("sensor", name),
		logRate:                 logRate,
		datasourceConfig:        map[string]*datasourceConfig{},
		datasource:              map[string]*datasource{},
		deviceRedeployTimer:     2 * time.Second,
		redeployDatasource:      make(chan string),
		statePrefix:             prefix,
		heartbeatInterval:       10 * time.Second, // default in case it is not set
		failureRetryBackoffBase: time.Minute,
		failureRetryBackoffMax:  1 * time.Hour,
		credResolver:            passthroughResolver,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.validateOptions()
	return s
}
