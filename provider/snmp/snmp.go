// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmp

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/sirupsen/logrus"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

const (
	snmpEntPhysicalClass               = ".1.3.6.1.2.1.47.1.1.1.1.5"
	snmpEntPhysicalDescr               = ".1.3.6.1.2.1.47.1.1.1.1.2"
	snmpEntPhysicalMfgName             = ".1.3.6.1.2.1.47.1.1.1.1.12"
	snmpEntPhysicalModelName           = ".1.3.6.1.2.1.47.1.1.1.1.13"
	snmpEntPhysicalSerialNum           = ".1.3.6.1.2.1.47.1.1.1.1.11"
	snmpEntPhysicalSoftwareRev         = ".1.3.6.1.2.1.47.1.1.1.1.10"
	snmpEntPhysicalTable               = ".1.3.6.1.2.1.47.1.1.1.1"
	snmpHrSystemUptime                 = ".1.3.6.1.2.1.25.1.1.0"
	snmpIfTable                        = ".1.3.6.1.2.1.2.2"
	snmpIfDescr                        = ".1.3.6.1.2.1.2.2.1.2"
	snmpIfType                         = ".1.3.6.1.2.1.2.2.1.3"
	snmpIfMtu                          = ".1.3.6.1.2.1.2.2.1.4"
	snmpIfName                         = ".1.3.6.1.2.1.31.1.1.1.1"
	snmpIfAdminStatus                  = ".1.3.6.1.2.1.2.2.1.7"
	snmpIfHCInOctets                   = ".1.3.6.1.2.1.31.1.1.1.6"
	snmpIfHCInBroadcastPkts            = ".1.3.6.1.2.1.31.1.1.1.9"
	snmpIfHCInMulticastPkts            = ".1.3.6.1.2.1.31.1.1.1.8"
	snmpIfHCInUcastPkts                = ".1.3.6.1.2.1.31.1.1.1.7"
	snmpIfHCOutOctets                  = ".1.3.6.1.2.1.31.1.1.1.10"
	snmpIfHCOutUcastPkts               = ".1.3.6.1.2.1.31.1.1.1.11"
	snmpIfHCOutBroadcastPkts           = ".1.3.6.1.2.1.31.1.1.1.13"
	snmpIfHCOutMulticastPkts           = ".1.3.6.1.2.1.31.1.1.1.12"
	snmpIfInDiscards                   = ".1.3.6.1.2.1.2.2.1.13"
	snmpIfInErrors                     = ".1.3.6.1.2.1.2.2.1.14"
	snmpIfInOctets                     = ".1.3.6.1.2.1.2.2.1.10"
	snmpIfInUcastPkts                  = ".1.3.6.1.2.1.2.2.1.11"
	snmpIfInUnknownProtos              = ".1.3.6.1.2.1.2.2.1.15"
	snmpIfOperStatus                   = ".1.3.6.1.2.1.2.2.1.8"
	snmpIfOutDiscards                  = ".1.3.6.1.2.1.2.2.1.19"
	snmpIfOutErrors                    = ".1.3.6.1.2.1.2.2.1.20"
	snmpIfOutOctets                    = ".1.3.6.1.2.1.2.2.1.16"
	snmpIfOutUcastPkts                 = ".1.3.6.1.2.1.2.2.1.17"
	snmpIfXTable                       = ".1.3.6.1.2.1.31.1.1"
	snmpLldpLocalSystemData            = ".1.0.8802.1.1.2.1.3"
	snmpLldpLocChassisID               = ".1.0.8802.1.1.2.1.3.2.0"
	snmpLldpLocChassisIDSubtype        = ".1.0.8802.1.1.2.1.3.1.0"
	snmpLldpLocPortID                  = ".1.0.8802.1.1.2.1.3.7.1.3"
	snmpLldpLocPortDesc                = ".1.0.8802.1.1.2.1.3.7.1.4"
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
	snmpLldpV2LocPortDesc              = ".1.3.111.2.802.1.1.13.1.3.7.1.4"
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
	snmpSysUpTimeInstance              = ".1.3.6.1.2.1.1.3.0"
)

var snmpErrMap = map[gosnmp.SNMPError]string{
	0: "noError", 1: "tooBig", 2: "noSuchName", 3: "badValue",
	4: "readOnly", 5: "genErr", 6: "noAccess", 7: "wrongType",
	8: "wrongLength", 9: "wrongEncoding", 10: "wrongValue",
	11: "noCreation", 12: "inconsistentValue", 13: "resourceUnavailable",
	14: "commitFailed", 15: "undoFailed", 16: "authorizationError",
	17: "notWritable", 18: "inconsistentName",
}

