// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	pmojo "arista/provider/mojo"
	"errors"
)

// The mojo device should never be created directly, so it needn't
// register.

type mojo struct {
	mojoProvider provider.GNMIOpenConfigProvider
	deviceID     string
	initialized  bool
}

func (m *mojo) Type() device.Type {
	return device.Target
}

var errMjUninitialized = errors.New("Mojo device cannot return " +
	"providers until initialized")

func (m *mojo) CheckAlive() (bool, error) {
	if !m.initialized {
		return false, errMjUninitialized
	}
	return true, nil
}

func (m *mojo) DeviceID() (string, error) {
	if !m.initialized {
		return "", errMjUninitialized
	}
	return m.deviceID, nil
}

func (m *mojo) Providers() ([]provider.Provider, error) {
	if !m.initialized {
		return nil, errMjUninitialized
	}
	return []provider.Provider{m.mojoProvider}, nil
}

// NewMojo returns a new Mojo Device.
func NewMojo(deviceID string,
	deviceUpdateChan chan *pmojo.ManagedDevice) device.Device {
	mj := &mojo{}
	mj.deviceID = deviceID
	mj.mojoProvider = pmojo.NewMojoProvider(deviceUpdateChan)
	mj.initialized = true
	return mj
}
