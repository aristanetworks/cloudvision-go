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
	"github.com/aristanetworks/cloudvision-go/log"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	datasourceLastErrorKey = pgnmi.PathFromString("last-error")
	undefinedPathElem      = &gnmi.PathElem{Name: "___undefined___"}
)

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
	clientFactory func(gnmi.GNMIClient, *Info) cvclient.CVClient

	// Current running config. Config changes require a datasource restart.
	config *datasourceConfig

	redeployTimer *time.Timer
	cancel        context.CancelFunc
	runDone       chan struct{}

	info     *Info
	grpcc    *grpc.ClientConn
	cvClient cvclient.CVClient

	// Holds datasource/state/sensor[id=sensor]/source[name=datasource]
	statePrefix *gnmi.Path
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

func (d *datasource) handleDatasourceError(ctx context.Context, e error) error {
	d.log.Error(e)
	return d.submitDatasourceUpdates(ctx,
		pgnmi.Update(datasourceLastErrorKey, agnmi.TypedValue(e.Error())))
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
	d.stop(ctx)
	d.config = cfg.clone()

	if !cfg.enabled {
		d.log.Infof("Datasource %v disabled", d.config.name)
		return nil
	}
	d.log.Info("Starting run")

	// Prepare device to execute based on datasource config
	info, err := NewDeviceInfo(&Config{
		Device:  cfg.typ,
		Options: mergeOpts(cfg.option, cfg.credential),
	})
	if err != nil {
		return d.handleDatasourceError(ctx, err)
	}
	d.cvClient = d.clientFactory(d.gnmic, info)
	d.info = info

	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel
	done := make(chan struct{})
	d.runDone = done

	go func() {
		defer close(done)
		err := d.Run(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			if pubErr := d.handleDatasourceError(ctx,
				fmt.Errorf("stopped running: %w", err)); pubErr != nil {
				d.log.Errorf("Unable to publish failure %v: %v", err, pubErr)
			}
		}
	}()

	return nil
}

// ctx is used to prevent this from blocking indefinitely if the datasource fails to stop the
// underlying device.
func (d *datasource) stop(ctx context.Context) {
	if d.redeployTimer != nil {
		_ = d.redeployTimer.Stop()
	}
	if d.cancel != nil {
		d.cancel()
	}
	select {
	case <-d.runDone:
	case <-ctx.Done():
	}
}

// Run executes the datasource
func (d *datasource) Run(ctx context.Context) error {
	errg, ctx := errgroup.WithContext(ctx)

	// Create device connection object.
	dc := &deviceConn{
		cvClient: d.cvClient,
		grpcConn: d.grpcc,
		info:     d.info,
		ctx:      ctx,
	}

	// Register the device before starting providers. If we can't reach
	// the device right now, we should return an error rather than
	// considering it added.
	if err := dc.cvClient.SendDeviceMetadata(dc.ctx); err != nil {
		return fmt.Errorf("Error sending device metadata for device "+
			"%q (%s): %w", dc.info.ID, dc.info.Config.Device, err)
	}

	// Start providers.
	if err := dc.runProviders(); err != nil {
		return fmt.Errorf("Error starting providers for device %q (%s): %w",
			dc.info.ID, dc.info.Config.Device, err)
	}

	// Send metadata updates.
	// XXX TODO: send datasource state updates
	errg.Go(func() error {
		err := dc.sendPeriodicUpdates()
		if err != nil {
			log.Log(dc.info.Device).Errorf("Error updating device metadata: %v", err)
		}
		return err
	})

	// Start manager, maybe.
	if manager, ok := dc.info.Device.(Manager); ok {
		inv := &sensorInventory{
			device: []Device{},
		}
		errg.Go(func() error {
			err := manager.Manage(dc.ctx, inv)
			if err != nil {
				log.Log(dc.info.Device).Errorf("Error in manager.Manage: %v", err)
			}
			return err
		})
	}

	dc.group.Wait()
	return errg.Wait()
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
	id    string
	gnmic gnmi.GNMIClient
	grpcc *grpc.ClientConn

	redeployDatasource chan string
	datasourceConfig   map[string]*datasourceConfig
	datasource         map[string]*datasource
	clientFactory      func(gnmi.GNMIClient, *Info) cvclient.CVClient

	deviceRedeployTimer time.Duration

	log *logrus.Entry
}

// SensorOption is used to configure the Sensor
type SensorOption func(m *Sensor)

