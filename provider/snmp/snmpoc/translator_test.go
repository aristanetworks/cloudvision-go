package snmpoc

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/pdu"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

var (
	// OIDs
	sysNameOID        = "1.3.6.1.2.1.1.5.0"
	lldpLocSysNameOID = "1.0.8802.1.1.2.1.3.3.0"
	ifInOctetsOID     = "1.3.6.1.2.1.2.2.1.10"
	ifHCInOctetsOID   = "1.3.6.1.2.1.31.1.1.1.6"
	ifDescrOID        = "1.3.6.1.2.1.2.2.1.2"

	// hostnames, etc.
	hostnameAbcd     = "abcd"
	hostnameAbcdLldp = "abcd-lldp"
	intf1            = "intf1"
	intf2            = "intf2"

	// PDUs
	sysNamePDU        = newPdu(sysNameOID, gosnmp.OctetString, hostnameAbcd)
	lldpLocSysNamePDU = newPdu(lldpLocSysNameOID, gosnmp.OctetString,
		hostnameAbcdLldp)

	// gNMI paths
	hostnamePath = pgnmi.Path("system", "state", "hostname")
	inOctetsPath = pgnmi.Path("interfaces",
		pgnmi.ListWithKey("interface", "name", "*"),
		"state", "counters", "in-octets")
	intfNamePath = pgnmi.Path("interfaces",
		pgnmi.ListWithKey("interface", "name", "*"), "name")
)

func newPdu(name string, t gosnmp.Asn1BER, val interface{}) *gosnmp.SnmpPDU {
	if s, ok := val.(string); ok {
		val = s
	}
	return &gosnmp.SnmpPDU{
		Name:  name,
		Type:  t,
		Value: val,
	}
}

func strval(s interface{}) *gnmi.TypedValue {
	switch t := s.(type) {
	case string:
		if t == "" {
			return nil
		}
		return pgnmi.Strval(t)
	case []byte:
		return pgnmi.Strval(string(t))
	}
	return nil
}

func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToUint64(u); err == nil {
		return pgnmi.Uintval(v)
	}
	return nil
}

func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return pgnmi.Update(path, val)
}

// genric system/state/hostname mapper
func hostname(store pdu.Store, oid string) ([]*gnmi.Update, error) {
	pdu, err := store.GetScalar(oid)
	if err != nil {
		return nil, err
	} else if pdu == nil {
		return nil, nil
	}

	path := pgnmi.Path("system", "state", "hostname")
	return []*gnmi.Update{update(path, strval(pdu.Value))}, nil
}

// primary system/state/hostname mapper
func hostnameSysName(ss smi.Store, ps pdu.Store, kv KVStore) ([]*gnmi.Update, error) {
	return hostname(ps, sysNameOID)
}

// secondary system/state/hostname mapper
func hostnameLldp(ss smi.Store, ps pdu.Store, kv KVStore) ([]*gnmi.Update, error) {
	return hostname(ps, lldpLocSysNameOID)
}

var hostnameErr = errors.New("Test error")

func hostnameError(ss smi.Store, ps pdu.Store, kv KVStore) ([]*gnmi.Update, error) {
	return nil, hostnameErr
}

type typedValue func(interface{}) *gnmi.TypedValue

func intfPath(kvs KVStore, ifIndex string) *gnmi.Path {
	m := kvs.Get("intfPath")
	if m == nil {
		panic("intfPath does not exist in Translator")
	}

	mp := m.(map[string]*gnmi.Path)
	if p, ok := mp[ifIndex]; ok {
		return p
	}
	return nil
}

func setIntfPath(kvs KVStore, ifIndex string, path *gnmi.Path) {
	m := kvs.Get("intfPath")
	if m == nil {
		kvs.Set("intfPath", map[string]*gnmi.Path{})
		m = kvs.Get("intfPath")
	}

	mp := m.(map[string]*gnmi.Path)
	mp[ifIndex] = path
}

// interface/interface[name=name] mapper
func ifDescrMapper(ss smi.Store, kvs KVStore,
	pdus []*gosnmp.SnmpPDU) ([]*gnmi.Update, error) {
	updates := []*gnmi.Update{}
	for _, p := range pdus {
		ifIndex := pdu.IndexValues(ss, p)[0]
		ifDescr := p.Value.(string)
		setIntfPath(kvs, ifIndex, pgnmi.Path("interfaces",
			pgnmi.ListWithKey("interface", "name", ifDescr)))
		updates = append(updates,
			update(pgnmi.IntfPath(ifDescr, "name"), strval(ifDescr)))
	}
	return updates, nil
}

