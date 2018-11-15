// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

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
// knowable before it's instantiated: its name, its factory function,
// and the options it supports.
type managerInfo struct {
	name    string
	options map[string]Option
	creator ManagerCreator
}

var managerMap = map[string]managerInfo{}
var managerInUse *managerInfo

// InitManager takes relevant information about a manager and does initial setup for that manager.
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

// RegisterManager registers a function that can create a new Manager
// of the given name.
func RegisterManager(name string, creator ManagerCreator, options map[string]Option) {
	managerMap[name] = managerInfo{
		creator: creator,
		options: options,
		name:    name,
	}
}

// setManagerInUse sets the current manager in use. This is separated from CreateManager so that
// we can print out help messages using -help of a specific manager if we fail to correctly
// configure the manager.
func setManagerInUse(name string) error {
	manager, ok := managerMap[name]
	if !ok {
		return fmt.Errorf("Manager %s doesn't exist", name)
	}

	managerInUse = &manager
	return nil
}

// UnregisterManager removes a manager from the registry.
func UnregisterManager(name string) {
	delete(managerMap, name)
}

type basicManager struct {
	device Device
}

func (m *basicManager) Manage(inventory Inventory) error {
	id, err := m.device.DeviceID()
	if err != nil {
		return err
	}
	err = inventory.Add(id, m.device)
	if err != nil {
		return fmt.Errorf("Error in adding device: %s", err)
	}
	return nil
}

func newBasicManager(device Device) *basicManager {
	return &basicManager{device: device}
}

func transformCreator(creator Creator) ManagerCreator {
	return func(options map[string]string) (Manager, error) {
		dev, err := creator(options)
		if err != nil {
			return nil, err
		}
		return newBasicManager(dev), nil
	}
}
