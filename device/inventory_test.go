// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"reflect"
	"testing"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	v1client "github.com/aristanetworks/cloudvision-go/device/cvclient/v1"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"

	"github.com/openconfig/gnmi/proto/gnmi"
)

func TestInventoryBasic(t *testing.T) {
	processor := func(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		return nil, nil
	}
	inventory := NewInventory(context.Background(), pgnmi.NewSimpleGNMIClient(processor),
		func(gc gnmi.GNMIClient, i *Info) cvclient.CVClient {
			return v1client.NewV1Client(gc, i.ID, false)
		})
	expectedDevice := testDevice{}
	deviceID := "dummy"
	err := inventory.Add(&Info{Device: expectedDevice, ID: deviceID})
	if err != nil {
		t.Fatal(err)
	}
	actualDevice, err := inventory.Get(deviceID)
	if err != nil {
		t.Fatalf("Device '%s' not found in inventory: %v", deviceID, err)
	}
	if !reflect.DeepEqual(expectedDevice, actualDevice.Device) {
		t.Fatalf("Added device different from retrieved device\nAdded: %v\nRetrieved: %v",
			expectedDevice, actualDevice.Device)
	}
	err = inventory.Delete(deviceID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = inventory.Get(deviceID)
	if err == nil {
		t.Fatalf("Device '%s' is found in inventory after deletion", deviceID)
	}
}
