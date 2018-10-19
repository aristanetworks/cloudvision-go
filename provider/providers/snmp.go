// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/gopenconfig/eos"
	"arista/gopenconfig/eos/converter"
	"arista/gopenconfig/eos/mapping"
	"arista/gopenconfig/model/node"
	"arista/provider"
	"arista/schema"
	"arista/types"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aristanetworks/glog"
	"github.com/soniah/gosnmp"
)

type snmp struct {
	provider.ReadOnly
	ready            chan struct{}
	done             chan struct{}
	errc             chan error
	ctx              context.Context
	cancel           context.CancelFunc
	interfaceIndex   map[string]string
	lldpLocPortIndex map[string]string
	address          string
	community        string
}

// Has read/write interface been established?
var connected bool

// Time we last heard back from target
var lastAlive time.Time

var pollInt time.Duration

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
	if time.Since(lastAlive) < pollInt {
		return true, nil
	}
	_, err := SNMPGetByOID(snmpSysUpTime)
	return true, err
}

func (s *snmp) WaitForNotification() {
	<-s.ready
}

func (s *snmp) Stop() {
	<-s.ready
	gosnmp.Default.Conn.Close()
	close(s.done)
}

// return base OID, index
func oidIndex(oid string) (string, string, error) {
	finalDotPos := strings.LastIndex(oid, ".")
	if finalDotPos < 0 {
		return "", "", fmt.Errorf("oid '%s' does not match expected format", oid)
	}
	return oid[:finalDotPos], oid[(finalDotPos + 1):], nil
}

