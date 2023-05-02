// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmp

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider/mock"
	"github.com/gosnmp/gosnmp"
)

// deviceIDTestCase describes a test of the SNMP DeviceID method: the
// mocked SNMP responses and an expected device ID.
type deviceIDTestCase struct {
	name      string
	responses map[string][]*gosnmp.SnmpPDU
	expected  string
}

// mockget and mockwalk are the SNMP get and walk routines used for
// injecting mocked SNMP data into the polling routines.
func mockget(oids []string,
	responses map[string][]*gosnmp.SnmpPDU) (*gosnmp.SnmpPacket, error) {
	if len(oids) != 1 {
		return nil, fmt.Errorf("Expected one OID, got %d", len(oids))
	}
	oid := oids[0]
	r, ok := responses[oid]
	if !ok {
		return nil, fmt.Errorf("mockget saw unexpected OID %s", oid)
	}
	pkt := &gosnmp.SnmpPacket{}
	for _, pdu := range r {
		pkt.Variables = append(pkt.Variables, *pdu)
	}
	// If the response map is empty, return a noSuchObject.
	if len(r) == 0 {
		pkt.Variables = []gosnmp.SnmpPDU{*pdu(oid, gosnmp.NoSuchObject, nil)}
	}
	return pkt, nil
}
func mockwalk(oid string, walker gosnmp.WalkFunc,
	responses map[string][]*gosnmp.SnmpPDU) error {
	pdus, ok := responses[oid]
	if !ok {
		return fmt.Errorf("mockwalk saw unexpected OID %s", oid)
	}
	for _, p := range pdus {
		if err := walker(*p); err != nil {
			return err
		}
	}
	return nil
}

// Call DeviceID, feeding in the mocked PDUs in tc.responses. Check
// that it returns the expected device ID.
func mockDeviceID(t *testing.T, s *Snmp, tc deviceIDTestCase) {
	s.getter = func(oids []string) (*gosnmp.SnmpPacket, error) {
		return mockget(oids, tc.responses)
	}
	s.walker = func(oid string, walker gosnmp.WalkFunc) error {
		return mockwalk(oid, walker, tc.responses)
	}
	s.now = func() time.Time {
		return time.Unix(1554954972, 0)
	}
	s.monitor = mock.NewMockMonitor()
	s.deviceID = ""
	did, err := s.DeviceID(context.Background())
	if err != nil {
		t.Fatalf("Error in DeviceID: %v", err)
	}
	if did != tc.expected {
		t.Fatalf("Device IDs not equal. Expected %v, got %v", tc.expected, did)
	}
}

var entPhysClassAristaResponse = `
.1.3.6.1.2.1.47.1.1.1.1.5.1 = INTEGER: chassis(3)
.1.3.6.1.2.1.47.1.1.1.1.5.100002001 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002002 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002003 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002004 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002005 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002006 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002051 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002052 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002053 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002054 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002055 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002056 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002101 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002102 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002103 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002105 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002106 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002151 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002152 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002153 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002154 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002155 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100002156 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100601000 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100601100 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100601110 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.5.100601111 = INTEGER: sensor(8)
.1.3.6.1.2.1.47.1.1.1.1.5.100601120 = INTEGER: fan(7)
`

