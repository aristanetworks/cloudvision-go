// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package snmp

import (
	"bytes"
	"cloudvision-go/provider"
	pgnmi "cloudvision-go/provider/gnmi"
	"cloudvision-go/provider/openconfig"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aristanetworks/glog"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

// return base OID, index
func oidIndex(oid string) (string, string, error) {
	finalDotPos := strings.LastIndex(oid, ".")
	if finalDotPos < 0 {
		return "", "", fmt.Errorf("oid '%s' does not match expected format", oid)
	}
	return oid[:finalDotPos], oid[(finalDotPos + 1):], nil
}

const (
	snmpEntPhysicalClass               = ".1.3.6.1.2.1.47.1.1.1.1.5"
	snmpEntPhysicalDescr               = ".1.3.6.1.2.1.47.1.1.1.1.2"
	snmpEntPhysicalMfgName             = ".1.3.6.1.2.1.47.1.1.1.1.12"
	snmpEntPhysicalModelName           = ".1.3.6.1.2.1.47.1.1.1.1.13"
	snmpEntPhysicalSerialNum           = ".1.3.6.1.2.1.47.1.1.1.1.11"
	snmpEntPhysicalSoftwareRev         = ".1.3.6.1.2.1.47.1.1.1.1.10"
	snmpEntPhysicalTable               = ".1.3.6.1.2.1.47.1.1.1.1"
	snmpIfTable                        = ".1.3.6.1.2.1.2.2"
	snmpIfDescr                        = ".1.3.6.1.2.1.2.2.1.2"
	snmpIfType                         = ".1.3.6.1.2.1.2.2.1.3"
	snmpIfMtu                          = ".1.3.6.1.2.1.2.2.1.4"
	snmpIfAdminStatus                  = ".1.3.6.1.2.1.2.2.1.7"
	snmpIfInBroadcastPkts              = ".1.3.6.1.2.1.31.1.1.1.3"
	snmpIfInDiscards                   = ".1.3.6.1.2.1.2.2.1.13"
	snmpIfInErrors                     = ".1.3.6.1.2.1.2.2.1.14"
	snmpIfInMulticastPkts              = ".1.3.6.1.2.1.31.1.1.1.2"
	snmpIfInOctets                     = ".1.3.6.1.2.1.2.2.1.10"
	snmpIfInUcastPkts                  = ".1.3.6.1.2.1.2.2.1.11"
	snmpIfInUnknownProtos              = ".1.3.6.1.2.1.2.2.1.15"
	snmpIfOperStatus                   = ".1.3.6.1.2.1.2.2.1.8"
	snmpIfOutBroadcastPkts             = ".1.3.6.1.2.1.31.1.1.1.5"
	snmpIfOutDiscards                  = ".1.3.6.1.2.1.2.2.1.19"
	snmpIfOutErrors                    = ".1.3.6.1.2.1.2.2.1.20"
	snmpIfOutMulticastPkts             = ".1.3.6.1.2.1.31.1.1.1.4"
	snmpIfOutOctets                    = ".1.3.6.1.2.1.2.2.1.16"
	snmpIfOutUcastPkts                 = ".1.3.6.1.2.1.2.2.1.17"
	snmpIfXTable                       = ".1.3.6.1.2.1.31.1.1"
	snmpLldpLocalSystemData            = ".1.0.8802.1.1.2.1.3"
	snmpLldpLocChassisID               = ".1.0.8802.1.1.2.1.3.2.0"
	snmpLldpLocChassisIDSubtype        = ".1.0.8802.1.1.2.1.3.1.0"
	snmpLldpLocPortID                  = ".1.0.8802.1.1.2.1.3.7.1.3"
	snmpLldpLocPortTable               = ".1.0.8802.1.1.2.1.3.7"
	snmpLldpLocSysDesc                 = ".1.0.8802.1.1.2.1.3.4.0"
	snmpLldpLocSysName                 = ".1.0.8802.1.1.2.1.3.3.0"
	snmpLldpRemPortID                  = ".1.0.8802.1.1.2.1.4.1.1.7"
	snmpLldpRemPortIDSubtype           = ".1.0.8802.1.1.2.1.4.1.1.6"
	snmpLldpRemChassisID               = ".1.0.8802.1.1.2.1.4.1.1.5"
	snmpLldpRemChassisIDSubtype        = ".1.0.8802.1.1.2.1.4.1.1.4"
	snmpLldpRemSysDesc                 = ".1.0.8802.1.1.2.1.4.1.1.10"
	snmpLldpRemSysName                 = ".1.0.8802.1.1.2.1.4.1.1.9"
	snmpLldpRemTable                   = ".1.0.8802.1.1.2.1.4.1"
	snmpLldpStatistics                 = ".1.0.8802.1.1.2.1.2"
	snmpLldpStatsRxPortFramesTotal     = ".1.0.8802.1.1.2.1.2.7.1.4"
	snmpLldpStatsRxPortTable           = ".1.0.8802.1.1.2.1.2.7"
	snmpLldpStatsRxPortTLVsDiscard     = ".1.0.8802.1.1.2.1.2.7.1.5"
	snmpLldpStatsRxPortTLVsUnrecog     = ".1.0.8802.1.1.2.1.2.7.1.6"
	snmpLldpStatsRxPortFramesDiscard   = ".1.0.8802.1.1.2.1.2.7.1.2"
	snmpLldpStatsRxPortFramesErrors    = ".1.0.8802.1.1.2.1.2.7.1.3"
	snmpLldpStatsTxPortFramesTotal     = ".1.0.8802.1.1.2.1.2.6.1.2"
	snmpLldpStatsTxPortTable           = ".1.0.8802.1.1.2.1.2.6"
	snmpLldpV2LocalSystemData          = ".1.3.111.2.802.1.1.13.1.3"
	snmpLldpV2LocChassisID             = ".1.3.111.2.802.1.1.13.1.3.2.0"
	snmpLldpV2LocChassisIDSubtype      = ".1.3.111.2.802.1.1.13.1.3.1.0"
	snmpLldpV2LocPortID                = ".1.3.111.2.802.1.1.13.1.3.7.1.3"
	snmpLldpV2LocPortTable             = ".1.3.111.2.802.1.1.13.1.3.7"
	snmpLldpV2LocSysDesc               = ".1.3.111.2.802.1.1.13.1.3.4.0"
	snmpLldpV2LocSysName               = ".1.3.111.2.802.1.1.13.1.3.3.0"
	snmpLldpV2RemPortID                = ".1.3.111.2.802.1.1.13.1.4.1.1.8"
	snmpLldpV2RemPortIDSubtype         = ".1.3.111.2.802.1.1.13.1.4.1.1.7"
	snmpLldpV2RemChassisID             = ".1.3.111.2.802.1.1.13.1.4.1.1.6"
	snmpLldpV2RemChassisIDSubtype      = ".1.3.111.2.802.1.1.13.1.4.1.1.5"
	snmpLldpV2RemSysDesc               = ".1.3.111.2.802.1.1.13.1.4.1.1.11"
	snmpLldpV2RemSysName               = ".1.3.111.2.802.1.1.13.1.4.1.1.10"
	snmpLldpV2RemTable                 = ".1.3.111.2.802.1.1.13.1.4.1"
	snmpLldpV2Statistics               = ".1.3.111.2.802.1.1.13.1.2"
	snmpLldpV2StatsRxPortFramesTotal   = ".1.3.111.2.802.1.1.13.1.2.7.1.5"
	snmpLldpV2StatsRxPortTable         = ".1.3.111.2.802.1.1.13.1.2.7"
	snmpLldpV2StatsRxPortTLVsDiscard   = "1.3.111.2.802.1.1.13.1.2.7.1.6"
	snmpLldpV2StatsRxPortTLVsUnrecog   = ".1.3.111.2.802.1.1.13.1.2.7.1.7"
	snmpLldpV2StatsRxPortFramesDiscard = ".1.3.111.2.802.1.1.13.1.2.7.1.3"
	snmpLldpV2StatsRxPortFramesErrors  = ".1.3.111.2.802.1.1.13.1.2.7.1.4"
	snmpLldpV2StatsTxPortFramesTotal   = ".1.3.111.2.802.1.1.13.1.2.6.1.3"
	snmpLldpV2StatsTxPortTable         = ".1.3.111.2.802.1.1.13.1.2.6"
	snmpSysName                        = ".1.3.6.1.2.1.1.5.0"
	snmpSysUpTime                      = ".1.3.6.1.2.1.1.3.0"
)

// Less typing: gNMI type helpers.
func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return pgnmi.Update(path, val)
}

