// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmpoc

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/pdu"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/gosnmp/gosnmp"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// A Mapper contains some logic for producing gNMI updates based on the
// contents of a pdu.Store and a mapper data cache.
type Mapper = func(smi.Store, pdu.Store, *sync.Map, Logger) ([]*gnmi.Update, error)

// A ValueProcessor takes an arbitrary value and returns a
// gnmi.TypedValue, possibly doing additional processing first.
type ValueProcessor func(interface{}) *gnmi.TypedValue

// Generic value processors
func strval(s interface{}) *gnmi.TypedValue {
	t := sanitizedString(s)
	if t == "" {
		return nil
	}
	return pgnmi.Strval(t)
}

// BytesToSanitizedString removes all but ASCII characters 32-126 to
// keep the JSON unmarshaler happy.
func BytesToSanitizedString(b []byte) string {
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

func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToUint64(u); err == nil {
		return pgnmi.Uintval(v)
	}
	return nil
}

func sanitizedString(s interface{}) string {
	switch t := s.(type) {
	case string:
		return t
	case []byte:
		return BytesToSanitizedString(t)
	case int64, int32, int16, int8, int:
		return strconv.Itoa(s.(int))
	case uint64, uint32, uint16, uint8, uint:
		return strconv.FormatUint(s.(uint64), 10)
	}
	return ""
}

func intval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToInt64(u); err == nil {
		return pgnmi.Intval(v)
	}
	return nil
}

func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return pgnmi.Update(path, val)
}

func scalarMapper(ps pdu.Store, path string, oid string,
	vp ValueProcessor) (*gnmi.Update, error) {
	pdu, err := ps.GetScalar(oid)
	if err != nil {
		return nil, err
	} else if pdu == nil {
		return nil, nil
	}

	val := vp(pdu.Value)
	if val == nil {
		return nil, nil
	}

	return update(pgnmi.PathFromString(path), val), nil
}

func scalarMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		md *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		u, err := scalarMapper(ps, path, oid, vp)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, nil
		}
		return []*gnmi.Update{u}, nil
	}
}

func getTabular(ps pdu.Store, oid string) ([]*gosnmp.SnmpPDU, error) {
	pdus, err := ps.GetTabular(oid)
	if err != nil {
		return nil, err
	} else if len(pdus) == 0 {
		return nil, nil
	}

	return pdus, nil
}

// SNMPErrCodes maps the SNMP error codes to their corresponding
// text descriptions.
var SNMPErrCodes = map[gosnmp.SNMPError]string{
	0: "noError", 1: "tooBig", 2: "noSuchName", 3: "badValue",
	4: "readOnly", 5: "genErr", 6: "noAccess", 7: "wrongType",
	8: "wrongLength", 9: "wrongEncoding", 10: "wrongValue",
	11: "noCreation", 12: "inconsistentValue", 13: "resourceUnavailable",
	14: "commitFailed", 15: "undoFailed", 16: "authorizationError",
	17: "notWritable", 18: "inconsistentName",
}

// Some implementations will return a hostname only, while others
// will return a fully qualified domain name.
func processSysName(s interface{}) *gnmi.TypedValue {
	sysName := sanitizedString(s)
	return strval(strings.SplitN(sysName, ".", 2)[0])
}

func processDomainName(s interface{}) *gnmi.TypedValue {
	sysName := sanitizedString(s)
	ss := strings.SplitN(sysName, ".", 2)
	if len(ss) > 1 {
		return strval(ss[1])
	}
	return nil
}

var now = time.Now

// Get boot-time by subtracting the target device's uptime from
// the Collector's current time. This isn't really correct--the
// boot time we produce is a blend of UNIX time on the Collector
// and UNIX time on the target device (which may not be in sync),
// and the target device's time may be recorded well before the
// Collector's. There doesn't seem to be a great way to get the
// device's time using SNMP. Assuming the two devices are
// roughly in sync, though, this shouldn't be disastrous.
// Convert it to nanoseconds as expected by the openconfig model
func processBootTime(x interface{}) *gnmi.TypedValue {
	t, err := provider.ToInt64(x)
	if err != nil {
		return nil
	}
	if t == 0 {
		return nil
	}
	return intval(now().UnixNano() - t*10000000)
}

func setIndexStringMappings(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, oid, indexName, mappingName string) error {
	_, ok := mapperData.Load(mappingName)
	if ok {
		return nil
	}

	mp := map[string]string{}

	pdus, err := getTabular(ps, oid)
	if err != nil {
		return err
	}

	for _, p := range pdus {
		indexVal, err := pdu.IndexValueByName(ss, p, indexName)
		if err != nil {
			return err
		}
		val := sanitizedString(p.Value)
		mp[indexVal] = val
	}

	mapperData.Store(mappingName, mp)
	return nil
}

func getMapperStringMapping(mapperData *sync.Map,
	mappingName, key string) (string, error) {
	m, ok := mapperData.Load(mappingName)
	if !ok {
		return "", fmt.Errorf("No mapping for '%s' in mapperData", mappingName)
	}
	mp := m.(map[string]string)
	return mp[key], nil
}

// interface helpers
func getIntfName(mapperData *sync.Map, ifIndex string) (string, error) {
	return getMapperStringMapping(mapperData, "intfName", ifIndex)
}

func setIntfNames(ss smi.Store, ps pdu.Store, mapperData *sync.Map) error {
	return setIndexStringMappings(ss, ps, mapperData, "ifDescr",
		"ifIndex", "intfName")
}

// generic mapper for PDUs from ifTable
func ifTableMapper(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger, path string,
	oid string, vp ValueProcessor) ([]*gnmi.Update, error) {
	pdus, err := getTabular(ps, oid)
	if err != nil || pdus == nil {
		return nil, err
	}

	updates := []*gnmi.Update{}
	for _, p := range pdus {
		indexVals, err := pdu.IndexValues(ss, p)
		if err != nil {
			logger.Errorf("failed to index value with error: %v", err)
			continue
		}
		ifIndex := indexVals[0]
		ifDescr, err := getIntfName(mapperData, ifIndex)
		if err != nil || ifDescr == "" {
			if err = setIntfNames(ss, ps, mapperData); err != nil {
				return nil, err
			}
		}
		ifDescr, err = getIntfName(mapperData, ifIndex)
		if err != nil {
			return nil, err
		}
		if ifDescr == "" {
			return nil, fmt.Errorf("No ifDescr for ifIndex '%s'", ifIndex)
		}
		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, ifDescr))
		processedVal := vp(p.Value)
		// do not send any update if value processor returns nil
		if processedVal != nil {
			updates = append(updates, update(fullPath, vp(p.Value)))
		}
	}
	return updates, nil
}

func ifTableMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return ifTableMapper(ss, ps, mapperData, logger, path, oid, vp)
	}
}

// Build a map from ipAddress -> ifDescr.
func buildIPAddrMap(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger) error {
	m, ok := mapperData.Load("ipAddrIntfName")
	if !ok {
		mapperData.Store("ipAddrIntfName", make(map[string]string))
		m, ok = mapperData.Load("ipAddrIntfName")
		if !ok {
			return errors.New("failed to load ipAddrIntfName from mapperData")
		}
	}

	mp := m.(map[string]string)

	// Build the ifIndex -> ifDescr mapping
	err := setIndexStringMappings(ss, ps, mapperData, "ifDescr",
		"ifIndex", "intfName")
	if err != nil {
		return err
	}

	n, ok := mapperData.Load("intfName")
	if !ok {
		return errors.New("Failed to load ipAdEntAddr from mapperData")
	}
	ifIdxToIfNameMap := n.(map[string]string)

	// Build the ipAddress -> ifIndex mapping
	_, ok = mapperData.Load("ipAddressIfIndex")
	if !ok {
		ipMap := map[string]string{}
		pdus, err := getTabular(ps, "ipAddressIfIndex")
		if err != nil {
			return err
		}

		for _, p := range pdus {
			ipAddrIdxVal, err := pdu.IndexValues(ss, p)
			if err != nil {
				logger.Errorf("failed to index value with error: %v", err)
				continue
			}
			if ipAddrIdxVal[0] != "1" && ipAddrIdxVal[0] != "2" {
				// when the ipaddress type is ipv4z(3) or ipv6z(4) or dns(5) ignore these addresses
				continue
			}
			if ipAddrIdxVal[1] == "" {
				continue
			}
			val := sanitizedString(p.Value)
			ipMap[ipAddrIdxVal[1]] = val
		}

		mapperData.Store("ipAddressIfIndex", ipMap)
	}

	i, ok := mapperData.Load("ipAddressIfIndex")
	if !ok {
		return errors.New("failed to load ipAddressIfIndex from mapperData")
	}
	ipAddrToIfIdxMap := i.(map[string]string)

	// map[ip address] -> interface name
	for ipAddr, idx := range ipAddrToIfIdxMap {
		mp[ipAddr] = ifIdxToIfNameMap[idx]
	}

	mapperData.Store("ipAddrIntfName", mp)
	return nil
}

func getInterfaceFromIPAddr(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, ipAddr string, logger Logger) (string, error) {
	_, ok := mapperData.Load("ipAddrIntfName")
	if !ok {
		if err := buildIPAddrMap(ss, ps, mapperData, logger); err != nil {
			return "", err
		}
	}

	m, ok := mapperData.Load("ipAddrIntfName")
	if !ok {
		return "", fmt.Errorf("ipAddrIntfName doesn't exist in mapperData")
	}
	mp := m.(map[string]string)

	intfName, ok := mp[ipAddr]
	if !ok {
		return "", nil
	}

	return intfName, nil
}

// generic mapper for PDUs from ipAddressTable
func ifSubIntfIPMapper(ss smi.Store, ps pdu.Store, mapperData *sync.Map, logger Logger, path string,
	oid string, vp ValueProcessor) ([]*gnmi.Update, error) {

	isIPv4 := strings.Contains(path, "ipv4")

	pdus, err := getTabular(ps, oid)
	if err != nil || pdus == nil {
		return nil, err
	}

	updates := []*gnmi.Update{}
	for _, p := range pdus {
		ipAddrIdxVal, err := pdu.IndexValues(ss, p)
		if err != nil {
			logger.Errorf("failed to index value with error: %v", err)
			continue
		}

		// An example oid is .1.3.6.1.2.1.4.34.1.3.1.4.172.31.26.139, where
		// .1.3.6.1.2.1.4.34.1.3 represents ipAddressIfIndex,
		// 1 represents ipAddressAddrType (ipv4 in this case)
		// 4 represents the next number of bytes to be considered as ip address
		// 172.31.26.139 represents ipAddressAddr which is the ip address.
		//
		// .1.3.6.1.2.1.4.34.1.3.2.16.253.122.98.159.82.164.119.119.0.0.0.0.0.31.26.139, where
		// 2 represents ipAddressAddrType (ipv6 in this case)
		// 16 represents the next number of bytes to be considered as ip address
		// 253.122.98.159.82.164.119.119.0.0.0.0.0.31.26.139 represents ipAddressAddr which is the
		// ip address(fd:7a:62:9f:52:a4:77:77:00:00:00:00:00:1f:1a:8b).
		if ipAddrIdxVal[0] != "1" && ipAddrIdxVal[0] != "2" {
			// when the ipaddress type is ipv4z(3) or ipv6z(4) or dns(5) ignore these addresses
			continue
		}

		if isIPv4 && ipAddrIdxVal[0] == "2" {
			// when path is ipv4 but oid is of ipv6 address
			continue
		}

		if !isIPv4 && ipAddrIdxVal[0] == "1" {
			// when path is ipv6 but oid is of ipv4 address
			continue
		}

		// Get interface name corresponding to this ipaddress.
		intfName, err := getInterfaceFromIPAddr(ss, ps, mapperData, ipAddrIdxVal[1], logger)
		if err != nil {
			return nil, err
		} else if intfName == "" {
			continue
		}

		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, intfName, 0, ipAddrIdxVal[1]))
		updates = append(updates, update(fullPath, strval(ipAddrIdxVal[1])))
	}
	return updates, nil
}

func ifSubIntfIPMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return ifSubIntfIPMapper(ss, ps, mapperData, logger, path, oid, vp)
	}
}

