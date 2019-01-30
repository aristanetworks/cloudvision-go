// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package managers

import (
	"fmt"

	"github.com/aristanetworks/cloudvision-go/device"
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
