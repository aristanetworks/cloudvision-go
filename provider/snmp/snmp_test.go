// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package snmp

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

// pollTestCase describes a test of an SNMP polling routine: the routine
// to test, the mocked SNMP responses for the poller, and the gNMI
// SetRequest we expect the poller to return.
type pollTestCase struct {
	name      string
	pollFn    func() (*gnmi.SetRequest, error)
	responses map[string][]gosnmp.SnmpPDU //OID -> PDUs
	expected  *gnmi.SetRequest
}

// deviceIDTestCase describes a test of the SNMP DeviceID method: the
// mocked SNMP responses and an expected device ID.
type deviceIDTestCase struct {
	name      string
	responses map[string][]gosnmp.SnmpPDU
	expected  string
}

// mockget and mockwalk are the SNMP get and walk routines used for
// injecting mocked SNMP data into the polling routines.
func mockget(oids []string,
	responses map[string][]gosnmp.SnmpPDU) (*gosnmp.SnmpPacket, error) {
	if len(oids) != 1 {
		return nil, fmt.Errorf("Expected one OID, got %d", len(oids))
	}
	oid := oids[0]
	r, ok := responses[oid]
	if !ok {
		return nil, fmt.Errorf("mockget saw unexpected OID %s", oid)
	}
	pkt := &gosnmp.SnmpPacket{
		Variables: r,
	}
	return pkt, nil
}
func mockwalk(oid string, walker gosnmp.WalkFunc,
	responses map[string][]gosnmp.SnmpPDU) error {
	pdus, ok := responses[oid]
	if !ok {
		return fmt.Errorf("mockwalk saw unexpected OID %s", oid)
	}
	for _, p := range pdus {
		if err := walker(p); err != nil {
			return err
		}
	}
	return nil
}

// Perform a poll with tc.pollFn, feeding it the mocked PDUs in
// tc.responses. Check that the poller returns the expected updates.
func mockPoll(t *testing.T, s *Snmp, tc pollTestCase) {
	// Set the provider's getter and walker to use the provided data.
	s.getter = func(oids []string) (*gosnmp.SnmpPacket, error) {
		return mockget(oids, tc.responses)
	}
	s.walker = func(oid string, walker gosnmp.WalkFunc) error {
		return mockwalk(oid, walker, tc.responses)
	}

	sr, err := tc.pollFn()
	if err != nil {
		t.Fatalf("Error in pollFn: %v", err)
	}
	if !reflect.DeepEqual(sr, tc.expected) {
		for i, d := range sr.Delete {
			if i >= len(tc.expected.Delete) {
				t.Fatalf("Got unexpected delete: %v", d)
			} else if !reflect.DeepEqual(d, tc.expected.Delete[i]) {
				t.Fatalf("Expected delete: %v, Got: %v", tc.expected.Delete[i], d)
			}
		}
		if len(tc.expected.Delete) > len(sr.Delete) {
			t.Fatalf("Expected but did not see delete: %v",
				tc.expected.Delete[len(sr.Delete)])
		}
		for i, r := range sr.Replace {
			if i >= len(tc.expected.Replace) {
				t.Fatalf("Got unexpected replace: %v", r)
			} else if !reflect.DeepEqual(r, tc.expected.Replace[i]) {
				t.Fatalf("Expected: %v, Got: %v", tc.expected.Replace[i], r)
			}
		}
		if len(tc.expected.Replace) > len(sr.Replace) {
			t.Fatalf("Expected but did not see replace: %v",
				tc.expected.Replace[len(sr.Replace)])
		}
		t.Fatalf("SetRequests not equal. Expected %v\nGot %v", tc.expected, sr)
	}
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
	did, err := s.DeviceID()
	if err != nil {
		t.Fatalf("Error in DeviceID: %v", err)
	}
	if did != tc.expected {
		t.Fatalf("Device IDs not equal. Expected %v, got %v", tc.expected, did)
	}
}

// PDU creation wrapper.
func pdu(name string, t gosnmp.Asn1BER, val interface{}) gosnmp.SnmpPDU {
	return gosnmp.SnmpPDU{
		Name:  name,
		Type:  t,
		Value: val,
	}
}

// SNMP PDU types of interest.
const (
	octstr  = gosnmp.OctetString
	counter = gosnmp.Counter32
	integer = gosnmp.Integer
)