func alternateIntfName(intfName string) string {
	if strings.Contains(intfName, "FortyGigabitEthernet") {
		intfName = strings.Replace(intfName, "FortyGigabitEthernet", "Fo", 1)
	} else if strings.Contains(intfName, "Fo") {
		intfName = strings.Replace(intfName, "Fo", "FortyGigabitEthernet", 1)
	} else if strings.Contains(intfName, "TenGigabitEthernet") {
		intfName = strings.Replace(intfName, "TenGigabitEthernet", "Te", 1)
	} else if strings.Contains(intfName, "Te") {
		intfName = strings.Replace(intfName, "Te", "TenGigabitEthernet", 1)
	} else if strings.Contains(intfName, "TwentyFiveGigE") {
		intfName = strings.Replace(intfName, "TwentyFiveGigE", "Twe", 1)
	} else if strings.Contains(intfName, "Twe") {
		intfName = strings.Replace(intfName, "Twe", "TwentyFiveGigE", 1)
	} else if strings.Contains(intfName, "TwoGigabitEthernet") {
		intfName = strings.Replace(intfName, "TwoGigabitEthernet", "Tw", 1)
	} else if strings.Contains(intfName, "Tw") {
		intfName = strings.Replace(intfName, "Tw", "TwoGigabitEthernet", 1)
	} else if strings.Contains(intfName, "FiveGigabitEthernet") {
		intfName = strings.Replace(intfName, "FiveGigabitEthernet", "Fi", 1)
	} else if strings.Contains(intfName, "Fi") {
		intfName = strings.Replace(intfName, "Fi", "FiveGigabitEthernet", 1)
	} else if strings.Contains(intfName, "GigabitEthernet") {
		intfName = strings.Replace(intfName, "GigabitEthernet", "Gi", 1)
	} else if strings.Contains(intfName, "Gi") {
		intfName = strings.Replace(intfName, "Gi", "GigabitEthernet", 1)
	} else if strings.Contains(intfName, "Ethernet") {
		intfName = strings.Replace(intfName, "Ethernet", "Eth", 1)
	} else if strings.Contains(intfName, "Eth") {
		intfName = strings.Replace(intfName, "Eth", "Ethernet", 1)
	} else if strings.Contains(intfName, "Mgmt") {
		intfName = strings.Replace(intfName, "Mgmt", "Management", 1)
	} else if strings.Contains(intfName, "Management") {
		intfName = strings.Replace(intfName, "Management", "Mgmt", 1)
	}
	return intfName
}

// Build a map from lldpLocPortNum -> ifDescr.
func buildLldpLocPortNumMap(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger) error {
	m, ok := mapperData.Load("lldpLocPortNum")
	if !ok {
		mapperData.Store("lldpLocPortNum", make(map[string]string))
		m, ok = mapperData.Load("lldpLocPortNum")
		if !ok {
			return errors.New("Failed to load lldpLocPortNum from mapperData")
		}
	}

	mp := m.(map[string]string)
	ifDescrs, err := getTabular(ps, "ifDescr")
	if err != nil {
		return fmt.Errorf("buildLldpLocPortNumMap: %s", err)
	}
	ifDescrMap := make(map[string]bool)
	for _, p := range ifDescrs {
		ifDescrMap[string(p.Value.([]byte))] = true
	}

	// XXX NOTE: The RFC says lldpLocPortDesc should have the
	// same value as a corresponding ifDescr object, but in
	// practice it seems more common to be have lldpLocPortId
	// equal to an ifDescr object, and lldpLocPortDesc is all
	// over the map--sometimes empty, sometimes set to
	// ifAlias. So just use whichever one matches ifDescr.
	for _, oid := range []string{"lldpLocPortId", "lldpV2LocPortId",
		"lldpLocPortDesc", "lldpV2LocPortDesc"} {
		pdus, err := getTabular(ps, oid)
		if err != nil {
			return fmt.Errorf("buildLldpLocPortNumMap: %s", err)
		}
		for _, p := range pdus {
			indexVals, err := pdu.IndexValues(ss, p)
			if err != nil {
				logger.Errorf("failed to index value with error: %v", err)
				continue
			}
			portNum := indexVals[0]
			intfName := string(p.Value.([]byte))
			if _, ok := ifDescrMap[intfName]; !ok {
				// We've seen some implementations where the lldpLocPortTable interface
				// name is an abbreviation of the the ifTable name.
				intfName = alternateIntfName(intfName)
				if _, ok = ifDescrMap[intfName]; !ok {
					continue
				}
			}
			mp[portNum] = intfName
		}

		// If we haven't built up the full mapping, keep trying.
		if len(mp) == len(ifDescrMap) {
			return nil
		}
	}

	return nil
}

func getInterfaceFromLldpPortNum(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, port string, logger Logger) (string, error) {
	_, ok := mapperData.Load("lldpLocPortNum")
	if !ok {
		if err := buildLldpLocPortNumMap(ss, ps, mapperData, logger); err != nil {
			return "", err
		}
	}
	m, ok := mapperData.Load("lldpLocPortNum")
	if !ok {
		return "", fmt.Errorf("lldpLocPortNum doesn't exist in mapperData")
	}
	mp := m.(map[string]string)
	intfName, ok := mp[port]
	if !ok {
		// XXX NOTE: We probably want to ignore failures here, but this
		// shouldn't really happen. Maybe log?
		return "", nil
	}
	return intfName, nil
}

// bpsToOCEthernetSpeeds maps bits per second(received from `ifSpeed` snmp oid) values
// to their corresponding openconfig ETHERNET_SPEED identities
var bpsToOCEthernetSpeeds = map[uint64]string{
	1e7:    "SPEED_10MB",
	1e8:    "SPEED_100MB",
	1e9:    "SPEED_1GB",
	2.5e9:  "SPEED_2500MB",
	5e9:    "SPEED_5GB",
	1e10:   "SPEED_10GB",
	2.5e10: "SPEED_25GB",
	4e10:   "SPEED_40GB",
	5e10:   "SPEED_50GB",
	1e11:   "SPEED_100GB",
	2e11:   "SPEED_200GB",
	4e11:   "SPEED_400GB",
	6e11:   "SPEED_600GB",
	8e11:   "SPEED_800GB",
}

func ifHighSpeedStrVal(s interface{}) *gnmi.TypedValue {
	return speedStrVal(s, true)
}

func ifSpeedStrVal(s interface{}) *gnmi.TypedValue {
	return speedStrVal(s, false)
}

const bpsPerMbps = 1e6

func speedStrVal(s interface{}, isHighSpeed bool) *gnmi.TypedValue {
	unknownSpeed := pgnmi.Strval("SPEED_UNKNOWN")

	bps, err := provider.ToUint64(s)
	if err != nil {
		return unknownSpeed
	}

	// highSpeed values will be in mbps, convert to bps before checking in map
	if isHighSpeed {
		bps = bps * bpsPerMbps
	}

	ethSpeed, found := bpsToOCEthernetSpeeds[bps]
	if !found {
		return unknownSpeed
	}
	return pgnmi.Strval(ethSpeed)
}