// Return the SNMP interface type string corresponding to the
// interface type number we receive. OpenConfig uses the same strings
// as SNMP.
func ifType(t int) string {
	ifTypeMap := []string{"", "other", "regular1822", "hdh1822", "ddnX25",
		"rfc877x25", "ethernetCsmacd", "iso88023Csmacd", "iso88024TokenBus",
		"iso88025TokenRing", "iso88026Man", "starLan", "proteon10Mbit",
		"proteon80Mbit", "hyperchannel", "fddi", "lapb", "sdlc", "ds1",
		"e1", "basicISDN", "primaryISDN", "propPointToPointSerial", "ppp",
		"softwareLoopback", "eon", "ethernet3Mbit", "nsip", "slip", "ultra",
		"ds3", "sip", "frameRelay", "rs232", "para", "arcnet", "arcnetPlus",
		"atm", "miox25", "sonet", "x25ple", "iso88022llc", "localTalk",
		"smdsDxi", "frameRelayService", "v35", "hssi", "hippi", "modem",
		"aal5", "sonetPath", "sonetVT", "smdsIcip", "propVirtual",
		"propMultiplexor", "ieee80212", "fibreChannel", "hippiInterface",
		"frameRelayInterconnect", "aflane8023", "aflane8025", "cctEmul",
		"fastEther", "isdn", "v11", "v36", "g703at64k", "g703at2mb", "qllc",
		"fastEtherFX", "channel", "ieee80211", "ibm370parChan", "escon",
		"dlsw", "isdns", "isdnu", "lapd", "ipSwitch", "rsrb", "atmLogical",
		"ds0", "ds0Bundle", "bsc", "async", "cnr", "iso88025Dtr", "eplrs",
		"arap", "propCnls", "hostPad", "termPad", "frameRelayMPI", "x213",
		"adsl", "radsl", "sdsl", "vdsl", "iso88025CRFPInt", "myrinet",
		"voiceEM", "voiceFXO", "voiceFXS", "voiceEncap", "voiceOverIp",
		"atmDxi", "atmFuni", "atmIma", "pppMultilinkBundle", "ipOverCdlc",
		"ipOverClaw", "stackToStack", "virtualIpAddress", "mpc",
		"ipOverAtm", "iso88025Fiber", "tdlc", "gigabitEthernet", "hdlc",
		"lapf", "v37", "x25mlp", "x25huntGroup", "transpHdlc", "interleave",
		"fast", "ip", "docsCableMaclayer", "docsCableDownstream",
		"docsCableUpstream", "a12MppSwitch", "tunnel", "coffee", "ces",
		"atmSubInterface", "l2vlan", "l3ipvlan", "l3ipxvlan",
		"digitalPowerline", "mediaMailOverIp", "dtm", "dcn", "ipForward",
		"msdsl", "ieee1394", "if-gsn", "dvbRccMacLayer", "dvbRccDownstream",
		"dvbRccUpstream", "atmVirtual", "mplsTunnel", "srp", "voiceOverAtm",
		"voiceOverFrameRelay", "idsl", "compositeLink", "ss7SigLink",
		"propWirelessP2P", "frForward", "rfc1483", "usb", "ieee8023adLag",
		"bgppolicyaccounting", "frf16MfrBundle", "h323Gatekeeper",
		"h323Proxy", "mpls", "mfSigLink", "hdsl2", "shdsl", "ds1FDL", "pos",
		"dvbAsiIn", "dvbAsiOut", "plc", "nfas", "tr008", "gr303RDT",
		"gr303IDT", "isup", "propDocsWirelessMaclayer",
		"propDocsWirelessDownstream", "propDocsWirelessUpstream",
		"hiperlan2", "propBWAp2Mp", "sonetOverheadChannel",
		"digitalWrapperOverheadChannel", "aal2", "radioMAC", "atmRadio",
		"imt", "mvl", "reachDSL", "frDlciEndPt", "atmVciEndPt",
		"opticalChannel", "opticalTransport", "propAtm", "voiceOverCable",
		"infiniband", "teLink", "q2931", "virtualTg", "sipTg", "sipSig",
		"docsCableUpstreamChannel", "econet", "pon155", "pon622", "bridge",
		"linegroup", "voiceEMFGD", "voiceFGDEANA", "voiceDID",
		"mpegTransport", "sixToFour", "gtp", "pdnEtherLoop1",
		"pdnEtherLoop2", "opticalChannelGroup", "homepna", "gfp",
		"ciscoISLvlan", "actelisMetaLOOP", "fcipLink", "rpr", "qam", "lmp",
		"cblVectaStar", "docsCableMCmtsDownstream", "adsl2",
		"macSecControlledIF", "macSecUncontrolledIF", "aviciOpticalEther",
		"atmbond", "voiceFGDOS", "mocaVersion1", "ieee80216WMAN",
		"adsl2plus", "dvbRcsMacLayer", "dvbTdm", "dvbRcsTdma", "x86Laps",
		"wwanPP", "wwanPP2", "voiceEBS", "ifPwType", "ilan", "pip",
		"aluELP", "gpon", "vdsl2", "capwapDot11Profile", "capwapDot11Bss",
		"capwapWtpVirtualRadio", "bits", "docsCableUpstreamRfPort",
		"cableDownstreamRfPort", "vmwareVirtualNic", "ieee802154", "otnOdu",
		"otnOtu", "ifVfiType", "g9981", "g9982", "g9983", "aluEpon",
		"aluEponOnu", "aluEponPhysicalUni", "aluEponLogicalLink",
		"aluGponOnu", "aluGponPhysicalUni", "vmwareNicTeam"}
	if t >= len(ifTypeMap) || t == 0 {
		t = 1
	}
	return ifTypeMap[t]
}

// Return the OpenConfig interface status string corresponding
// to the SNMP interface status.
func ifStatus(status int) string {
	switch uint32(status) {
	case eos.IntfOperUp().EnumValue():
		return converter.IntfOperStatusUp
	case eos.IntfOperDown().EnumValue():
		return converter.IntfOperStatusDown
	case eos.IntfOperTesting().EnumValue():
		return converter.IntfOperStatusTesting
	case eos.IntfOperUnknown().EnumValue():
		return converter.IntfOperStatusUnknown
	case eos.IntfOperDormant().EnumValue():
		return converter.IntfOperStatusDormant
	case eos.IntfOperNotPresent().EnumValue():
		return converter.IntfOperStatusNotPresent
	case eos.IntfOperLowerLayerDown().EnumValue():
		return converter.IntfOperStatusLowerLayerDown
	}
	return ""
}

