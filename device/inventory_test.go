// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"reflect"
	"testing"

	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/sync/errgroup"
)

func TestInventoryBasic(t *testing.T) {
	processor := func(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		return nil, nil
	}
	group, ctx := errgroup.WithContext(context.Background())
	inventory := NewInventory(ctx, group, pgnmi.NewSimpleGNMIClient(processor))
	expectedDevice := testDevice{}
	deviceID := "dummy"
	err := inventory.Add(deviceID, expectedDevice)
	if err != nil {
		t.Fatal(err)
	}
	actualDevice, ok := inventory.Get(deviceID)
	if !ok {
		t.Fatalf("Device '%s' not found in inventory", deviceID)
	}
	if !reflect.DeepEqual(expectedDevice, actualDevice) {
		t.Fatalf("Added device different from retrieved device\nAdded: %v\nRetrieved: %v",
			expectedDevice, actualDevice)
	}
	err = inventory.Delete(deviceID)
	if err != nil {
		t.Fatal(err)
	}
	_, ok = inventory.Get(deviceID)
	if ok {
		t.Fatalf("Device '%s' is found in inventory after deletion", deviceID)
	}
}
