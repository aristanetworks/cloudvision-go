// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/aristanetworks/glog"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

// An Inventory maintains a set of devices.
type Inventory interface {
	Add(key string, device Device) error
	Delete(key string) error
	Get(key string) (Device, bool)
}

// deviceConn contains a device and its gNMI connections.
type deviceConn struct {
	device            Device
	deviceID          string
	version           string
	deviceType        string
	ctx               context.Context
	cancel            context.CancelFunc
	rawGNMIClient     gnmi.GNMIClient
	wrappedGNMIClient *gNMIClientWrapper
	providerGroup     *errgroup.Group
}

// inventory implements the Inventory interface.
type inventory struct {
	ctx            context.Context
	group          *errgroup.Group
	gnmiServerAddr string
	devices        map[string]*deviceConn
}

func startGNMIClient(serverAddr string) (gnmi.GNMIClient, error) {
	if serverAddr == "" {
		return nil, fmt.Errorf("Invalid gNMI server address '%v'", serverAddr)
	}
	return agnmi.Dial(&agnmi.Config{Addr: serverAddr})
}

func (dc *deviceConn) sendPeriodicUpdates() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-dc.ctx.Done():
			return nil
		case <-ticker.C:
			if dc.deviceID == "" {
				var err error
				dc.deviceID, err = dc.device.DeviceID()
				if err != nil {
					return err
				}
			}
			if dc.deviceType == "" {
				dc.deviceType = dc.device.Type().String()
			}
			ctx := metadata.AppendToOutgoingContext(dc.ctx,
				deviceTypeMetadata, dc.deviceType,
				deviceLivenessMetadata, "true",
				collectorVersionMetadata, version.Version)
			dc.wrappedGNMIClient.Set(ctx, &gnmi.SetRequest{})
		}
	}
}

func (dc *deviceConn) handleErrors() error {
	return dc.providerGroup.Wait()
}

// Add adds a device to the inventory, opens up any gNMI connections
// required by the device's providers, and then starts its providers.
func (i *inventory) Add(key string, device Device) error {
	if _, ok := i.devices[key]; ok {
		return nil
	}

	dc := &deviceConn{device: device}
	ctx, cancel := context.WithCancel(i.ctx)
	dc.providerGroup, dc.ctx = errgroup.WithContext(ctx)
	dc.cancel = cancel

	i.devices[key] = dc

	providers, err := device.Providers()
	if err != nil {
		return err
	}

	dc.rawGNMIClient, err = startGNMIClient(i.gnmiServerAddr)
	if err != nil {
		return err
	}
	dc.wrappedGNMIClient = newGNMIClientWrapper(dc.rawGNMIClient,
		key, false)

	for _, p := range providers {
		pt, ok := p.(provider.GNMIProvider)
		if !ok {
			return errors.New("unexpected provider type; need GNMIProvider")
		}

		pt.InitGNMI(newGNMIClientWrapper(dc.rawGNMIClient, key, pt.OpenConfig()))

		// Watch for provider errors in the provider errgroup and
		// propagate them up to the inventory errgroup.
		i.group.Go(func() error {
			return dc.handleErrors()
		})

		// Start the providers.
		dc.providerGroup.Go(func() error {
			return p.Run(dc.ctx)
		})
	}

	// Send periodic updates of device-level metadata.
	i.group.Go(func() error {
		return dc.sendPeriodicUpdates()
	})

	glog.V(2).Infof("Added device %s", key)
	return nil
}

func (i *inventory) Delete(key string) error {
	dc, ok := i.devices[key]
	if !ok {
		return nil
	}

	// Cancel the device context and delete the device from the device
	// map. We don't have to worry about propagating errors up to the
	// inventory errgroup, since handleErrors will do that. We just need
	// to make sure this device's providers are finished before deleting
	// the device.
	dc.cancel()
	_ = dc.providerGroup.Wait()
	delete(i.devices, key)
	glog.V(2).Infof("Deleted device %s", key)
	return nil
}

func (i *inventory) Get(key string) (Device, bool) {
	d, ok := i.devices[key]
	if !ok {
		return nil, ok
	}
	return d.device, ok
}

// NewInventory creates an Inventory.
func NewInventory(ctx context.Context, group *errgroup.Group,
	gnmiServerAddr string) Inventory {
	return &inventory{
		ctx:            ctx,
		devices:        make(map[string]*deviceConn),
		group:          group,
		gnmiServerAddr: gnmiServerAddr,
	}
}
