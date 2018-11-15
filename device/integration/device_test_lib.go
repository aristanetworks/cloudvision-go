// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devicetest

import (
	"arista/device"
	"testing"
)

func init() {
	device.RegisterDevice("test", device.NewTestDevice, device.TestDeviceOptions)
}

// XXX_jcr: The device test runner (RunDeviceTest) has to be outside the
// device package because packages imported by Collector (such as device)
// cannot import package testing in files that aren't named *_test.go,
// but if it's named *_test.go it can't export symbols. This is also why
// there's both a device/device_test.go and device/test/device_test.go.

// RunDeviceTest creates a device and fails on an unepected error.
func RunDeviceTest(t *testing.T, deviceName string,
	deviceConfig map[string]string, pluginDir string, shouldPass bool) {

	err := device.Init(pluginDir, deviceName, nil, nil)
	if err != nil && shouldPass {
		t.Fatalf("Unexpected error in device.Init: %s", err)
	}
	_, err = device.CreateManager(deviceConfig)
	if err != nil && shouldPass {
		t.Fatalf("Unexpected error in device.Create: %s", err)
	}

	if err == nil && !shouldPass {
		t.Fatal("Expected error but got none")
	}
	device.Delete()
}
