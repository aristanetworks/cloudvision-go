// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/provider"
	"arista/provider/openconfig"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/aristanetworks/glog"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

type snmp struct {
	errc             chan error
	ctx              context.Context
	client           gnmi.GNMIClient
	interfaceIndex   map[string]string
	lldpLocPortIndex map[string]string
	address          string
	community        string
	lock             sync.Mutex // gosnmp can't handle parallel gets
	initialized      bool
}

// Has read/write interface been established?
var connected bool

// Time we last heard back from target
var lastAlive time.Time

var pollInterval time.Duration

func snmpNetworkInit() error {
	if connected {
		return nil
	}
	err := gosnmp.Default.Connect()
	if err == nil {
		connected = true
	}
	return err
}

// SNMPGetByOID returns the value at oid.
func SNMPGetByOID(oid string) (string, error) {
	oids := []string{oid}
	err := snmpNetworkInit()
	if err != nil {
		return "", err
	}

	// Ask for object
	result, err := gosnmp.Default.Get(oids)
	if err != nil {
		return "", err
	}

	lastAlive = time.Now()

	// Retrieve it from results
	for _, v := range result.Variables {
		switch v.Type {
		case gosnmp.OctetString:
			return string(v.Value.([]byte)), nil
		default:
			return gosnmp.ToBigInt(v.Value).String(), nil
		}
	}

	return "", errors.New("How did we get here?")
}

// SNMPDeviceID returns the device ID
func SNMPDeviceID() (string, error) {
	return SNMPGetByOID(snmpEntPhysicalSerialNum)
}

// SNMPCheckAlive checks if device is still alive if poll interval has passed.
func SNMPCheckAlive() (bool, error) {
	if time.Since(lastAlive) < pollInterval {
		return true, nil
	}
	_, err := SNMPGetByOID(snmpSysUpTime)
	return true, err
}

func (s *snmp) stop() {
	gosnmp.Default.Conn.Close()
}

// return base OID, index
func oidIndex(oid string) (string, string, error) {
	finalDotPos := strings.LastIndex(oid, ".")
	if finalDotPos < 0 {
		return "", "", fmt.Errorf("oid '%s' does not match expected format", oid)
	}
	return oid[:finalDotPos], oid[(finalDotPos + 1):], nil
}

const (
	snmpEntPhysicalSerialNum         = ".1.3.6.1.2.1.47.1.1.1.1.11.1"
	snmpSysName                      = ".1.3.6.1.2.1.1.5.0"
	snmpIfTable                      = ".1.3.6.1.2.1.2.2"
	snmpIfXTable                     = ".1.3.6.1.2.1.31.1.1"
	snmpIfDescr                      = ".1.3.6.1.2.1.2.2.1.2"
	snmpIfType                       = ".1.3.6.1.2.1.2.2.1.3"
	snmpIfMtu                        = ".1.3.6.1.2.1.2.2.1.4"
	snmpIfAdminStatus                = ".1.3.6.1.2.1.2.2.1.7"
	snmpIfOperStatus                 = ".1.3.6.1.2.1.2.2.1.8"
	snmpIfInOctets                   = ".1.3.6.1.2.1.2.2.1.10"
	snmpIfInUcastPkts                = ".1.3.6.1.2.1.2.2.1.11"
	snmpIfInMulticastPkts            = ".1.3.6.1.2.1.31.1.1.1.2"
	snmpIfInBroadcastPkts            = ".1.3.6.1.2.1.31.1.1.1.3"
	snmpIfInDiscards                 = ".1.3.6.1.2.1.2.2.1.13"
	snmpIfInErrors                   = ".1.3.6.1.2.1.2.2.1.14"
	snmpIfInUnknownProtos            = ".1.3.6.1.2.1.2.2.1.15"
	snmpIfOutOctets                  = ".1.3.6.1.2.1.2.2.1.16"
	snmpIfOutUcastPkts               = ".1.3.6.1.2.1.2.2.1.17"
	snmpIfOutMulticastPkts           = ".1.3.6.1.2.1.31.1.1.1.4"
	snmpIfOutBroadcastPkts           = ".1.3.6.1.2.1.31.1.1.1.5"
	snmpIfOutDiscards                = ".1.3.6.1.2.1.2.2.1.19"
	snmpIfOutErrors                  = ".1.3.6.1.2.1.2.2.1.20"
	snmpLldpLocalSystemData          = ".1.0.8802.1.1.2.1.3"
	snmpLldpLocPortTable             = ".1.0.8802.1.1.2.1.3.7"
	snmpLldpRemTable                 = ".1.0.8802.1.1.2.1.4.1"
	snmpLldpStatistics               = ".1.0.8802.1.1.2.1.2"
	snmpLldpStatsTxPortTable         = ".1.0.8802.1.1.2.1.2.6"
	snmpLldpStatsRxPortTable         = ".1.0.8802.1.1.2.1.2.7"
	snmpLldpLocChassisID             = ".1.0.8802.1.1.2.1.3.2.0"
	snmpLldpLocChassisIDSubtype      = ".1.0.8802.1.1.2.1.3.1.0"
	snmpLldpLocSysName               = ".1.0.8802.1.1.2.1.3.3.0"
	snmpLldpLocSysDesc               = ".1.0.8802.1.1.2.1.3.4.0"
	snmpLldpLocPortID                = ".1.0.8802.1.1.2.1.3.7.1.3"
	snmpLldpRemPortID                = ".1.0.8802.1.1.2.1.4.1.1.7"
	snmpLldpRemPortIDSubtype         = ".1.0.8802.1.1.2.1.4.1.1.6"
	snmpLldpRemChassisID             = ".1.0.8802.1.1.2.1.4.1.1.5"
	snmpLldpRemChassisIDSubtype      = ".1.0.8802.1.1.2.1.4.1.1.4"
	snmpLldpRemSysName               = ".1.0.8802.1.1.2.1.4.1.1.9"
	snmpLldpRemSysDesc               = ".1.0.8802.1.1.2.1.4.1.1.10"
	snmpLldpStatsTxPortFramesTotal   = ".1.0.8802.1.1.2.1.2.6.1.2"
	snmpLldpStatsRxPortFramesDiscard = ".1.0.8802.1.1.2.1.2.7.1.2"
	snmpLldpStatsRxPortFramesErrors  = ".1.0.8802.1.1.2.1.2.7.1.3"
	snmpLldpStatsRxPortFramesTotal   = ".1.0.8802.1.1.2.1.2.7.1.4"
	snmpLldpStatsRxPortTLVsDiscard   = ".1.0.8802.1.1.2.1.2.7.1.5"
	snmpLldpStatsRxPortTLVsUnrecog   = ".1.0.8802.1.1.2.1.2.7.1.6"
	snmpSysUpTime                    = ".1.3.6.1.2.1.1.3.0"
)