func strval(s interface{}) *gnmi.TypedValue {
	t, ok := s.(string)
	if ok {
		return pgnmi.Strval(t)
	}
	u, ok := s.([]byte)
	if ok {
		// Remove newlines. OpenConfig will reject multiline strings.
		ss := strings.Replace(string(u), "\n", " ", -1)
		return pgnmi.Strval(ss)
	}
	glog.Fatalf("Unexpected type in strval: %T", s)
	return nil
}

func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToUint64(u); err == nil {
		return pgnmi.Uintval(v)
	}
	return nil
}

// Snmp contains everything needed to implement an SNMP provider.
type Snmp struct {
	errc   chan error
	client gnmi.GNMIClient

	// interfaceName is a map of interface name (as discovered in ifTable) -> true.
	// It's used so that we don't include inactive interfaces we see in
	// snmpLldpLocPortTable.
	interfaceName map[string]bool

	// lldpV2 indicates whether to use LLDP-V2-MIB.
	lldpV2 bool

	address   string
	community string

	// gosnmp can't handle parallel gets.
	lock sync.Mutex

	pollInterval time.Duration
	lastAlive    time.Time
	initialized  bool

	// Alternative Walk() and Get() for mock testing.
	getter func([]string) (*gosnmp.SnmpPacket, error)
	walker func(string, gosnmp.WalkFunc) error
}

