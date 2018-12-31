// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package mojo

import (
	"arista/mwm/apiclient"
	"arista/provider"
	pgnmi "arista/provider/gnmi"
	"context"
	"fmt"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type mojo struct {
	ctx              context.Context
	client           gnmi.GNMIClient
	deviceUpdateChan chan *apiclient.ManagedDevice
	errChan          chan error
	initialized      bool
}

func (m *mojo) InitGNMIOpenConfig(client gnmi.GNMIClient) {
	m.client = client
	m.initialized = true
}

// DeviceIDFromMac returns a Mojo device ID from a MAC address.
func DeviceIDFromMac(mac string) string {
	return strings.Replace(mac, ":", "", -1)
}

func (m *mojo) handleDeviceUpdate(deviceUpdate *apiclient.ManagedDevice) {
	updates := []*gnmi.Update{
		pgnmi.Update(pgnmi.PlatformComponentConfigPath("chassis", "name"),
			pgnmi.Strval("chassis")),
		pgnmi.Update(pgnmi.PlatformComponentPath("chassis", "name"),
			pgnmi.Strval("chassis")),
		pgnmi.Update(pgnmi.PlatformComponentStatePath("chassis", "name"),
			pgnmi.Strval("chassis")),
	}
	if deviceUpdate.Model != "" {
		updates = append(updates, pgnmi.Update(
			pgnmi.PlatformComponentStatePath("chassis", "hardware-version"),
			pgnmi.Strval(deviceUpdate.Model)))
	}
	if deviceUpdate.FirmwareVersion != "" {
		updates = append(updates, pgnmi.Update(
			pgnmi.PlatformComponentStatePath("chassis", "software-version"),
			pgnmi.Strval(deviceUpdate.FirmwareVersion)))
	}
	if deviceUpdate.Macaddress != "" {
		updates = append(updates, pgnmi.Update(
			pgnmi.Path("system", "state", "hostname"),
			pgnmi.Strval(DeviceIDFromMac(deviceUpdate.Macaddress))))
	}
	if deviceUpdate.UpSinceTimestamp != 0 {
		// XXX_jcr: This is seconds, I think?
		updates = append(updates, pgnmi.Update(
			pgnmi.Path("system", "state", "boot-time"),
			pgnmi.Intval(100*deviceUpdate.UpSinceTimestamp)))
	}

	sr := &gnmi.SetRequest{
		Delete: []*gnmi.Path{
			pgnmi.Path("system", "state"),
			pgnmi.Path("components")},
		Replace: updates,
	}
	if _, err := m.client.Set(m.ctx, sr); err != nil {
		m.errChan <- err
	}
}

func (m *mojo) Run(ctx context.Context) error {
	if !m.initialized {
		return fmt.Errorf("mojo provider is uninitialized")
	}

	m.ctx = ctx

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-m.deviceUpdateChan:
			m.handleDeviceUpdate(update)
		case err := <-m.errChan:
			return err
		}
	}
}

// NewMojoProvider returns a Mojo provider.
func NewMojoProvider(ch chan *apiclient.ManagedDevice) provider.GNMIOpenConfigProvider {
	return &mojo{
		deviceUpdateChan: ch,
		errChan:          make(chan error),
	}
}
