// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// +build integration,device

package devicetest

import (
	"testing"

	"github.com/aristanetworks/cloudvision-go/device"
)

var pluginDir = "../plugins/build"
var testPluginName = "test"

func TestDevice(t *testing.T) {
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, "", true)
}

// XXX NOTE: This test has to be in a different package than the other
// tests because of https://github.com/golang/go/issues/17928 ("cannot
// load a plugin from a test where the plugin includes the tested
// package"). The plugin also has to be outside the test package
// because plugins have to be in a "main" package.

// Test creating a device with a basic plugin.
func TestPlugin(t *testing.T) {
	device.Unregister(testPluginName)
	RunDeviceTest(t, testPluginName, device.TestDeviceConfig, pluginDir, true)
}