func intfPath(intfName string, elems ...interface{}) node.Path {
	p := []interface{}{"interfaces", "interface", intfName}
	return node.NewPath(append(p, elems...)...)
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

// Given an incoming PDU, update the appropriate interface state.
func (s *snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU) error {
	// Get/set interface name from index. If there's no mapping, just return and
	// wait for the mapping to show up.
	baseOid, index, err := oidIndex(pdu.Name)
	if err != nil {
		return err
	}
	intfName, ok := s.interfaceIndex[index]
	if !ok && baseOid != snmpIfDescr {
		return nil
	} else if !ok && baseOid == snmpIfDescr {
		intfName = string(pdu.Value.([]byte))
		s.interfaceIndex[index] = intfName
	}

	statePath := intfPath(intfName, "state")
	countersPath := intfPath(intfName, "state", "counters")

	err = nil
	switch baseOid {
	case snmpIfDescr:
		err = OpenConfigUpdateLeaf(s.ctx, statePath, "name",
			string(pdu.Value.([]byte)))
	case snmpIfType:
		err = OpenConfigUpdateLeaf(s.ctx, statePath, "type",
			ifType(pdu.Value.(int)))
	case snmpIfMtu:
		err = OpenConfigUpdateLeaf(s.ctx, statePath, "mtu",
			uint16(pdu.Value.(int)))
	case snmpIfAdminStatus:
		err = OpenConfigUpdateLeaf(s.ctx, statePath, "admin-status",
			ifStatus(pdu.Value.(int)))
	case snmpIfOperStatus:
		err = OpenConfigUpdateLeaf(s.ctx, statePath, "oper-status",
			ifStatus(pdu.Value.(int)))
	case snmpIfInOctets:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-octets",
			uint64(pdu.Value.(uint)))
	case snmpIfInUcastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-unicast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfInMulticastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-multicast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfInBroadcastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-broadcast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfInDiscards:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-discards",
			uint64(pdu.Value.(uint)))
	case snmpIfInErrors:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-errors",
			uint64(pdu.Value.(uint)))
	case snmpIfInUnknownProtos:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "in-unknown-protos",
			uint64(pdu.Value.(uint)))
	case snmpIfOutOctets:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-octets",
			uint64(pdu.Value.(uint)))
	case snmpIfOutUcastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-unicast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfOutMulticastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-multicast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfOutBroadcastPkts:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-broadcast-pkts",
			uint64(pdu.Value.(uint)))
	case snmpIfOutDiscards:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-discards",
			uint64(pdu.Value.(uint)))
	case snmpIfOutErrors:
		err = OpenConfigUpdateLeaf(s.ctx, countersPath, "out-errors",
			uint64(pdu.Value.(uint)))
	}
	// default: ignore update
	return err
}

func (s *snmp) updateInterfaces() error {
	// XXX_jcr: We still need to add code for understanding deletes.
	intfWalk := func(data gosnmp.SnmpPDU) error {
		return s.handleInterfacePDU(data)
	}

	// ifTable
	if err := gosnmp.Default.Walk(snmpIfTable, intfWalk); err != nil {
		return err
	}

	// ifXTable
	return gosnmp.Default.Walk(snmpIfXTable, intfWalk)
}

func (s *snmp) updateSystemState() error {
	sysName, err := SNMPGetByOID(snmpSysName)
	if err != nil {
		return err
	}
	hostname := strings.Split(sysName, ".")[0]
	domainName := strings.Join(strings.Split(sysName, ".")[1:], ".")

	err = OpenConfigUpdateLeaf(s.ctx, node.NewPath("system", "state"),
		"hostname", hostname)
	if err != nil {
		return err
	}
	return OpenConfigUpdateLeaf(s.ctx, node.NewPath("system", "state"),
		"domain-name", domainName)
}