func macAddrStrVal(s interface{}) *gnmi.TypedValue {
	sByte, ok := s.([]byte)
	if !ok {
		return nil
	}

	macAddr := MacFromBytes(sByte)
	if macAddr == "" {
		return nil
	}

	return pgnmi.Strval(macAddr)
}

// regex pattern for mac-address as specified in openconfig
var macAddrRegex = regexp.MustCompile(`^[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}$`)

// MacFromBytes returns a MAC address from a string or hex byte string.
func MacFromBytes(s []byte) string {
	// string case
	if len(s) == 17 {
		return string(s)
	}

	// else assume hex string
	mac := net.HardwareAddr(s).String()
	if !macAddrRegex.MatchString(mac) {
		return ""
	}
	return mac
}

// IPFromBytes returns an IP address from a string or byte string.
func IPFromBytes(s []byte) string {
	// bytes case
	if len(s) == 5 && int(s[0]) == 1 {
		// IPv4
		return net.IPv4(s[1], s[2], s[3], s[4]).String()
	} else if len(s) == 17 && int(s[0]) == 2 {
		// IPv6
		return net.IP(s[1:]).String()
	}

	return string(s)
}

func processChassisID(p *gosnmp.SnmpPDU, subtype int) (interface{}, error) {
	switch subtype {
	case 4:
		return MacFromBytes(p.Value.([]byte)), nil
	case 5:
		return IPFromBytes(p.Value.([]byte)), nil
	}
	// Trim any NULL bytes from the value.
	return string(bytes.Trim(p.Value.([]byte), "\x00")), nil
}

func processPortID(p *gosnmp.SnmpPDU, subtype int) (interface{}, error) {
	switch subtype {
	case 3:
		return MacFromBytes(p.Value.([]byte)), nil
	case 4:
		return IPFromBytes(p.Value.([]byte)), nil
	}
	return string(p.Value.([]byte)), nil
}

func lldpChassisIDMapper(ss smi.Store, ps pdu.Store, mapperData *sync.Map,
	logger Logger, path, idOid, subtypeOid string,
	vp ValueProcessor) ([]*gnmi.Update, error) {
	pcid, err := ps.GetScalar(idOid)
	if err != nil || pcid == nil {
		return nil, err
	}
	pst, err := ps.GetScalar(subtypeOid)
	if err != nil {
		return nil, err
	}

	// We have seen implementations where walking lldpLocalSystemData
	// doesn't return lldpLocChassisIdSubtype. Assume it's a macAddress
	// and see if it works.
	chassisSubtypeMacAddress := int(4)
	v := chassisSubtypeMacAddress
	if pst != nil {
		v = pst.Value.(int)
	}
	cid, err := processChassisID(pcid, v)
	if err != nil || cid == "" {
		return nil, err
	}
	return []*gnmi.Update{update(pgnmi.PathFromString(path), vp(cid))}, nil
}

func lldpChassisIDFn(path, idOid, subtypeOid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return lldpChassisIDMapper(ss, ps, mapperData, logger, path,
			idOid, subtypeOid, vp)
	}
}

func lldpLocPortTableMapper(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger, path string, oid string,
	vp ValueProcessor) ([]*gnmi.Update, error) {
	pdus, err := getTabular(ps, oid)
	if err != nil || pdus == nil {
		return nil, err
	}

	updates := []*gnmi.Update{}

	for _, p := range pdus {
		indexVals, err := pdu.IndexValues(ss, p)
		if err != nil {
			logger.Errorf("failed to index value with error: %v", err)
			continue
		}
		lldpPortNum := indexVals[0]
		// Get interface name corresponding to this port number.
		intfName, err := getInterfaceFromLldpPortNum(ss, ps, mapperData, lldpPortNum, logger)
		if err != nil {
			return nil, err
		} else if intfName == "" {
			continue
		}

		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, intfName))

		// If we're mapping to interface name, just use the intfName to make sure
		// they match.
		val := vp(p.Value)
		if strings.HasSuffix(path, "name") {
			val = vp(intfName)
		}
		updates = append(updates, update(fullPath, val))
	}
	return updates, nil
}

func lldpLocPortTableMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return lldpLocPortTableMapper(ss, ps, mapperData, logger, path, oid, vp)
	}
}

func processLldpRemTableVal(ps pdu.Store, p *gosnmp.SnmpPDU, oid,
	lldpPortNum, lldpRemIndex string) (interface{}, error) {

	var pdus []*gosnmp.SnmpPDU
	var err error

	// Special-case chassis and port IDs, which need a subtype to be
	// interpreted.
	portNumOid := "lldpRemLocalPortNum"
	remIdxOid := "lldpRemIndex"
	v2, err := regexp.MatchString("V2", oid)
	if err != nil {
		return nil, err
	}
	if v2 {
		portNumOid = "lldpV2RemLocalIfIndex"
		remIdxOid = "lldpV2RemIndex"
	}

	lldpPortNumIdx := pdu.Index{Name: portNumOid, Value: lldpPortNum}
	lldpRemIdx := pdu.Index{Name: remIdxOid, Value: lldpRemIndex}
	switch oid {
	case "lldpRemChassisId":
		pdus, err = ps.GetTabular("lldpRemChassisIdSubtype", lldpPortNumIdx, lldpRemIdx)
	case "lldpV2RemChassisId":
		pdus, err = ps.GetTabular("lldpV2RemChassisIdSubtype", lldpPortNumIdx, lldpRemIdx)
	case "lldpRemPortId":
		pdus, err = ps.GetTabular("lldpRemPortIdSubtype", lldpPortNumIdx, lldpRemIdx)
	case "lldpV2RemPortId":
		pdus, err = ps.GetTabular("lldpV2RemPortIdSubtype", lldpPortNumIdx, lldpRemIdx)
	default:
		return p.Value, nil
	}

	if err != nil {
		return nil, err
	}
	if len(pdus) == 0 {
		return nil, fmt.Errorf("Expected PDUs for subtype of %s, got none", oid)
	}

	v := pdus[0].Value.(int)
	if oid == "lldpRemChassisId" || oid == "lldpV2RemChassisId" {
		return processChassisID(p, v)
	} else if oid == "lldpRemPortId" || oid == "lldpV2RemPortId" {
		return processPortID(p, v)
	}
	return "", fmt.Errorf("processLldpRemTableVal shouldn't get here, oid = %s", oid)
}

