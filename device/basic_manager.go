// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import "fmt"

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

// NewBasicManager creates a simple manager around a given device.
func NewBasicManager(device Device) Manager {
	return &basicManager{device: device}
}
