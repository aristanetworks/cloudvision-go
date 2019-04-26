// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

// An Inventory maintains a set of devices.
type Inventory interface {
	Add(device *Info) error
	Delete(key string) error
	Get(key string) (Device, bool)
	Update(devices []*Info) error
}

// deviceConn contains a device and its gNMI connections.
type deviceConn struct {
	deviceInfo        *Info
	ctx               context.Context
	cancel            context.CancelFunc
	rawGNMIClient     gnmi.GNMIClient
	wrappedGNMIClient *gNMIClientWrapper
	group             sync.WaitGroup
}

// inventory implements the Inventory interface.
type inventory struct {
	ctx           context.Context
	rawGNMIClient gnmi.GNMIClient
	devices       map[string]*deviceConn
	lock          sync.Mutex
}

func (dc *deviceConn) sendPeriodicUpdates() error {
	ticker := time.NewTicker(time.Second)
	ctx := metadata.AppendToOutgoingContext(dc.ctx,
		collectorVersionMetadata, version.Version)
	if _, ok := dc.deviceInfo.Device.(Manager); ok {
		// ManagementSystem is a system managing other devices which itself
		// shouldn't be treated as an actual streaming device in CloudVision.
		ctx = metadata.AppendToOutgoingContext(ctx,
			deviceTypeMetadata, "managementSystem")
	} else {
		// Target is an ordinary device streaming to CloudVision.
		ctx = metadata.AppendToOutgoingContext(ctx,
			deviceTypeMetadata, "target")
	}
	dc.wrappedGNMIClient.Set(ctx, &gnmi.SetRequest{})
	for {
		select {
		case <-dc.ctx.Done():
			return nil
		case <-ticker.C:
			if alive, err := dc.deviceInfo.Device.Alive(); err == nil {
				if alive {
					ctx := metadata.AppendToOutgoingContext(dc.ctx,
						deviceLivenessMetadata, "true")
					dc.wrappedGNMIClient.Set(ctx, &gnmi.SetRequest{})
				} else {
					did, _ := dc.deviceInfo.Device.DeviceID()
					log.Log(dc.deviceInfo.Device).Infof("Device %s is not alive", did)
				}
			} else {
				return err
			}
		}
	}
}

// Add adds a device to the inventory, opens up any gNMI connections
// required by the device's providers, and then starts its providers.
func (i *inventory) Add(device *Info) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.devices[device.ID]; ok {
		return nil
	}

	dc := &deviceConn{deviceInfo: device}
	dc.ctx, dc.cancel = context.WithCancel(i.ctx)

	i.devices[device.ID] = dc

	providers, err := device.Device.Providers()
	if err != nil {
		return err
	}
	dc.rawGNMIClient = i.rawGNMIClient
	dc.wrappedGNMIClient = newGNMIClientWrapper(dc.rawGNMIClient, nil,
		device.ID, device.Config.Address, false)

	logFileName := device.ID + ".log"
	err = log.InitLogging(logFileName, device.Device)
	if err != nil {
		return fmt.Errorf("Error setting up logging for device %s: %v", device.ID, err)
	}

	for _, p := range providers {
		err = log.InitLogging(logFileName, p)
		if err != nil {
			return fmt.Errorf("Error setting up logging for provider %#v: %v", p, err)
		}

		pt, ok := p.(provider.GNMIProvider)
		if !ok {
			return errors.New("unexpected provider type; need GNMIProvider")
		}

		pt.InitGNMI(newGNMIClientWrapper(dc.rawGNMIClient, pt,
			device.ID, device.Config.Address, pt.OpenConfig()))

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

	// Send periodic updates of device-level metadata.
	dc.group.Add(1)
	go func() {
		err := dc.sendPeriodicUpdates()
		if err != nil {
			log.Log(device.Device).Errorf("Error updating device metadata: %v", err)
		}
		dc.group.Done()
	}()

	log.Log(device.Device).Infof("Added device %s", device.ID)
	return nil
}

func (i *inventory) Delete(key string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	dc, ok := i.devices[key]
	if !ok {
		return nil
	}

	// Cancel the device context and delete the device from the device
	// map. We need to make sure this device's providers are finished
	// before deleting the device.
	dc.cancel()
	dc.group.Wait()
	delete(i.devices, key)
	log.Log(dc.deviceInfo.Device).Infof("Deleted device %s", key)
	return nil
}

func (i *inventory) Get(key string) (Device, bool) {
	i.lock.Lock()
	defer i.lock.Unlock()
	d, ok := i.devices[key]
	if !ok {
		return nil, ok
	}
	return d.deviceInfo.Device, ok
}

func (i *inventory) Update(devices []*Info) error {
	idToDevice := map[string]Device{}
	for _, info := range devices {
		err := i.Add(info)
		if err != nil {
			return err
		}
		idToDevice[info.ID] = info.Device
	}
	for id := range i.devices {
		if _, ok := idToDevice[id]; !ok {
			err := i.Delete(id)
			if err != nil {
				return err
			}
		}
	}
	if DeviceIDFile != "" {
		err := DumpDeviceIDs(devices, DeviceIDFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewInventory creates an Inventory.
func NewInventory(ctx context.Context, gnmiClient gnmi.GNMIClient) Inventory {
	inv := &inventory{
		ctx:           ctx,
		devices:       make(map[string]*deviceConn),
		rawGNMIClient: gnmiClient,
	}
	return inv
}
