// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"context"
	"testing"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/device/gen"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

type testDevice struct{}

func (td testDevice) Alive() (bool, error) {
	return true, nil
}

func (td testDevice) DeviceID() (string, error) {
	return "aaa", nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

type testManager struct{}

func (tm testManager) Alive() (bool, error) {
	return true, nil
}

func (tm testManager) DeviceID() (string, error) {
	return "bbb", nil
}

func (tm testManager) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (tm testManager) Manage(inv device.Inventory) error {
	return inv.Add(&device.Info{Device: testDevice{}, ID: "aaa"})
}

// newTestManager returns a dummy device for testing.
func newTestManager(map[string]string) (device.Device, error) {
	return testManager{}, nil
}

func TestGRPCServer(t *testing.T) {
	ctx := context.Background()
	inventory, err := device.NewInventory(ctx, pgnmi.NewSimpleGNMIClient(
		func(context.Context, *gnmi.SetRequest) (*gnmi.SetResponse, error) {
			return nil, nil
		}), "")
	if err != nil {
		t.Fatal(err)
	}
	grpcServer, listener, err := newGRPCServer("localhost:0", inventory)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			t.Error(err)
		}
	}()
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := gen.NewDeviceInventoryClient(conn)
	device.Register("test", newTestManager, nil)
	_, err = client.Add(context.Background(), &gen.AddRequest{
		DeviceConfig: &gen.DeviceConfig{
			DeviceType: "test",
		}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), &gen.GetRequest{
		DeviceID: "aaa",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), &gen.GetRequest{
		DeviceID: "bbb",
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.List(context.Background(), &gen.ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.DeviceInfos) != 2 {
		t.Fatalf("Expect one device in inventory but got %d", len(resp.DeviceInfos))
	}
	_, err = client.Delete(context.Background(), &gen.DeleteRequest{DeviceID: "aaa"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), &gen.GetRequest{
		DeviceID: "aaa",
	})
	if err == nil {
		t.Fatalf("Device is found in inventory after deletion")
	}
	_, err = client.Delete(context.Background(), &gen.DeleteRequest{DeviceID: "bbb"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), &gen.GetRequest{
		DeviceID: "bbb",
	})
	if err == nil {
		t.Fatalf("Device is found in inventory after deletion")
	}
}
