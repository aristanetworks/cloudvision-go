// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package main

import (
	"github.com/aristanetworks/cloudvision-go/device"
)

var testPluginName = "test"

func init() {
	device.Register(testPluginName, device.NewTestDevice,
		device.TestDeviceOptions)
}

// appease go install ./...
func main() {
}
