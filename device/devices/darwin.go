// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	pdarwin "arista/provider/darwin"
	"os/exec"
	"strings"
)

func init() {
	device.RegisterDevice("darwin", NewDarwinDevice, make(map[string]device.Option))
}

type darwinDevice struct {
	isAlive  bool
	deviceID string
	provider provider.EOSProvider
}

// NewDarwinDevice instantiates a Mac device.
func NewDarwinDevice(options map[string]string) (device.Device, error) {
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
	device.isAlive = true
	device.deviceID = strings.TrimSuffix(string(out), "\n")
	device.provider = pdarwin.NewDarwinProvider()
	return &device, nil
}

func (d *darwinDevice) CheckAlive() (bool, error) {
	return d.isAlive, nil
}

func (d *darwinDevice) DeviceID() (string, error) {
	return d.deviceID, nil
}

func (d *darwinDevice) Providers() []provider.Provider {
	return []provider.Provider{d.provider}
}

func (d *darwinDevice) Type() device.Type {
	return device.Target
}