func (s *Snmp) snmpNetworkInit() error {
	if s.initialized {
		return nil
	}
	return gosnmp.Default.Connect()
}

func (s *Snmp) getByOID(oids []string) (*gosnmp.SnmpPacket, error) {
	if s.getter == nil {
		return nil, errors.New("SNMP getter not set")
	}
	pkt, err := s.getter(oids)
	if err != nil {
		s.lastAlive = time.Now()
	}
	return pkt, err
}

// getStringByOID does a Get on the specified OID, an octet string, and
// returns the result as a string.
func (s *Snmp) getStringByOID(oid string) (string, error) {
	pkt, err := s.getByOID([]string{oid})
	if err != nil {
		return "", err
	}
	if len(pkt.Variables) == 0 {
		return "", fmt.Errorf("No variables in SNMP packet for OID %s", oid)
	}
	v := pkt.Variables[0]
	if v.Type != gosnmp.OctetString {
		return "", fmt.Errorf("PDU type for OID %s is not octet string", oid)
	}
	return string(v.Value.([]byte)), nil
}

func (s *Snmp) walk(rootOid string, walkFn gosnmp.WalkFunc) error {
	if s.walker == nil {
		return errors.New("SNMP walker not set")
	}
	err := s.walker(rootOid, walkFn)
	if err != nil {
		s.lastAlive = time.Now()
	}
	return err
}

// DeviceID returns the device ID.
func (s *Snmp) DeviceID() (string, error) {
	serial := ""
	var done bool
	chassisIndex := ""
	var snmpEntPhysicalClassTypeChassis = 3

	// Get the serial number corresponding to the index whose class
	// type is chassis(3).
	entPhysicalWalk := func(data gosnmp.SnmpPDU) error {
		if done {
			return nil
		}
		baseOid, index, err := oidIndex(data.Name)
		if err != nil {
			return err
		}
		// If the physical class is "chassis", this is the index we want.
		if baseOid == snmpEntPhysicalClass {
			if data.Value.(int) == snmpEntPhysicalClassTypeChassis {
				chassisIndex = index
			}
		}
		if baseOid == snmpEntPhysicalSerialNum {
			// Take the first non-empty serial number as a backup, in
			// case there isn't a chassis serial number.
			if serial == "" {
				serial = string(data.Value.([]byte))
			}
			if index == chassisIndex {
				serial = string(data.Value.([]byte))
				done = true
			}
		}

		return nil
	}

	if err := s.walk(snmpEntPhysicalTable, entPhysicalWalk); err != nil {
		return "", err
	}
	if serial == "" {
		return "", errors.New("Failed to get serial number")
	}
	return serial, nil
}

// CheckAlive checks if device is still alive if poll interval has passed.
func (s *Snmp) CheckAlive() (bool, error) {
	if time.Since(s.lastAlive) < s.pollInterval {
		return true, nil
	}
	_, err := s.getByOID([]string{snmpSysUpTime})
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *Snmp) stop() {
	gosnmp.Default.Conn.Close()
}

