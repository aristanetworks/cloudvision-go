// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package v2

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"

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
		if checkValues && !reflect.DeepEqual(v, u.Val) {
			return fmt.Errorf("unexpected value for leaf %s, expected: %v, got: %v", p, v, u.Val)
		}
	}
	return nil
}

func verifyMetadataLeaves(r *gnmi.SetRequest, c *v2Client) error {
	expData := map[string]interface{}{
		"/device-metadata/state/metadata/type":              pgnmi.Strval(c.deviceType),
		"/device-metadata/state/metadata/collector-version": pgnmi.Strval(versionString),
		"/device-metadata/state/metadata/ip-addr":           pgnmi.Strval(testDevice{}.IPAddr()),
	}
	return verifyUpdates(r, expData, true)
}

type testDevice struct{}

func (td testDevice) DeviceID() (string, error) {
	return "mycontroller", nil
}

func (td testDevice) Alive() (bool, error) {
	return true, nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return []provider.Provider{}, nil
}

func (td testDevice) Type() string {
	return ""
}

func (td testDevice) IPAddr() string {
	return "192.168.5.10"
}

func TestMetadataRequest(t *testing.T) {
	c := NewV2Client(nil, testDevice{}).(*v2Client)
	r := c.metadataRequest()
	if err := verifyMetadataLeaves(r, c); err != nil {
		t.Logf("Error verifying leaves in set request: %v", err)
		t.Fail()
	}
}

func verifyHeartbeatLeaves(r *gnmi.SetRequest, c *v2Client) error {
	expData := map[string]interface{}{
		"/device-metadata/state/metadata/last-seen": true,
	}
	return verifyUpdates(r, expData, false)
}

func TestHeartbeatRequest(t *testing.T) {
	c := &v2Client{
		deviceID: "myswitch",
	}
	r := c.heartbeatRequest()
	if err := verifyHeartbeatLeaves(r, c); err != nil {
		t.Logf("Error verifying leaves in set request: %v", err)
		t.Fail()
	}
}