// Split the final index off an OID and return it along with the remaining OID.
func oidSplitEnd(oid string) (string, string, error) {
	finalDotPos := strings.LastIndex(oid, ".")
	if finalDotPos < 0 {
		return "", "", fmt.Errorf("oid '%s' does not match expected format", oid)
	}
	return oid[:finalDotPos], oid[(finalDotPos + 1):], nil
}

// Less typing: gNMI type helpers.
func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return pgnmi.Update(path, val)
}

func appendUpdates(base []*gnmi.Update, updates ...*gnmi.Update) []*gnmi.Update {
	for _, update := range updates {
		if update.Val != nil {
			base = append(base, update)
		}
	}
	return base
}

// Remove all but ASCII characters 32-126 to keep the JSON
// unmarshaler happy.
func bytesToSanitizedString(b []byte) string {
	out := make([]byte, len(b))
	j := 0
	for i := 0; i < len(b); i++ {
		if b[i] < 32 || b[i] > 126 {
			continue
		}

		// Replace square brackets with parentheses.
		c := b[i]
		if c == '[' {
			c = '('
		} else if c == ']' {
			c = ')'
		}

		out[j] = c
		j++
	}
	return string(out[:j])
}

func strval(s interface{}) *gnmi.TypedValue {
	switch t := s.(type) {
	case string:
		if t == "" {
			return nil
		}
		return pgnmi.Strval(t)
	case []byte:
		str := bytesToSanitizedString(t)
		if str == "" {
			return nil
		}
		return pgnmi.Strval(str)
	default:
		logrus.Fatalf("Unexpected type in strval: %T", s)
	}
	return nil
}

func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToUint64(u); err == nil {
		return pgnmi.Uintval(v)
	}
	return nil
}

func intval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToInt64(u); err == nil {
		return pgnmi.Intval(v)
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

	// Whether the device implements ifXTable.
	ifX bool

	// lldpV2 indicates whether to use LLDP-V2-MIB.
	lldpV2 bool

	lldpLocChassisIDSubtype string

	gsnmp *gosnmp.GoSNMP // gosnmp object
	mock  bool           // if true, don't do any network init

	// Lock for synchronizing access to state produced by SNMP polls.
	providerLock sync.Mutex
	// gosnmp can't handle parallel gets, so we also need to lock
	// access to its connection object.
	connectionLock sync.Mutex

	pollInterval time.Duration
	lastAlive    time.Time
	initialized  bool
	deviceID     string

	// Alternative Walk() and Get() for mock testing.
	getter func([]string) (*gosnmp.SnmpPacket, error)
	walker func(string, gosnmp.WalkFunc) error

	// Alternative time.Now() for mock testing.
	now func() time.Time
}