// Less typing:
func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return GNMIUpdate(path, val)
}

// toInt64 converts an interface to a int64.
func toInt64(valIntf interface{}) (int64, error) {
	var val int64
	switch t := valIntf.(type) {
	case int:
		val = int64(t)
	case int8:
		val = int64(t)
	case int16:
		val = int64(t)
	case int32:
		val = int64(t)
	case int64:
		val = t
	case uint:
		val = int64(t)
	case uint8:
		val = int64(t)
	case uint16:
		val = int64(t)
	case uint32:
		val = int64(t)
	case uint64:
		if t > math.MaxInt64 {
			return 0, fmt.Errorf("could not convert to int64, %d larger than max of %d",
				t, uint64(math.MaxInt64))
		}
		val = int64(t)
	default:
		return 0, fmt.Errorf("update contained value of unexpected type %T", valIntf)
	}
	return val, nil
}

// toUint64 converts an interface to a uint64.
func toUint64(valIntf interface{}) (uint64, error) {
	var val uint64
	switch t := valIntf.(type) {
	case int, int8, int16, int32, int64:
		v, e := toInt64(t)
		if e != nil {
			return 0, e
		}
		if v < 0 {
			return 0, fmt.Errorf("value %d cannot be converted to uint as it is negative", v)
		}
		val = uint64(v)
	case uint:
		val = uint64(t)
	case uint8:
		val = uint64(t)
	case uint16:
		val = uint64(t)
	case uint32:
		val = uint64(t)
	case uint64:
		val = t
	default:
		return 0, fmt.Errorf("update contained value of unexpected type %T", valIntf)
	}
	return val, nil
}

func strval(s interface{}) *gnmi.TypedValue {
	t, ok := s.(string)
	if ok {
		return GNMIStrval(t)
	}
	u, ok := s.([]byte)
	if ok {
		return GNMIStrval(string(u))
	}
	glog.Fatalf("Unexpected type in strval: %T", s)
	return nil
}
func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := toUint64(u); err == nil {
		return GNMIUintval(v)
	}
	return nil
}

