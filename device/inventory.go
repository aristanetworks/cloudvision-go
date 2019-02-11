// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"

	"github.com/aristanetworks/cloudvision-go/provider"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// An Inventory maintains a set of devices.
type Inventory interface {
	Add(key string, device Device) error
	Delete(key string) error
	Get(key string) (Device, bool)
}

// deviceConn contains a device and its gNMI connections.
type deviceConn struct {
	device        Device
	ctx           context.Context
	cancel        context.CancelFunc
	gnmiClient    gnmi.GNMIClient
	gnmiOCClient  gnmi.GNMIClient
	providerGroup *errgroup.Group
}

type inventory struct {
	ctx              context.Context
	group            *errgroup.Group
	gnmiServerAddr   string
	gnmiOCServerAddr string
	devices          map[string]*deviceConn
}

func startGNMIClient(serverAddr string) (gnmi.GNMIClient, error) {
	if serverAddr == "" {
		return nil, fmt.Errorf("Invalid gNMI server address '%v'", serverAddr)
	}
	return agnmi.Dial(&agnmi.Config{Addr: serverAddr})
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

	for _, p := range providers {
		pt, ok := p.(provider.GNMIProvider)
		if !ok {
			return errors.New("unexpected provider type; need GNMIProvider")
		}

		// Initialize the provider, depending on whether it wants OpenConfig
		// type-checking.
		if pt.OpenConfig() {
			if dc.gnmiOCClient == nil {
				dc.gnmiOCClient, err = startGNMIClient(i.gnmiOCServerAddr)
				if err != nil {
					return err
				}
			}
			pt.InitGNMI(dc.gnmiOCClient)
		} else {
			if dc.gnmiClient == nil {
				dc.gnmiClient, err = startGNMIClient(i.gnmiServerAddr)
				if err != nil {
					return err
				}
			}
			pt.InitGNMI(dc.gnmiClient)
		}

		// Watch for provider errors in the provider errgroup and
		// propagate them up to the inventory errgroup.
		i.group.Go(func() error {
			return i.handleErrors(dc)
		})

		// Start the providers.
		dc.providerGroup.Go(func() error {
			return p.Run(dc.ctx)
		})
	}
	return nil
}

func (i *inventory) handleErrors(dc *deviceConn) error {
	return dc.providerGroup.Wait()
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
	gnmiServerAddr, gnmiOCServerAddr string) Inventory {
	return &inventory{
		ctx:              ctx,
		devices:          make(map[string]*deviceConn),
		group:            group,
		gnmiServerAddr:   gnmiServerAddr,
		gnmiOCServerAddr: gnmiOCServerAddr,
	}
}