// Given an incoming PDU, update the appropriate interface state.
func (s *Snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU,
	interfaceIndex map[string]string) ([]*gnmi.Update, error) {
	// Get/set interface name from index. If there's no mapping, just return and
	// wait for the mapping to show up.
	baseOid, index, err := oidIndex(pdu.Name)
	if err != nil {
		return nil, err
	}
	intfName, ok := interfaceIndex[index]
	if !ok && baseOid != snmpIfDescr {
		return nil, nil
	} else if !ok && baseOid == snmpIfDescr {
		intfName = string(pdu.Value.([]byte))
		interfaceIndex[index] = intfName
		s.interfaceName[intfName] = true
	}

	var u *gnmi.Update
	switch baseOid {
	case snmpIfDescr:
		u = update(pgnmi.IntfStatePath(intfName, "name"),
			strval(pdu.Value))
	case snmpIfType:
		u = update(pgnmi.IntfStatePath(intfName, "type"),
			strval(openconfig.InterfaceType(pdu.Value.(int))))
	case snmpIfMtu:
		u = update(pgnmi.IntfStatePath(intfName, "mtu"),
			uintval(pdu.Value))
	case snmpIfAdminStatus:
		u = update(pgnmi.IntfStatePath(intfName, "admin-status"),
			strval(openconfig.IntfAdminStatus(pdu.Value.(int))))
	case snmpIfOperStatus:
		u = update(pgnmi.IntfStatePath(intfName, "oper-status"),
			strval(openconfig.IntfOperStatus(pdu.Value.(int))))
	case snmpIfInOctets:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-octets"),
			uintval(pdu.Value))
	case snmpIfInUcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-unicast-pkts"),
			uintval(pdu.Value))
	case snmpIfInMulticastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfInBroadcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-broadcast-pkts"),
			uintval(pdu.Value))
	case snmpIfInDiscards:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-discards"),
			uintval(pdu.Value))
	case snmpIfInErrors:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-errors"),
			uintval(pdu.Value))
	case snmpIfInUnknownProtos:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-unknown-protos"),
			uintval(pdu.Value))
	case snmpIfOutOctets:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-octets"),
			uintval(pdu.Value))
	case snmpIfOutUcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-unicast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutMulticastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutBroadcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-broadcast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutDiscards:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-discards"),
			uintval(pdu.Value))
	case snmpIfOutErrors:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-errors"),
			uintval(pdu.Value))
	default:
		// default: ignore update
		return nil, nil
	}

	updates := []*gnmi.Update{u}
	// When we get a name, add name, config/name, state/name.
	if baseOid == snmpIfDescr {
		updates = append(updates,
			update(pgnmi.IntfPath(intfName, "name"), strval(intfName)),
			update(pgnmi.IntfConfigPath(intfName, "name"), strval(intfName)))
	}
	return updates, nil
}

func (s *Snmp) updateInterfaces() (*gnmi.SetRequest, error) {
	// interfaceIndex is a map of SNMP interface index -> name for this poll.
	interfaceIndex := make(map[string]string)

	s.lock.Lock()
	// Clear interfaceName map for each new poll. It should be
	// protected by the lock, because updateLldp needs it, too. :(
	s.interfaceName = make(map[string]bool)
	defer s.lock.Unlock()

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	intfWalk := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleInterfacePDU(data, interfaceIndex)
		if err != nil {
			return err
		}
		updates = append(updates, u...)
		return nil
	}

	// ifTable
	if err := s.walk(snmpIfTable, intfWalk); err != nil {
		return nil, err
	}

	// ifXTable
	if err := s.walk(snmpIfXTable, intfWalk); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{pgnmi.Path("interfaces", "interface")}
	setReq.Replace = updates
	return setReq, nil
}

// Some implementations will return a hostname only, while others
// will return a fully qualified domain name. splitSysName returns
// the hostname and the domain if it exists.
func splitSysName(sysName string) (string, string) {
	ss := strings.Split(sysName, ".")
	hn := ss[0]
	var dn string
	if len(ss) > 1 {
		dn = strings.Join(ss[1:], ".")
	}
	return hn, dn
}