// Return the OpenConfig chassis ID type string corresponding to the
// SNMP chassis ID type value.
func lldpChassisIDType(id int) string {
	switch uint32(id) {
	case eos.LLDPChassisIDChassisComponent().EnumValue():
		return mapping.ChassisIDTypeChassisComponent
	case eos.LLDPChassisIDInterfaceAlias().EnumValue():
		return mapping.ChassisIDTypeInterfaceAlias
	case eos.LLDPChassisIDPortComponent().EnumValue():
		return mapping.ChassisIDTypePortComponent
	case eos.LLDPChassisIDMacAddress().EnumValue():
		return mapping.ChassisIDTypeMACAddress
	case eos.LLDPChassisIDNetworkAddress().EnumValue():
		return mapping.ChassisIDTypeNetworkAddress
	case eos.LLDPChassisIDInterfaceName().EnumValue():
		return mapping.ChassisIDTypeInterfaceName
	case eos.LLDPChassisIDLocal().EnumValue():
		return mapping.ChassisIDTypeLocal
	}
	return ""
}

// Return the OpenConfig port ID type string corresponding to the
// SNMP port ID type value.
func lldpPortIDType(id int) string {
	switch uint32(id) {
	case eos.LLDPPidInterfaceAlias().EnumValue():
		return mapping.PortIDTypeInterfaceAlias
	case eos.LLDPPidPortComponent().EnumValue():
		return mapping.PortIDTypePortComponent
	case eos.LLDPPidMacAddress().EnumValue():
		return mapping.PortIDTypeMACAddress
	case eos.LLDPPidNetworkAddress().EnumValue():
		return mapping.PortIDTypeNetworkAddress
	case eos.LLDPPidInterfaceName().EnumValue():
		return mapping.PortIDTypeInterfaceName
	case eos.LLDPPidAgentCircuitID().EnumValue():
		return mapping.PortIDTypeAgentCircuitID
	case eos.LLDPPidLocal().EnumValue():
		return mapping.PortIDTypeLocal
	}
	return ""
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

// LLDP paths of interest.
func lldpStatePath() node.Path {
	return node.NewPath("lldp", "state")
}
func lldpIntfStatePath(intfName string) node.Path {
	return node.NewPath("lldp", "interfaces", "interface", intfName, "state")
}
func lldpIntfCountersPath(intfName string) node.Path {
	return node.NewPath("lldp", "interfaces", "interface", intfName, "state", "counters")
}
func lldpNeighborStatePath(intfName string, id string) node.Path {
	return node.NewPath("lldp", "interfaces", "interface", intfName,
		"neighbors", "neighbor", id, "state")
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

func (s *snmp) handleLldpPDU(pdu gosnmp.SnmpPDU) error {
	// Split OID into parts.
	locIndex, remoteID, baseOid, err := processLldpOid(pdu.Name)
	if err != nil {
		return err
	}

	// If we haven't yet seen this local interface, add it to our list.
	intfName := ""
	var ok bool
	if locIndex != "" {
		intfName, ok = s.lldpLocPortIndex[locIndex]
		if !ok {
			if baseOid != snmpLldpLocPortID {
				return nil
			}
			intfName = string(pdu.Value.([]byte))
			s.lldpLocPortIndex[locIndex] = intfName
		}
	}

	// If we haven't yet seen this remote system, add its ID.
	if remoteID != "" {
		if err := OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName,
			remoteID), "id", remoteID); err != nil {
			return err
		}
	}

	err = nil
	switch baseOid {
	case snmpLldpLocChassisID:
		err = OpenConfigUpdateLeaf(s.ctx, lldpStatePath(),
			"chassis-id", macFromBytes(pdu.Value.([]byte)))
	case snmpLldpLocChassisIDSubtype:
		err = OpenConfigUpdateLeaf(s.ctx, lldpStatePath(),
			"chassis-id-type", lldpChassisIDType(pdu.Value.(int)))
	case snmpLldpLocSysName:
		err = OpenConfigUpdateLeaf(s.ctx, lldpStatePath(),
			"system-name", string(pdu.Value.([]byte)))
	case snmpLldpLocSysDesc:
		err = OpenConfigUpdateLeaf(s.ctx, lldpStatePath(),
			"system-description", string(pdu.Value.([]byte)))
	case snmpLldpLocPortID:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfStatePath(intfName),
			"name", intfName)
	case snmpLldpStatsTxPortFramesTotal:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"frame-out", uint64(pdu.Value.(uint)))
	case snmpLldpStatsRxPortFramesDiscard:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"frame-discard", uint64(pdu.Value.(uint)))
	case snmpLldpStatsRxPortFramesErrors:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"frame-error-in", uint64(pdu.Value.(uint)))
	case snmpLldpStatsRxPortFramesTotal:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"frame-in", uint64(pdu.Value.(uint)))
	case snmpLldpStatsRxPortTLVsDiscard:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"tlv-discard", uint64(pdu.Value.(uint)))
	case snmpLldpStatsRxPortTLVsUnrecog:
		err = OpenConfigUpdateLeaf(s.ctx, lldpIntfCountersPath(intfName),
			"tlv-unknown", uint64(pdu.Value.(uint)))
	case snmpLldpRemPortID:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"port-id", string(pdu.Value.([]byte)))
	case snmpLldpRemPortIDSubtype:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"port-id-type", lldpPortIDType(pdu.Value.(int)))
	case snmpLldpRemChassisID:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"chassis-id", macFromBytes(pdu.Value.([]byte)))
	case snmpLldpRemChassisIDSubtype:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"chassis-id-type", lldpChassisIDType(pdu.Value.(int)))
	case snmpLldpRemSysName:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"system-name", string(pdu.Value.([]byte)))
	case snmpLldpRemSysDesc:
		err = OpenConfigUpdateLeaf(s.ctx, lldpNeighborStatePath(intfName, remoteID),
			"system-description", string(pdu.Value.([]byte)))
	}
	return err
}