func (s *Snmp) snmpNetworkInit() error {
	if s.initialized || s.mock {
		return nil
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()
	err := s.gsnmp.Connect()

	s.initialized = err == nil
	return err
}

func (s *Snmp) get(oid string) (*gosnmp.SnmpPacket, error) {
	log.Log(s).Debugf("get (OID = %s)", oid)
	if s.getter == nil {
		return nil, errors.New("SNMP getter not set")
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()
	pkt, err := s.getter([]string{oid})
	log.Log(s).Debugf("get complete (OID = %s): pkt = %v, err = %v", oid, pkt, err)
	if err != nil {
		return nil, err
	}

	// Handle packet errors.
	if pkt.Error != gosnmp.NoError {
		errstr, ok := snmpErrMap[pkt.Error]
		if !ok {
			errstr = "Unknown error"
		}
		return nil, fmt.Errorf("Error in packet (%v): %v", pkt, errstr)
	}

	s.lastAlive = s.now()

	return pkt, err
}

func oidExists(pdu gosnmp.SnmpPDU) bool {
	return pdu.Type != gosnmp.NoSuchObject && pdu.Type != gosnmp.NoSuchInstance
}

func (s *Snmp) getFirstPDU(oid string) (*gosnmp.SnmpPDU, error) {
	pkt, err := s.get(oid)
	if err != nil {
		return nil, err
	}
	if len(pkt.Variables) == 0 {
		return nil, fmt.Errorf("No variables in SNMP packet for OID %s", oid)
	}
	return &pkt.Variables[0], err
}

// getString does a Get on the specified OID, an octet string, and
// returns the result as a string.
func (s *Snmp) getString(oid string) (string, error) {
	pdu, err := s.getFirstPDU(oid)

	// Accept a noSuchObject or noSuchInstance, but otherwise, if it's not
	// an octet string, something went wrong.
	if err != nil || !oidExists(*pdu) {
		return "", err
	}
	if pdu.Type != gosnmp.OctetString {
		return "", fmt.Errorf("Variable type in PDU for OID %s is not octet string", oid)
	}

	return string(pdu.Value.([]byte)), nil
}

func (s *Snmp) walk(rootOid string, walkFn gosnmp.WalkFunc) error {
	log.Log(s).Debugf("walk (OID = %s)", rootOid)
	if s.walker == nil {
		return errors.New("SNMP walker not set")
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()
	err := s.walker(rootOid, walkFn)
	if err != nil {
		return err
	}
	log.Log(s).Debugf("walk complete (OID = %s)", rootOid)
	s.lastAlive = s.now()
	return err
}

var errStopWalk = errors.New("stop walk")

func (s *Snmp) getSerialNumber() (string, error) {
	serial := ""
	var done bool
	chassisIndex := ""
	var snmpEntPhysicalClassTypeChassis = 3

	// Get the serial number corresponding to the index whose class
	// type is chassis(3).
	log.Log(s).Tracef("getSerialNumber")
	entPhysicalWalk := func(data gosnmp.SnmpPDU) error {
		// If we're finished, throw a pseudo-error to indicate to the
		// walker that no more walking is required.
		if done {
			return errStopWalk
		}
		baseOid, index, err := oidSplitEnd(data.Name)
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
			if index == chassisIndex && string(data.Value.([]byte)) != "" {
				serial = string(data.Value.([]byte))
				done = true
			}
		}

		return nil
	}

	if err := s.walk(snmpEntPhysicalClass, entPhysicalWalk); err != nil {
		return "", err
	}
	if err := s.walk(snmpEntPhysicalSerialNum, entPhysicalWalk); err != nil {
		if err != errStopWalk {
			return "", err
		}
	}
	log.Log(s).Tracef("getSerialNumber complete (serial = %v)", serial)
	return serial, nil
}

func (s *Snmp) getChassisID() (string, error) {
	log.Log(s).Tracef("getChassisID")
	pdu, err := s.getFirstPDU(snmpLldpLocChassisIDSubtype)
	if err != nil || !oidExists(*pdu) {
		return "", err
	}

	subtype := openconfig.LLDPChassisIDType(pdu.Value.(int))
	pkt, err := s.getFirstPDU(snmpLldpLocChassisID)
	if err != nil {
		return "", err
	}
	log.Log(s).Tracef("getChassisID (chassisID = %v)", chassisID(pkt.Value.([]byte), subtype))
	return chassisID(pkt.Value.([]byte), subtype), nil
}

// DeviceID returns the device ID.
func (s *Snmp) DeviceID() (string, error) {
	log.Log(s).Trace("Snmp.DeviceID")
	if err := s.snmpNetworkInit(); err != nil {
		return "", fmt.Errorf("Error connecting to device: %v", err)
	}

	if s.deviceID != "" {
		return s.deviceID, nil
	}

	did, err := s.getSerialNumber()
	if err != nil {
		return did, err
	} else if did != "" {
		s.deviceID = did
		return did, nil
	}

	did, err = s.getChassisID()
	if err != nil {
		return did, err
	} else if did != "" {
		s.deviceID = did
		return did, nil
	}

	// The device didn't give us a serial number. Use the device
	// address instead. It's not great but better than nothing.
	log.Log(s).Infof("Failed to retrieve serial number for device '%s'; "+
		"using address for device ID", s.gsnmp.Target)
	s.deviceID = s.gsnmp.Target
	return s.gsnmp.Target, nil
}

// Alive checks if device is still alive if poll interval has passed.
func (s *Snmp) Alive() (bool, error) {
	log.Log(s).Debugf("Alive")
	if err := s.snmpNetworkInit(); err != nil {
		return false, fmt.Errorf("Error connecting to device: %v", err)
	}
	if time.Since(s.lastAlive) < s.pollInterval {
		return true, nil
	}
	_, err := s.get(snmpSysUpTimeInstance)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *Snmp) stop() {
	if !s.mock {
		s.gsnmp.Conn.Close()
	}
}

// Given an incoming PDU, update the appropriate interface state.
func (s *Snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU,
	interfaceIndex map[string]string) ([]*gnmi.Update, error) {
	// Get/set interface name from index. If there's no mapping, just return and
	// wait for the mapping to show up.
	baseOid, index, err := oidSplitEnd(pdu.Name)
	log.Log(s).Tracef("handleInterfacePDU (OID = %s)", pdu.Name)
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
			strval("iana-if-type:"+openconfig.InterfaceType(pdu.Value.(int))))
	case snmpIfMtu:
		v, err := provider.ToUint64(pdu.Value)
		if err != nil {
			return nil, nil
		}
		if v > math.MaxUint16 {
			v = math.MaxUint16
		}
		u = update(pgnmi.IntfStatePath(intfName, "mtu"), uintval(v))
	case snmpIfName:
		s.ifX = true
	case snmpIfAdminStatus:
		u = update(pgnmi.IntfStatePath(intfName, "admin-status"),
			strval(openconfig.IntfAdminStatus(pdu.Value.(int))))
	case snmpIfOperStatus:
		u = update(pgnmi.IntfStatePath(intfName, "oper-status"),
			strval(openconfig.IntfOperStatus(pdu.Value.(int))))
	case snmpIfInOctets:
		if !s.ifX {
			u = update(pgnmi.IntfStateCountersPath(intfName, "in-octets"),
				uintval(pdu.Value))
		}
	case snmpIfHCInOctets:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-octets"),
			uintval(pdu.Value))
	case snmpIfInUcastPkts:
		if !s.ifX {
			u = update(pgnmi.IntfStateCountersPath(intfName, "in-unicast-pkts"),
				uintval(pdu.Value))
		}
	case snmpIfHCInUcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-unicast-pkts"),
			uintval(pdu.Value))
	case snmpIfHCInMulticastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfHCInBroadcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "in-broadcast-pkts"),
			uintval(pdu.Value))
	case snmpIfHCOutMulticastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfHCOutBroadcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-broadcast-pkts"),
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
		if !s.ifX {
			u = update(pgnmi.IntfStateCountersPath(intfName, "out-octets"),
				uintval(pdu.Value))
		}
	case snmpIfHCOutOctets:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-octets"),
			uintval(pdu.Value))
	case snmpIfOutUcastPkts:
		if !s.ifX {
			u = update(pgnmi.IntfStateCountersPath(intfName, "out-unicast-pkts"),
				uintval(pdu.Value))
		}
	case snmpIfHCOutUcastPkts:
		u = update(pgnmi.IntfStateCountersPath(intfName, "out-unicast-pkts"),
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

	var updates []*gnmi.Update
	if u != nil {
		updates = appendUpdates(updates, u)
	}
	// When we get a name, add name, config/name, state/name.
	if baseOid == snmpIfDescr {
		updates = appendUpdates(updates,
			update(pgnmi.IntfPath(intfName, "name"), strval(intfName)),
			update(pgnmi.IntfConfigPath(intfName, "name"), strval(intfName)))
	}
	return updates, nil
}

