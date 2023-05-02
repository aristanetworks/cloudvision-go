// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/device/gen"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testDevice struct{}

func (td testDevice) Alive(ctx context.Context) (bool, error) {
	return true, nil
}

func (td testDevice) DeviceID(ctx context.Context) (string, error) {
	return "aaa", nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (td testDevice) Type() string {
	return ""
}

func (td testDevice) IPAddr(ctx context.Context) (string, error) {
	return "192.168.1.2", nil
}

type testManager struct{}

func (tm testManager) Alive(ctx context.Context) (bool, error) {
	return true, nil
}

func (tm testManager) DeviceID(ctx context.Context) (string, error) {
	return "bbb", nil
}

func (tm testManager) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (tm testManager) Type() string {
	return ""
}

func (tm testManager) Manage(ctx context.Context, inv device.Inventory) error {
	return inv.Add(&device.Info{
		Device: testDevice{},
		ID:     "aaa",
		Config: &device.Config{},
	})
}

func (tm testManager) IPAddr(ctx context.Context) (string, error) {
	return "192.168.0.123", nil
}

// newTestManager returns a dummy device for testing.
func newTestManager(ctx context.Context, opts map[string]string,
	monitor provider.Monitor) (device.Device, error) {
	return testManager{}, nil
}

func TestGRPCServer(t *testing.T) {
	ctx := context.Background()
	inventory := device.NewInventory(ctx, pgnmi.NewSimpleGNMIClient(
		func(context.Context, *gnmi.SetRequest) (*gnmi.SetResponse, error) {
			return nil, nil
		}), newCVClient)
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
	conn, err := grpc.Dial(listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func TestLoggingFormat(t *testing.T) {
	lvl := "info"
	dir := t.TempDir()
	logLevel = &lvl
	logDir = &dir
	// Mostly ensuring we don't crash
	initLogging()

	// Log looks like
	// time="2023-04-17T15:52:12-07:00" level=info msg="info 123" file="libmain/libmain_test.go:163"
	logrus.Infof("info %d", 123)
}

func TestWatchConfig(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		name          string
		content       string
		expectConfigs []*device.Config
	}{
		{
			name: "one config",
			content: `
- Name: cfg1
  Device: type1
  Options:
    opt1: val1
    opt2: val2`,
			expectConfigs: []*device.Config{
				{
					Name:   "cfg1",
					Device: "type1",
					Options: map[string]string{
						"opt1": "val1",
						"opt2": "val2",
					},
				},
			},
		},
		{
			name: "two new unnamed configs and delete previous config",
			content: `
- Device: type2
  Options:
    opt1: val1
    opt2: val2
- Device: type3
  Options:
    opt1: val1
    opt2: val2`,
			expectConfigs: []*device.Config{
				{
					Name:   "auto-datasource-000",
					Device: "type2",
					Options: map[string]string{
						"opt1": "val1",
						"opt2": "val2",
					},
				},
				{
					Name:   "auto-datasource-001",
					Device: "type3",
					Options: map[string]string{
						"opt1": "val1",
						"opt2": "val2",
					},
				},
				device.NewDeletedConfig("cfg1"),
			},
		},
	}

	file := filepath.Join(t.TempDir(), "test.yml")
	// The initial file should not be seen as an update.
	initialContent := `
- Device: should-not-read
  Options:
    opt1: val1
    opt2: val2
`
	if err := os.WriteFile(file, []byte(initialContent), 0666); err != nil {
		t.Fatal(err)
	}

	configCh := make(chan *device.Config)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		defer close(configCh)
		return watchConfig(ctx, nil, file, configCh, 30*time.Millisecond)
	})

	<-time.After(40 * time.Millisecond) // need to wait watcher to be ready

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := func() error {
				if err := os.WriteFile(file, []byte(tc.content), 0666); err != nil {
					return err
				}

				timer := time.NewTimer(5 * time.Second)
				for {
					select {
					case <-timer.C:
						select {
						case cfg := <-configCh:
							return fmt.Errorf("Unexpected config update: %v", cfg)
						default:
							if len(tc.expectConfigs) > 0 {
								return fmt.Errorf("Did not match %v! Plus %d others",
									tc.expectConfigs[0].Name, len(tc.expectConfigs)-1)
							}
							return nil
						}
					case cfg, ok := <-configCh:
						if !ok {
							continue
						}
						if len(tc.expectConfigs) == 0 {
							return fmt.Errorf("Unexpected config update: %v", cfg)
						}
						expect := tc.expectConfigs[0]
						if !reflect.DeepEqual(cfg, expect) {
							return fmt.Errorf("Config mismatch, expected:\n%v\ngot:\n%v",
								expect, cfg)
						}
						tc.expectConfigs = tc.expectConfigs[1:] // wait for next
						if len(tc.expectConfigs) == 0 {
							timer.Reset(31 * time.Millisecond)
						}
					}
				}
			}()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