// generic mapper for PDUs from ifTable
func ifTableMapper(ss smi.Store, ps pdu.Store, kvs KVStore, oid string,
	tv typedValue) ([]*gnmi.Update, error) {
	pdus, err := ps.GetTabular(oid)
	if err != nil {
		return nil, err
	} else if len(pdus) == 0 {
		return nil, nil
	}

	// ifDescr is our index, so handle it separately
	o := ss.GetObject(pdus[0].Name)
	if o == nil {
		return nil, fmt.Errorf("GetObject failed on OID %s", oid)
	} else if o.Name == "ifDescr" {
		return ifDescrMapper(ss, kvs, pdus)
	}

	updates := []*gnmi.Update{}
	for _, p := range pdus {
		ifIndex := pdu.IndexValues(ss, p)[0]
		path := intfPath(kvs, ifIndex)
		if path == nil {
			return nil, fmt.Errorf("No path for ifIndex %s", ifIndex)
		}
		fullPath := pgnmi.PathAppend(path, "state", "counters", "in-octets")
		updates = append(updates, update(fullPath, tv(p.Value)))
	}
	return updates, nil
}

// interfaces/interface[name=name] mapper
func intfName(ss smi.Store, ps pdu.Store, kvs KVStore) ([]*gnmi.Update, error) {
	return ifTableMapper(ss, ps, kvs, ifDescrOID, strval)
}

// secondary interfaces/interface[name=name]/state/counters/in-octets mapper
func inOctets32(ss smi.Store, ps pdu.Store, kvs KVStore) ([]*gnmi.Update, error) {
	return ifTableMapper(ss, ps, kvs, ifInOctetsOID, uintval)
}

// primary interfaces/interface[name=name]/state/counters/in-octets mapper
func inOctets64(ss smi.Store, ps pdu.Store, kvs KVStore) ([]*gnmi.Update, error) {
	return ifTableMapper(ss, ps, kvs, ifHCInOctetsOID, uintval)
}

func checkUpdates(t *testing.T, updates, expected []*gnmi.Update) {
	if len(updates) != len(expected) {
		t.Fatalf("Expected %d updates; got %d", len(expected), len(updates))
	}
	for _, expectedUpdate := range expected {
		found := false
		for _, update := range updates {
			if reflect.DeepEqual(update, expectedUpdate) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected update not found: %s", expectedUpdate)
		}
	}
}

func checkError(t *testing.T, err, expectedErr error) {
	if err == nil && expectedErr != nil {
		t.Fatalf("Expected but didn't get error: %s", expectedErr)
	} else if err != nil && expectedErr == nil {
		t.Fatalf("Unexpected error: %s", err)
	} else if err != nil {
		if err.Error() != expectedErr.Error() {
			t.Fatalf("Expected error '%s', got '%s'", expectedErr, err)
		}
	}
}

func runTranslatorTest(t *testing.T, trans Translator, ss smi.Store, ps pdu.Store,
	kvs KVStore, tc *translatorTestCase) {
	// add data to PDU store
	if tc.clearPDUs {
		if err := ps.Clear(); err != nil {
			t.Fatalf("Error in Clear: %s", err)
		}
	}
	for _, p := range tc.pdus {
		if err := ps.Add(p); err != nil {
			t.Fatalf("Error in Add: %s", err)
		}
	}

	// clear/add mappings
	if tc.clear {
		trans.(*translator).mappings = make(map[string][]Mapper)
		trans.(*translator).successfulMappings = make(map[string]Mapper)
	}
	for _, m := range tc.mappings {
		if err := trans.AddMapping(m.path, m.mapper); err != nil {
			t.Fatalf("AddMapping error: %s", err)
		}
	}

	// produce updates and check output
	updates, err := trans.Updates(tc.updatePaths)
	checkError(t, err, tc.expectedErr)
	checkUpdates(t, updates, tc.expectedUpdates)
}

type mapping struct {
	path   *gnmi.Path
	mapper Mapper
}

type translatorTestCase struct {
	name            string
	clearPDUs       bool
	pdus            []*gosnmp.SnmpPDU
	clear           bool
	mappings        []mapping
	updatePaths     []*gnmi.Path
	expectedUpdates []*gnmi.Update
	expectedErr     error
}

