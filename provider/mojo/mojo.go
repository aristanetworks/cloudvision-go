// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package mojo

import (
	"arista/provider"
	pgnmi "arista/provider/gnmi"
	"context"
	"fmt"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// ManagedDevice encapsulates an MWM managed device.
type ManagedDevice struct {
	BoxID           int64  `json:"boxId"`
	Name            string `json:"name"`
	Macaddress      string `json:"macaddress"`
	Model           string `json:"model"`
	SoftwareVersion string `json:"softwareVersion"`
	IPAddress       string `json:"ipAddress"`
	Active          bool   `json:"active"`
	UpSince         uint64 `json:"upSince"`
}

type mojo struct {
	ctx              context.Context
	client           gnmi.GNMIClient
	deviceUpdateChan chan *ManagedDevice
	errChan          chan error
	initialized      bool
}

func (m *mojo) InitGNMIOpenConfig(client gnmi.GNMIClient) {
	m.client = client
	m.initialized = true
}

func platformComponentPath(leafName string) *gnmi.Path {
	return pgnmi.Path("components",
		pgnmi.ListWithKey("component", "name", "chassis"), leafName)
}
func platformComponentConfigPath(leafName string) *gnmi.Path {
	return pgnmi.Path("components",
		pgnmi.ListWithKey("component", "name", "chassis"), "config", leafName)
}
func platformComponentStatePath(leafName string) *gnmi.Path {
	return pgnmi.Path("components",
		pgnmi.ListWithKey("component", "name", "chassis"), "state", leafName)
}

// DeviceIDFromMac returns a Mojo device ID from a MAC address.
func DeviceIDFromMac(mac string) string {
	return strings.Replace(mac, ":", "", -1)
}

func (m *mojo) handleDeviceUpdate(deviceUpdate *ManagedDevice) {
	updates := []*gnmi.Update{
		pgnmi.Update(platformComponentConfigPath("name"), pgnmi.Strval("chassis")),
		pgnmi.Update(platformComponentPath("name"), pgnmi.Strval("chassis")),
		pgnmi.Update(platformComponentStatePath("name"), pgnmi.Strval("chassis")),
	}
	if deviceUpdate.Model != "" {
		updates = append(updates, pgnmi.Update(
			platformComponentStatePath("hardware-version"),
			pgnmi.Strval(deviceUpdate.Model)))
	}
	if deviceUpdate.SoftwareVersion != "" {
		updates = append(updates, pgnmi.Update(
			platformComponentStatePath("software-version"),
			pgnmi.Strval(deviceUpdate.SoftwareVersion)))
	}
	if deviceUpdate.Macaddress != "" {
		updates = append(updates, pgnmi.Update(
			pgnmi.Path("system", "state", "hostname"),
			pgnmi.Strval(DeviceIDFromMac(deviceUpdate.Macaddress))))
	}
	if deviceUpdate.UpSince != 0 {
		// XXX_jcr: This is seconds, I think?
		updates = append(updates, pgnmi.Update(
			pgnmi.Path("system", "state", "boot-time"),
			pgnmi.Uintval(100*deviceUpdate.UpSince)))
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
func NewMojoProvider(deviceUpdateChan chan *ManagedDevice) provider.GNMIOpenConfigProvider {
	return &mojo{
		deviceUpdateChan: deviceUpdateChan,
		errChan:          make(chan error),
	}
}
