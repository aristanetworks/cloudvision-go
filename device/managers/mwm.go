// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package managers

import (
	"arista/device"
	"arista/device/devices"
	"arista/mwm/apiclient"
	pmojo "arista/provider/mojo"
	"net/http"
	"time"
)

func init() {
	options := map[string]device.Option{
		"address": device.Option{
			Description: "MWM address",
			Required:    true,
		},
		"authTimeout": device.Option{
			Description: "MWM authentication timeout, with unit suffix (s/m/h)",
			Default:     "1h",
		},
		"pollInterval": device.Option{
			Description: "Polling interval, with unit suffix (s/m/h)",
			Default:     "30s",
		},
		"username": device.Option{
			Description: "MWM username",
			Required:    true,
		},
		"password": device.Option{
			Description: "MWM password",
			Required:    true,
		},
	}
	device.RegisterManager("mwm", newMWM, options)
}

type mwm struct {
	address          string
	authTimeout      time.Duration
	client           *http.Client
	deviceUpdateChan map[string]chan *apiclient.ManagedDevice
	inventory        device.Inventory
	lastAuthTime     time.Time
	password         string
	pollInterval     time.Duration
	session          *apiclient.APISession
	username         string
}

func (m *mwm) addDevice(d *apiclient.ManagedDevice) (device.Device, error) {
	ch := make(chan *apiclient.ManagedDevice)
	did := pmojo.DeviceIDFromMac(d.Macaddress)
	md := devices.NewMojo(did, ch)
	if err := m.inventory.Add(did, md); err != nil {
		return nil, err
	}
	m.deviceUpdateChan[did] = ch
	return md, nil
}

func (m *mwm) removeDevice(key string) error {
	if err := m.inventory.Delete(key); err != nil {
		return err
	}
	delete(m.deviceUpdateChan, key)
	return nil
}

func (m *mwm) handleDeviceUpdate(d *apiclient.ManagedDevice) error {
	did := pmojo.DeviceIDFromMac(d.Macaddress)
	_, ok := m.inventory.Get(did)
	if !ok {
		if !d.IsActive {
			return nil
		}
		_, err := m.addDevice(d)
		if err != nil {
			return err
		}
	}

	if !d.IsActive {
		if err := m.removeDevice(did); err != nil {
			return err
		}
		return nil
	}

	m.deviceUpdateChan[did] <- d
	return nil
}

func (m *mwm) pollMWM() error {
	// Log in if the previous login has expired.
	t := time.Now()
	nextAuthNeeded := m.lastAuthTime.Add(m.authTimeout)
	if t.Add(m.pollInterval).After(nextAuthNeeded) {
		_, err := m.session.Open(m.username, m.password, m.authTimeout)
		if err != nil {
			return err
		}
		m.lastAuthTime = t
	}

	devices, err := m.session.GetManagedDevices(apiclient.Filter{})
	if err != nil {
		return err
	}

	for _, d := range devices {
		if err := m.handleDeviceUpdate(d); err != nil {
			return err
		}
	}
	return nil
}

func (m *mwm) Manage(inventory device.Inventory) error {
	m.inventory = inventory

	var err error
	m.session, err = apiclient.NewAPISession(m.address)
	if err != nil {
		return err
	}

	if err = m.pollMWM(); err != nil {
		return err
	}

	tick := time.NewTicker(m.pollInterval)
	defer tick.Stop()
	for range tick.C {
		if err := m.pollMWM(); err != nil {
			return err
		}
	}

	return nil
}

func newMWM(options map[string]string) (device.Manager, error) {
	m := &mwm{}
	var err error

	m.address, err = device.GetAddressOption("address", options)
	if err != nil {
		return nil, err
	}

	m.authTimeout, err = device.GetDurationOption("authTimeout", options)
	if err != nil {
		return nil, err
	}

	m.password, err = device.GetStringOption("password", options)
	if err != nil {
		return nil, err
	}

	m.pollInterval, err = device.GetDurationOption("pollInterval", options)
	if err != nil {
		return nil, err
	}

	m.username, err = device.GetStringOption("username", options)
	if err != nil {
		return nil, err
	}

	m.deviceUpdateChan = make(map[string]chan *apiclient.ManagedDevice)

	return m, nil
}