func (s *Snmp) updateInterfaces() ([]*gnmi.SetRequest, error) {
	// interfaceIndex is a map of SNMP interface index -> name for this poll.
	interfaceIndex := make(map[string]string)
	s.providerLock.Lock()
	defer s.providerLock.Unlock()
	log.Log(s).Tracef("updateInterfaces")
	// Clear interfaceName map for each new poll. It should be
	// protected by the lock, because updateLldp needs it, too. :(
	s.interfaceName = make(map[string]bool)

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	intfWalk := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleInterfacePDU(data, interfaceIndex)
		if err != nil {
			return err
		}
		updates = appendUpdates(updates, u...)
		return nil
	}

	// ifTable
	if err := s.walk(snmpIfTable, intfWalk); err != nil {
		return nil, err
	}

	log.Log(s).Tracef("updateInterfaces finished ifTable")

	// ifXTable
	if err := s.walk(snmpIfXTable, intfWalk); err != nil {
		return nil, err
	}

	log.Log(s).Tracef("updateInterfaces finished ifXTable")

	setReq.Delete = []*gnmi.Path{pgnmi.Path("interfaces", "interface")}
	setReq.Replace = updates
	log.Log(s).Tracef("updateInterfaces produced %d updates", len(updates))
	return []*gnmi.SetRequest{setReq}, nil
}