var basicEntPhysicalTableResponse = `
.1.3.6.1.2.1.47.1.1.1.1.2.1 = STRING: DCS-7504 Chassis
.1.3.6.1.2.1.47.1.1.1.1.2.100002001 = STRING: Supervisor Slot 1
.1.3.6.1.2.1.47.1.1.1.1.2.100002101 = STRING: DCS-7500E-SUP Supervisor Module
.1.3.6.1.2.1.47.1.1.1.1.2.100601110 = STRING: Fan Tray 1 Fan 1
.1.3.6.1.2.1.47.1.1.1.1.5.1 = INTEGER: chassis(3)
.1.3.6.1.2.1.47.1.1.1.1.5.100002001 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002101 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100601110 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.7.1 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100002101 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.8.1 = STRING: 02.00
.1.3.6.1.2.1.47.1.1.1.1.8.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.8.100002101 = STRING: 02.02
.1.3.6.1.2.1.47.1.1.1.1.8.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.1 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.100002101 = STRING: 4.21.0F
.1.3.6.1.2.1.47.1.1.1.1.10.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.1 = STRING: JSH11420017
.1.3.6.1.2.1.47.1.1.1.1.11.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002101 = STRING: JPE15200157
.1.3.6.1.2.1.47.1.1.1.1.11.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.12.1 = STRING: Arista Networks
.1.3.6.1.2.1.47.1.1.1.1.12.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.12.100002101 = STRING: Arista Networks
.1.3.6.1.2.1.47.1.1.1.1.12.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.13.1 = STRING: DCS-7504
.1.3.6.1.2.1.47.1.1.1.1.13.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.13.100002101 = STRING: DCS-7500E-SUP
.1.3.6.1.2.1.47.1.1.1.1.13.100601110 = STRING:
`

// snmpwalk responses for six interfaces, four interface types.
var basicIfTableResponse = `
.1.3.6.1.2.1.2.2.1.1.3001 = INTEGER: 3001
.1.3.6.1.2.1.2.2.1.1.3002 = INTEGER: 3002
.1.3.6.1.2.1.2.2.1.1.999011 = INTEGER: 999011
.1.3.6.1.2.1.2.2.1.1.1000001 = INTEGER: 1000001
.1.3.6.1.2.1.2.2.1.1.2001610 = INTEGER: 2001610
.1.3.6.1.2.1.2.2.1.1.5000000 = INTEGER: 5000000
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.2.2.1.2.999011 = STRING: Management1/1
.1.3.6.1.2.1.2.2.1.2.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.2.2.1.2.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.2.2.1.2.5000000 = STRING: Loopback0
.1.3.6.1.2.1.2.2.1.3.3001 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.3002 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.999011 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.1000001 = INTEGER: 161
.1.3.6.1.2.1.2.2.1.3.2001610 = INTEGER: 136
.1.3.6.1.2.1.2.2.1.3.5000000 = INTEGER: 24
.1.3.6.1.2.1.2.2.1.6.3001 = STRING: 74:83:ef:f:6b:6d
.1.3.6.1.2.1.2.2.1.6.3002 = STRING: 74:83:ef:f:6b:6e
.1.3.6.1.2.1.2.2.1.6.999011 = STRING: 0:1c:73:d6:22:c7
.1.3.6.1.2.1.2.2.1.6.1000001 = STRING: 0:1c:73:3c:a2:d3
.1.3.6.1.2.1.2.2.1.6.2001610 = STRING: 0:1c:73:3:13:36
.1.3.6.1.2.1.2.2.1.6.5000000 = STRING:
.1.3.6.1.2.1.2.2.1.10.3001 = Counter32: 1193032336
.1.3.6.1.2.1.2.2.1.10.3002 = Counter32: 3571250484
.1.3.6.1.2.1.2.2.1.10.999011 = Counter32: 102453638
.1.3.6.1.2.1.2.2.1.10.1000001 = Counter32: 3997687336
.1.3.6.1.2.1.2.2.1.10.2001610 = Counter32: 0
.1.3.6.1.2.1.2.2.1.10.5000000 = Counter32: 0
`

var basicIfXTableResponse = `
.1.3.6.1.2.1.31.1.1.1.1.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.31.1.1.1.1.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.31.1.1.1.1.999011 = STRING: Management1/1
.1.3.6.1.2.1.31.1.1.1.1.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.31.1.1.1.1.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.31.1.1.1.1.5000000 = STRING: Loopback0
.1.3.6.1.2.1.31.1.1.1.4.3001 = Counter32: 1303028
.1.3.6.1.2.1.31.1.1.1.4.3002 = Counter32: 5498034
.1.3.6.1.2.1.31.1.1.1.4.999011 = Counter32: 210209
.1.3.6.1.2.1.31.1.1.1.4.1000001 = Counter32: 17708654
.1.3.6.1.2.1.31.1.1.1.4.2001610 = Counter32: 0
.1.3.6.1.2.1.31.1.1.1.4.5000000 = Counter32: 0
`

var basicLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.1 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.2 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.1 = STRING: Management1/1
.1.0.8802.1.1.2.1.3.7.1.3.451 = STRING: Ethernet3/1
.1.0.8802.1.1.2.1.3.7.1.3.452 = STRING: Ethernet3/2
`

var lldpLocalSystemDataResponseStringID = `
.1.0.8802.1.1.2.1.3.2.0 = STRING: 50:87:89:a1:64:4f
`

var basicLldpRemTableResponse = `
.1.0.8802.1.1.2.1.4.1.1.4.0.1.1 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.451.3 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.451.4 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.5.0.1.1 = Hex-STRING: 00 1C 73 0C 97 60
.1.0.8802.1.1.2.1.4.1.1.5.0.451.3 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.451.4 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.6.0.1.1 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.451.3 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.451.4 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.7.0.1.1 = STRING: Ethernet41
.1.0.8802.1.1.2.1.4.1.1.7.0.451.3 = STRING: p255p1
.1.0.8802.1.1.2.1.4.1.1.7.0.451.4 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.9.0.1.1 = STRING: r1-rack1-tor1.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.9.0.451.3 = STRING: server123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.9.0.451.4 = STRING: server123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.10.0.1.1 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.4.1.1.10.0.451.3 = STRING: Linux x.y.z
.1.0.8802.1.1.2.1.4.1.1.10.0.451.4 = STRING: Linux x.y.z
`

var basicLldpStatisticsResponse = `
.1.0.8802.1.1.2.1.2.6.1.2.1 = Counter32: 210277
.1.0.8802.1.1.2.1.2.6.1.2.2 = Counter32: 0
.1.0.8802.1.1.2.1.2.6.1.2.4 = Counter32: 0
.1.0.8802.1.1.2.1.2.6.1.2.451 = Counter32: 210214
.1.0.8802.1.1.2.1.2.6.1.2.452 = Counter32: 207597
.1.0.8802.1.1.2.1.2.7.1.2.1 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.2 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.4 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.451 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.452 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.453 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.454 = Counter32: 0
`

var basicLldpV2IntfSetupResponse = `
.1.3.6.1.2.1.2.2.1.1.6 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.1.18 = INTEGER: 18
.1.3.6.1.2.1.2.2.1.1.19 = INTEGER: 19
.1.3.6.1.2.1.2.2.1.2.6 = STRING: ethernet1/1
.1.3.6.1.2.1.2.2.1.2.18 = STRING: ethernet1/13
.1.3.6.1.2.1.2.2.1.2.19 = STRING: ethernet1/14
`

var basicLldpV2LocalSystemDataResponse = `
.1.3.111.2.802.1.1.13.1.3.1.0 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.3.2.0 = STRING: 24:0b:0a:00:70:98
.1.3.111.2.802.1.1.13.1.3.3.0 = STRING: firewall337-PAN3060
.1.3.111.2.802.1.1.13.1.3.4.0 = STRING: Palo Alto Networks 3000 series firewall
.1.3.111.2.802.1.1.13.1.3.5.0 = Hex-STRING: E8 00
.1.3.111.2.802.1.1.13.1.3.6.0 = Hex-STRING: C0 00
.1.3.111.2.802.1.1.13.1.3.7.1.2.6 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.2.18 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.2.19 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.3.6 = STRING: ethernet1/1
.1.3.111.2.802.1.1.13.1.3.7.1.3.18 = STRING: ethernet1/13
.1.3.111.2.802.1.1.13.1.3.7.1.3.19 = STRING: ethernet1/14
`

var basicLldpV2RemTableResponse = `
.1.3.111.2.802.1.1.13.1.4.1.1.5.0.18.1.1 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.4.1.1.5.0.19.1.2 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.4.1.1.6.0.18.1.1 = STRING: 28:99:3a:bf:26:46
.1.3.111.2.802.1.1.13.1.4.1.1.6.0.19.1.2 = STRING: 28:99:3a:bf:23:f6
.1.3.111.2.802.1.1.13.1.4.1.1.7.0.18.1.1 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.4.1.1.7.0.19.1.2 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.4.1.1.8.0.18.1.1 = STRING: Ethernet46
.1.3.111.2.802.1.1.13.1.4.1.1.8.0.19.1.2 = STRING: Ethernet46
.1.3.111.2.802.1.1.13.1.4.1.1.10.0.18.1.1 = STRING: switch123.sjc.aristanetworks.com
.1.3.111.2.802.1.1.13.1.4.1.1.10.0.19.1.2 = STRING: switch124.sjc.aristanetworks.com
.1.3.111.2.802.1.1.13.1.4.1.1.11.0.18.1.1 = STRING: Arista Networks EOS version x.y.z
.1.3.111.2.802.1.1.13.1.4.1.1.11.0.19.1.2 = STRING: Arista Networks EOS version x.y.z
.1.3.111.2.802.1.1.13.1.4.1.1.12.0.18.1.1 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.12.0.19.1.2 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.13.0.18.1.1 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.13.0.19.1.2 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.14.0.18.1.1 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.14.0.19.1.2 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.15.0.18.1.1 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.15.0.19.1.2 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.2.1.3.1.9.49.53.48.46.48.46.48.46.49 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.2.1.3.1.9.49.53.48.46.48.46.48.46.50 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.3.1.2.127 = STRING: 0x00120f0427d8
`

var basicLldpV2StatisticsResponse = `
.1.3.111.2.802.1.1.13.1.2.6.1.3.18.1 = Counter32: 118331
.1.3.111.2.802.1.1.13.1.2.6.1.3.19.1 = Counter32: 118329
.1.3.111.2.802.1.1.13.1.2.7.1.5.18.2 = Counter32: 118219
.1.3.111.2.802.1.1.13.1.2.7.1.5.19.3 = Counter32: 118194
`

// Include another interface, Ethernet3/3, that's inactive.
var inactiveIntfLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.1 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.2 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.1 = STRING: Management1/1
.1.0.8802.1.1.2.1.3.7.1.3.451 = STRING: Ethernet3/1
.1.0.8802.1.1.2.1.3.7.1.3.452 = STRING: Ethernet3/2
.1.0.8802.1.1.2.1.3.7.1.3.453 = STRING: Ethernet3/3
`

