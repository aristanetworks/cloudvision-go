// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"cloudvision-go/provider"
)

// TestDeviceOptions is a set of test options
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

type testDevice struct{}

// TestDeviceID is the device ID used for retriving the device from the inventory.
var TestDeviceID = "0a0a.0a0a.0a0a"

func (td testDevice) CheckAlive() (bool, error) {
	return true, nil
}

func (td testDevice) DeviceID() (string, error) {
	return TestDeviceID, nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (td testDevice) Type() Type {
	return Target
}

// NewTestDevice returns a dummy device for testing.
func NewTestDevice(map[string]string) (Device, error) {
	return testDevice{}, nil
}

// TestDeviceConfig is a map of test config options.
var TestDeviceConfig = map[string]string{
	"a": "abc",
	"b": "stuff",
}

type testManager struct{}

func (td testManager) Manage(inventory Inventory) error {
	err := inventory.Delete(TestDeviceID)
	if err != nil {
		return err
	}
	return nil
}

// NewTestManager returns a dummy manager for testing.
func NewTestManager(map[string]string) (Manager, error) {
	return testManager{}, nil
}
