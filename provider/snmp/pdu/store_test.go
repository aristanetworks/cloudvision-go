package pdu

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/soniah/gosnmp"
)

func pdusMatch(pdu, expected *gosnmp.SnmpPDU) bool {
	if pdu.Name != expected.Name {
		return false
	}
	if pdu.Type != expected.Type {
		return false
	}
	if pdu.Value != expected.Value {
		return false
	}
	return true
}

func checkPDUs(t *testing.T, pdus, expected []*gosnmp.SnmpPDU) {
	if len(pdus) != len(expected) {
		t.Fatalf("got %d PDUs, expected %d", len(pdus), len(expected))
	}
	received := make(map[string]*gosnmp.SnmpPDU)
	for _, p := range pdus {
		received[p.Name] = p
	}
	for _, e := range expected {
		r, ok := received[e.Name]
		if !ok {
			t.Fatalf("Did not receive expected PDU: %v", e)
		}
		if !pdusMatch(r, e) {
			t.Fatalf("expected %v, got %v", e, r)
		}
	}
}

func pdu(name string, t gosnmp.Asn1BER, val interface{}) *gosnmp.SnmpPDU {
	if s, ok := val.(string); ok {
		val = s
	}
	return &gosnmp.SnmpPDU{
		Name:  name,
		Type:  t,
		Value: val,
	}
}

type testGet struct {
	oid          string
	scalar       bool
	constraints  []Index
	expectedPDUs []*gosnmp.SnmpPDU
	err          error
}

