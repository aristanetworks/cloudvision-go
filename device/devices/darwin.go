// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	"arista/provider/providers"
	"io"
	"os/exec"
)

func init() {
	device.RegisterDevice("darwin", NewDarwinDevice)
}

type darwinDevice struct {
	name     string
	isAlive  bool
	deviceID string
	provider provider.Provider
}

// NewDarwinDevice gives a representation of device provider for our laptop
func NewDarwinDevice(configFile io.Reader) (device.Device, error) {
	profiler := exec.Command("system_profiler", "SPHardwareDataType")
	awk := exec.Command("awk", "/Serial/ {print $4}")
	pipe, err := profiler.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer pipe.Close()

	awk.Stdin = pipe

	err = profiler.Start()
	if err != nil {
		return nil, err
	}
	out, err := awk.Output()
	if err != nil {
		return nil, err
	}
	device := darwinDevice{}
	device.name = "MacBook"
	device.isAlive = true
	device.deviceID = string(out)
	device.provider = providers.NewDarwinProvider()
	return &device, nil
}

func (d *darwinDevice) Name() string {
	return d.name
}

func (d *darwinDevice) CheckAlive() bool {
	return d.isAlive
}

func (d *darwinDevice) DeviceID() string {
	return d.deviceID
}

func (d *darwinDevice) Providers() []provider.Provider {
	return []provider.Provider{d.provider}
}