func (s *Snmp) updateSystemState() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	setReq := new(gnmi.SetRequest)
	sysName, err := s.getStringByOID(snmpSysName)
	if err != nil {
		return nil, err
	}
	hostname, domainName := splitSysName(sysName)

	hn := update(pgnmi.Path("system", "state", "hostname"), strval(hostname))
	upd := []*gnmi.Update{hn}
	if domainName != "" {
		upd = append(upd,
			update(pgnmi.Path("system", "state", "domain-name"),
				strval(domainName)))
	}
	setReq.Replace = upd

	return setReq, nil
}

// There are three kinds of LLDP data: local general (non-port-specific),
// local per-port (comes with a local interface index), and remote
// (comes with a local interface index and remote port ID).
// processLldpOid extracts the relevant indices (if present) and returns
// them along with the true base OID.
func processLldpOid(oid string) (locIndex, remoteID,
	baseOid string, err error) {
	baseOid = oid

	// Local per-port
	if strings.HasPrefix(oid, snmpLldpStatsTxPortTable) ||
		strings.HasPrefix(oid, snmpLldpStatsRxPortTable) ||
		strings.HasPrefix(oid, snmpLldpLocPortTable) ||
		strings.HasPrefix(oid, snmpLldpV2LocPortTable) {
		baseOid, locIndex, err = oidIndex(oid)
		return
	}
	// Local per-port V2
	if strings.HasPrefix(oid, snmpLldpV2StatsTxPortTable) ||
		strings.HasPrefix(oid, snmpLldpV2StatsRxPortTable) {
		baseOid, _, err = oidIndex(oid) // remove lldpV2StatsTxDestMACAddress
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidIndex(baseOid)
		return
	}

	// Remote
	if strings.HasPrefix(oid, snmpLldpRemTable) {
		baseOid, remoteID, err = oidIndex(oid)
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidIndex(baseOid)
		if err != nil {
			return
		}
		baseOid, _, err = oidIndex(baseOid) // remove lldpRemTimeMark
		return
	}
	// Remote V2
	if strings.HasPrefix(oid, snmpLldpV2RemTable) {
		baseOid, remoteID, err = oidIndex(oid)
		if err != nil {
			return
		}
		baseOid, _, err = oidIndex(baseOid) // remove lldpV2RemLocalDestMACAddress
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidIndex(baseOid)
		if err != nil {
			return
		}
		baseOid, _, err = oidIndex(baseOid) // remove lldpRemTimeMark
		return
	}
	return
}

// Return MAC address from string or hex byte string.
func macFromBytes(s []byte) string {
	// string case
	if len(s) == 17 {
		return string(s)
	}

	// else assume hex string
	var t bytes.Buffer
	for i := 0; i < len(s); i++ {
		if i != 0 {
			t.WriteString(":")
		}
		t.WriteString(hex.EncodeToString(s[i : i+1]))
	}
	return t.String()
}

