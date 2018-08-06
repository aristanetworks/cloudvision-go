// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"fmt"
	"io"
)

// Creator is a function that when called returns a new instance of a Device
type Creator = func(io.Reader) (Device, error)

var deviceMap = map[string]Creator{}

// RegisterDevice registers a function that can create a new Device given by the name
func RegisterDevice(name string, creator Creator) {
	deviceMap[name] = creator
}

// CreateDevice returns a device from the registry
func CreateDevice(name string, config io.Reader) (Device, error) {
	creator, ok := deviceMap[name]

	if ok {
		return creator(config)
	}
	return nil, fmt.Errorf("Device %s doesn't exist", name)
}

// UnregisterDevice removes the device from the registry
func UnregisterDevice(name string) {
	delete(deviceMap, name)
}