// Given an incoming PDU, update the appropriate interface state.
func (s *snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU) ([]*gnmi.Update, error) {
	// Get/set interface name from index. If there's no mapping, just return and
	// wait for the mapping to show up.
	baseOid, index, err := oidIndex(pdu.Name)
	if err != nil {
		return nil, err
	}
	intfName, ok := s.interfaceIndex[index]
	if !ok && baseOid != snmpIfDescr {
		return nil, nil
	} else if !ok && baseOid == snmpIfDescr {
		intfName = string(pdu.Value.([]byte))
		s.interfaceIndex[index] = intfName
	}

	var u *gnmi.Update
	switch baseOid {
	case snmpIfDescr:
		u = update(GNMIIntfStatePath(intfName, "name"),
			strval(pdu.Value))
	case snmpIfType:
		u = update(GNMIIntfStatePath(intfName, "type"),
			strval(openconfig.InterfaceType(pdu.Value.(int))))
	case snmpIfMtu:
		u = update(GNMIIntfStatePath(intfName, "mtu"),
			uintval(pdu.Value))
	case snmpIfAdminStatus:
		u = update(GNMIIntfStatePath(intfName, "admin-status"),
			strval(openconfig.IntfAdminStatus(pdu.Value.(int))))
	case snmpIfOperStatus:
		u = update(GNMIIntfStatePath(intfName, "oper-status"),
			strval(openconfig.IntfOperStatus(pdu.Value.(int))))
	case snmpIfInOctets:
		u = update(GNMIIntfStateCountersPath(intfName, "in-octets"),
			uintval(pdu.Value))
	case snmpIfInUcastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "in-unicast-pkts"),
			uintval(pdu.Value))
	case snmpIfInMulticastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "in-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfInBroadcastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "in-broadcast-pkts"),
			uintval(pdu.Value))
	case snmpIfInDiscards:
		u = update(GNMIIntfStateCountersPath(intfName, "in-discards"),
			uintval(pdu.Value))
	case snmpIfInErrors:
		u = update(GNMIIntfStateCountersPath(intfName, "in-errors"),
			uintval(pdu.Value))
	case snmpIfInUnknownProtos:
		u = update(GNMIIntfStateCountersPath(intfName, "in-unknown-protos"),
			uintval(pdu.Value))
	case snmpIfOutOctets:
		u = update(GNMIIntfStateCountersPath(intfName, "out-octets"),
			uintval(pdu.Value))
	case snmpIfOutUcastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "out-unicast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutMulticastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "out-multicast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutBroadcastPkts:
		u = update(GNMIIntfStateCountersPath(intfName, "out-broadcast-pkts"),
			uintval(pdu.Value))
	case snmpIfOutDiscards:
		u = update(GNMIIntfStateCountersPath(intfName, "out-discards"),
			uintval(pdu.Value))
	case snmpIfOutErrors:
		u = update(GNMIIntfStateCountersPath(intfName, "out-errors"),
			uintval(pdu.Value))
	default:
		// default: ignore update
		return nil, nil
	}

	updates := []*gnmi.Update{u}
	// When we get a name, add name, config/name, state/name.
	if baseOid == snmpIfDescr {
		updates = append(updates,
			update(GNMIIntfPath(intfName, "name"), strval(intfName)),
			update(GNMIIntfConfigPath(intfName, "name"), strval(intfName)))
	}
	return updates, nil
}

func (s *snmp) updateInterfaces() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	intfWalk := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleInterfacePDU(data)
		if err != nil {
			return err
		}
		updates = append(updates, u...)
		return nil
	}

	// ifTable
	if err := gosnmp.Default.Walk(snmpIfTable, intfWalk); err != nil {
		return nil, err
	}

	// ifXTable
	if err := gosnmp.Default.Walk(snmpIfXTable, intfWalk); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{GNMIPath("interfaces", "interface")}
	setReq.Replace = updates
	return setReq, nil
}

func (s *snmp) updateSystemState() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	setReq := new(gnmi.SetRequest)
	sysName, err := SNMPGetByOID(snmpSysName)
	if err != nil {
		return nil, err
	}
	hostname := strings.Split(sysName, ".")[0]
	domainName := strings.Join(strings.Split(sysName, ".")[1:], ".")

	hn := update(GNMIPath("system", "state", "hostname"), strval(hostname))
	dn := update(GNMIPath("system", "state", "domain-name"),
		strval(domainName))
	setReq.Replace = []*gnmi.Update{hn, dn}

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
	if strings.HasPrefix(oid, snmpLldpStatsTxPortTable) ||
		strings.HasPrefix(oid, snmpLldpStatsRxPortTable) ||
		strings.HasPrefix(oid, snmpLldpLocPortTable) {
		baseOid, locIndex, err = oidIndex(oid)
		return
	} else if strings.HasPrefix(oid, snmpLldpRemTable) {
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
	return
}

// Return MAC address from hex byte string.
func macFromBytes(s []byte) string {
	// XXX_jcr: hex assumption is only right for MAC
	var t bytes.Buffer
	for i := 0; i < len(s); i++ {
		if i != 0 {
			t.WriteString(":")
		}
		t.WriteString(hex.EncodeToString(s[i : i+1]))
	}
	return t.String()
}