func lldpRemTableMapper(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger, path string, oid string,
	vp ValueProcessor) ([]*gnmi.Update, error) {
	pdus, err := getTabular(ps, oid)
	if err != nil || pdus == nil {
		return nil, err
	}

	updates := []*gnmi.Update{}
	lldpIDRegex, err := regexp.Compile("^/lldp.*state/id$")
	if err != nil {
		return nil, err
	}

	for _, p := range pdus {
		lldpPortNum := ""
		lldpRemIndex := ""
		for _, o := range []string{"lldpRemLocalPortNum", "lldpV2RemLocalIfIndex"} {
			lldpPortNum, _ = pdu.IndexValueByName(ss, p, o)
			if lldpPortNum != "" {
				break
			}
		}
		for _, o := range []string{"lldpRemIndex", "lldpV2RemIndex"} {
			lldpRemIndex, _ = pdu.IndexValueByName(ss, p, o)
			if lldpRemIndex != "" {
				break
			}
		}

		// get the interface name corresponding to this lldpRemLocalPortNum
		intfName, err := getInterfaceFromLldpPortNum(ss, ps, mapperData, lldpPortNum, logger)
		if err != nil {
			return nil, err
		} else if intfName == "" {
			// Log error but keep going.
			logger.Debugf("Failed to convert lldpRemLocaLPortNum '%s' to intf name", lldpPortNum)
			continue
		}
		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, intfName, lldpRemIndex))

		neighborID := lldpIDRegex.MatchString(path)
		if (oid == "lldpRemPortId" || oid == "lldpV2RemPortId") && neighborID {
			updates = append(updates, update(fullPath, vp(lldpRemIndex)))
			continue
		}

		v, err := processLldpRemTableVal(ps, p, oid, lldpPortNum, lldpRemIndex)
		if err != nil {
			return nil, err
		}
		val := vp(v)
		// If the interface lldp neighbour's system-name or system-description is not set,
		// the value processor returns nil, in which case no updates should be generated.
		if val != nil {
			updates = append(updates, update(fullPath, val))
		}
	}
	return updates, nil
}

func lldpRemTableMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return lldpRemTableMapper(ss, ps, mapperData, logger, path, oid, vp)
	}
}

func processEntPhysicalTableVal(oid string, pdu *gosnmp.SnmpPDU) (interface{}, error) {
	if oid != "entPhysicalClass" {
		return pdu.Value, nil
	}

	v := pdu.Value.(int)
	ocComponentTypeMap := map[int]string{
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
	class, ok := ocComponentTypeMap[v]
	if !ok {
		return nil, fmt.Errorf("Unexpected entPhysicalClass value %v", v)
	}
	if class != "" {
		class = "openconfig-platform-types:" + class
	}
	return class, nil
}

func entPhysicalMapper(ss smi.Store, ps pdu.Store,
	mapperData *sync.Map, logger Logger, path string, oid string,
	vp ValueProcessor) ([]*gnmi.Update, error) {
	pdus, err := getTabular(ps, oid)
	if err != nil || pdus == nil {
		return nil, err
	}

	updates := []*gnmi.Update{}
	nameRegex, err := regexp.Compile(".*/(name|id)$")
	if err != nil {
		return nil, err
	}

	for _, p := range pdus {
		indexVals, err := pdu.IndexValues(ss, p)
		if err != nil {
			logger.Errorf("failed to fetch index value with error: %v", err)
			continue
		}
		epi := indexVals[0]
		var v interface{}
		namePath := nameRegex.MatchString(path)
		if oid == "entPhysicalDescr" && namePath {
			v = epi
		} else {
			v, err = processEntPhysicalTableVal(oid, p)
			if err != nil {
				return nil, err
			}
		}
		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, epi))
		val := vp(v)
		if val != nil {
			updates = append(updates, update(fullPath, val))
		}
	}

	return updates, nil
}

func entPhysicalTableMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store,
		mapperData *sync.Map, logger Logger) ([]*gnmi.Update, error) {
		return entPhysicalMapper(ss, ps, mapperData, logger, path, oid, vp)
	}
}

// downcastUint returns the downsized uint value.
func downcastUint(x interface{}, max uint64) uint64 {
	v, err := provider.ToUint64(x)
	if err != nil {
		return 0
	}
	if v > max {
		v = max
	}
	return v
}

