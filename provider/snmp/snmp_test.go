// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package snmp

import (
	pgnmi "arista/provider/gnmi"
	"arista/provider/openconfig"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

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
		t.Fatalf("SetRequests not equal. Expected %v\nGot %v", tc.expected, sr)
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

func parsePDU(line string) (oid, pduTypeString, value string) {
	t := strings.Split(line, " = ")
	if len(t) < 2 {
		return "", "", ""
	}
	oid = t[0]
	t = strings.Split(t[1], ": ")
	pduTypeString = t[0]
	if len(t) >= 2 {
		value = t[1]
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
		errc:             make(chan error),
		interfaceIndex:   make(map[string]string),
		interfaceName:    make(map[string]bool),
		lldpLocPortIndex: make(map[string]string),
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockPoll(t, s, tc)
		})
	}
}