func (s *snmp) handleLldpPDU(pdu gosnmp.SnmpPDU) ([]*gnmi.Update, error) {
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
		intfName, ok = s.lldpLocPortIndex[locIndex]
		if !ok {
			// If we have the port ID, add it to the port index map.
			// Otherwise we can't do anything and should return.
			if baseOid != snmpLldpLocPortID {
				return nil, nil
			}
			intfName = string(pdu.Value.([]byte))
			s.lldpLocPortIndex[locIndex] = intfName
		}
	}

	// If we haven't yet seen this remote system, add its ID.
	if remoteID != "" {
		updates = append(updates,
			update(GNMILldpNeighborStatePath(intfName, remoteID, "id"),
				strval(remoteID)))
	}

	var u *gnmi.Update
	switch baseOid {
	case snmpLldpLocPortID:
		updates = append(updates,
			update(GNMILldpIntfConfigPath(intfName, "name"),
				strval(intfName)),
			update(GNMILldpIntfPath(intfName, "name"),
				strval(intfName)),
			update(GNMILldpIntfStatePath(intfName, "name"),
				strval(intfName)))
	case snmpLldpLocChassisID:
		u = update(GNMILldpStatePath("chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpLocChassisIDSubtype:
		u = update(GNMILldpStatePath("chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpLocSysName:
		u = update(GNMILldpStatePath("system-name"),
			strval(pdu.Value))
	case snmpLldpLocSysDesc:
		u = update(GNMILldpStatePath("system-description"),
			strval(pdu.Value))
	case snmpLldpStatsTxPortFramesTotal:
		u = update(GNMILldpIntfCountersPath(intfName, "frame-out"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesDiscard:
		u = update(GNMILldpIntfCountersPath(intfName, "frame-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesErrors:
		u = update(GNMILldpIntfCountersPath(intfName, "frame-error-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesTotal:
		u = update(GNMILldpIntfCountersPath(intfName, "frame-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsDiscard:
		u = update(GNMILldpIntfCountersPath(intfName, "tlv-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsUnrecog:
		u = update(GNMILldpIntfCountersPath(intfName, "tlv-unknown"),
			uintval(pdu.Value))
	case snmpLldpRemPortID:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "port-id"),
			strval(pdu.Value))
	case snmpLldpRemPortIDSubtype:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "port-id-type"),
			strval(openconfig.LLDPPortIDType(pdu.Value.(int))))
	case snmpLldpRemChassisID:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpRemChassisIDSubtype:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpRemSysName:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "system-name"),
			strval(pdu.Value))
	case snmpLldpRemSysDesc:
		u = update(GNMILldpNeighborStatePath(intfName, remoteID, "system-description"),
			strval(pdu.Value))
	}
	if u != nil {
		updates = append(updates, u)
	}
	return updates, nil
}

func (s *snmp) updateLldp() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	setReq := new(gnmi.SetRequest)
	updates := make([]*gnmi.Update, 0)
	updater := func(data gosnmp.SnmpPDU) error {
		u, err := s.handleLldpPDU(data)
		if err != nil {
			return err
		}
		if u != nil {
			updates = append(updates, u...)
		}
		return nil
	}
	if err := gosnmp.Default.Walk(snmpLldpLocalSystemData, updater); err != nil {
		return nil, err
	}

	if err := gosnmp.Default.Walk(snmpLldpRemTable, updater); err != nil {
		return nil, err
	}

	if err := gosnmp.Default.Walk(snmpLldpStatistics, updater); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{GNMIPath("lldp")}
	setReq.Replace = updates
	return setReq, nil
}

func (s *snmp) init() error {
	// Do SNMP networking setup.
	err := snmpNetworkInit()
	if err != nil {
		return fmt.Errorf("Error connecting to device: %v", err)
	}

	return nil
}

func (s *snmp) InitGNMI(client gnmi.GNMIClient) {
	s.client = client
	err := s.init()
	if err != nil {
		glog.Errorf("Error in initialization: %v", err)
		return
	}
	s.initialized = true
}

func (s *snmp) handleErrors(ctx context.Context) {
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

func (s *snmp) Run(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("SNMP provider is uninitialized")
	}

	// Do periodic state updates.
	go OpenConfigPollForever(ctx, s.client, pollInterval,
		s.updateSystemState, s.errc)
	go OpenConfigPollForever(ctx, s.client, pollInterval,
		s.updateInterfaces, s.errc)
	go OpenConfigPollForever(ctx, s.client, pollInterval,
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
	pollInt time.Duration) provider.GNMIProvider {
	gosnmp.Default.Target = address
	gosnmp.Default.Community = community
	pollInterval = pollInt
	gosnmp.Default.Timeout = 2 * pollInterval
	return &snmp{
		errc:             make(chan error),
		interfaceIndex:   make(map[string]string),
		lldpLocPortIndex: make(map[string]string),
		address:          address,
		community:        community,
	}
}
