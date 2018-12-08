// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package managers

import (
	"arista/device"
	"fmt"
)

// There's no need to register this manager, as we should never need to
// create it dynamically.

type basicManager struct {
	device device.Device
}

func (m *basicManager) Manage(inventory device.Inventory) error {
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

// NewBasicManager creates a simple manager around a given device.
func NewBasicManager(d device.Device) device.Manager {
	return &basicManager{device: d}
}
