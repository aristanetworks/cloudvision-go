// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package snmp

import (
	"arista/provider"
	pgnmi "arista/provider/gnmi"
	"arista/provider/openconfig"
	"bytes"
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
	return pgnmi.Update(path, val)
}

func strval(s interface{}) *gnmi.TypedValue {
	t, ok := s.(string)
	if ok {
		return pgnmi.Strval(t)
	}
	u, ok := s.([]byte)
	if ok {
		return pgnmi.Strval(string(u))
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

	// interfaceIndex is a map of SNMP interface index -> name.
	interfaceIndex map[string]string

	// interfaceName is a map of interface name (as discovered in ifTable) -> true.
	// It's used so that we don't include inactive interfaces we see in
	// snmpLldpLocPortTable.
	interfaceName map[string]bool

	// lldpLocPortIndex is a map of lldpLocPortNum -> lldpLocPortId.
	lldpLocPortIndex map[string]string

	// lldpRemoteID is a map of remote system ID -> true. It's used to
	// remember which remote IDs we've already seen in a given round of polling.
	lldpRemoteID map[string]bool

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
	return s.getStringByOID(snmpEntPhysicalSerialNum)
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
func (s *Snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU) ([]*gnmi.Update, error) {
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
	s.lock.Lock()
	// Clear interface index and name maps for each new poll. It should be
	// protected by the lock, because updateLldp needs it, too. :(
	s.interfaceIndex = make(map[string]string)
	s.interfaceName = make(map[string]bool)
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

func (s *Snmp) updateSystemState() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	setReq := new(gnmi.SetRequest)
	sysName, err := s.getStringByOID(snmpSysName)
	if err != nil {
		return nil, err
	}
	hostname := strings.Split(sysName, ".")[0]
	domainName := strings.Join(strings.Split(sysName, ".")[1:], ".")

	hn := update(pgnmi.Path("system", "state", "hostname"), strval(hostname))
	dn := update(pgnmi.Path("system", "state", "domain-name"),
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

func (s *Snmp) handleLldpPDU(pdu gosnmp.SnmpPDU) ([]*gnmi.Update, error) {
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
			// If we have the port ID AND this interface is in the interfaceIndex,
			// add it to the port index map. Otherwise we can't do anything and
			// should return.
			if baseOid != snmpLldpLocPortID {
				return nil, nil
			}
			intfName = string(pdu.Value.([]byte))
			if _, ok = s.interfaceName[intfName]; !ok {
				return nil, nil
			}
			s.lldpLocPortIndex[locIndex] = intfName
		}
	}

	// If we haven't yet seen this remote system, add its ID.
	if remoteID != "" {
		if _, ok = s.lldpRemoteID[remoteID]; !ok {
			updates = append(updates,
				update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "id"),
					strval(remoteID)))
			s.lldpRemoteID[remoteID] = true
		}
	}

	var u *gnmi.Update
	switch baseOid {
	case snmpLldpLocPortID:
		updates = append(updates,
			update(pgnmi.LldpIntfConfigPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfPath(intfName, "name"),
				strval(intfName)),
			update(pgnmi.LldpIntfStatePath(intfName, "name"),
				strval(intfName)))
	case snmpLldpLocChassisID:
		u = update(pgnmi.LldpStatePath("chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpLocChassisIDSubtype:
		u = update(pgnmi.LldpStatePath("chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpLocSysName:
		u = update(pgnmi.LldpStatePath("system-name"),
			strval(pdu.Value))
	case snmpLldpLocSysDesc:
		u = update(pgnmi.LldpStatePath("system-description"),
			strval(pdu.Value))
	case snmpLldpStatsTxPortFramesTotal:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-out"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesDiscard:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesErrors:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-error-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortFramesTotal:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "frame-in"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsDiscard:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "tlv-discard"),
			uintval(pdu.Value))
	case snmpLldpStatsRxPortTLVsUnrecog:
		u = update(pgnmi.LldpIntfCountersPath(intfName, "tlv-unknown"),
			uintval(pdu.Value))
	case snmpLldpRemPortID:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id"),
			strval(pdu.Value))
	case snmpLldpRemPortIDSubtype:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "port-id-type"),
			strval(openconfig.LLDPPortIDType(pdu.Value.(int))))
	case snmpLldpRemChassisID:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id"),
			strval(macFromBytes(pdu.Value.([]byte))))
	case snmpLldpRemChassisIDSubtype:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "chassis-id-type"),
			strval(openconfig.LLDPChassisIDType(pdu.Value.(int))))
	case snmpLldpRemSysName:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-name"),
			strval(pdu.Value))
	case snmpLldpRemSysDesc:
		u = update(pgnmi.LldpNeighborStatePath(intfName, remoteID, "system-description"),
			strval(pdu.Value))
	}
	if u != nil {
		updates = append(updates, u)
	}
	return updates, nil
}

func (s *Snmp) updateLldp() (*gnmi.SetRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.lldpRemoteID = make(map[string]bool)
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
	if err := s.walk(snmpLldpLocalSystemData, updater); err != nil {
		return nil, err
	}

	if err := s.walk(snmpLldpRemTable, updater); err != nil {
		return nil, err
	}

	if err := s.walk(snmpLldpStatistics, updater); err != nil {
		return nil, err
	}

	setReq.Delete = []*gnmi.Path{pgnmi.Path("lldp")}
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
		errc:             make(chan error),
		interfaceIndex:   make(map[string]string),
		interfaceName:    make(map[string]bool),
		lldpLocPortIndex: make(map[string]string),
		address:          address,
		community:        community,
		pollInterval:     pollInt,
		getter:           gosnmp.Default.Get,
		walker:           gosnmp.Default.BulkWalk,
	}
	if err := s.snmpNetworkInit(); err != nil {
		glog.Errorf("Error connecting to device: %v", err)
	}
	return s
}