// Some implementations will return a hostname only, while others
// will return a fully qualified domain name. splitSysName returns
// the hostname and the domain if it exists.
func splitSysName(sysName string) (string, string) {
	ss := append(strings.SplitN(sysName, ".", 2), "")
	return ss[0], ss[1]
}

func (s *Snmp) updateSystemState() ([]*gnmi.SetRequest, error) {
	s.providerLock.Lock()
	defer s.providerLock.Unlock()
	log.Log(s).Tracef("updateSystemState")
	setReq := new(gnmi.SetRequest)
	var upd []*gnmi.Update

	// get hostname
	sysName, err := s.getString(snmpSysName)
	if err != nil || sysName == "" {
		// Try lldpLocSysName if sysName isn't there.
		sysName, err = s.getString(snmpLldpLocSysName)
		if err != nil {
			return nil, err
		}
	}
	if sysName != "" {
		hostname, domainName := splitSysName(sysName)

		hn := update(pgnmi.Path("system", "state", "hostname"), strval(hostname))
		upd = appendUpdates(upd, hn)
		if domainName != "" {
			upd = appendUpdates(upd,
				update(pgnmi.Path("system", "state", "domain-name"),
					strval(domainName)))
		}
	}

	// Get boot-time by subtracting the target device's uptime from
	// the Collector's current time. This isn't really correct--the
	// boot time we produce is a blend of UNIX time on the Collector
	// and UNIX time on the target device (which may not be in sync),
	// and the target device's time may be recorded well before the
	// Collector's. There doesn't seem to be a great way to get the
	// device's time using SNMP. Assuming the two devices are
	// roughly in sync, though, this shouldn't be disastrous.
	var bootTime int64
	for _, uptimeOid := range []string{snmpHrSystemUptime, snmpSysUpTimeInstance} {
		if bootTime != 0 {
			break
		}
		pdu, err := s.getFirstPDU(uptimeOid)
		if err != nil {
			return nil, err
		}
		if oidExists(*pdu) {
			t, err := provider.ToInt64(pdu.Value)
			if err != nil {
				return nil, err
			}
			if t != 0 {
				bootTime = s.now().Unix() - t/100
			}
		}
	}
	if bootTime != 0 {
		upd = appendUpdates(upd, update(pgnmi.Path("system", "state", "boot-time"),
			intval(bootTime)))
	}

	// Didn't get anything useful. Don't return a SetRequest.
	if len(upd) == 0 {
		return nil, nil
	}

	setReq.Replace = upd

	log.Log(s).Tracef("updateSystemState produced %d updates", len(upd))

	return []*gnmi.SetRequest{setReq}, nil
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
		baseOid, locIndex, err = oidSplitEnd(oid)
		return
	}
	// Local per-port V2
	if strings.HasPrefix(oid, snmpLldpV2StatsTxPortTable) ||
		strings.HasPrefix(oid, snmpLldpV2StatsRxPortTable) {
		baseOid, _, err = oidSplitEnd(oid) // remove lldpV2StatsTxDestMACAddress
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidSplitEnd(baseOid)
		return
	}

	// Remote
	if strings.HasPrefix(oid, snmpLldpRemTable) {
		baseOid, remoteID, err = oidSplitEnd(oid)
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidSplitEnd(baseOid)
		if err != nil {
			return
		}
		baseOid, _, err = oidSplitEnd(baseOid) // remove lldpRemTimeMark
		return
	}
	// Remote V2
	if strings.HasPrefix(oid, snmpLldpV2RemTable) {
		baseOid, remoteID, err = oidSplitEnd(oid)
		if err != nil {
			return
		}
		baseOid, _, err = oidSplitEnd(baseOid) // remove lldpV2RemLocalDestMACAddress
		if err != nil {
			return
		}
		baseOid, locIndex, err = oidSplitEnd(baseOid)
		if err != nil {
			return
		}
		baseOid, _, err = oidSplitEnd(baseOid) // remove lldpRemTimeMark
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
	return net.HardwareAddr(s).String()
}

var chassisIDSubtypeMacAddress = openconfig.LLDPChassisIDType(4)

func chassisID(b []byte, subtype string) string {
	if subtype == chassisIDSubtypeMacAddress {
		return macFromBytes(b)
	}
	return string(b)
}

func (s *Snmp) locChassisID(b []byte) string {
	return chassisID(b, s.lldpLocChassisIDSubtype)
}