// WithSensorGNMIClient sets a gNMI client on the Sensor.
func WithSensorGNMIClient(c gnmi.GNMIClient) SensorOption {
	return func(s *Sensor) { s.gnmic = c }
}

// WithSensorGRPCConn sets a gRPC connection on the Sensor.
func WithSensorGRPCConn(c *grpc.ClientConn) SensorOption {
	return func(s *Sensor) { s.grpcc = c }
}

// WithSensorClientFactory sets a cvclient factory on the Sensor.
func WithSensorClientFactory(f func(gnmi.GNMIClient,
	*Info) cvclient.CVClient) SensorOption {
	return func(s *Sensor) { s.clientFactory = f }
}

func (s *Sensor) handleConfigUpdate(ctx context.Context,
	resp *gnmi.Notification, postSync bool, execute func(deviceFn func() error)) error {
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
				ds.stop(ctx)
				delete(s.datasource, name)
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
	closedCh := make(chan struct{})
	close(closedCh)
	runtime := &datasource{
		log:           s.log.WithField("datasource", name),
		sensorID:      s.id,
		clientFactory: s.clientFactory,
		gnmic:         s.gnmic,
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
		grpcc:   s.grpcc,
		runDone: closedCh,
		statePrefix: pgnmi.PathFromString(fmt.Sprintf(
			"datasource/state/sensor[id=%s]/source[name=%s]", s.id, name)),
	}
	_ = runtime.redeployTimer.Stop() // we don't want to run it yet.
	s.datasource[name] = runtime
	return runtime
}

func (s *Sensor) runDatasourceConfig(ctx context.Context, name string) error {

	// Copy the current configuration to send it to the device.
	// The device should keep its own copy so that we can compare against changes and
	// decide if we need to restart the datasource.
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

func (s *Sensor) cleanUpState(ctx context.Context, names map[string]struct{}) error {
	s.log.Debugf("cleanUpState found states: %v", names)

	var toDelete []*gnmi.Path
	for name := range names {
		// Only clean up datasources that are not configured
		if _, ok := s.datasourceConfig[name]; !ok {
			gnmiPath := pgnmi.PathFromString(fmt.Sprintf("source[name=%s]", name))
			toDelete = append(toDelete, gnmiPath)
		}
	}

	if len(toDelete) > 0 {
		s.log.Infof("cleaning up old states: %v", toDelete)
		prefix := pgnmi.PathFromString(
			fmt.Sprintf("/datasource/state/sensor[id=%s]", s.id),
		)
		prefix.Origin = "arista"
		prefix.Target = "cv"
		if _, err := s.gnmic.Set(ctx, &gnmi.SetRequest{
			Prefix: prefix,
			Delete: toDelete,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Run executes the sensor to start to manage datasources
func (s *Sensor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	runningDevices, _ := errgroup.WithContext(ctx)
	runDevice := func(deviceFn func() error) {
		runningDevices.Go(func() error {
			err := deviceFn()
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
			return nil
		})
	}

	sensorPathStr := "/datasource/%s/sensor[id=%s]"
	sensorConfigPath := fmt.Sprintf(sensorPathStr, "config", s.id)
	sensorStatePath := fmt.Sprintf(sensorPathStr, "state", s.id)

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

	// XXX TODO: send periodic state updates for sensor

	// Subscribe to config forever.
	handleConfigUpdate := func(ctx context.Context,
		notif *gnmi.Notification, postSync bool) error {
		return s.handleConfigUpdate(ctx, notif, postSync, runDevice)
	}
	handleConfigSync := func(ctx context.Context) error {
		if err := s.cleanUpState(ctx, stateDatasourceNames); err != nil {
			return err
		}
		// Run synced configs
		for name := range s.datasourceConfig {
			if err := s.runDatasourceConfig(ctx, name); err != nil {
				s.log.Errorf("resync run failed %s: %v", name, err)
			}
		}
		s.log.Info("config sync complete")
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
		cancel()
		s.log.Infof("config subscription returned: error: %v", err)
		_ = runningDevices.Wait()
		return err
	}

	return runningDevices.Wait()
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
	log := logrus.WithField("sensor", name)
	s := &Sensor{
		id:                  name,
		log:                 log,
		datasourceConfig:    map[string]*datasourceConfig{},
		datasource:          map[string]*datasource{},
		deviceRedeployTimer: 2 * time.Second,
		redeployDatasource:  make(chan string),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.validateOptions()
	return s
}
