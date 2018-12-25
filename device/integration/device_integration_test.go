// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

// +build integration,device

package devicetest

import (
	"cloudvision-go/device"
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

// Test creating a device with a basic plugin.
func TestPlugin(t *testing.T) {
	device.Unregister(testPluginName)
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, pluginDir, true)
}