var portIDSubtypeMacAddress = openconfig.LLDPPortIDType(3)

func portID(b []byte, subtype string) string {
	if subtype == portIDSubtypeMacAddress {
		return macFromBytes(b)
	}
	return string(b)
}

type remoteKey struct{ intfName, remoteID string }

// Data collected during a round of polling
type lldpSeen struct {
	// lldpLocPortNum -> lldpLocPortId
	locPortID map[string]string

	// Which intfName/remoteID pairs we've already seen in the round,
	// and their associated lldpChassisIdSubtypes.
	remoteID map[remoteKey]string

	// intfName/remoteID -> lldpRemPortIdSubtype
	remotePortID map[remoteKey]string

	// The OID from which we pulled interface names matching ifDescr.
	intfOid string
}

func (s *Snmp) handleLldpPDU(pdu gosnmp.SnmpPDU, seen *lldpSeen) ([]*gnmi.Update, error) {
	// Split OID into parts.
	locIndex, remoteID, baseOid, err := processLldpOid(pdu.Name)
	if err != nil {
		return nil, err
	}
	log.Log(s).Tracef("handleLldpPDU (OID = %s)", pdu.Name)

	// If we haven't yet seen this local interface, add it to our list.
	intfName := ""
	updates := make([]*gnmi.Update, 0)
	var ok bool
	if locIndex != "" {
		intfName, ok = seen.locPortID[locIndex]
		if !ok {
			// If this is an interface name AND this interface is in the
			// interfaceName map, add it to the port index map.
			// Otherwise we can't do anything and should return.
			if baseOid != snmpLldpLocPortID && baseOid != snmpLldpV2LocPortID &&
				baseOid != snmpLldpLocPortDesc && baseOid != snmpLldpV2LocPortDesc {
				return nil, nil
			}

			// XXX NOTE: The RFC says lldpLocPortDesc should have the
			// same value as a corresponding ifDescr object, but in
			// practice it seems more common to be have lldpLocPortId
			// equal to an ifDescr object, and lldpLocPortDesc is all
			// over the map--sometimes empty, sometimes set to
			// ifAlias. So just use whichever one matches ifDescr.
			if seen.intfOid != "" && baseOid != seen.intfOid {
				return nil, nil
			}
			intfName = string(pdu.Value.([]byte))
			if _, ok = s.interfaceName[intfName]; !ok {
				return nil, nil
			}
			seen.intfOid = baseOid
			seen.locPortID[locIndex] = intfName
		}
	}

	// If we haven't yet seen this remote system, add its ID.
	if remoteID != "" {
		if _, ok := seen.remoteID[remoteKey{intfName, remoteID}]; !ok {
			updates = appendUpdates(updates,
				update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "id"),
					strval(remoteID)))
			seen.remoteID[remoteKey{intfName, remoteID}] = ""
		}
	}

	var u *gnmi.Update
	switch baseOid {
	case seen.intfOid:
		// lldpLocPortID, lldpV2LocPortID, lldpLocPortDesc,
		// lldpV2LocPortDesc
		updates = appendUpdates(updates,
			update(pgnmi.LldpIntfConfigPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfStatePath(intfName, "name"),
				strval(intfName)))
	case snmpLldpLocChassisID, snmpLldpV2LocChassisID:
		v := s.locChassisID(pdu.Value.([]byte))
		u = update(pgnmi.LldpStatePath("chassis-id"), strval(v))
	case snmpLldpLocChassisIDSubtype, snmpLldpV2LocChassisIDSubtype:
		s.lldpLocChassisIDSubtype = openconfig.LLDPChassisIDType(pdu.Value.(int))
		u = update(pgnmi.LldpStatePath("chassis-id-type"),
			strval(s.lldpLocChassisIDSubtype))
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
		subtype := seen.remotePortID[remoteKey{intfName, remoteID}]
		v := portID(pdu.Value.([]byte), subtype)
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id"),
			strval(v))
	case snmpLldpRemPortIDSubtype, snmpLldpV2RemPortIDSubtype:
		v := openconfig.LLDPPortIDType(pdu.Value.(int))
		seen.remotePortID[remoteKey{intfName, remoteID}] = v
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id-type"),
			strval(v))
	case snmpLldpRemChassisID, snmpLldpV2RemChassisID:
		subtype := seen.remoteID[remoteKey{intfName, remoteID}]
		v := chassisID(pdu.Value.([]byte), subtype)
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id"),
			strval(v))
	case snmpLldpRemChassisIDSubtype, snmpLldpV2RemChassisIDSubtype:
		v := openconfig.LLDPChassisIDType(pdu.Value.(int))
		seen.remoteID[remoteKey{intfName, remoteID}] = v
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id-type"),
			strval(v))
	case snmpLldpRemSysName, snmpLldpV2RemSysName:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-name"),
			strval(pdu.Value))
	case snmpLldpRemSysDesc, snmpLldpV2RemSysDesc:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-description"),
			strval(pdu.Value))
	}
	if u != nil {
		updates = appendUpdates(updates, u)
	}
	return updates, nil
}