var entPhysSerialNumAristaResponse = `
.1.3.6.1.2.1.47.1.1.1.1.11.1 = STRING: JSH11420017
.1.3.6.1.2.1.47.1.1.1.1.11.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002002 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002003 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002004 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002005 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002006 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002051 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002052 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002053 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002054 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002055 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002056 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002101 = STRING: JPE15200157
.1.3.6.1.2.1.47.1.1.1.1.11.100002102 = STRING: JPE15200256
.1.3.6.1.2.1.47.1.1.1.1.11.100002103 = STRING: JPE17400037
.1.3.6.1.2.1.47.1.1.1.1.11.100002105 = STRING: JPE15214958
.1.3.6.1.2.1.47.1.1.1.1.11.100002106 = STRING: JPE13351729
.1.3.6.1.2.1.47.1.1.1.1.11.100002151 = STRING: JPE15253426
.1.3.6.1.2.1.47.1.1.1.1.11.100002152 = STRING: JPE15253614
.1.3.6.1.2.1.47.1.1.1.1.11.100002153 = STRING: JPE15253369
.1.3.6.1.2.1.47.1.1.1.1.11.100002154 = STRING: JPE15253523
.1.3.6.1.2.1.47.1.1.1.1.11.100002155 = STRING: JPE15253366
.1.3.6.1.2.1.47.1.1.1.1.11.100002156 = STRING: JPE15253556
.1.3.6.1.2.1.47.1.1.1.1.11.100601000 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100601100 = STRING: JPE15253426
.1.3.6.1.2.1.47.1.1.1.1.11.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100601111 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100601120 = STRING:
`

var entPhysClassN9KResponse = `
.1.3.6.1.2.1.47.1.1.1.1.5.10 = INTEGER: stack(11)
.1.3.6.1.2.1.47.1.1.1.1.5.22 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.23 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.149 = INTEGER: chassis(3)
.1.3.6.1.2.1.47.1.1.1.1.5.214 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.215 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.278 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.279 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.342 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.343 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.344 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.470 = INTEGER: powerSupply(6)
.1.3.6.1.2.1.47.1.1.1.1.5.471 = INTEGER: powerSupply(6)
.1.3.6.1.2.1.47.1.1.1.1.5.534 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.5.535 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.5.536 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.5.598 = INTEGER: other(1)
.1.3.6.1.2.1.47.1.1.1.1.5.5206 = INTEGER: port(10)
.1.3.6.1.2.1.47.1.1.1.1.5.5207 = INTEGER: port(10)
.1.3.6.1.2.1.47.1.1.1.1.5.5208 = INTEGER: port(10)
`

var entPhysSerialNumN9KResponse = `
.1.3.6.1.2.1.47.1.1.1.1.11.10 = STRING: SAL1817R822
.1.3.6.1.2.1.47.1.1.1.1.11.22 = STRING: SAL1817R822
.1.3.6.1.2.1.47.1.1.1.1.11.23 = STRING: SAL1807M59Z
.1.3.6.1.2.1.47.1.1.1.1.11.149 = STRING: SAL1817R822
.1.3.6.1.2.1.47.1.1.1.1.11.214 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.215 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.278 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.279 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.342 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.343 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.344 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.470 = STRING: DCH1815R075
.1.3.6.1.2.1.47.1.1.1.1.11.471 = STRING: DCH1815R07C
.1.3.6.1.2.1.47.1.1.1.1.11.534 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.535 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.536 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.598 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.5206 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.5207 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.5208 = STRING:
`

var entPhysClassDefaultResponse = `
`

var entPhysSerialNumDefaultResponse = `
.1.3.6.1.2.1.47.1.1.1.1.11.1 = STRING: JSH11420018
.1.3.6.1.2.1.47.1.1.1.1.11.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002002 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002003 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002004 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002005 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002006 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002051 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002052 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002053 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002054 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002055 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002056 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002101 = STRING: JPE15200157
.1.3.6.1.2.1.47.1.1.1.1.11.100002102 = STRING: JPE15200256
.1.3.6.1.2.1.47.1.1.1.1.11.100002103 = STRING: JPE17400037
.1.3.6.1.2.1.47.1.1.1.1.11.100002105 = STRING: JPE15214958
.1.3.6.1.2.1.47.1.1.1.1.11.100002106 = STRING: JPE13351729
.1.3.6.1.2.1.47.1.1.1.1.11.100002151 = STRING: JPE15253426
.1.3.6.1.2.1.47.1.1.1.1.11.100002152 = STRING: JPE15253614
`

var lldpChassisIDTypeDefaultResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
`

var lldpChassisIDDefaultResponse = `
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
`

var sysUpTimeInstanceDefaultResponse = `
.1.3.6.1.2.1.1.3.0 = Timeticks: (10643788) 1 day, 5:33:57.88
`

var sysUpTimeInstanceBadUsernameResponse = `
.1.3.6.1.6.3.15.1.1.3.0 = Counter32: 1
`

func TestDeviceID(t *testing.T) {
	s := &Snmp{
		mock:  true,
		gsnmp: &gosnmp.GoSNMP{Target: "1.2.3.4"},
	}
	for _, tc := range []deviceIDTestCase{
		{
			name: "deviceIDArista",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpEntPhysicalClass:     PDUsFromString(entPhysClassAristaResponse),
				snmpEntPhysicalSerialNum: PDUsFromString(entPhysSerialNumAristaResponse),
			},
			expected: "JSH11420017",
		},
		{
			name: "deviceIDN9K",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpEntPhysicalClass:     PDUsFromString(entPhysClassN9KResponse),
				snmpEntPhysicalSerialNum: PDUsFromString(entPhysSerialNumN9KResponse),
			},
			expected: "SAL1817R822",
		},
		{
			name: "noChassisFound",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpEntPhysicalClass:     PDUsFromString(entPhysClassDefaultResponse),
				snmpEntPhysicalSerialNum: PDUsFromString(entPhysSerialNumDefaultResponse),
			},
			expected: "JSH11420018",
		},
		{
			name: "lldpChassisID",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpEntPhysicalClass:        {},
				snmpEntPhysicalSerialNum:    {},
				snmpLldpLocChassisIDSubtype: PDUsFromString(lldpChassisIDTypeDefaultResponse),
				snmpLldpLocChassisID:        PDUsFromString(lldpChassisIDDefaultResponse),
			},
			expected: "00:1c:73:03:13:36",
		},
		{
			name: "badChassisIDType",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpEntPhysicalClass:     {},
				snmpEntPhysicalSerialNum: {},
				snmpLldpLocChassisIDSubtype: {{
					Name:  ".1.3.6.1.6.3.15.1.1.3.0",
					Type:  gosnmp.Counter32,
					Value: uint(12345),
				}},
				snmpLldpLocChassisID:          {},
				snmpLldpV2LocChassisIDSubtype: {},
				snmpLldpV2LocChassisID:        {},
			},
			expected: "1.2.3.4",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockDeviceID(t, s, tc)
		})
	}
}

func TestTestGet(t *testing.T) {
	s := &Snmp{
		mock:  true,
		gsnmp: &gosnmp.GoSNMP{Target: "1.2.3.4"},
	}

	type testGetTestCase struct {
		name        string
		responses   map[string][]*gosnmp.SnmpPDU
		expectedErr error
	}

	for _, tc := range []testGetTestCase{
		{
			name: "basic",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpSysUpTimeInstance: PDUsFromString(sysUpTimeInstanceDefaultResponse),
			},
		},
		{
			name: "bad username",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpSysUpTimeInstance: PDUsFromString(sysUpTimeInstanceBadUsernameResponse),
			},
			expectedErr: errors.New("unexpected response from SNMP server in " +
				"snmpNetworkInit: [{Value:1 Name:.1.3.6.1.6.3.15.1.1.3.0 " +
				"Type:Counter32}]"),
		},
		{
			name: "NoSuchObject",
			responses: map[string][]*gosnmp.SnmpPDU{
				snmpSysUpTimeInstance: nil,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s.getter = func(oids []string) (*gosnmp.SnmpPacket, error) {
				return mockget(oids, tc.responses)
			}
			err := s.doTestGet()
			if (err == nil) != (tc.expectedErr == nil) ||
				(err != nil && err.Error() != tc.expectedErr.Error()) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}
