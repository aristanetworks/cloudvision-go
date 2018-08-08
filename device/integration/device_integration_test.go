// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

// +build integration,device

package devicetest

import (
	"arista/device"
	"testing"
)

var pluginDir = "../plugins/build"
var testPluginName = "test"

func TestDevice(t *testing.T) {
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, "", true)
}

// XXX_jcr: This test has to be in a different package than the other
// tests because of https://github.com/golang/go/issues/17928 ("cannot
// load a plugin from a test where the plugin includes the tested
// package"). The plugin also has to be outside the test package
// because plugins have to be in a "main" package.

// Test creating an adapter with basic plugin
func TestPlugin(t *testing.T) {
	device.UnregisterDevice(testPluginName)
	// Adapter should fail with no pluginDir provided to make sure that there
	// is no built-in device registered with the same name
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, "", false)
	// After checking that no built-in device with the same name is registered,
	// now we can check if the plugin is actually registered
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, pluginDir, true)
}