func (s *Snmp) handleLldpPDU(pdu gosnmp.SnmpPDU,
	lldpLocPortIndex map[string]string,
	lldpRemoteID map[string]bool) ([]*gnmi.Update, error) {

	// Split OID into parts.
	locIndex, remoteID, baseOid, err := processLldpOid(pdu.Name)
	if err != nil {
		return nil, err
	}

	// If we haven't yet seen this local interface, add it to our list.
	intfName := ""
	updates := make([]*gnmi.Update, 0)
	var ok bool
	if locIndex != "" {
		intfName, ok = lldpLocPortIndex[locIndex]
		if !ok {
			// If we have the port ID AND this interface is in the interfaceName
			// map, add it to the port index map. Otherwise we can't do anything
			// and should return.
			if baseOid != snmpLldpLocPortID && baseOid != snmpLldpV2LocPortID {
				return nil, nil
			}
			intfName = string(pdu.Value.([]byte))
			if _, ok = s.interfaceName[intfName]; !ok {
				return nil, nil
			}
			lldpLocPortIndex[locIndex] = intfName
		}
	}

	// If we haven't yet seen this remote system, add its ID.
	if remoteID != "" {
		if _, ok = lldpRemoteID[remoteID]; !ok {
			updates = append(updates,
				update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "id"),
					strval(remoteID)))
			lldpRemoteID[remoteID] = true
		}
	}

	var u *gnmi.Update
	switch baseOid {
	case snmpLldpLocPortID, snmpLldpV2LocPortID:
		updates = append(updates,
			update(pgnmi.LldpIntfConfigPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfStatePath(intfName, "name"),
				strval(intfName)))
	case snmpLldpLocChassisID, snmpLldpV2LocChassisID:
		u = update(pgnmi.LldpStatePath("chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpLocChassisIDSubtype, snmpLldpV2LocChassisIDSubtype:
		u = update(pgnmi.LldpStatePath("chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpLocSysName, snmpLldpV2LocSysName:
		u = update(pgnmi.LldpStatePath("system-name"),
			strval(pdu.Value))
	case snmpLldpLocSysDesc, snmpLldpV2LocSysDesc:
		u = update(pgnmi.LldpStatePath("system-description"),
			strval(pdu.Value))
	case snmpLldpStatsTxPortFramesTotal, snmpLldpV2StatsTxPortFramesTotal:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-out"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesDiscard, snmpLldpV2StatsRxPortFramesDiscard:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesErrors, snmpLldpV2StatsRxPortFramesErrors:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-error-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesTotal, snmpLldpV2StatsRxPortFramesTotal:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsDiscard, snmpLldpV2StatsRxPortTLVsDiscard:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "tlv-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsUnrecog, snmpLldpV2StatsRxPortTLVsUnrecog:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "tlv-unknown"),
			uintval(pdu.Value))
	case snmpLldpRemPortID, snmpLldpV2RemPortID:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id"),
			strval(pdu.Value))
	case snmpLldpRemPortIDSubtype, snmpLldpV2RemPortIDSubtype:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id-type"),
			strval(openconfig.LLDPPortIDType(pdu.Value.(int))))
	case snmpLldpRemChassisID, snmpLldpV2RemChassisID:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpRemChassisIDSubtype, snmpLldpV2RemChassisIDSubtype:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpRemSysName, snmpLldpV2RemSysName:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-name"),
			strval(pdu.Value))
	case snmpLldpRemSysDesc, snmpLldpV2RemSysDesc:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-description"),
			strval(pdu.Value))
	}
	if u != nil {
		updates = append(updates, u)
	}
	return updates, nil
}

func (s *Snmp) updateLldp() (*gnmi.SetRequest, error) {
	// lldpLocPortIndex is a map of lldpLocPortNum -> lldpLocPortId.
	lldpLocPortIndex := make(map[string]string)

	// lldpRemoteID is a map of remote system ID -> true. It's used to
	// remember which remote IDs we've already seen in a given round of polling.
	lldpRemoteID := make(map[string]bool)

	s.lock.Lock()
	defer s.lock.Unlock()

	locSysData := snmpLldpLocalSystemData
	remTable := snmpLldpRemTable
	statsTable := snmpLldpStatistics
	if s.lldpV2 {
		locSysData = snmpLldpV2LocalSystemData
		remTable = snmpLldpV2RemTable
		statsTable = snmpLldpV2Statistics
	}

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	updater := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleLldpPDU(data, lldpLocPortIndex, lldpRemoteID)
		if err != nil {
			return err
		}
		if u != nil {
			updates = append(updates, u...)
		}
		return nil
	}

	if err := s.walk(locSysData, updater); err != nil {
		return nil, err
	}
	// XXX_jcr: Ultimately we'll want to add a proper mechanism for discovering which
	// MIBs the target device supports. For now just try a different version next time.
	if len(updates) == 0 {
		s.lldpV2 = !s.lldpV2
		return setReq, nil
	}

	if err := s.walk(remTable, updater); err != nil {
		return nil, err
	}

	if err := s.walk(statsTable, updater); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{pgnmi.Path("lldp")}
	setReq.Replace = updates
	return setReq, nil
}