type storeTestCase struct {
	name  string
	clear bool
	adds  []*gosnmp.SnmpPDU
	get   testGet
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

func runStoreTest(t *testing.T, s Store, tc storeTestCase) {
	if tc.clear {
		if err := s.Clear(); err != nil {
			t.Fatalf("Error in Clear: %s", err)
		}
	}

	for _, pdu := range tc.adds {
		if err := s.Add(pdu); err != nil {
			t.Fatalf("Error in Add: %s", err)
		}
	}

	if tc.get.scalar {
		pdu, err := s.GetScalar(tc.get.oid)
		checkError(t, err, tc.get.err)

		if len(tc.get.expectedPDUs) > 1 {
			t.Fatalf("GetScalar should have at most one expectedPDU, got %d",
				len(tc.get.expectedPDUs))
		}
		var pdus []*gosnmp.SnmpPDU
		if pdu != nil {
			pdus = []*gosnmp.SnmpPDU{pdu}
		}
		checkPDUs(t, pdus, tc.get.expectedPDUs)
		return
	}
	pdus, err := s.GetTabular(tc.get.oid, tc.get.constraints...)
	checkError(t, err, tc.get.err)
	checkPDUs(t, pdus, tc.get.expectedPDUs)
}

const (
	sysNameOid        = "1.3.6.1.2.1.1.5.0"
	ifDescrOid        = "1.3.6.1.2.1.2.2.1.2"
	lldpRemSysNameOid = "1.0.8802.1.1.2.1.4.1.1.9"
)

func ifDescrPDU(i, d string) *gosnmp.SnmpPDU {
	return pdu("1.3.6.1.2.1.2.2.1.2."+i, gosnmp.OctetString, d)
}

func sysNamePDU(n string) *gosnmp.SnmpPDU {
	return pdu("1.3.6.1.2.1.1.5.0", gosnmp.OctetString, n)
}

func lldpRemSysNamePDU(l, r, name string) *gosnmp.SnmpPDU {
	return pdu(fmt.Sprintf("1.0.8802.1.1.2.1.4.1.1.9.42.%s.%s", l, r),
		gosnmp.OctetString, name)
}

func TestStore(t *testing.T) {
	mibStore, err := smi.NewStore("../smi/mibs")
	if err != nil {
		t.Fatalf("Error creating smi.Store: %s", err)
	}

	store, err := NewStore(mibStore)
	if err != nil {
		t.Fatalf("Error creating store: %s", err)
	}

	for _, tc := range []storeTestCase{
		{
			name: "basic scalar",
			adds: []*gosnmp.SnmpPDU{sysNamePDU("device123.x.y.z")},
			get: testGet{
				oid:    sysNameOid,
				scalar: true,
				expectedPDUs: []*gosnmp.SnmpPDU{
					sysNamePDU("device123.x.y.z"),
				},
			},
		},
		{
			name:  "basic tabular",
			clear: true,
			adds: []*gosnmp.SnmpPDU{
				ifDescrPDU("1", "intf1"),
				ifDescrPDU("2", "intf2"),
			},
			get: testGet{
				oid: ifDescrOid,
				expectedPDUs: []*gosnmp.SnmpPDU{
					ifDescrPDU("1", "intf1"),
					ifDescrPDU("2", "intf2"),
				},
			},
		},
		{
			name: "text OID",
			get: testGet{
				oid: "ifDescr",
				expectedPDUs: []*gosnmp.SnmpPDU{
					ifDescrPDU("1", "intf1"),
					ifDescrPDU("2", "intf2"),
				},
			},
		},
		{
			name: "tabular with constraint",
			get: testGet{
				oid: ifDescrOid,
				constraints: []Index{
					Index{
						Name:  "ifIndex",
						Value: "2",
					},
				},
				expectedPDUs: []*gosnmp.SnmpPDU{ifDescrPDU("2", "intf2")},
			},
		},
		{
			name: "numeric constraint OID",
			get: testGet{
				oid: ifDescrOid,
				constraints: []Index{
					Index{
						Name:  "1.3.6.1.2.1.2.2.1.1",
						Value: "2",
					},
				},
				expectedPDUs: []*gosnmp.SnmpPDU{ifDescrPDU("2", "intf2")},
			},
		},
		{
			name:  "multi-index tabular",
			clear: true,
			adds: []*gosnmp.SnmpPDU{
				lldpRemSysNamePDU("1", "1", "d1-1"),
				lldpRemSysNamePDU("1", "2", "d1-2"),
				lldpRemSysNamePDU("2", "1", "d2-1"),
				lldpRemSysNamePDU("2", "2", "d2-2"),
			},
			get: testGet{
				oid: lldpRemSysNameOid,
				expectedPDUs: []*gosnmp.SnmpPDU{
					lldpRemSysNamePDU("1", "1", "d1-1"),
					lldpRemSysNamePDU("1", "2", "d1-2"),
					lldpRemSysNamePDU("2", "1", "d2-1"),
					lldpRemSysNamePDU("2", "2", "d2-2"),
				},
			},
		},
		{
			name: "multi-index tabular with one constraint",
			get: testGet{
				oid: lldpRemSysNameOid,
				constraints: []Index{
					Index{
						Name:  "lldpRemLocalPortNum",
						Value: "2",
					},
				},
				expectedPDUs: []*gosnmp.SnmpPDU{
					lldpRemSysNamePDU("2", "1", "d2-1"),
					lldpRemSysNamePDU("2", "2", "d2-2"),
				},
			},
		},
		{
			name: "multi-index tabular with two constraints",
			get: testGet{
				oid: lldpRemSysNameOid,
				constraints: []Index{
					Index{
						Name:  "lldpRemLocalPortNum",
						Value: "2",
					},
					Index{
						Name:  "lldpRemIndex",
						Value: "1",
					},
				},
				expectedPDUs: []*gosnmp.SnmpPDU{
					lldpRemSysNamePDU("2", "1", "d2-1"),
				},
			},
		},
		{
			name: "constraint sorting",
			get: testGet{
				oid: lldpRemSysNameOid,
				constraints: []Index{
					Index{
						Name:  "lldpRemIndex",
						Value: "1",
					},
					Index{
						Name:  "lldpRemLocalPortNum",
						Value: "2",
					},
				},
				expectedPDUs: []*gosnmp.SnmpPDU{
					lldpRemSysNamePDU("2", "1", "d2-1"),
				},
			},
		},
		{
			name:  "PDU not added scalar",
			clear: true,
			get: testGet{
				oid:          sysNameOid,
				scalar:       true,
				expectedPDUs: nil,
			},
		},
		{
			name:  "PDU not added tabular",
			clear: true,
			get: testGet{
				oid:          ifDescrOid,
				expectedPDUs: nil,
			},
		},
		{
			name: "too many constraints",
			get: testGet{
				oid: ifDescrOid,
				constraints: []Index{
					Index{
						Name:  "ifIndex",
						Value: "2",
					},
					Index{
						Name:  "bogus",
						Value: "whatever",
					},
				},
				err: errors.New("2 constraints is more than 1 indexes"),
			},
		},
		{
			name: "constraint not found in MIB store",
			get: testGet{
				oid: ifDescrOid,
				constraints: []Index{
					Index{
						Name:  "ifBogus",
						Value: "2",
					},
				},
				err: errors.New("Index 'ifBogus' not found in MIB store"),
			},
		},
		{
			name: "invalid constraint index",
			get: testGet{
				oid: ifDescrOid,
				constraints: []Index{
					Index{
						Name:  "ifAdminStatus",
						Value: "2",
					},
				},
				err: fmt.Errorf("Invalid constraint 'ifAdminStatus' for OID %s", ifDescrOid),
			},
		},
		{
			name: "no such scalar object",
			get: testGet{
				oid: "1.2.3.4.5.0",
				err: errors.New("No corresponding object in MIB store for OID 1.2.3.4.5.0"),
			},
		},
		{
			name: "no such tabular object",
			get: testGet{
				oid: "1.2.3.4.5",
				err: errors.New("No corresponding object in MIB store for OID 1.2.3.4.5"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runStoreTest(t, store, tc)
		})
	}
}