func (s *Snmp) updateLldp() ([]*gnmi.SetRequest, error) {
	s.providerLock.Lock()
	defer s.providerLock.Unlock()
	log.Log(s).Tracef("updateLldp")

	// Unset local chassis ID subtype.
	s.lldpLocChassisIDSubtype = ""

	locSysData := snmpLldpLocalSystemData
	remTable := snmpLldpRemTable
	statsRoot := snmpLldpStatistics
	if s.lldpV2 {
		locSysData = snmpLldpV2LocalSystemData
		remTable = snmpLldpV2RemTable
		statsRoot = snmpLldpV2Statistics
	}

	seen := &lldpSeen{
		locPortID:    make(map[string]string),
		remoteID:     make(map[remoteKey]string),
		remotePortID: make(map[remoteKey]string),
	}

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	updater := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleLldpPDU(data, seen)
		if err != nil {
			return err
		}
		if u != nil {
			updates = appendUpdates(updates, u...)
		}
		return nil
	}

	if err := s.walk(locSysData, updater); err != nil {
		return nil, err
	}
	log.Log(s).Tracef("updateLldp finished lldpLocalSystemData")
	// XXX NOTE: Ultimately we'll want to add a proper mechanism for discovering which
	// MIBs the target device supports. Here we could just request lldpV2LocSysName
	// to see if the device supports V2. But for now just try a different version
	// next time.
	if len(updates) == 0 {
		s.lldpV2 = !s.lldpV2
		return []*gnmi.SetRequest{setReq}, nil
	}

	if err := s.walk(remTable, updater); err != nil {
		return nil, err
	}
	log.Log(s).Tracef("updateLldp finished lldpRemoteSystemsData")

	if err := s.walk(statsRoot, updater); err != nil {
		return nil, err
	}
	log.Log(s).Tracef("updateLldp finished lldpStatistics")

	setReq.Delete = []*gnmi.Path{pgnmi.Path("lldp")}
	setReq.Replace = updates
	log.Log(s).Tracef("updateLldp produced %d updates", len(updates))
	return []*gnmi.SetRequest{setReq}, nil
}

// It's necessary to run updateInterfaces before updateLldp, since the
// lldp model depends on the interfaces being there already.
func (s *Snmp) updateInterfacesAndLldp() ([]*gnmi.SetRequest, error) {
	intfSR, err := s.updateInterfaces()
	if err != nil {
		return nil, err
	}
	lldpSR, err := s.updateLldp()
	if err != nil {
		return nil, err
	}
	return append(intfSR, lldpSR...), nil
}

func (s *Snmp) handleEntityMibPDU(pdu gosnmp.SnmpPDU,
	entityIndexMap map[string]bool) ([]*gnmi.Update, error) {

	baseOid, index, err := oidSplitEnd(pdu.Name)
	if err != nil {
		return nil, err
	}
	log.Log(s).Tracef("handleEntityMibPDU (OID = %s)", pdu.Name)

	updates := make([]*gnmi.Update, 0)
	if _, ok := entityIndexMap[index]; !ok {
		entityIndexMap[index] = true
		updates = appendUpdates(updates,
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
			updates = appendUpdates(updates,
				update(pgnmi.PlatformComponentStatePath(index, "type"),
					strval("openconfig-platform-types:"+class)))
		}
	case snmpEntPhysicalDescr:
		updates = appendUpdates(updates,
			update(pgnmi.PlatformComponentStatePath(index, "description"),
				strval(pdu.Value)))
	case snmpEntPhysicalMfgName:
		updates = appendUpdates(updates,
			update(pgnmi.PlatformComponentStatePath(index, "mfg-name"),
				strval(pdu.Value)))
	case snmpEntPhysicalSerialNum:
		updates = appendUpdates(updates,
			update(pgnmi.PlatformComponentStatePath(index, "serial-no"),
				strval(pdu.Value)))
	case snmpEntPhysicalSoftwareRev:
		updates = appendUpdates(updates,
			update(pgnmi.PlatformComponentStatePath(index, "software-version"),
				strval(pdu.Value)))
	case snmpEntPhysicalModelName:
		updates = appendUpdates(updates,
			update(pgnmi.PlatformComponentStatePath(index, "hardware-version"),
				strval(pdu.Value)))
	}
	return updates, nil
}

