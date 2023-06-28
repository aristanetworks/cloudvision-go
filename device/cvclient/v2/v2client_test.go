// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package v2

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/provider"

	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func verifyUpdates(r *gnmi.SetRequest, expData map[string]interface{},
	checkValues bool) error {
	if len(r.Update) != len(expData) {
		return fmt.Errorf("number of updates (%d) does not match expected value (%d)",
			len(r.Update), len(expData))
	}
	for _, u := range r.Update {
		p := agnmi.StrPath(agnmi.JoinPaths(r.Prefix, u.Path))
		v, ok := expData[p]
		if !ok {
			return fmt.Errorf("unexpected leaf in update: %s, value: %v", p, v)
		}
		// override time value to 42
		if strings.HasSuffix(p, "last-seen") {
			u.Val = agnmi.TypedValue(int(42))
		}
		if checkValues && !reflect.DeepEqual(v, u.Val) {
			return fmt.Errorf("unexpected value for leaf %s, expected: %+v, got: %v", p, v, u.Val)
		}
	}
	return nil
}

func verifyMetadataLeaves(r *gnmi.SetRequest, c *v2Client) error {
	ip, _ := testDevice{}.IPAddr(context.Background())
	integ := agnmi.TypedValue(c.info.Config.Device)
	expData := map[string]interface{}{
		"/device-metadata/state/metadata/type":              agnmi.TypedValue(c.deviceType),
		"/device-metadata/state/metadata/source-type":       integ,
		"/device-metadata/state/metadata/collector-version": agnmi.TypedValue(versionString),
		"/device-metadata/state/metadata/ip-addr":           agnmi.TypedValue(ip),
		"/device-metadata/state/metadata/managed-device-status": agnmi.
			TypedValue(string(c.info.Status)),
	}
	return verifyUpdates(r, expData, true)
}

type testDevice struct{}

func (td testDevice) DeviceID(ctx context.Context) (string, error) {
	return "mycontroller", nil
}

func (td testDevice) Alive(ctx context.Context) (bool, error) {
	return true, nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return []provider.Provider{}, nil
}

func (td testDevice) Type() string {
	return ""
}

func (td testDevice) IPAddr(ctx context.Context) (string, error) {
	return "192.168.5.10", nil
}

func TestMetadataRequest(t *testing.T) {
	c := NewV2Client(nil,
		&device.Info{
			Device: testDevice{},
			Config: &device.Config{
				Device: "test",
			},
			Status: device.StatusRemoved,
		}).(*v2Client)
	r := c.metadataRequest(context.Background())
	if err := verifyMetadataLeaves(r, c); err != nil {
		t.Logf("Error verifying leaves in set request: %v", err)
		t.Fail()
	}
}

func TestHeartbeatRequest(t *testing.T) {
	for _, tc := range []struct {
		typ            string
		managedDevices []string
		expectLeaves   map[string]any
	}{
		{
			typ: "",
			expectLeaves: map[string]any{
				"/device-metadata/state/metadata/last-seen": agnmi.TypedValue(int(42)),
			},
		},
		{
			typ: NetworkElement,
			expectLeaves: map[string]any{
				"/device-metadata/state/metadata/last-seen": agnmi.TypedValue(int(42)),
			},
		},
		{
			typ: WirelessAP,
			expectLeaves: map[string]any{
				"/device-metadata/state/metadata/last-seen": agnmi.TypedValue(int(42)),
			},
		},
		{
			typ: DeviceManager,
			expectLeaves: map[string]any{
				"/device-metadata/state/metadata/last-seen": agnmi.TypedValue(int(42)),
			},
		},
		{
			typ:            DeviceManager,
			managedDevices: []string{"1", "2"},
			expectLeaves: map[string]any{
				"/device-metadata/state/metadata/last-seen": agnmi.TypedValue(int(42)),
				"/device-metadata/state/metadata/managed-devices": &gnmi.TypedValue{
					Value: &gnmi.TypedValue_LeaflistVal{
						LeaflistVal: &gnmi.ScalarArray{
							Element: []*gnmi.TypedValue{
								agnmi.TypedValue("1"),
								agnmi.TypedValue("2")},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.typ, func(t *testing.T) {
			c := &v2Client{
				deviceID:   "myswitch",
				deviceType: tc.typ,
			}
			if tc.managedDevices != nil {
				c.SetManagedDevices(tc.managedDevices)
			}
			r := c.heartbeatRequest()
			if err := verifyUpdates(r, tc.expectLeaves, true); err != nil {
				t.Fatalf("Error verifying leaves in set request: %v", err)
			}

			r = c.heartbeatRequest()
			// second heartbeat will not send managed devices again
			delete(tc.expectLeaves, "/device-metadata/state/metadata/managed-devices")
			if err := verifyUpdates(r, tc.expectLeaves, true); err != nil {
				t.Fatalf("Error verifying leaves in set request: %v", err)
			}
		})
	}
}
