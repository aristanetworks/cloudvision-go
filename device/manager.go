// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
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

// CreateManager takes a config map, sanitizes the provided config, and
// returns a manager.
func CreateManager(name string, config map[string]string) (Manager, error) {
	managerInfo, ok := managerMap[name]
	if !ok {
		return nil, fmt.Errorf("Manager '%v' not found", name)
	}
	sanitizedConfig, err := sanitizedOptions(managerInfo.options, config)
	if err != nil {
		return nil, err
	}

	return managerInfo.creator(sanitizedConfig)
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

// RegisteredManagers returns the names of all registered managers.
func RegisteredManagers() (keys []string) {
	for k := range managerMap {
		keys = append(keys, k)
	}
	return
}

// ManagerOptionHelp returns the options and associated help strings of
// the manager in use.
func ManagerOptionHelp(managerName string) (map[string]string, error) {
	managerInfo, ok := managerMap[managerName]
	if !ok {
		return nil, fmt.Errorf("Manager '%v' not found", managerName)
	}
	return helpDesc(managerInfo.options), nil
}