func (s *Snmp) handleEntityMibPDU(pdu gosnmp.SnmpPDU,
	entityIndexMap map[string]bool) ([]*gnmi.Update, error) {

	baseOid, index, err := oidIndex(pdu.Name)
	if err != nil {
		return nil, err
	}

	updates := make([]*gnmi.Update, 0)
	if _, ok := entityIndexMap[index]; !ok {
		entityIndexMap[index] = true
		updates = append(updates,
			update(pgnmi.PlatformComponentConfigPath(index, "name"),
				strval(index)),
			update(pgnmi.PlatformComponentPath(index, "name"),
				strval(index)),
			update(pgnmi.PlatformComponentStatePath(index, "name"),
				strval(index)),
			update(pgnmi.PlatformComponentStatePath(index, "id"),
				strval(index)))
	}

	switch baseOid {
	case snmpEntPhysicalClass:
		v := pdu.Value.(int)
		// OpenConfig's OPENCONFIG_HARDWARE_COMPONENT type identities don't
		// map perfectly to SNMP's PhysicalClass values. If we see a
		// PhysicalClass value of other(1), unknown(2), container(5), or
		// module(9), just leave the type blank.
		snmpOpenConfigComponentTypeMap := map[int]string{
			1:  "",
			2:  "",
			3:  "CHASSIS",
			4:  "BACKPLANE",
			5:  "",
			6:  "POWER_SUPPLY",
			7:  "FAN",
			8:  "SENSOR",
			9:  "",
			10: "PORT",
			11: "",
			12: "CPU",
		}
		class, ok := snmpOpenConfigComponentTypeMap[v]
		if !ok {
			return nil, fmt.Errorf("Unexpected PhysicalClass value %v", v)
		}
		if class != "" {
			updates = append(updates,
				update(pgnmi.PlatformComponentStatePath(index, "type"), strval(class)))
		}
	case snmpEntPhysicalDescr:
		updates = append(updates,
			update(pgnmi.PlatformComponentStatePath(index, "description"),
				strval(pdu.Value)))
	case snmpEntPhysicalMfgName:
		updates = append(updates,
			update(pgnmi.PlatformComponentStatePath(index, "mfg-name"),
				strval(pdu.Value)))
	case snmpEntPhysicalSerialNum:
		updates = append(updates,
			update(pgnmi.PlatformComponentStatePath(index, "serial-no"),
				strval(pdu.Value)))
	case snmpEntPhysicalSoftwareRev:
		updates = append(updates,
			update(pgnmi.PlatformComponentStatePath(index, "software-version"),
				strval(pdu.Value)))
	case snmpEntPhysicalModelName:
		updates = append(updates,
			update(pgnmi.PlatformComponentStatePath(index, "hardware-version"),
				strval(pdu.Value)))
	}
	return updates, nil
}

func (s *Snmp) updatePlatform() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	entityIndexMap := make(map[string]bool)
	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)

	updater := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleEntityMibPDU(data, entityIndexMap)
		if err != nil {
			return err
		}
		updates = append(updates, u...)
		return nil
	}

	if err := s.walk(snmpEntPhysicalTable, updater); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{pgnmi.Path("components")}
	setReq.Replace = updates
	return setReq, nil
}

// InitGNMIOpenConfig initializes the Snmp provider with a gNMI client.
func (s *Snmp) InitGNMIOpenConfig(client gnmi.GNMIClient) {
	s.client = client
	s.initialized = true
}

func (s *Snmp) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-s.errc:
			// XXX_jcr: We should probably return for some errors.
			// Others we can't return from. For example, an LLDP poll
			// may fail if it takes place after an interface change that
			// hasn't yet showed up in an interface poll, since LLDP
			// interfaces also have to be present in interfaces/interface.
			glog.Errorf("Failure in gNMI stream: %v", err)
		}
	}
}

// Run sets the Snmp provider running and returns only on error.
func (s *Snmp) Run(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("SNMP provider is uninitialized")
	}

	// Do periodic state updates.
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updateSystemState, s.errc)
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updatePlatform, s.errc)
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updateInterfaces, s.errc)
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updateLldp, s.errc)

	// Watch for errors.
	s.handleErrors(ctx)

	s.stop()
	return nil
}

// NewSNMPProvider returns a new SNMP provider for the device at 'address'
// using a community value for authentication and pollInterval for rate
// limiting requests.
func NewSNMPProvider(address string, community string,
	pollInt time.Duration) provider.GNMIOpenConfigProvider {
	gosnmp.Default.Target = address
	gosnmp.Default.Community = community
	gosnmp.Default.Timeout = 2 * pollInt
	s := &Snmp{
		errc:          make(chan error),
		interfaceName: make(map[string]bool),
		address:       address,
		community:     community,
		pollInterval:  pollInt,
		getter:        gosnmp.Default.Get,
		walker:        gosnmp.Default.BulkWalk,
	}
	if err := s.snmpNetworkInit(); err != nil {
		glog.Errorf("Error connecting to device: %v", err)
	}
	return s
}
