// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"errors"
	"fmt"
)

// A Manager manages a device inventory, adding and deleting
// devices as appropriate.
type Manager interface {
	Manage(inventory Inventory) error
}

// ManagerCreator returns a new instance of a Manager.
type ManagerCreator = func(map[string]string) (Manager, error)

// managerInfo contains all the information about a device manager that's
// knowable before it's instantiated: its name, its factory function, and
// the options it supports.
type managerInfo struct {
	name    string
	options map[string]Option
	creator ManagerCreator
}

var managerMap = map[string]managerInfo{}
var managerInUse *managerInfo

// setManagerInUse sets the current manager in use. This is separated
// from CreateManager so that we can print out help messages using
// -help of a specific manager if we fail to correctly configure the
// manager.
func setManagerInUse(name string) error {
	manager, ok := managerMap[name]
	if !ok {
		return fmt.Errorf("Manager %s doesn't exist", name)
	}

	managerInUse = &manager
	return nil
}

// InitManager takes relevant information about a manager and does initial
// setup for that manager.
func InitManager(pluginDir, name string, creator *ManagerCreator,
	managerOpt map[string]Option) error {
	if creator != nil {
		RegisterManager(name, *creator, managerOpt)
	}
	err := loadPlugins(pluginDir)
	if err != nil {
		return fmt.Errorf("Failure in device.loadPlugins: %v", err)
	}
	err = setManagerInUse(name)
	if err != nil {
		return fmt.Errorf("Failure in device.setManagerInUse: %s", err)
	}
	return nil
}

// CreateManager takes a config map, sanitizes the provided config, and
// returns a manager from the current manager in use initialized with the
// sanitized config.
func CreateManager(config map[string]string) (Manager, error) {

	if managerInUse == nil {
		return nil, errors.New("No manager or device in use")
	}

	sanitizedConfig, err := sanitizedOptions(managerInUse.options, config)
	if err != nil {
		return nil, err
	}

	return managerInUse.creator(sanitizedConfig)
}

// RegisterManager registers a function that can create a new Manager of
// the given name.
func RegisterManager(name string, creator ManagerCreator, options map[string]Option) {
	managerMap[name] = managerInfo{
		creator: creator,
		options: options,
		name:    name,
	}
}

// UnregisterManager removes a manager from the registry.
func UnregisterManager(name string) {
	delete(managerMap, name)
}

// ManagerName returns the name of the current manager in use if any.
func ManagerName() string {
	if managerInUse == nil {
		return ""
	}
	return (*managerInUse).name
}

// Delete clears the manager currently in use.
func Delete() {
	managerInUse = nil
}

// ManagerOptionHelp returns the options and associated help strings of
// the manager in use.
func ManagerOptionHelp() (map[string]string, error) {
	if managerInUse == nil {
		return nil, errors.New("No manager in use")
	}
	return helpDesc(managerInUse.options), nil
}
