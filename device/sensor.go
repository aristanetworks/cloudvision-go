// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/version"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

type datasourceConfig struct {
	name       string
	typ        string
	enabled    bool
	option     map[string]string
	credential map[string]string
}

func (c datasourceConfig) clone() *datasourceConfig {
	return &datasourceConfig{
		name:       c.name,
		typ:        c.typ,
		enabled:    c.enabled,
		option:     cloneMap(c.option),
		credential: cloneMap(c.credential),
	}
}

func (c datasourceConfig) String() string {
	return fmt.Sprintf("name: %s, typ: %s, enabled: %t, option: %v",
		c.name, c.typ, c.enabled, c.option)
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

	// Current running config. Config changes require a datasource restart.
	config *datasourceConfig

	redeployTimer *time.Timer
	cancel        context.CancelFunc
	execGroup     *errgroup.Group

	info     *Info
	grpcc    *grpc.ClientConn
	cvClient cvclient.CVClient

	// Holds datasource/state/sensor[id=sensor]/source[name=datasource]
	statePrefix       *gnmi.Path
	heartbeatInterval time.Duration
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
		d.log.Errorf("Failed to publish error: %v", err)
	}
}

func (d *datasource) scheduleRestart(in time.Duration) {
	// If was already triggered we are ok in triggering it again
	_ = d.redeployTimer.Reset(in)
}