func TestTranslator(t *testing.T) {
	mibStore, err := smi.NewStore("../smi/mibs")
	if err != nil {
		t.Fatalf("Error in smi.NewStore: %s", err)
	}

	pduStore, err := pdu.NewStore(mibStore)
	if err != nil {
		t.Fatalf("Error in pdu.NewStore: %s", err)
	}

	kvStore := NewKVStore()
	trans := NewTranslator(pduStore, mibStore, kvStore)
	for _, tc := range []translatorTestCase{
		{
			name: "simple",
			pdus: []*gosnmp.SnmpPDU{
				sysNamePDU,
			},
			mappings: []mapping{
				mapping{
					path:   hostnamePath,
					mapper: hostnameSysName,
				},
			},
			updatePaths: []*gnmi.Path{
				hostnamePath,
			},
			expectedUpdates: []*gnmi.Update{
				update(hostnamePath, strval(hostnameAbcd)),
			},
		},
		{
			name: "no such mapping",
			updatePaths: []*gnmi.Path{
				inOctetsPath,
			},
			expectedErr: fmt.Errorf("No mapping supplied for path %v",
				inOctetsPath),
		},
		{
			name:      "backup mapping",
			clearPDUs: true,
			clear:     true,
			pdus: []*gosnmp.SnmpPDU{
				lldpLocSysNamePDU,
			},
			mappings: []mapping{
				mapping{
					path:   hostnamePath,
					mapper: hostnameSysName,
				},
				mapping{
					path:   hostnamePath,
					mapper: hostnameLldp,
				},
			},
			updatePaths: []*gnmi.Path{
				hostnamePath,
			},
			expectedUpdates: []*gnmi.Update{
				update(hostnamePath, strval(hostnameAbcdLldp)),
			},
		},
		{
			name:      "no updates",
			clearPDUs: true,
			mappings: []mapping{
				mapping{
					path:   hostnamePath,
					mapper: hostnameSysName,
				},
			},
			updatePaths: []*gnmi.Path{
				hostnamePath,
			},
			expectedUpdates: nil,
		},
		{
			name:  "mapper error",
			clear: true,
			mappings: []mapping{
				mapping{
					path:   hostnamePath,
					mapper: hostnameError,
				},
			},
			updatePaths: []*gnmi.Path{
				hostnamePath,
			},
			expectedErr: hostnameErr,
		},
		{
			name: "basic tabular",
			pdus: []*gosnmp.SnmpPDU{
				newPdu(ifDescrOID+".1", gosnmp.OctetString, intf1),
				newPdu(ifDescrOID+".2", gosnmp.OctetString, intf2),
			},
			mappings: []mapping{
				mapping{
					path:   intfNamePath,
					mapper: intfName,
				},
			},
			updatePaths: []*gnmi.Path{
				intfNamePath,
			},
			expectedUpdates: []*gnmi.Update{
				update(pgnmi.IntfPath(intf1, "name"), strval(intf1)),
				update(pgnmi.IntfPath(intf2, "name"), strval(intf2)),
			},
		},
		{
			name: "in-octets",
			pdus: []*gosnmp.SnmpPDU{
				newPdu(ifInOctetsOID+".1", gosnmp.Counter32, 11),
				newPdu(ifInOctetsOID+".2", gosnmp.Counter32, 22),
			},
			mappings: []mapping{
				mapping{
					path:   inOctetsPath,
					mapper: inOctets32,
				},
			},
			updatePaths: []*gnmi.Path{
				inOctetsPath,
			},
			expectedUpdates: []*gnmi.Update{
				update(pgnmi.IntfStateCountersPath(intf1, "in-octets"), uintval(11)),
				update(pgnmi.IntfStateCountersPath(intf2, "in-octets"), uintval(22)),
			},
		},
		{
			name:      "in-octets-64",
			clearPDUs: true,
			clear:     true,
			pdus: []*gosnmp.SnmpPDU{
				newPdu(ifDescrOID+".1", gosnmp.OctetString, intf1),
				newPdu(ifDescrOID+".2", gosnmp.OctetString, intf2),
				newPdu(ifHCInOctetsOID+".1", gosnmp.Counter64, 111),
				newPdu(ifHCInOctetsOID+".2", gosnmp.Counter64, 222),
			},
			mappings: []mapping{
				mapping{
					path:   intfNamePath,
					mapper: intfName,
				},
				mapping{
					path:   inOctetsPath,
					mapper: inOctets64,
				},
				mapping{
					path:   inOctetsPath,
					mapper: inOctets32,
				},
			},
			updatePaths: []*gnmi.Path{
				intfNamePath,
				inOctetsPath,
			},
			expectedUpdates: []*gnmi.Update{
				update(pgnmi.IntfPath(intf1, "name"), strval(intf1)),
				update(pgnmi.IntfPath(intf2, "name"), strval(intf2)),
				update(pgnmi.IntfStateCountersPath(intf1, "in-octets"), uintval(111)),
				update(pgnmi.IntfStateCountersPath(intf2, "in-octets"), uintval(222)),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runTranslatorTest(t, trans, mibStore, pduStore, kvStore, &tc)
		})
	}
}

// - check that once a mapping is successful, we keep using it
// - check that we try all paths supplied
// - non-string type
// - Path with no such key
