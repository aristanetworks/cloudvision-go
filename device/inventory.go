// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"

	"github.com/openconfig/gnmi/proto/gnmi"

	"google.golang.org/grpc"
)

const heartbeatInterval = 10 * time.Second

// An Inventory maintains a set of devices.
type Inventory interface {
	Add(deviceInfo *Info) error
	Delete(key string) error
	Get(key string) (*Info, error)
	List() []*Info
}

// InventoryOption configures how we create the Inventory.
type InventoryOption func(*inventory)

// deviceConn contains a device and its gNMI/gRPC connections.
type deviceConn struct {
	info     *Info
	ctx      context.Context
	cancel   context.CancelFunc
	cvClient cvclient.CVClient
	grpcConn *grpc.ClientConn
	group    sync.WaitGroup
}

// inventory implements the Inventory interface.
type inventory struct {
	ctx            context.Context
	rawGNMIClient  gnmi.GNMIClient
	grpcConn       *grpc.ClientConn
	grpcServerAddr string
	grpcConnector  GRPCConnector // Connector to get gRPC connection
	devices        map[string]*deviceConn
	lock           sync.Mutex
	clientFactory  func(gnmi.GNMIClient, *Info) cvclient.CVClient
}

func (dc *deviceConn) sendPeriodicUpdates() error {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	did, _ := dc.info.Device.DeviceID()

	logger := log.Log(dc.info.Device)
	wasFailing := false // used to only log once when device is unhealthy and back alive

	for {
		select {
		case <-dc.ctx.Done():
			return nil
		case <-ticker.C:
			alive, err := dc.info.Device.Alive()
			if err == nil && alive {
				if wasFailing {
					logger.Infof("Device %s is back alive", did)
					wasFailing = false
				}
				if err := dc.cvClient.SendHeartbeat(dc.ctx, alive); err != nil {
					// Don't give up if an update fails for some reason.
					logger.Infof("Error sending periodic update for device %v: %v",
						did, err)
				}
			} else {
				if !wasFailing {
					logger.Infof("Device %s is not alive, err: %v", did, err)
					wasFailing = true
				}
			}
		}
	}
}

func (i *inventory) newDeviceConn(info *Info) (*deviceConn, error) {
	dc := &deviceConn{
		cvClient: i.clientFactory(i.rawGNMIClient, info),
		info:     info,
	}

	// Take any metadata associated with the device context.
	if info.Context != nil {
		dc.ctx, dc.cancel = context.WithCancel(info.Context)
	} else {
		dc.ctx, dc.cancel = context.WithCancel(i.ctx)
	}

	// i.grpcConnector is set,
	// only if grpcServerAddr is provided
	if i.grpcConnector != nil {
		cc := GRPCConnectorConfig{
			dc.info.ID,
		}
		conn, err := i.grpcConnector.Connect(dc.ctx, i.grpcConn, i.grpcServerAddr, cc)
		if err != nil {
			return nil, fmt.Errorf("gRPC connection to device %v failed: %w", cc.DeviceID, err)
		}
		dc.grpcConn = conn
	}
	return dc, nil
}

func (dc *deviceConn) runProviders() error {
	providers, err := dc.info.Device.Providers()
	if err != nil {
		return err
	}
	logFileName := dc.info.ID + ".log"
	err = log.InitLogging(logFileName, dc.info.Device)
	if err != nil {
		return fmt.Errorf("Error setting up logging for device %s: %v", dc.info.ID, err)
	}

	for _, p := range providers {
		err = log.InitLogging(logFileName, p)
		if err != nil {
			return fmt.Errorf("Error setting up logging for provider %#v: %v", p, err)
		}

		switch pt := p.(type) {
		case provider.GNMIProvider:
			pt.InitGNMI(dc.cvClient.ForProvider(pt))
		case provider.GRPCProvider:
			pt.InitGRPC(dc.grpcConn)
		default:
			return fmt.Errorf("unexpected provider type %T", p)
		}

		// Start the providers.
		dc.group.Add(1)
		go func(p provider.Provider) {
			err := p.Run(dc.ctx)
			if err != nil {
				log.Log(p).Errorf("Provider exiting with error %v", err)
			}
			dc.group.Done()
		}(p)
	}
	return nil
}

