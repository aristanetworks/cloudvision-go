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
	err := inventory.Add(&Info{
		Config: &Config{},
		Device: expectedDevice,
		ID:     deviceID,
	})
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

type invTestCase struct {
	name string
	del  string // del takes priority
	add  *Info
}

func runInventoryTest(t *testing.T, inv Inventory, tc *invTestCase) {
	if tc.del != "" {
		if err := inv.Delete(tc.del); err != nil {
			t.Fatalf("failed in inventory.Delete [%s]: %s", tc.del, err)
		}
		if _, err := inv.Get(tc.del); err == nil {
			t.Fatalf("did not delete %q", tc.del)
		}
	} else if tc.add != nil {
		if err := inv.Add(tc.add); err != nil {
			t.Fatalf("failed in inventory.Add [%v]: %s", tc.add, err)
		}
		info, err := inv.Get(tc.add.ID)
		if err != nil {
			t.Fatalf("failed in inventory.Get [%s]: %s", tc.add.ID, err)
		}
		if info.String() != tc.add.String() {
			t.Fatalf("did not update Info in inventory.Add (got: %s, want: %s)",
				info.String(), tc.add.String())
		}
	}
}

func TestInventory(t *testing.T) {
	deviceInfo1 := &Info{
		Config: &Config{
			Device: "bogus",
		},
		Device: &testDevice{
			deviceID: "deviceone",
		},
		ID: "deviceone",
	}
	deviceInfo1Alt := &Info{
		Config: &Config{
			Device: "stuff",
		},
		Device: &testDevice{
			deviceID: "deviceone",
		},
		ID: "deviceone",
	}
	deviceInfo2 := &Info{
		Config: &Config{
			Device: "whatever",
		},
		Device: &testDevice{
			deviceID: "devicetwo",
		},
		ID: "devicetwo",
	}

	tests := []invTestCase{
		{
			name: "add first device",
			add:  deviceInfo1,
		},
		{
			name: "add second device",
			add:  deviceInfo2,
		},
		{
			name: "delete second device",
			del:  deviceInfo2.ID,
		},
		{
			name: "update first device config",
			add:  deviceInfo1Alt,
		},
	}

	processor := func(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		return nil, nil
	}
	inv := NewInventory(context.Background(), pgnmi.NewSimpleGNMIClient(processor),
		func(gc gnmi.GNMIClient, i *Info) cvclient.CVClient {
			return v1client.NewV1Client(gc, i.ID, false)
		})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runInventoryTest(t, inv, &tc)
		})
	}
}
