package snmpoc

import (
	"fmt"

	"github.com/aristanetworks/cloudvision-go/provider/snmp/pdu"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// KVStore is a simple key-value store. Its keys are strings and its
// values are interface{}.
type KVStore interface {
	Get(string) interface{}
	Set(string, interface{})
}

// A Mapper contains some logic for producing gNMI updates based on the
// contents of a pdu.Store and KVStore.
type Mapper = func(smi.Store, pdu.Store, KVStore) ([]*gnmi.Update, error)

// Translator defines an interface for producing translations from a
// set of received SNMP PDUs to a set of gNMI updates.
type Translator interface {
	// AddMapping associates a mapping with a gNMI path.
	AddMapping(*gnmi.Path, Mapper) error

	// Updates looks for mappings associated with the provided path and
	// tries to produce a set of updates using those mappings.
	Updates([]*gnmi.Path) ([]*gnmi.Update, error)
}

// NewTranslator returns a Translator.
func NewTranslator(pduStore pdu.Store, mibStore smi.Store, kvs KVStore) Translator {
	return &translator{
		pduStore:           pduStore,
		mibStore:           mibStore,
		kvStore:            kvs,
		mappings:           make(map[string][]Mapper),
		successfulMappings: make(map[string]Mapper),
	}
}

func NewKVStore() KVStore {
	return &kvstore{
		data: make(map[string]interface{}),
	}
}

type kvstore struct {
	data map[string]interface{}
}

func (s *kvstore) Set(k string, v interface{}) {
	s.data[k] = v
}

func (s *kvstore) Get(k string) interface{} {
	if v, ok := s.data[k]; ok {
		return v
	}
	return nil
}

// translator implements Translator.
type translator struct {
	pduStore           pdu.Store
	mibStore           smi.Store
	mappings           map[string][]Mapper
	successfulMappings map[string]Mapper
	kvStore            KVStore
}

// AddMapping associates a mapping function with a path. The translator
// will use the provided mapper to produce updates for the given path.
func (t *translator) AddMapping(path *gnmi.Path, m Mapper) error {
	ps := path.String()
	if _, ok := t.mappings[ps]; !ok {
		t.mappings[ps] = []Mapper{}
	}
	t.mappings[ps] = append([]Mapper{m}, t.mappings[ps]...)
	return nil
}

// Updates produces updates for the provided set of paths.
func (t *translator) Updates(paths []*gnmi.Path) ([]*gnmi.Update, error) {
	updates := []*gnmi.Update{}
	for _, path := range paths {
		// If we have a mapping that already worked, use it.
		if mapping, ok := t.successfulMappings[path.String()]; ok {
			u, err := mapping(t.mibStore, t.pduStore, t.kvStore)
			if err != nil {
				return nil, err
			}
			updates = append(updates, u...)
		}

		// Otherwise, try each mapping in order.
		mappings, ok := t.mappings[path.String()]
		if !ok {
			return nil, fmt.Errorf("No mapping supplied for path %v", path)
		}
		for _, mapping := range mappings {
			u, err := mapping(t.mibStore, t.pduStore, t.kvStore)
			if err != nil {
				return nil, err
			} else if len(u) == 0 {
				continue
			}
			t.successfulMappings[path.String()] = mapping
			updates = append(updates, u...)
			break
		}
	}
	return updates, nil
}
