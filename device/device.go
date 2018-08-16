// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"arista/flag"
	"arista/provider"
	"bytes"
	"errors"
	"fmt"
)

// Option defines a command-line option accepted by a device.
type Option struct {
	Description string
	Default     string
	Required    bool
}

// A Device knows how to interact with a specific device.
type Device interface {
	CheckAlive() (bool, error)
	DeviceID() (string, error)
	Providers() []provider.Provider
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
var deviceInUse *deviceInfo

func Name() string {
	if deviceInUse == nil {
		return ""
	}
	return (*deviceInUse).name
}

// sanitizedOptions takes the map of device option keys and values
// passed in at the command line and checks it against the device's
// exported list of accepted options, returning an error if there
// are inappropriate or missing options.
func sanitizedOptions(device *deviceInfo,
	config map[string]string) (map[string]string, error) {
	if device == nil {
		return nil, fmt.Errorf("Nil deviceInfo")
	}

	options := device.options
	sopt := make(map[string]string)

	// Check whether the user gave us bad options.
	for k, v := range config {
		_, ok := options[k]
		if !ok {
			return nil, fmt.Errorf("Bad option '%s' for device '%s'", k, device.name)
		}
		sopt[k] = v
	}

	// Check that all required options were specified, and fill in
	// any others with defaults.
	for k, v := range options {
		_, found := sopt[k]
		if v.Required && !found {
			return nil, fmt.Errorf("Required option '%s' not provided", k)
		}
		if !found {
			sopt[k] = v.Default
		}
	}

	return sopt, nil
}

// RegisterDevice registers a function that can create a new Device
// of the given name.
func RegisterDevice(name string, creator Creator, options map[string]Option) {
	deviceMap[name] = deviceInfo{
		name:    name,
		creator: creator,
		options: options,
	}
}

// CreateDevice takes a device name and config map, sanitizes the
// provided config, and returns a device from the registry initialized
// with the sanitized config.
func CreateDevice(name string, config map[string]string) (Device, error) {
	di, ok := deviceMap[name]
	if !ok {
		return nil, fmt.Errorf("Device %s doesn't exist", name)
	}

	// Set device in use
	deviceInUse = &di

	sanitizedConfig, err := sanitizedOptions(&di, config)
	if err != nil {
		return nil, err
	}

	return di.creator(sanitizedConfig)
}

// DeleteDevice clears the registry of any created devices.
func DeleteDevice() {
	deviceInUse = nil
}

// UnregisterDevice removes a device from the registry.
func UnregisterDevice(name string) {
	delete(deviceMap, name)
}

// Create map of option key to description.
func helpDesc(devInfo deviceInfo) map[string]string {
	hd := make(map[string]string)

	for k, v := range devInfo.options {
		desc := v.Description
		// Add default if there's a non-empty one.
		if v.Default != "" {
			desc = desc + " (default " + v.Default + ")"
		}
		hd[k] = desc
	}
	return hd
}

// Return deviceInUse's help string.
func help(deviceInUse deviceInfo) string {
	b := new(bytes.Buffer)
	hd := helpDesc(deviceInUse)
	// Don't print out device separator if the device has no options.
	if len(hd) == 0 {
		return ""
	}
	flag.FormatOptions(b, "Help options for device "+deviceInUse.name+":", hd)
	return b.String()
}

// AddHelp adds the deviceInUse's options to flag.Usage.
func AddHelp() error {
	if deviceInUse == nil {
		return errors.New("No device in use")
	}

	h := help(*deviceInUse)
	flag.AddHelp("", h)
	return nil
}