func (s *Snmp) updatePlatform() ([]*gnmi.SetRequest, error) {
	s.providerLock.Lock()
	defer s.providerLock.Unlock()
	log.Log(s).Tracef("updatePlatform")

	entityIndexMap := make(map[string]bool)
	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)

	updater := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleEntityMibPDU(data, entityIndexMap)
		if err != nil {
			return err
		}
		updates = appendUpdates(updates, u...)
		return nil
	}

	if err := s.walk(snmpEntPhysicalTable, updater); err != nil {
		return nil, err
	}
	log.Log(s).Tracef("updatePlatform finished entPhysicalTable")

	setReq.Delete = []*gnmi.Path{pgnmi.Path("components")}
	setReq.Replace = updates
	log.Log(s).Tracef("updatePlatform produced %d updates", len(updates))
	return []*gnmi.SetRequest{setReq}, nil
}

// InitGNMI initializes the Snmp provider with a gNMI client.
func (s *Snmp) InitGNMI(client gnmi.GNMIClient) {
	s.client = client
}

// OpenConfig indicates that this provider wants OpenConfig
// type-checking.
func (s *Snmp) OpenConfig() bool {
	return true
}

func (s *Snmp) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Log(s).Infof("SNMP provider is finished")
			return
		case err := <-s.errc:
			// XXX NOTE: We should probably return for some errors.
			// Others we can't return from. For example, an LLDP poll
			// may fail if it takes place after an interface change that
			// hasn't yet showed up in an interface poll, since LLDP
			// interfaces also have to be present in interfaces/interface.
			log.Log(s).Errorf("Failure in SNMP -> gNMI stream: %v", err)
		}
	}
}

// Run sets the Snmp provider running and returns only on error.
func (s *Snmp) Run(ctx context.Context) error {
	if s.client == nil {
		return errors.New("Run called before InitGNMI")
	}
	log.Log(s).Debugf("Run")

	if err := s.snmpNetworkInit(); err != nil {
		return fmt.Errorf("Error connecting to device: %v", err)
	}
	log.Log(s).Debugf("gosnmp.Connect complete")

	// Do periodic state updates.
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updateSystemState, s.errc)
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updatePlatform, s.errc)
	go pgnmi.PollForever(ctx, s.client, s.pollInterval,
		s.updateInterfacesAndLldp, s.errc)

	// Watch for errors.
	s.handleErrors(ctx)

	s.stop()
	return nil
}

// V3Params contains options related to SNMPv3.
type V3Params struct {
	SecurityModel gosnmp.SnmpV3SecurityModel
	Level         gosnmp.SnmpV3MsgFlags
	UsmParams     *gosnmp.UsmSecurityParameters
}

// NewSNMPProvider returns a new SNMP provider for the device at 'address'
// using a community value for authentication and pollInterval for rate
// limiting requests.
func NewSNMPProvider(address string, port uint16, community string,
	pollInt time.Duration, version gosnmp.SnmpVersion,
	v3Params *V3Params, verbose bool, mock bool) provider.GNMIProvider {
	gsnmp := &gosnmp.GoSNMP{
		Port:               port,
		Version:            version,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
		Target:             address,
		Community:          community,
		Timeout:            time.Duration(2) * time.Second,
		Logger:             nil,
		MaxRepetitions:     12,
	}
	if v3Params != nil {
		gsnmp.MsgFlags = v3Params.Level
		gsnmp.SecurityModel = v3Params.SecurityModel
		gsnmp.SecurityParameters = v3Params.UsmParams
	}

	s := &Snmp{
		gsnmp:         gsnmp,
		errc:          make(chan error),
		interfaceName: make(map[string]bool),
		pollInterval:  pollInt,
		mock:          mock,
		getter:        gsnmp.Get,
		walker:        gsnmp.BulkWalk,
		now:           time.Now,
	}
	return s
}
