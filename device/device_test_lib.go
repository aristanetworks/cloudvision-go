// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"arista/provider"
)

type TestDevice struct {
}

var TestDeviceOptions = map[string]Option{
	"a": Option{
		Description: "option a is a required option",
		Default:     "",
		Required:    true,
	},
	"b": Option{
		Description: "option b is not required",
		Default:     "stuff",
		Required:    false,
	},
}

func (td TestDevice) CheckAlive() (bool, error) {
	return true, nil
}

func (td TestDevice) DeviceID() (string, error) {
	return "0a0a.0a0a.0a0a", nil
}

func (td TestDevice) Providers() []provider.Provider {
	return nil
}

// NewTestDevice returns a dummy device for testing.
func NewTestDevice(map[string]string) (Device, error) {
	return TestDevice{}, nil
}

var TestDeviceConfig = map[string]string{
	"a": "abc",
	"b": "stuff",
}
