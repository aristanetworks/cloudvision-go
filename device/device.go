// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"

	"github.com/aristanetworks/cloudvision-go/provider"
)

// Type is used to distinguish between target devices and management systems.
type Type int

const (
	// Target is an ordinary device streaming to CloudVision.
	Target Type = 0
	// ManagementSystem is a system managing other devices which itself
	// shouldn't be treated an actual streaming device in CloudVision.
	ManagementSystem Type = 1
)

// A Device knows how to interact with a specific device.
type Device interface {
	Type() Type
	Alive() (bool, error)
	DeviceID() (string, error)
	Providers() ([]provider.Provider, error)
}

// String converts a device Type enum to its string value.
func (t Type) String() string {
	return map[Type]string{
		Target:           "target",
		ManagementSystem: "managementSystem",
	}[t]
}

// Creator returns a new instance of a Device.
type Creator = func(map[string]string) (Device, error)

// deviceInfo contains all the information about a device that's
// knowable before it's instantiated: its name, its factory function,
// and the options it supports.
type deviceInfo struct {
	name    string
	creator Creator
	options map[string]Option
}

var deviceMap = map[string]deviceInfo{}

// Register registers a function that can create a new Device
// of the given name.
func Register(name string, creator Creator, options map[string]Option) {
	deviceMap[name] = deviceInfo{
		name:    name,
		creator: creator,
		options: options,
	}
}

// Unregister removes a device from the registry.
func Unregister(name string) {
	delete(deviceMap, name)
}

// Registered returns a list of registered device names.
func Registered() (keys []string) {
	for k := range deviceMap {
		keys = append(keys, k)
	}
	return
}

// Create takes a config map, sanitizes the provided config, and returns
// a Device.
func Create(name string, config map[string]string) (Device, error) {
	deviceInfo, ok := deviceMap[name]
	if !ok {
		return nil, fmt.Errorf("Device '%v' not found", name)
	}
	sanitizedConfig, err := sanitizedOptions(deviceInfo.options, config)
	if err != nil {
		return nil, err
	}
	return deviceInfo.creator(sanitizedConfig)
}

// OptionHelp returns the options and associated help strings of the
// specified device.
func OptionHelp(deviceName string) (map[string]string, error) {
	deviceInfo, ok := deviceMap[deviceName]
	if !ok {
		return nil, fmt.Errorf("Device '%v' not found", deviceName)
	}
	return helpDesc(deviceInfo.options), nil
}