// Add adds a device to the inventory, opens up any gNMI connections
// required by the device's providers, and then starts its providers.
func (i *inventory) Add(info *Info) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	if info.ID == "" {
		return errors.New("ID in device.Info cannot be empty")
	} else if info.Config == nil {
		return errors.New("Config in device.Info cannot be empty")
	}
	if dev, ok := i.devices[info.ID]; ok {
		log.Log(info.Device).Debugf("Replacing device %s (type %s)",
			info.ID, info.Config.Device)
		dev.cancel()
		delete(i.devices, info.ID)
	}

	// Create device connection object.
	dc, err := i.newDeviceConn(info)
	if err != nil {
		return err
	}

	// Register the device before starting providers. If we can't reach
	// the device right now, we should return an error rather than
	// considering it added.
	if err := dc.cvClient.SendDeviceMetadata(dc.ctx); err != nil {
		return fmt.Errorf("Error sending device metadata for device "+
			"%q (%s): %w", info.ID, info.Config.Device, err)
	}

	// Start providers.
	if err := dc.runProviders(); err != nil {
		return fmt.Errorf("Error starting providers for device %q (%s): %w",
			info.ID, info.Config.Device, err)
	}

	// We're connected to the device, have told CloudVision about the
	// device, and are streaming the device's data now, so add the
	// device to the inventory.
	i.devices[info.ID] = dc

	// Send periodic updates of device-level metadata.
	if !info.Config.NoStream {
		dc.group.Add(1)
		go func() {
			err := dc.sendPeriodicUpdates()
			if err != nil {
				log.Log(info.Device).Errorf("Error updating device metadata: %v", err)
			}
			dc.group.Done()
		}()
	}

	if manager, ok := info.Device.(Manager); ok {
		dc.group.Add(1)
		go func() {
			err := manager.Manage(dc.ctx, i)
			if err != nil {
				log.Log(info.Device).Errorf("Error in manager.Manage: %v", err)
			}
			dc.group.Done()
		}()
	}

	log.Log(info.Device).Infof("Added device %q (%s)", info.ID,
		info.Config.Device)
	return nil
}

func (i *inventory) Delete(key string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	if key == "" {
		return fmt.Errorf("key in inventory.Delete cannot be empty")
	}
	dc, ok := i.devices[key]
	if !ok {
		return nil
	}

	// Cancel the device context and delete the device from the device
	// map. We need to make sure this device's providers are finished
	// before deleting the device. We also need to make sure Manager device
	// has manage go routine closed too.
	dc.cancel()
	dc.group.Wait()
	delete(i.devices, key)
	log.Log(dc.info.Device).Infof("Deleted device %s", key)
	return nil
}

func (i *inventory) Get(key string) (*Info, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if key == "" {
		return nil, fmt.Errorf("key in inventory.Get cannot be empty")
	}
	d, ok := i.devices[key]
	if !ok {
		return nil, fmt.Errorf("Device %s not found", key)
	}
	return d.info, nil
}

func (i *inventory) List() []*Info {
	var ret []*Info
	for _, conn := range i.devices {
		ret = append(ret, conn.info)
	}
	return ret
}

// WithGNMIClient sets a gNMI client on the Inventory.
func WithGNMIClient(c gnmi.GNMIClient) InventoryOption {
	return func(i *inventory) {
		i.rawGNMIClient = c
	}
}

// WithGRPCConn sets a gRPC connection on the Inventory.
func WithGRPCConn(c *grpc.ClientConn) InventoryOption {
	return func(i *inventory) {
		i.grpcConn = c
	}
}

// WithGRPCServerAddr sets a gRPC connection on the Inventory.
func WithGRPCServerAddr(addr string) InventoryOption {
	return func(i *inventory) {
		i.grpcServerAddr = addr
	}
}

// WithGRPCConnector sets a gRPC connector on the Inventory.
func WithGRPCConnector(c GRPCConnector) InventoryOption {
	return func(i *inventory) {
		i.grpcConnector = c
	}
}

// WithClientFactory sets a client factory on the Inventory.
func WithClientFactory(
	f func(gnmi.GNMIClient, *Info) cvclient.CVClient) InventoryOption {
	return func(i *inventory) {
		i.clientFactory = f
	}
}

// NewInventoryWithOptions creates an Inventory with the supplied options.
func NewInventoryWithOptions(ctx context.Context,
	options ...InventoryOption) Inventory {
	inv := &inventory{
		ctx:     ctx,
		devices: make(map[string]*deviceConn),
	}
	for _, opt := range options {
		opt(inv)
	}
	return inv
}

// NewInventory creates an Inventory.
// Deprecated: Use NewInventoryWithOptions instead.
func NewInventory(ctx context.Context, gnmiClient gnmi.GNMIClient,
	clientFactory func(gnmi.GNMIClient, *Info) cvclient.CVClient) Inventory {
	return NewInventoryWithOptions(ctx,
		WithGNMIClient(gnmiClient),
		WithClientFactory(clientFactory))
}