// mappers
var (
	// /interfaces
	interfacePath              = "/interfaces/interface[name=%s]/"
	interfaceEthernetPath      = interfacePath + "ethernet/"
	interfaceStatePath         = interfacePath + "state/"
	interfaceConfigPath        = interfacePath + "config/"
	interfaceSubIntfPath       = interfacePath + "subinterfaces/subinterface[index=%d]/"
	interfaceCounterPath       = interfaceStatePath + "counters/"
	interfaceEthernetStatePath = interfaceEthernetPath + "state/"
	interfaceName              = ifTableMapperFn(interfacePath+"name",
		"ifDescr", strval)
	interfaceStateName = ifTableMapperFn(interfaceStatePath+"name",
		"ifDescr", strval)
	interfaceConfigName = ifTableMapperFn(interfaceConfigPath+"name",
		"ifDescr", strval)
	interfaceMtu = ifTableMapperFn(interfaceStatePath+"mtu",
		"ifMtu", func(x interface{}) *gnmi.TypedValue {
			// Make MTU values at most uint16_max, since the SNMP ifMtu
			// type is int32 while the OpenConfig equivalent is uint16.
			return uintval(downcastUint(x, math.MaxUint16))
		})
	interfaceType = ifTableMapperFn(interfaceStatePath+"type",
		"ifType", func(x interface{}) *gnmi.TypedValue {
			return strval("iana-if-type:" + openconfig.InterfaceType(x.(int)))
		})
	interfaceAdminStatus = ifTableMapperFn(interfaceStatePath+"admin-status",
		"ifAdminStatus", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.IntfAdminStatus(x.(int)))
		})
	interfaceOperStatus = ifTableMapperFn(interfaceStatePath+"oper-status",
		"ifOperStatus", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.IntfOperStatus(x.(int)))
		})
	interfaceIndex = ifTableMapperFn(interfaceStatePath+"ifindex",
		"ifIndex", func(x interface{}) *gnmi.TypedValue {
			// Make ifIndex values at most uint32.
			return uintval(downcastUint(x, math.MaxUint32))
		})
	interfaceInOctets64 = ifTableMapperFn(interfaceCounterPath+"in-octets",
		"ifHCInOctets", uintval)
	interfaceInOctets32 = ifTableMapperFn(interfaceCounterPath+"in-octets",
		"ifInOctets", uintval)
	interfaceInUnicastPkts64 = ifTableMapperFn(interfaceCounterPath+"in-unicast-pkts",
		"ifHCInUcastPkts", uintval)
	interfaceInUnicastPkts32 = ifTableMapperFn(interfaceCounterPath+"in-unicast-pkts",
		"ifInUcastPkts", uintval)
	interfaceInMulticastPkts = ifTableMapperFn(interfaceCounterPath+"in-multicast-pkts",
		"ifHCInMulticastPkts", uintval)
	interfaceInBroadcastPkts = ifTableMapperFn(interfaceCounterPath+"in-broadcast-pkts",
		"ifHCInBroadcastPkts", uintval)
	interfaceOutMulticastPkts = ifTableMapperFn(interfaceCounterPath+"out-multicast-pkts",
		"ifHCOutMulticastPkts", uintval)
	interfaceOutBroadcastPkts = ifTableMapperFn(interfaceCounterPath+"out-broadcast-pkts",
		"ifHCOutBroadcastPkts", uintval)
	interfaceInDiscards = ifTableMapperFn(interfaceCounterPath+"in-discards",
		"ifInDiscards", uintval)
	interfaceInErrors = ifTableMapperFn(interfaceCounterPath+"in-errors",
		"ifInErrors", uintval)
	interfaceInUnknownProtos = ifTableMapperFn(interfaceCounterPath+"in-unknown-protos",
		"ifInUnknownProtos", uintval)
	interfaceOutOctets64 = ifTableMapperFn(interfaceCounterPath+"out-octets",
		"ifHCOutOctets", uintval)
	interfaceOutOctets32 = ifTableMapperFn(interfaceCounterPath+"out-octets",
		"ifOutOctets", uintval)
	interfaceOutUnicastPkts64 = ifTableMapperFn(interfaceCounterPath+"out-unicast-pkts",
		"ifHCOutUcastPkts", uintval)
	interfaceOutUnicastPkts32 = ifTableMapperFn(interfaceCounterPath+"out-unicast-pkts",
		"ifOutUcastPkts", uintval)
	interfaceOutDiscards = ifTableMapperFn(interfaceCounterPath+"out-discards",
		"ifOutDiscards", uintval)
	interfaceOutErrors = ifTableMapperFn(interfaceCounterPath+"out-errors",
		"ifOutErrors", uintval)
	interfaceHighSpeed = ifTableMapperFn(interfaceEthernetStatePath+"port-speed",
		"ifHighSpeed", ifHighSpeedStrVal)
	interfaceSpeed = ifTableMapperFn(interfaceEthernetStatePath+"port-speed",
		"ifSpeed", ifSpeedStrVal)
	interfacePhysAddress = ifTableMapperFn(interfaceEthernetStatePath+"mac-address",
		"ifPhysAddress", macAddrStrVal)
	interfaceSubIntfIPv4 = ifSubIntfIPMapperFn(
		interfaceSubIntfPath+"ipv4/addresses/address[ip=%s]/ip", "ipAddressIfIndex", strval)
	interfaceSubIntfIPv6 = ifSubIntfIPMapperFn(
		interfaceSubIntfPath+"ipv6/addresses/address[ip=%s]/ip", "ipAddressIfIndex", strval)

	// /system/state
	systemStatePath     = "/system/state/"
	systemStateHostname = scalarMapperFn(systemStatePath+"hostname",
		"sysName.0", processSysName)
	systemStateHostnameLldp = scalarMapperFn(systemStatePath+"hostname",
		"lldpLocSysName.0", processSysName)
	systemStateDomainName = scalarMapperFn(systemStatePath+"domain-name",
		"sysName.0", processDomainName)
	systemStateBootTime64 = scalarMapperFn(systemStatePath+"boot-time",
		"hrSystemUptime.0", processBootTime)
	systemStateBootTime32 = scalarMapperFn(systemStatePath+"boot-time",
		"sysUpTimeInstance", processBootTime)

	// /lldp
	lldpPath                        = "/lldp/"
	lldpStatePath                   = lldpPath + "state/"
	lldpInterfacePath               = lldpPath + "interfaces/interface[name=%s]/"
	lldpInterfaceStatePath          = lldpInterfacePath + "state/"
	lldpInterfaceCountersPath       = lldpInterfaceStatePath + "counters/"
	lldpInterfaceConfigPath         = lldpInterfacePath + "config/"
	lldpInterfaceNeighborsPath      = lldpInterfacePath + "neighbors/neighbor[id=%s]/"
	lldpInterfaceNeighborsStatePath = lldpInterfaceNeighborsPath + "state/"
	lldpChassisID                   = lldpChassisIDFn(lldpStatePath+"chassis-id",
		"lldpLocChassisId", "lldpLocChassisIdSubtype", strval)
	lldpV2ChassisID = lldpChassisIDFn(lldpStatePath+"chassis-id",
		"lldpV2LocChassisId", "lldpV2LocChassisIdSubtype", strval)
	lldpChassisIDType = scalarMapperFn(lldpStatePath+"chassis-id-type",
		"lldpLocChassisIdSubtype", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.LLDPChassisIDType(x.(int)))
		})
	lldpV2ChassisIDType = scalarMapperFn(lldpStatePath+"chassis-id-type",
		"lldpV2LocChassisIdSubtype", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.LLDPChassisIDType(x.(int)))
		})
	lldpSystemName = scalarMapperFn(lldpStatePath+"system-name",
		"lldpLocSysName", strval)
	lldpV2SystemName = scalarMapperFn(lldpStatePath+"system-name",
		"lldpV2LocSysName", strval)
	lldpSystemDescription = scalarMapperFn(lldpStatePath+"system-description",
		"lldpLocSysDesc", strval)
	lldpV2SystemDescription = scalarMapperFn(lldpStatePath+"system-description",
		"lldpV2LocSysDesc", strval)
	lldpInterfaceName = lldpLocPortTableMapperFn(lldpInterfacePath+"name",
		"lldpLocPortId", strval)
	lldpV2InterfaceName = lldpLocPortTableMapperFn(lldpInterfacePath+"name",
		"lldpV2LocPortId", strval)
	lldpInterfaceStateName = lldpLocPortTableMapperFn(lldpInterfaceStatePath+"name",
		"lldpLocPortId", strval)
	lldpV2InterfaceStateName = lldpLocPortTableMapperFn(lldpInterfaceStatePath+"name",
		"lldpV2LocPortId", strval)
	lldpInterfaceConfigName = lldpLocPortTableMapperFn(lldpInterfaceConfigPath+"name",
		"lldpLocPortId", strval)
	lldpV2InterfaceConfigName = lldpLocPortTableMapperFn(lldpInterfaceConfigPath+"name",
		"lldpV2LocPortId", strval)
	lldpInterfaceFrameOut = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"frame-out",
		"lldpStatsTxPortFramesTotal", uintval)
	lldpV2InterfaceFrameOut = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"frame-out",
		"lldpV2StatsTxPortFramesTotal", uintval)
	lldpInterfaceFrameDiscard = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"frame-discard",
		"lldpStatsRxPortFramesDiscardedTotal", uintval)
	lldpV2InterfaceFrameDiscard = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+
		"frame-discard", "lldpV2StatsRxPortFramesDiscardedTotal", uintval)
	lldpInterfaceFrameErrorIn = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+
		"frame-error-in", "lldpStatsRxPortFramesErrors", uintval)
	lldpV2InterfaceFrameErrorIn = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+
		"frame-error-in", "lldpV2StatsRxPortFramesErrors", uintval)
	lldpInterfaceFrameIn = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"frame-in",
		"lldpStatsRxPortFramesTotal", uintval)
	lldpV2InterfaceFrameIn = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"frame-in",
		"lldpV2StatsRxPortFramesTotal", uintval)
	lldpInterfaceTlvDiscard = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"tlv-discard",
		"lldpStatsRxPortTLVsDiscardedTotal", uintval)
	lldpV2InterfaceTlvDiscard = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"tlv-discard",
		"lldpV2StatsRxPortTLVsDiscardedTotal", uintval)
	lldpInterfaceTlvUnknown = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"tlv-unknown",
		"lldpStatsRxPortTLVsUnrecognizedTotal", uintval)
	lldpV2InterfaceTlvUnknown = lldpLocPortTableMapperFn(lldpInterfaceCountersPath+"tlv-unknown",
		"lldpV2StatsRxPortTLVsUnrecognizedTotal", uintval)
	lldpInterfaceNeighborPortID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+"port-id",
		"lldpRemPortId", strval)
	lldpV2InterfaceNeighborPortID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+"port-id",
		"lldpV2RemPortId", strval)
	lldpInterfaceNeighborID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+"id",
		"lldpRemPortId", strval)
	lldpV2InterfaceNeighborID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+"id",
		"lldpV2RemPortId", strval)
	lldpInterfaceNeighborPortIDType = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"port-id-type", "lldpRemPortIdSubtype", func(x interface{}) *gnmi.TypedValue {
		return strval(openconfig.LLDPPortIDType(x.(int)))
	})
	lldpV2InterfaceNeighborPortIDType = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"port-id-type", "lldpV2RemPortIdSubtype", func(x interface{}) *gnmi.TypedValue {
		return strval(openconfig.LLDPPortIDType(x.(int)))
	})
	lldpInterfaceNeighborChassisID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"chassis-id", "lldpRemChassisId", strval)
	lldpV2InterfaceNeighborChassisID = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"chassis-id", "lldpV2RemChassisId", strval)
	lldpInterfaceNeighborChassisIDType = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"chassis-id-type", "lldpRemChassisIdSubtype", func(x interface{}) *gnmi.TypedValue {
		return strval(openconfig.LLDPChassisIDType(x.(int)))
	})
	lldpV2InterfaceNeighborChassisIDType = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"chassis-id-type", "lldpV2RemChassisIdSubtype", func(x interface{}) *gnmi.TypedValue {
		return strval(openconfig.LLDPChassisIDType(x.(int)))
	})
	lldpInterfaceNeighborSystemName = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"system-name", "lldpRemSysName", strval)
	lldpV2InterfaceNeighborSystemName = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"system-name", "lldpV2RemSysName", strval)
	lldpInterfaceNeighborSystemDescription = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"system-description", "lldpRemSysDesc", strval)
	lldpV2InterfaceNeighborSystemDescription = lldpRemTableMapperFn(lldpInterfaceNeighborsStatePath+
		"system-description", "lldpV2RemSysDesc", strval)

	// /components
	componentPath       = "/components/component[name=%s]/"
	componentStatePath  = componentPath + "state/"
	componentConfigPath = componentPath + "config/"
	componentName       = entPhysicalTableMapperFn(componentPath+"name",
		"entPhysicalDescr", strval)
	componentStateName = entPhysicalTableMapperFn(componentStatePath+"name",
		"entPhysicalDescr", strval)
	componentConfigName = entPhysicalTableMapperFn(componentConfigPath+"name",
		"entPhysicalDescr", strval)
	componentID = entPhysicalTableMapperFn(componentStatePath+"id",
		"entPhysicalDescr", strval)
	componentType = entPhysicalTableMapperFn(componentStatePath+"type",
		"entPhysicalClass", strval)
	componentDescription = entPhysicalTableMapperFn(componentStatePath+"description",
		"entPhysicalDescr", strval)
	componentMfgName = entPhysicalTableMapperFn(componentStatePath+"mfg-name",
		"entPhysicalMfgName", strval)
	componentSerialNo = entPhysicalTableMapperFn(componentStatePath+"serial-no",
		"entPhysicalSerialNum", strval)
	componentSoftwareVersion = entPhysicalTableMapperFn(componentStatePath+"software-version",
		"entPhysicalSoftwareRev", strval)
	componentModelName = entPhysicalTableMapperFn(componentStatePath+"part-no",
		"entPhysicalModelName", strval)
	componentHardwareVersion = entPhysicalTableMapperFn(componentStatePath+"hardware-version",
		"entPhysicalHardwareRev", strval)
)

