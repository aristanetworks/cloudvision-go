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

// The Type is used to distinguish between a normal device and a management system.
type Type int

const (
	// Target is an ordinary device streaming to CloudVision.
	Target Type = 0
	// ManagementSystem is a system managing other devices which itself shouldn't be
	// treated an actual streaming device in CloudVision.
	ManagementSystem Type = 1
)

// A Device knows how to interact with a specific device.
type Device interface {
	Type() Type
	CheckAlive() (bool, error)
	DeviceID() (string, error)
	Providers() []provider.Provider
}

// String converts a device Type enum to it's string value.
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

// ManagerName returns the name of the current manager in use if any.
func ManagerName() string {
	if managerInUse == nil {
		return ""
	}
	return (*managerInUse).name
}

// sanitizedOptions takes the map of device option keys and values
// passed in at the command line and checks it against the device's
// exported list of accepted options, returning an error if there
// are inappropriate or missing options.
func sanitizedOptions(manager *managerInfo, config map[string]string) (map[string]string, error) {
	if manager == nil {
		return nil, fmt.Errorf("Nil deviceInfo")
	}

	options := manager.options
	sopt := make(map[string]string)

	// Check whether the user gave us bad options.
	for k, v := range config {
		_, ok := options[k]
		if !ok {
			return nil, fmt.Errorf("Bad option '%s' for manager '%s'", k, manager.name)
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

// setDeviceInUse sets the current device in use. This is separated from CreateDevice so that
// we can print out help messages using -help of a specific device if we fail to correctly
// configure the device.
func setDeviceInUse(name string) error {
	di, ok := deviceMap[name]
	if !ok {
		return fmt.Errorf("Device %s doesn't exist", name)
	}

	managerInUse = &managerInfo{
		name:    di.name,
		options: di.options,
		creator: transformCreator(di.creator),
	}
	return nil
}

// CreateManager takes a config map, sanitizes the provided config, and
// returns a manager from the current manager in use initialized with the sanitized config.
func CreateManager(config map[string]string) (Manager, error) {

	if managerInUse == nil {
		return nil, errors.New("No manager in use")
	}

	sanitizedConfig, err := sanitizedOptions(managerInUse, config)
	if err != nil {
		return nil, err
	}

	return managerInUse.creator(sanitizedConfig)
}

// Init takes relevant information about a device and does initial setup for that device.
func Init(pluginDir, deviceName string, creator *Creator,
	deviceOpt map[string]Option) error {

	if creator != nil {
		RegisterDevice(deviceName, *creator, deviceOpt)
	}
	err := loadPlugins(pluginDir)
	if err != nil {
		return fmt.Errorf("Failure in device.loadPlugins: %v", err)
	}
	err = setDeviceInUse(deviceName)
	if err != nil {
		return fmt.Errorf("Failure in device.setDeviceInUse: %s", err)
	}
	return nil
}

// Delete clears the manager currently in use.
func Delete() {
	managerInUse = nil
}

// UnregisterDevice removes a device from the registry.
func UnregisterDevice(name string) {
	delete(deviceMap, name)
}

// Create map of option key to description.
func helpDesc(options map[string]Option) map[string]string {
	hd := make(map[string]string)

	for k, v := range options {
		desc := v.Description
		// Add default if there's a non-empty one.
		if v.Default != "" {
			desc = desc + " (default " + v.Default + ")"
		}
		hd[k] = desc
	}
	return hd
}

// Return managerInUse's help string.
func help(options map[string]Option, name string) string {
	b := new(bytes.Buffer)
	hd := helpDesc(options)
	// Don't print out device separator if the device has no options.
	if len(hd) == 0 {
		return ""
	}
	flag.FormatOptions(b, "Help options for device/manager "+name+":", hd)
	return b.String()
}

// AddHelp adds the deviceInUse's options to flag.Usage.
func AddHelp() error {
	if managerInUse == nil {
		return errors.New("No manager in use")
	}

	h := help(managerInUse.options, managerInUse.name)
	flag.AddHelp("", h)
	return nil
}