func (d *datasource) deploy(ctx context.Context, cfg *datasourceConfig) error {
	d.log.Tracef("Trying to deploy config: %v. Current: %v", cfg, d.config)

	// We don't need to re-run if we already are running the same config.
	if cfg.equals(d.config) {
		d.log.Info("Deploy requested with the same config, ignoring request.")
		return nil
	}

	// Start processing new config
	d.stop()
	d.config = cfg.clone()

	if !cfg.enabled {
		d.log.Infof("Datasource %v disabled", d.config.name)
		if err := d.submitDatasourceUpdates(ctx,
			pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(false))); err != nil {
			d.log.Error("Failed to publish initial status:", err)
		}
		return nil
	}
	d.log.Info("Starting run")

	ctx, cancel := context.WithCancel(ctx)
	d.execGroup, ctx = errgroup.WithContext(ctx)
	d.cancel = cancel

	// Start running the device in the background
	d.execGroup.Go(func() error {
		err := d.Run(ctx)

		// Handle return error by pushing it to datasource last-error state
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) {
			return err
		}

		err = fmt.Errorf("Datasource stopped unexpectedly: %w", err)
		d.handleDatasourceError(ctx, err)
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

// Run executes the datasource
func (d *datasource) Run(ctx context.Context) (err error) {
	// Submit initial status information initially as the next operations can be slow
	ts := time.Now().UnixNano()
	if err := d.submitDatasourceUpdates(ctx,
		// last-seen is not sent as to not qualify the datasource as streaming yet
		pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
		pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue(d.config.typ)),
		pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(ts)),
		pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
	); err != nil {
		return err
	}

	// Prepare device to execute based on datasource config
	info, err := NewDeviceInfo(&Config{
		Device:  d.config.typ,
		Options: mergeOpts(d.config.option, d.config.credential),
	})
	if err != nil {
		return err
	}
	d.cvClient = d.clientFactory(d.gnmic, info)
	d.info = info

	deviceID, err := d.info.Device.DeviceID()
	var updates []*gnmi.Update
	if err != nil {
		errStr := fmt.Sprintf("failed to determine DeviceID: %v", err)
		updates = append(updates, pgnmi.Update(lastErrorKey, agnmi.TypedValue(errStr)))
	} else if len(deviceID) > 0 {
		updates = append(updates, pgnmi.Update(lastSeenKey, agnmi.TypedValue(ts)))
		updates = append(updates, pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue(deviceID)))
	}

	if err := d.submitDatasourceUpdates(ctx, updates...); err != nil {
		return fmt.Errorf("failed to publish startup status: %w", err)
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

	// Start manager, maybe.
	if manager, ok := d.info.Device.(Manager); ok {
		inv := &sensorInventory{
			device: []Device{},
		}
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

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			alive, err := d.info.Device.Alive()
			if err == nil && alive {
				if wasFailing {
					d.log.Info("Device is back alive")
					wasFailing = false
				}
				ts := time.Now().UnixNano()
				if err := d.submitDatasourceUpdates(ctx,
					pgnmi.Update(lastSeenKey, agnmi.TypedValue(ts))); err != nil {
					d.log.Error("Publish status failed:", err)
				}
				if err := d.cvClient.SendHeartbeat(ctx, alive); err != nil {
					// Don't give up if an update fails for some reason.
					d.log.Errorf("Error sending heartbeat: %v", err)
				}
			} else if !wasFailing {
				d.handleDatasourceError(ctx, fmt.Errorf("Device not alive: %w", err))
				wasFailing = true
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

type sensorInventory struct {
	device []Device
}

func (s *sensorInventory) Add(deviceInfo *Info) error {
	// TODO
	return nil
}

func (s *sensorInventory) Delete(key string) error {
	// TODO
	return nil
}

func (s *sensorInventory) Get(key string) (*Info, error) {
	// TODO
	return nil, nil
}

func (s *sensorInventory) List() []*Info {
	// TODO
	return nil
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

	redeployDatasource chan string
	datasourceConfig   map[string]*datasourceConfig
	datasource         map[string]*datasource
	clientFactory      func(gnmi.GNMIClient, *Info) cvclient.CVClient

	deviceRedeployTimer time.Duration

	log         *logrus.Entry
	statePrefix *gnmi.Path
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

// WithSensorStandaloneStatus sets the stanalone status
func WithSensorStandaloneStatus(standalone bool) SensorOption {
	return func(s *Sensor) { s.standalone = standalone }
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
	for _, p := range resp.Delete {
		fullPath := pgnmi.PathJoin(resp.Prefix, p)
		leafName := fullPath.Elem[len(fullPath.Elem)-1].Name
		if leafName == "id" {
			// If the sensor itself has been deleted, log but leave
			// the datasource deletion/addition until the relevant
			// leaf nodes are deleted.
			s.log.Infof("Sensor deleted: %v", fullPath)
		} else if leafName != "name" {
			continue
		}
		if name := datasourceFromPath(fullPath); name != "" {
			if ds, ok := s.datasource[name]; ok {
				s.log.Infof("Config removed: %s", name)
				ds.stop()
				delete(s.datasource, name)
				_, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
					Delete: []*gnmi.Path{
						ds.statePrefix,
					},
				})
				if err != nil {
					s.log.Errorf("Failed to delete state for %s: %v", name, err)
				}
			}
			delete(s.datasourceConfig, name)
		}
	}

	// For each updated datasource, update the datasource config but
	// hold off on restarting the datasource.
	dsUpdated := map[string]struct{}{}
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
		dsUpdated[name] = struct{}{}
		dscfg, ok := s.datasourceConfig[name]
		if !ok {
			dscfg = &datasourceConfig{
				name:       name,
				option:     map[string]string{},
				credential: map[string]string{},
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

		switch elem := elemNext(); elem.Name {
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
		}
	}

	// If we are still on sync phase, nothing else to do.
	if !postSync {
		return nil
	}

	// Reset redeploy timer so we have a few moments to aggregate more changes
	for name := range dsUpdated {
		ds := s.getDatasource(ctx, name)
		// We don't care if it already ran, we want to reschedule a run after the config changes.
		ds.scheduleRestart(s.deviceRedeployTimer)
	}

	return nil
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
		config: &datasourceConfig{
			name: name,
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
	}
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

	s.log.Infof("resyncing state. Deleting: %v", toDelete)
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

	// Subscribe to config forever.
	handleConfigUpdate := func(ctx context.Context,
		notif *gnmi.Notification, postSync bool) error {
		return s.handleConfigUpdate(ctx, notif, postSync)
	}
	handleConfigSync := func(ctx context.Context) error {
		// Update sensor and datasources state based on existing configs
		if err := s.syncState(ctx, stateDatasourceNames); err != nil {
			return err
		}
		// Run synced configs
		for name := range s.datasourceConfig {
			if err := s.runDatasourceConfig(ctx, name); err != nil {
				s.log.Errorf("resync run failed %s: %v", name, err)
			}
		}
		s.log.Info("config sync complete")
		go s.heartbeatLoop(ctx)
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
		s.log.Infof("config subscription returned: error: %v", err)
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
func NewSensor(name string, opts ...SensorOption) *Sensor {
	prefix := pgnmi.PathFromString(fmt.Sprintf("datasource/state/sensor[id=%s]", name))
	prefix.Origin = "arista"
	prefix.Target = "cv"

	s := &Sensor{
		id:                  name,
		log:                 logrus.WithField("sensor", name),
		datasourceConfig:    map[string]*datasourceConfig{},
		datasource:          map[string]*datasource{},
		deviceRedeployTimer: 2 * time.Second,
		redeployDatasource:  make(chan string),
		statePrefix:         prefix,
		heartbeatInterval:   10 * time.Second, // default in case it is not set
	}
	for _, opt := range opts {
		opt(s)
	}
	s.validateOptions()
	return s
}
