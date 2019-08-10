package smi

import (
	"reflect"
	"testing"
)

type parserTestCase struct {
	name           string
	oid            string
	expectedObject *Object
}

func checkEqual(t *testing.T, o, exp *Object) {
	if (o == nil) != (exp == nil) {
		t.Fatalf("Got one nil")
	}
	if exp == nil {
		return
	}
	if o.Access != exp.Access {
		t.Fatalf("Expected Access: %s. Got: %s", exp.Access.String(),
			o.Access.String())
	}
	if o.Description != exp.Description {
		t.Fatalf("Expected Description: %s. Got: %s", exp.Description,
			o.Description)
	}
	if !reflect.DeepEqual(o.Indexes, exp.Indexes) {
		t.Fatalf("Expected Indexes: %v. Got: %v", exp.Indexes, o.Indexes)
	}
	if o.Kind != exp.Kind {
		t.Fatalf("Expected Kind: %v. Got: %v", exp.Kind, o.Kind)
	}
	if o.Name != exp.Name {
		t.Fatalf("Expected Name: %s. Got %s", exp.Name, o.Name)
	}
	if o.Oid != exp.Oid {
		t.Fatalf("Expected OID: %s. Got %s", exp.Oid, o.Oid)
	}
	if o.Status != exp.Status {
		t.Fatalf("Expected Status: %s. Got %s", exp.Status, o.Status)
	}
	if o.Module != exp.Module {
		t.Fatalf("Expected Module: %s. Got %s", exp.Module, o.Module)
	}
}

func runParserTest(t *testing.T, store Store, tc parserTestCase) {
	o := store.GetObject(tc.oid)
	checkEqual(t, o, tc.expectedObject)
}