func (s *snmp) updateLLDP() error {
	walker := func(data gosnmp.SnmpPDU) error {
		return s.handleLldpPDU(data)
	}
	if err := gosnmp.Default.Walk(snmpLldpLocalSystemData, walker); err != nil {
		return err
	}

	if err := gosnmp.Default.Walk(snmpLldpRemTable, walker); err != nil {
		return err
	}

	return gosnmp.Default.Walk(snmpLldpStatistics, walker)
}

func (s *snmp) init(ch chan<- types.Notification) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	// Do SNMP networking setup.
	err := snmpNetworkInit()
	if err != nil {
		return fmt.Errorf("Error connecting to device: %v", err)
	}

	// Set up notifying data tree.
	s.ctx, err = OpenConfigNotifyingTree(ctx, ch, s.errc)
	if err != nil {
		return err
	}

	close(s.ready)

	return nil
}

func (s *snmp) Run(schema *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	// Do necessary setup.
	err := s.init(ch)
	if err != nil {
		glog.Infof("Error in initialization: %v", err)
		return
	}

	// Do periodic state updates
	tick := time.NewTicker(pollInt)
	defer tick.Stop()
	defer s.cancel()
	for {
		select {
		case <-tick.C:
			if err := s.updateSystemState(); err != nil {
				glog.Infof("Failure in updateSystemState: %v", err)
				return
			}
			if err := s.updateInterfaces(); err != nil {
				glog.Infof("Failure in updateInterfaces: %s", err)
			}
			if err := s.updateLLDP(); err != nil {
				glog.Infof("Failure in updateLLDP: %v", err)
			}
		case <-s.done:
			return
		case err := <-s.errc:
			glog.Errorf("Failure in gNMI stream: %v", err)
			return
		}
	}
}

// NewSNMPProvider returns a new SNMP provider for the device at 'address'
// using a community value for authentication and pollInterval for rate
// limiting requests.
func NewSNMPProvider(address string, community string,
	pollInterval time.Duration) provider.Provider {
	gosnmp.Default.Target = address
	gosnmp.Default.Community = community
	gosnmp.Default.Timeout = 30 * time.Second
	pollInt = pollInterval
	return &snmp{
		ready:            make(chan struct{}),
		done:             make(chan struct{}),
		errc:             make(chan error),
		interfaceIndex:   make(map[string]string),
		lldpLocPortIndex: make(map[string]string),
		address:          address,
		community:        community,
	}
}