var entPhysTableAristaResponse = `
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

var entPhysTableN9KResponse = `
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

var entPhysTableDefaultResponse = `
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

func parsePDU(line string) (oid, pduTypeString, value string) {
	t := strings.Split(line, " = ")
	if len(t) < 2 {
		return "", "", ""
	}
	oid = t[0]
	t = strings.Split(t[1], ": ")
	pduTypeString = t[0]
	if len(t) >= 2 {
		// Handle case where value is of format "chassis(3)".
		s := strings.Split(strings.Split(t[1], ")")[0], "(")
		value = s[0]
		if len(s) > 1 {
			value = s[1]
		}
	} else {
		pduTypeString = strings.Split(t[0], ":")[0]
	}
	return oid, pduTypeString, value
}

// Convert a set of formatted SNMP responses (as returned by
// snmpwalk -v3 -O ne <target> <oid>) to PDUs.
func pdusFromString(s string) []gosnmp.SnmpPDU {
	pdus := make([]gosnmp.SnmpPDU, 0)
	for _, line := range strings.Split(s, "\n") {
		oid, pduTypeString, val := parsePDU(line)
		if oid == "" {
			continue
		}

		var pduType gosnmp.Asn1BER
		var value interface{}
		switch pduTypeString {
		case "INTEGER":
			pduType = integer
			v, _ := strconv.ParseInt(val, 10, 32)
			value = int(v)
		case "STRING":
			pduType = octstr
			value = []byte(val)
		case "Hex-STRING":
			pduType = octstr
			s := strings.Replace(val, " ", "", -1)
			value, _ = hex.DecodeString(strings.Replace(s, " ", "", -1))
		case "Counter32":
			pduType = counter
			v, _ := strconv.ParseUint(val, 10, 32)
			value = uint(v)
		default:
			panic("Shouldn't get here")
		}
		pdus = append(pdus, pdu(oid, pduType, value))
	}
	return pdus
}

func TestSnmp(t *testing.T) {
	s := &Snmp{
		errc:          make(chan error),
		interfaceName: make(map[string]bool),
	}
	for _, tc := range []pollTestCase{
		{
			name:   "updateSystemStateBasic",
			pollFn: s.updateSystemState,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpSysName: []gosnmp.SnmpPDU{
					pdu(snmpSysName, octstr, []byte("device123.sjc.aristanetworks.com")),
				},
			},
			expected: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					update(pgnmi.Path("system", "state", "hostname"), strval("device123")),
					update(pgnmi.Path("system", "state", "domain-name"),
						strval("sjc.aristanetworks.com")),
				},
			},
		},
		{
			name:   "updatePlatformBasic",
			pollFn: s.updatePlatform,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpEntPhysicalTable: pdusFromString(basicEntPhysicalTableResponse),
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("components")},
				Replace: []*gnmi.Update{
					update(pgnmi.PlatformComponentConfigPath("1", "name"),
						strval("1")),
					update(pgnmi.PlatformComponentPath("1", "name"),
						strval("1")),
					update(pgnmi.PlatformComponentStatePath("1", "name"),
						strval("1")),
					update(pgnmi.PlatformComponentStatePath("1", "id"),
						strval("1")),
					update(pgnmi.PlatformComponentStatePath("1", "description"),
						strval("DCS-7504 Chassis")),
					update(pgnmi.PlatformComponentConfigPath("100002001", "name"),
						strval("100002001")),
					update(pgnmi.PlatformComponentPath("100002001", "name"),
						strval("100002001")),
					update(pgnmi.PlatformComponentStatePath("100002001", "name"),
						strval("100002001")),
					update(pgnmi.PlatformComponentStatePath("100002001", "id"),
						strval("100002001")),
					update(pgnmi.PlatformComponentStatePath("100002001", "description"),
						strval("Supervisor Slot 1")),
					update(pgnmi.PlatformComponentConfigPath("100002101", "name"),
						strval("100002101")),
					update(pgnmi.PlatformComponentPath("100002101", "name"),
						strval("100002101")),
					update(pgnmi.PlatformComponentStatePath("100002101", "name"),
						strval("100002101")),
					update(pgnmi.PlatformComponentStatePath("100002101", "id"),
						strval("100002101")),
					update(pgnmi.PlatformComponentStatePath("100002101", "description"),
						strval("DCS-7500E-SUP Supervisor Module")),
					update(pgnmi.PlatformComponentConfigPath("100601110", "name"),
						strval("100601110")),
					update(pgnmi.PlatformComponentPath("100601110", "name"),
						strval("100601110")),
					update(pgnmi.PlatformComponentStatePath("100601110", "name"),
						strval("100601110")),
					update(pgnmi.PlatformComponentStatePath("100601110", "id"),
						strval("100601110")),
					update(pgnmi.PlatformComponentStatePath("100601110", "description"),
						strval("Fan Tray 1 Fan 1")),
					update(pgnmi.PlatformComponentStatePath("1", "type"),
						strval("CHASSIS")),
					update(pgnmi.PlatformComponentStatePath("100601110", "type"),
						strval("FAN")),
					update(pgnmi.PlatformComponentStatePath("1", "software-version"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("100002001", "software-version"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("100002101", "software-version"),
						strval("4.21.0F")),
					update(pgnmi.PlatformComponentStatePath("100601110", "software-version"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("1", "serial-no"),
						strval("JSH11420017")),
					update(pgnmi.PlatformComponentStatePath("100002001", "serial-no"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("100002101", "serial-no"),
						strval("JPE15200157")),
					update(pgnmi.PlatformComponentStatePath("100601110", "serial-no"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("1", "mfg-name"),
						strval("Arista Networks")),
					update(pgnmi.PlatformComponentStatePath("100002001", "mfg-name"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("100002101", "mfg-name"),
						strval("Arista Networks")),
					update(pgnmi.PlatformComponentStatePath("100601110", "mfg-name"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("1", "hardware-version"),
						strval("DCS-7504")),
					update(pgnmi.PlatformComponentStatePath("100002001", "hardware-version"),
						strval("")),
					update(pgnmi.PlatformComponentStatePath("100002101", "hardware-version"),
						strval("DCS-7500E-SUP")),
					update(pgnmi.PlatformComponentStatePath("100601110", "hardware-version"),
						strval("")),
				},
			},
		},
		{
			name:   "updateInterfacesBasic",
			pollFn: s.updateInterfaces,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpIfTable:  pdusFromString(basicIfTableResponse),
				snmpIfXTable: pdusFromString(basicIfXTableResponse),
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("interfaces", "interface")},
				Replace: []*gnmi.Update{
					update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
					update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
					update(pgnmi.IntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
					update(pgnmi.IntfStatePath("Management1/1", "name"), strval("Management1/1")),
					update(pgnmi.IntfPath("Management1/1", "name"), strval("Management1/1")),
					update(pgnmi.IntfConfigPath("Management1/1", "name"), strval("Management1/1")),
					update(pgnmi.IntfStatePath("Port-Channel1", "name"), strval("Port-Channel1")),
					update(pgnmi.IntfPath("Port-Channel1", "name"), strval("Port-Channel1")),
					update(pgnmi.IntfConfigPath("Port-Channel1", "name"), strval("Port-Channel1")),
					update(pgnmi.IntfStatePath("Vlan1610", "name"), strval("Vlan1610")),
					update(pgnmi.IntfPath("Vlan1610", "name"), strval("Vlan1610")),
					update(pgnmi.IntfConfigPath("Vlan1610", "name"), strval("Vlan1610")),
					update(pgnmi.IntfStatePath("Loopback0", "name"), strval("Loopback0")),
					update(pgnmi.IntfPath("Loopback0", "name"), strval("Loopback0")),
					update(pgnmi.IntfConfigPath("Loopback0", "name"), strval("Loopback0")),
					update(pgnmi.IntfStatePath("Ethernet3/1", "type"), strval("ethernetCsmacd")),
					update(pgnmi.IntfStatePath("Ethernet3/2", "type"), strval("ethernetCsmacd")),
					update(pgnmi.IntfStatePath("Management1/1", "type"), strval("ethernetCsmacd")),
					update(pgnmi.IntfStatePath("Port-Channel1", "type"), strval("ieee8023adLag")),
					update(pgnmi.IntfStatePath("Vlan1610", "type"), strval("l3ipvlan")),
					update(pgnmi.IntfStatePath("Loopback0", "type"), strval("softwareLoopback")),
					update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-octets"),
						uintval(uint(1193032336))),
					update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-octets"),
						uintval(uint(3571250484))),
					update(pgnmi.IntfStateCountersPath("Management1/1", "in-octets"),
						uintval(uint(102453638))),
					update(pgnmi.IntfStateCountersPath("Port-Channel1", "in-octets"),
						uintval(uint(3997687336))),
					update(pgnmi.IntfStateCountersPath("Vlan1610", "in-octets"), uintval(0)),
					update(pgnmi.IntfStateCountersPath("Loopback0", "in-octets"), uintval(0)),
					update(pgnmi.IntfStateCountersPath("Ethernet3/1", "out-multicast-pkts"),
						uintval(1303028)),
					update(pgnmi.IntfStateCountersPath("Ethernet3/2", "out-multicast-pkts"),
						uintval(5498034)),
					update(pgnmi.IntfStateCountersPath("Management1/1", "out-multicast-pkts"),
						uintval(210209)),
					update(pgnmi.IntfStateCountersPath("Port-Channel1", "out-multicast-pkts"),
						uintval(17708654)),
					update(pgnmi.IntfStateCountersPath("Vlan1610", "out-multicast-pkts"),
						uintval(0)),
					update(pgnmi.IntfStateCountersPath("Loopback0", "out-multicast-pkts"),
						uintval(0)),
				},
			},
		},
		{
			name:   "updateLldpBasic",
			pollFn: s.updateLldp,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpLldpLocalSystemData: pdusFromString(basicLldpLocalSystemDataResponse),
				snmpLldpRemTable:        pdusFromString(basicLldpRemTableResponse),
				snmpLldpStatistics:      pdusFromString(basicLldpStatisticsResponse),
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("lldp")},
				Replace: []*gnmi.Update{
					update(pgnmi.LldpStatePath("chassis-id-type"),
						strval(openconfig.LLDPChassisIDType(4))),
					update(pgnmi.LldpStatePath("chassis-id"), strval("00:1c:73:03:13:36")),
					update(pgnmi.LldpStatePath("system-name"),
						strval("device123.sjc.aristanetworks.com")),
					update(pgnmi.LldpStatePath("system-description"),
						strval("Arista Networks EOS version x.y.z")),
					update(pgnmi.LldpIntfConfigPath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfPath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfStatePath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfConfigPath("Ethernet3/1", "name"),
						strval("Ethernet3/1")),
					update(pgnmi.LldpIntfPath("Ethernet3/1", "name"),
						strval("Ethernet3/1")),
					update(pgnmi.LldpIntfStatePath("Ethernet3/1", "name"),
						strval("Ethernet3/1")),
					update(pgnmi.LldpIntfConfigPath("Ethernet3/2", "name"),
						strval("Ethernet3/2")),
					update(pgnmi.LldpIntfPath("Ethernet3/2", "name"),
						strval("Ethernet3/2")),
					update(pgnmi.LldpIntfStatePath("Ethernet3/2", "name"),
						strval("Ethernet3/2")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "id"),
						strval("1")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "chassis-id-type"),
						strval("MAC_ADDRESS")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "id"), strval("3")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id-type"),
						strval("MAC_ADDRESS")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "id"), strval("4")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "chassis-id-type"),
						strval("MAC_ADDRESS")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "chassis-id"),
						strval("00:1c:73:0c:97:60")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id"),
						strval("02:82:9b:3e:e5:fa")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "chassis-id"),
						strval("02:82:9b:3e:e5:fa")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "port-id-type"),
						strval("INTERFACE_NAME")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id-type"),
						strval("INTERFACE_NAME")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "port-id-type"),
						strval("INTERFACE_NAME")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "port-id"),
						strval("Ethernet41")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id"),
						strval("p255p1")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "port-id"),
						strval("macvlan-bond0")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "system-name"),
						strval("r1-rack1-tor1.sjc.aristanetworks.com")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "system-name"),
						strval("server123.sjc.aristanetworks.com")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "system-name"),
						strval("server123.sjc.aristanetworks.com")),
					update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "system-description"),
						strval("Arista Networks EOS version x.y.z")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "system-description"),
						strval("Linux x.y.z")),
					update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "system-description"),
						strval("Linux x.y.z")),
					update(pgnmi.LldpIntfCountersPath("Management1/1", "frame-out"),
						uintval(210277)),
					update(pgnmi.LldpIntfCountersPath("Ethernet3/1", "frame-out"), uintval(210214)),
					update(pgnmi.LldpIntfCountersPath("Ethernet3/2", "frame-out"), uintval(207597)),
					update(pgnmi.LldpIntfCountersPath("Management1/1", "frame-discard"),
						uintval(0)),
					update(pgnmi.LldpIntfCountersPath("Ethernet3/1", "frame-discard"), uintval(0)),
					update(pgnmi.LldpIntfCountersPath("Ethernet3/2", "frame-discard"), uintval(0)),
				},
			},
		},
		{
			name:   "updateLldpOmitInactiveIntfs",
			pollFn: s.updateLldp,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpLldpLocalSystemData: pdusFromString(inactiveIntfLldpLocalSystemDataResponse),
				snmpLldpRemTable:        []gosnmp.SnmpPDU{},
				snmpLldpStatistics:      []gosnmp.SnmpPDU{},
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("lldp")},
				Replace: []*gnmi.Update{
					update(pgnmi.LldpStatePath("chassis-id-type"),
						strval(openconfig.LLDPChassisIDType(4))),
					update(pgnmi.LldpStatePath("chassis-id"), strval("00:1c:73:03:13:36")),
					update(pgnmi.LldpStatePath("system-name"),
						strval("device123.sjc.aristanetworks.com")),
					update(pgnmi.LldpStatePath("system-description"),
						strval("Arista Networks EOS version x.y.z")),
					update(pgnmi.LldpIntfConfigPath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfPath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfStatePath("Management1/1", "name"),
						strval("Management1/1")),
					update(pgnmi.LldpIntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.LldpIntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.LldpIntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
					update(pgnmi.LldpIntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
					update(pgnmi.LldpIntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
					update(pgnmi.LldpIntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
				},
			},
		},
		{
			name:   "updateLldpStringChassisID",
			pollFn: s.updateLldp,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpLldpLocalSystemData: pdusFromString(lldpLocalSystemDataResponseStringID),
				snmpLldpRemTable:        []gosnmp.SnmpPDU{},
				snmpLldpStatistics:      []gosnmp.SnmpPDU{},
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("lldp")},
				Replace: []*gnmi.Update{
					update(pgnmi.LldpStatePath("chassis-id"), strval("50:87:89:a1:64:4f")),
				},
			},
		},
		{
			name:   "updateSystemStateHostnameOnly",
			pollFn: s.updateSystemState,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpSysName: []gosnmp.SnmpPDU{
					pdu(snmpSysName, octstr, []byte("deviceABC")),
				},
			},
			expected: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					update(pgnmi.Path("system", "state", "hostname"), strval("deviceABC")),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockPoll(t, s, tc)
		})
	}
}

func TestSnmpLldpV2(t *testing.T) {
	s := &Snmp{
		errc:          make(chan error),
		interfaceName: make(map[string]bool),
		lldpV2:        true,
	}
	for _, tc := range []pollTestCase{
		{
			name:   "lldpV2IntfSetup",
			pollFn: s.updateInterfaces,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpIfTable:  pdusFromString(basicLldpV2IntfSetupResponse),
				snmpIfXTable: []gosnmp.SnmpPDU{},
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("interfaces", "interface")},
				Replace: []*gnmi.Update{
					update(pgnmi.IntfStatePath("ethernet1/1", "name"), strval("ethernet1/1")),
					update(pgnmi.IntfPath("ethernet1/1", "name"), strval("ethernet1/1")),
					update(pgnmi.IntfConfigPath("ethernet1/1", "name"), strval("ethernet1/1")),
					update(pgnmi.IntfStatePath("ethernet1/13", "name"), strval("ethernet1/13")),
					update(pgnmi.IntfPath("ethernet1/13", "name"), strval("ethernet1/13")),
					update(pgnmi.IntfConfigPath("ethernet1/13", "name"), strval("ethernet1/13")),
					update(pgnmi.IntfStatePath("ethernet1/14", "name"), strval("ethernet1/14")),
					update(pgnmi.IntfPath("ethernet1/14", "name"), strval("ethernet1/14")),
					update(pgnmi.IntfConfigPath("ethernet1/14", "name"), strval("ethernet1/14")),
				},
			},
		},
		{
			name:   "updateLldpV2Basic",
			pollFn: s.updateLldp,
			responses: map[string][]gosnmp.SnmpPDU{
				snmpLldpV2LocalSystemData: pdusFromString(basicLldpV2LocalSystemDataResponse),
				snmpLldpV2RemTable:        pdusFromString(basicLldpV2RemTableResponse),
				snmpLldpV2Statistics:      pdusFromString(basicLldpV2StatisticsResponse),
			},
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{pgnmi.Path("lldp")},
				Replace: []*gnmi.Update{
					update(pgnmi.LldpStatePath("chassis-id-type"),
						strval(openconfig.LLDPChassisIDType(4))),
					update(pgnmi.LldpStatePath("chassis-id"), strval("24:0b:0a:00:70:98")),
					update(pgnmi.LldpStatePath("system-name"),
						strval("firewall337-PAN3060")),
					update(pgnmi.LldpStatePath("system-description"),
						strval("Palo Alto Networks 3000 series firewall")),
					update(pgnmi.LldpIntfConfigPath("ethernet1/1", "name"),
						strval("ethernet1/1")),
					update(pgnmi.LldpIntfPath("ethernet1/1", "name"),
						strval("ethernet1/1")),
					update(pgnmi.LldpIntfStatePath("ethernet1/1", "name"),
						strval("ethernet1/1")),
					update(pgnmi.LldpIntfConfigPath("ethernet1/13", "name"),
						strval("ethernet1/13")),
					update(pgnmi.LldpIntfPath("ethernet1/13", "name"),
						strval("ethernet1/13")),
					update(pgnmi.LldpIntfStatePath("ethernet1/13", "name"),
						strval("ethernet1/13")),
					update(pgnmi.LldpIntfConfigPath("ethernet1/14", "name"),
						strval("ethernet1/14")),
					update(pgnmi.LldpIntfPath("ethernet1/14", "name"),
						strval("ethernet1/14")),
					update(pgnmi.LldpIntfStatePath("ethernet1/14", "name"),
						strval("ethernet1/14")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "id"),
						strval("1")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "chassis-id-type"),
						strval("MAC_ADDRESS")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "id"), strval("2")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "chassis-id-type"),
						strval("MAC_ADDRESS")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "chassis-id"),
						strval("28:99:3a:bf:26:46")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "chassis-id"),
						strval("28:99:3a:bf:23:f6")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "port-id-type"),
						strval("INTERFACE_NAME")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "port-id-type"),
						strval("INTERFACE_NAME")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "port-id"),
						strval("Ethernet46")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "port-id"),
						strval("Ethernet46")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "system-name"),
						strval("switch123.sjc.aristanetworks.com")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "system-name"),
						strval("switch124.sjc.aristanetworks.com")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "system-description"),
						strval("Arista Networks EOS version x.y.z")),
					update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "system-description"),
						strval("Arista Networks EOS version x.y.z")),
					update(pgnmi.LldpIntfCountersPath("ethernet1/13", "frame-out"),
						uintval(118331)),
					update(pgnmi.LldpIntfCountersPath("ethernet1/14", "frame-out"),
						uintval(118329)),
					update(pgnmi.LldpIntfCountersPath("ethernet1/13", "frame-in"),
						uintval(118219)),
					update(pgnmi.LldpIntfCountersPath("ethernet1/14", "frame-in"), uintval(118194)),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockPoll(t, s, tc)
		})
	}
}

func TestDeviceID(t *testing.T) {
	s := &Snmp{
		errc:          make(chan error),
		interfaceName: make(map[string]bool),
	}
	for _, tc := range []deviceIDTestCase{
		{
			name: "deviceIDArista",
			responses: map[string][]gosnmp.SnmpPDU{
				snmpEntPhysicalTable: pdusFromString(entPhysTableAristaResponse),
			},
			expected: "JSH11420017",
		},
		{
			name: "deviceIDN9K",
			responses: map[string][]gosnmp.SnmpPDU{
				snmpEntPhysicalTable: pdusFromString(entPhysTableN9KResponse),
			},
			expected: "SAL1817R822",
		},
		{
			name: "noChassisFound",
			responses: map[string][]gosnmp.SnmpPDU{
				snmpEntPhysicalTable: pdusFromString(entPhysTableDefaultResponse),
			},
			expected: "JSH11420018",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockDeviceID(t, s, tc)
		})
	}
}