func TestParser(t *testing.T) {
	store, err := NewStore("mibs")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []parserTestCase{
		{
			name:           "nonsense",
			oid:            "1.2.3.4.5.6",
			expectedObject: nil,
		},
		{
			name: "interfaces",
			oid:  "interfaces",
			expectedObject: &Object{
				Kind:   KindObject,
				Module: "IF-MIB",
				Name:   "interfaces",
				Oid:    "1.3.6.1.2.1.2",
			},
		},
		{
			name: "interfaces numeric",
			oid:  "1.3.6.1.2.1.2",
			expectedObject: &Object{
				Kind:   KindObject,
				Module: "IF-MIB",
				Name:   "interfaces",
				Oid:    "1.3.6.1.2.1.2",
			},
		},
		{
			name: "interfaces numeric with leading '.'",
			oid:  ".1.3.6.1.2.1.2",
			expectedObject: &Object{
				Kind:   KindObject,
				Module: "IF-MIB",
				Name:   "interfaces",
				Oid:    "1.3.6.1.2.1.2",
			},
		},
		{
			name: "IF-MIB::interfaces",
			oid:  "IF-MIB::interfaces",
			expectedObject: &Object{
				Kind:   KindObject,
				Module: "IF-MIB",
				Name:   "interfaces",
				Oid:    "1.3.6.1.2.1.2",
			},
		},
		{
			name: "ifTable",
			oid:  "ifTable",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "A list of interface entries. The number " +
					"of entries is given by the value of ifNumber.",
				Kind:   KindTable,
				Module: "IF-MIB",
				Name:   "ifTable",
				Oid:    "1.3.6.1.2.1.2.2",
				Status: StatusCurrent,
			},
		},
		{
			name: "ifEntry",
			oid:  "ifEntry",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "An entry containing management information " +
					"applicable to a particular interface.",
				Indexes: []string{"ifIndex"},
				Kind:    KindRow,
				Module:  "IF-MIB",
				Name:    "ifEntry",
				Oid:     "1.3.6.1.2.1.2.2.1",
				Status:  StatusCurrent,
			},
		},
		{
			name: "ifIndex",
			oid:  "ifIndex",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "A unique value, greater than zero, for each " +
					"interface. It is recommended that values are assigned " +
					"contiguously starting from 1. The value for each " +
					"interface sub-layer must remain constant at least " +
					"from one re-initialization of the entity's network " +
					"management system to the next re- initialization.",
				Kind:   KindColumn,
				Module: "IF-MIB",
				Name:   "ifIndex",
				Oid:    "1.3.6.1.2.1.2.2.1.1",
				Status: StatusCurrent,
			},
		},
		{
			name: "ifAdminStatus",
			oid:  "ifAdminStatus",
			expectedObject: &Object{
				Access: AccessReadWrite,
				Description: "The desired state of the interface. The " +
					"testing(3) state indicates that no operational " +
					"packets can be passed. When a managed system " +
					"initializes, all interfaces start with ifAdminStatus " +
					"in the down(2) state. As a result of either explicit " +
					"management action or per configuration information " +
					"retained by the managed system, ifAdminStatus is " +
					"then changed to either the up(1) or testing(3) " +
					"states (or remains in the down(2) state).",
				Kind:   KindColumn,
				Module: "IF-MIB",
				Name:   "ifAdminStatus",
				Oid:    "1.3.6.1.2.1.2.2.1.7",
				Status: StatusCurrent,
			},
		},
		{
			name: "ifInOctets",
			oid:  "ifInOctets",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The total number of octets received on the " +
					"interface, including framing characters. " +
					"Discontinuities in the value of this counter can " +
					"occur at re-initialization of the management system, " +
					"and at other times as indicated by the value of " +
					"ifCounterDiscontinuityTime.",
				Kind:   KindColumn,
				Module: "IF-MIB",
				Name:   "ifInOctets",
				Oid:    "1.3.6.1.2.1.2.2.1.10",
				Status: StatusCurrent,
			},
		},
		{
			name: "ifXTable",
			oid:  "ifXTable",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "A list of interface entries. The number of " +
					"entries is given by the value of ifNumber. This " +
					"table contains additional objects for the interface " +
					"table.",
				Kind:   KindTable,
				Module: "IF-MIB",
				Name:   "ifXTable",
				Oid:    "1.3.6.1.2.1.31.1.1",
				Status: StatusCurrent,
			},
		},
		{
			name: "ifXEntry",
			oid:  "ifXEntry",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "An entry containing additional management " +
					"information applicable to a particular interface.",
				Indexes: []string{"ifIndex"},
				Kind:    KindRow,
				Module:  "IF-MIB",
				Name:    "ifXEntry",
				Oid:     "1.3.6.1.2.1.31.1.1.1",
				Status:  StatusCurrent,
			},
		},
		{
			name: "ifHCInOctets",
			oid:  "1.3.6.1.2.1.31.1.1.1.6",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The total number of octets received on the " +
					"interface, including framing characters. This object " +
					"is a 64-bit version of ifInOctets. Discontinuities " +
					"in the value of this counter can occur at re-" +
					"initialization of the management system, and at " +
					"other times as indicated by the value of " +
					"ifCounterDiscontinuityTime.",
				Kind:   KindColumn,
				Module: "IF-MIB",
				Name:   "ifHCInOctets",
				Oid:    "1.3.6.1.2.1.31.1.1.1.6",
				Status: StatusCurrent,
			},
		},
		{
			name: "LLDP-MIB::lldpLocChassisId",
			oid:  "LLDP-MIB::lldpLocChassisId",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The string value used to identify the " +
					"chassis component associated with the local system.",
				Kind:   KindScalar,
				Module: "LLDP-MIB",
				Name:   "lldpLocChassisId",
				Oid:    "1.0.8802.1.1.2.1.3.2",
				Status: StatusCurrent,
			},
		},
		{
			name: "lldpLocChassisId.0",
			oid:  "1.0.8802.1.1.2.1.3.2.0",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The string value used to identify the " +
					"chassis component associated with the local system.",
				Kind:   KindScalar,
				Module: "LLDP-MIB",
				Name:   "lldpLocChassisId",
				Oid:    "1.0.8802.1.1.2.1.3.2",
				Status: StatusCurrent,
			},
		},
		{
			name: "lldpRemEntry",
			oid:  "lldpRemEntry",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "Information about a particular physical " +
					"network connection. Entries may be created and " +
					"deleted in this table by the agent, if a physical " +
					"topology discovery process is active.",
				Indexes: []string{"lldpRemTimeMark",
					"lldpRemLocalPortNum", "lldpRemIndex"},
				Kind:   KindRow,
				Module: "LLDP-MIB",
				Name:   "lldpRemEntry",
				Oid:    "1.0.8802.1.1.2.1.4.1.1",
				Status: StatusCurrent,
			},
		},
		{
			name: "lldpRemChassisId",
			oid:  "lldpRemChassisId",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The string value used to identify the " +
					"chassis component associated with the remote system.",
				Kind:   KindColumn,
				Module: "LLDP-MIB",
				Name:   "lldpRemChassisId",
				Oid:    "1.0.8802.1.1.2.1.4.1.1.5",
				Status: StatusCurrent,
			},
		},
		{
			name: "lldpRemChassisId with indexes",
			oid:  "1.0.8802.1.1.2.1.4.1.1.5.1.2.3",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The string value used to identify the " +
					"chassis component associated with the remote system.",
				Kind:   KindColumn,
				Module: "LLDP-MIB",
				Name:   "lldpRemChassisId",
				Oid:    "1.0.8802.1.1.2.1.4.1.1.5",
				Status: StatusCurrent,
			},
		},
		{
			name:           "lldpRemChassisId with too many indexes",
			oid:            "1.0.8802.1.1.2.1.4.1.1.5.1.2.3.4",
			expectedObject: nil,
		},
		{
			name:           "lldpRemChassisId with too few indexes",
			oid:            "1.0.8802.1.1.2.1.4.1.1.5.1.2",
			expectedObject: nil,
		},
		{
			name: "entPhysicalTable",
			oid:  "ENTITY-MIB::entPhysicalTable",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "This table contains one row per physical " +
					"entity. There is always at least one row for an " +
					"'overall' physical entity.",
				Kind:   KindTable,
				Module: "ENTITY-MIB",
				Name:   "entPhysicalTable",
				Oid:    "1.3.6.1.2.1.47.1.1.1",
				Status: StatusCurrent,
			},
		},
		{
			name: "entPhysicalEntry",
			oid:  "entPhysicalEntry",
			expectedObject: &Object{
				Access: AccessNotAccessible,
				Description: "Information about a particular physical " +
					"entity. Each entry provides objects " +
					"(entPhysicalDescr, entPhysicalVendorType, and " +
					"entPhysicalClass) to help an NMS identify and " +
					"characterize the entry and objects " +
					"(entPhysicalContainedIn and entPhysicalParentRelPos) " +
					"to help an NMS relate the particular entry to other " +
					"entries in this table.",
				Indexes: []string{"entPhysicalIndex"},
				Kind:    KindRow,
				Module:  "ENTITY-MIB",
				Name:    "entPhysicalEntry",
				Oid:     "1.3.6.1.2.1.47.1.1.1.1",
				Status:  StatusCurrent,
			},
		},
		{
			name: "hrSystemUptime",
			oid:  "hrSystemUptime",
			expectedObject: &Object{
				Access: AccessReadOnly,
				Description: "The amount of time since this host was last " +
					"initialized. Note that this is different from " +
					"sysUpTime in the SNMPv2-MIB [RFC1907] because " +
					"sysUpTime is the uptime of the network management " +
					"portion of the system.",
				Kind:   KindScalar,
				Module: "HOST-RESOURCES-MIB",
				Name:   "hrSystemUptime",
				Oid:    "1.3.6.1.2.1.25.1.1",
				Status: StatusCurrent,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runParserTest(t, store, tc)
		})
	}
}