var defaultMappings = map[string][]Mapper{
	// interface
	"/interfaces/interface[name=name]/name":               {interfaceName},
	"/interfaces/interface[name=name]/state/name":         {interfaceStateName},
	"/interfaces/interface[name=name]/config/name":        {interfaceConfigName},
	"/interfaces/interface[name=name]/state/mtu":          {interfaceMtu},
	"/interfaces/interface[name=name]/state/type":         {interfaceType},
	"/interfaces/interface[name=name]/state/admin-status": {interfaceAdminStatus},
	"/interfaces/interface[name=name]/state/oper-status":  {interfaceOperStatus},
	"/interfaces/interface[name=name]/state/ifindex":      {interfaceIndex},
	"/interfaces/interface[name=name]/state/in-octets": {interfaceInOctets64,
		interfaceInOctets32},
	"/interfaces/interface[name=name]/state/in-unicast-pkts": {interfaceInUnicastPkts64,
		interfaceInUnicastPkts32},
	"/interfaces/interface[name=name]/state/in-multicast-pkts":  {interfaceInMulticastPkts},
	"/interfaces/interface[name=name]/state/in-broadcast-pkts":  {interfaceInBroadcastPkts},
	"/interfaces/interface[name=name]/state/out-multicast-pkts": {interfaceOutMulticastPkts},
	"/interfaces/interface[name=name]/state/out-broadcast-pkts": {interfaceOutBroadcastPkts},
	"/interfaces/interface[name=name]/state/in-discards":        {interfaceInDiscards},
	"/interfaces/interface[name=name]/state/in-errors":          {interfaceInErrors},
	"/interfaces/interface[name=name]/state/in-unknown-protos":  {interfaceInUnknownProtos},
	"/interfaces/interface[name=name]/state/out-octets": {interfaceOutOctets64,
		interfaceOutOctets32},
	"/interfaces/interface[name=name]/state/out-unicast-pkts": {interfaceOutUnicastPkts64,
		interfaceOutUnicastPkts32},
	"/interfaces/interface[name=name]/state/out-discards":         {interfaceOutDiscards},
	"/interfaces/interface[name=name]/state/out-errors":           {interfaceOutErrors},
	"/interfaces/interface[name=name]/ethernet/state/mac-address": {interfacePhysAddress},
	"/interfaces/interface[name=name]/ethernet/state/port-speed": {interfaceHighSpeed,
		interfaceSpeed},
	"/interfaces/interface[name=name]/subinterfaces/subinterface[index=index]/ipv4/addresses" +
		"/address/[ip=ip]/ip": {interfaceSubIntfIPv4},
	"/interfaces/interface[name=name]/subinterfaces/subinterface[index=index]/ipv6/addresses" +
		"/address/[ip=ip]/ip": {interfaceSubIntfIPv6},

	// system
	"/system/state/hostname":    {systemStateHostname, systemStateHostnameLldp},
	"/system/state/domain-name": {systemStateDomainName},
	"/system/state/boot-time":   {systemStateBootTime64, systemStateBootTime32},

	//// platform
	"/components/component[name=name]/name":                   {componentName},
	"/components/component[name=name]/state/name":             {componentStateName},
	"/components/component[name=name]/config/name":            {componentConfigName},
	"/components/component[name=name]/state/id":               {componentID},
	"/components/component[name=name]/state/type":             {componentType},
	"/components/component[name=name]/state/description":      {componentDescription},
	"/components/component[name=name]/state/mfg-name":         {componentMfgName},
	"/components/component[name=name]/state/serial-no":        {componentSerialNo},
	"/components/component[name=name]/state/software-version": {componentSoftwareVersion},
	"/components/component[name=name]/state/part-no":          {componentModelName},
	"/components/component[name=name]/state/hardware-version": {componentHardwareVersion},

	//// lldp
	"/lldp/state/chassis-id":                     {lldpChassisID, lldpV2ChassisID},
	"/lldp/state/chassis-id-type":                {lldpChassisIDType, lldpV2ChassisIDType},
	"/lldp/state/system-name":                    {lldpSystemName, lldpV2SystemName},
	"/lldp/state/system-description":             {lldpSystemDescription, lldpV2SystemDescription},
	"/lldp/interfaces/interface[name=name]/name": {lldpInterfaceName, lldpV2InterfaceName},
	"/lldp/interfaces/interface[name=name]/config/" +
		"name": {lldpInterfaceConfigName, lldpV2InterfaceConfigName},
	"/lldp/interfaces/interface[name=name]/state/" +
		"name": {lldpInterfaceStateName, lldpV2InterfaceStateName},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/frame-out": {lldpInterfaceFrameOut, lldpV2InterfaceFrameOut},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/frame-discard": {lldpInterfaceFrameDiscard, lldpV2InterfaceFrameDiscard},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/frame-error-in": {lldpInterfaceFrameErrorIn, lldpV2InterfaceFrameErrorIn},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/frame-in": {lldpInterfaceFrameIn, lldpV2InterfaceFrameIn},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/tlv-discard": {lldpInterfaceTlvDiscard, lldpV2InterfaceTlvDiscard},
	"/lldp/interfaces/interface[name=name]/state/" +
		"counters/tlv-unknown": {lldpInterfaceTlvUnknown, lldpV2InterfaceTlvUnknown},
	"/lldp/interfaces/interface[name=name]/neighbors/neighbor[id=id]/state/id": {
		lldpInterfaceNeighborID, lldpV2InterfaceNeighborID},
	"/lldp/interfaces/interface[name=name]/neighbors/" +
		"neighbor[id=id]/state/port-id": {lldpInterfaceNeighborPortID,
		lldpV2InterfaceNeighborPortID},
	"/lldp/interfaces/interface[name=name]/neighbors/" +
		"neighbor[id=id]/state/port-id-type": {lldpInterfaceNeighborPortIDType,
		lldpV2InterfaceNeighborPortIDType},
	"/lldp/interfaces/interface[name=name]/neighbors/" +
		"neighbor[id=id]/state/chassis-id": {lldpInterfaceNeighborChassisID,
		lldpV2InterfaceNeighborChassisID},
	"/lldp/interfaces/interface[name=name]/neighbors/" +
		"neighbor[id=id]/state/chassis-id-type": {lldpInterfaceNeighborChassisIDType,
		lldpV2InterfaceNeighborChassisIDType},
	"/lldp/interfaces/interface[name=name]/neighbors/neighbor[id=id]/state/" +
		"system-name": {lldpInterfaceNeighborSystemName,
		lldpV2InterfaceNeighborSystemName},
	"/lldp/interfaces/interface[name=name]/neighbors/neighbor[id=id]/state/" +
		"system-description": {lldpInterfaceNeighborSystemDescription,
		lldpV2InterfaceNeighborSystemDescription},
}

// DefaultMappings returns an ordered set of mappings per supported path.
func DefaultMappings() map[string][]Mapper {
	dm := make(map[string][]Mapper)
	for k, v := range defaultMappings {
		dm[k] = v
	}
	return dm
}
