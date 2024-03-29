// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmpoc

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
)

// This file contains functionality useful for producing a dump from
// a series of SNMP requests and responses, or for reading such a
// dump. The format of responses is what you'd see running
// `snmpwalk -O ne`, so:
//
// <OID> = <type>: <value>
//
// For example:
//
// .1.3.6.1.2.1.47.1.1.1.1.13.156025601 = STRING: Ucd90120
//
// The request format is:
//
// <request-type>: <OID>
//
// For example:
//
// WALK: .1.3.6.1.2.1.47.1.1.1.1

// SNMP PDU types of interest.
const (
	octstr              = gosnmp.OctetString
	counter             = gosnmp.Counter32
	counter64           = gosnmp.Counter64
	integer             = gosnmp.Integer
	gauge32             = gosnmp.Gauge32
	timeticks           = gosnmp.TimeTicks // nolint: deadcode
	octstrTypeString    = "STRING"
	hexstrTypeString    = "Hex-STRING"
	integerTypeString   = "INTEGER"
	counterTypeString   = "Counter32"
	counter64TypeString = "Counter64"
	gauge32TypeString   = "Gauge32"
	getString           = "GET"       // nolint: deadcode
	walkString          = "WALK"      // nolint: deadcode
	timeticksString     = "Timeticks" // nolint: deadcode

	oidPrefixIfPhysAddr = ".1.3.6.1.2.1.2.2.1.6."
)

// A pduParser takes OID, pdu type string, and value parsed from a dumped PDU and returns a
// gosnmp.SnmpPDU
type pduParser func(oid, pduTypeString, val string) *gosnmp.SnmpPDU

// PDU creation wrapper.
func PDU(name string, t gosnmp.Asn1BER, val interface{}) *gosnmp.SnmpPDU {
	return &gosnmp.SnmpPDU{
		Name:  name,
		Type:  t,
		Value: val,
	}
}

// Get OID, type, and value from a dumped PDU.
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
		value = strings.Trim(value, "\"")
	} else {
		pduTypeString = strings.Split(t[0], ":")[0]
	}
	return oid, pduTypeString, value
}

// customPDUParsers maps an oid prefix to the custom pdu parsing logic required for
// getting a gosnmp.SnmpPdu from oids starting with that prefix
var customPDUParsers = map[string]pduParser{
	oidPrefixIfPhysAddr: ifPhysAddrPDUParser,
}

// pduFromString returns a PDU from a string representation of a PDU.
// If it sees something it doesn't like, it returns nil.
func pduFromString(s string) *gosnmp.SnmpPDU {
	oid, pduTypeString, val := parsePDU(s)
	if oid == "" {
		return nil
	}

	// some values can't be generically parsed using the type switch on the `pduTypeString` below,
	// use custom parsers in such cases
	for oidPrefix, customParser := range customPDUParsers {
		if strings.HasPrefix(oid, oidPrefix) {
			return customParser(oid, pduTypeString, val)
		}
	}

	var pduType gosnmp.Asn1BER
	var value interface{}
	switch pduTypeString {
	case integerTypeString:
		pduType = integer
		v, _ := strconv.ParseInt(val, 10, 32)
		value = int(v)
	case octstrTypeString:
		pduType = octstr
		value = []byte(val)
	case hexstrTypeString:
		pduType = octstr
		s := strings.Replace(val, " ", "", -1)
		value, _ = hex.DecodeString(s)
	case counterTypeString:
		pduType = counter
		v, _ := strconv.ParseUint(val, 10, 32)
		value = uint(v)
	case counter64TypeString:
		pduType = counter64
		v, _ := strconv.ParseUint(val, 10, 64)
		value = v
	case gauge32TypeString:
		pduType = gauge32
		v, _ := strconv.ParseUint(val, 10, 32)
		value = v
	default:
		return nil
	}
	return PDU(oid, pduType, value)
}

// PDUsFromString converts a set of formatted SNMP responses (as
// returned by snmpwalk -v3 -O ne <target> <oid>) to PDUs.
func PDUsFromString(s string) []*gosnmp.SnmpPDU {
	pdus := make([]*gosnmp.SnmpPDU, 0)
	for _, line := range strings.Split(s, "\n") {
		pdu := pduFromString(line)
		if pdu != nil {
			pdus = append(pdus, pdu)
		}
	}
	return pdus
}

// ifPhysAddrPDUParser takes OID, pdu type string , and value parsed from a dumped PDU string for
// ifPhysAddress and returns a parsed gosnmp.SnmpPDU
func ifPhysAddrPDUParser(oid, pduTypeString, val string) *gosnmp.SnmpPDU {
	return PDU(oid, octstr, macAddrToByteHWAddr(val))
}

func macAddrToByteHWAddr(mac string) []byte {
	groups := strings.Split(mac, ":")
	result := make([]byte, len(groups))
	for i, group := range groups {
		groupUint, err := strconv.ParseUint(group, 16, 8)
		if err != nil {
			return nil
		}
		result[i] = byte(groupUint)
	}
	return result
}
