// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package smi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// A Store contains the SMI parse tree and allows users to query
// for objects.
type Store interface {
	GetObject(oid string) *Object
	GetOids() map[string]*Object
	GetNames() map[string]*Object
	GetModules() map[string]*Module
}

// JSONStoreObject defines the JSONStoreObject json structure.
type JSONStoreObject struct {
	Access      Access
	Description string
	Indexes     []string `json:"indexes,omitempty"`
	Kind        Kind
	Module      string
	Name        string
	Oid         string
	Status      Status
	// Parent and Children fields are  stored as
	// string, string-slice in the json instead of
	// storing the pointer object(s) for optimization and
	// to avoid marshalling cyclic data structure.
	// The parsed store from NewStore is used to populate
	// these fields.
	Parent   string   `json:"parent,omitempty"`
	Children []string `json:"children,omitempty"`
}

// JSONStore defines the JSONStore json structure.
type JSONStore struct {
	ObjectMap map[string]*JSONStoreObject `json:"jsonstoreobject,omitempty"`
	Modules   map[string]*Module          `json:"jsonmodule,omitempty"`
}

// store implements the Store interface.
type store struct {
	lock    *sync.RWMutex
	modules map[string]*Module
	oids    map[string]*Object
	names   map[string]*Object
	known   map[string]*Object
}

// NewStore returns a Store.
func NewStore(files ...string) (Store, error) {
	parseModules, err := parseFiles(files...)
	if err != nil {
		return nil, err
	}

	store := &store{
		lock:    &sync.RWMutex{},
		modules: make(map[string]*Module),
		oids:    make(map[string]*Object),
		names:   make(map[string]*Object),
		known:   make(map[string]*Object),
	}

	// After initially building the parse tree, there are certain
	// fixes we have to make that are easier to do once the store
	// already exists, such as resolving OIDs and certain indexes.
	resolvedModules := map[string]bool{}
	for moduleName, pm := range parseModules {
		if err := resolveModule(moduleName, store, parseModules,
			resolvedModules); err != nil {
			return nil, err
		}
		store.modules[moduleName] = createModule(pm)
	}
	return store, nil
}

// MergeJSONWithStore merges smi store data from a JSON file with a store object.
func MergeJSONWithStore(mergejsonbytes []byte, baseStore Store) (Store, error) {
	astore, ok := baseStore.(*store)
	if !ok {
		return nil, fmt.Errorf("error in casting %T to store object", baseStore)
	}
	var jsrStore JSONStore
	err := json.Unmarshal(mergejsonbytes, &jsrStore)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling json mib file: %v", err)
	}

	cstore := &store{
		lock:    &sync.RWMutex{},
		modules: make(map[string]*Module),
		oids:    make(map[string]*Object),
		names:   make(map[string]*Object),
		known:   make(map[string]*Object),
	}

	cstore.lock.Lock()
	// create store object from unmarshalled json file.
	for oid, jStoreObj := range jsrStore.ObjectMap {
		cstore.oids[oid] = &Object{
			Access:      jStoreObj.Access,
			Description: jStoreObj.Description,
			Kind:        jStoreObj.Kind,
			Module:      jStoreObj.Module,
			Name:        jStoreObj.Name,
			Oid:         jStoreObj.Name,
			Status:      jStoreObj.Status,
		}
		cstore.oids[oid].Indexes = make([]string, len(jStoreObj.Indexes))
		copy(cstore.oids[oid].Indexes, jStoreObj.Indexes)
		cstore.names[jStoreObj.Name] = cstore.oids[oid]
	}

	// Populate the Parent nodes in the store from the ObjectMap from the
	// unmarshalled json file.
	for oid, jStoreObj := range jsrStore.ObjectMap {
		parentOid := jsrStore.ObjectMap[oid].Parent
		parentObj := cstore.oids[parentOid]
		if parentOid != "" {
			if parentObj == nil {
				return nil, fmt.Errorf(
					"error in MergeJSONWithStore: Parent object is nil for oid: %s", parentOid)
			}
			cstore.oids[oid].Parent = parentObj
			cstore.names[jStoreObj.Name].Parent = parentObj
		}
	}

	// Populate the Children nodes in the store from the ObjectMap from the
	// unmarshalled json file.
	for oid, jStoreObj := range jsrStore.ObjectMap {
		for _, chi := range jStoreObj.Children {
			childrenObj := cstore.oids[chi]
			cstore.oids[oid].Children = append(cstore.oids[oid].Children, childrenObj)
			cstore.names[jStoreObj.Name].Children = append(
				cstore.names[jStoreObj.Name].Children, childrenObj)
		}
	}
	cstore.modules = jsrStore.Modules
	astore.lock.Lock()
	// merge jsonstore object cstore with astore.
	for oid, jStoreObj := range cstore.modules {
		_, ok := astore.modules[oid]
		if !ok {
			astore.modules[oid] = jStoreObj
		}
	}
	for k, v := range cstore.oids {
		_, ok := astore.oids[k]
		if !ok {
			astore.oids[k] = v
		}
	}
	for k, v := range cstore.names {
		_, ok := astore.names[k]
		if !ok {
			astore.names[k] = v
		}
	}
	astore.lock.Unlock()
	cstore.lock.Unlock()

	return astore, nil
}

// CheckStoresEqual checks for equality between two SMI stores and returns
// boolean for equality and error.
func CheckStoresEqual(storeA Store, storeB Store) (bool, error) {
	storeAObj, ok := storeA.(*store)
	if !ok {
		return false, fmt.Errorf("error in type conversion to store object")
	}
	storeBObj, ok := storeA.(*store)
	if !ok {
		return false, fmt.Errorf("error in type conversion to store object")
	}
	if !reflect.DeepEqual(storeAObj, storeBObj) {
		return false, fmt.Errorf(
			"mismatch in storeA and storeB objects: storeA: %v, \nstoreB: %v",
			storeAObj, storeB)
	}
	return true, nil
}

func (s *store) checkKnown(oid string) *Object {
	s.lock.RLock()
	o := s.known[oid]
	s.lock.RUnlock()
	return o
}

func (s *store) updateKnown(oid string, o *Object) {
	s.lock.Lock()
	s.known[oid] = o
	s.lock.Unlock()
}

// GetObject takes a text or numeric object identifier and returns
// the corresponding parsed Object, if one exists.
func (s *store) GetObject(oid string) *Object {
	// First check the cache.
	if o := s.checkKnown(oid); o != nil {
		return o
	}
	origOid := oid

	// Remove module name from text OID
	ss := strings.Split(oid, "::")
	if len(ss) >= 2 {
		oid = ss[1]
	}

	if strings.Contains(oid, ".") {
		// Remove leading "." if there is one.
		if oid[0] == '.' {
			oid = oid[1:]
		}

		var o *Object
		var ok bool

		// Start removing possible index values from the OID.
		// If we find an object with a matching OID and the right
		// number of indexes, return that.
		ss = strings.Split(oid, ".")
		for i := len(ss); i > 0; i-- {
			shortenedOid := strings.Join(ss[:i], ".")
			// Try it as a numeric OID first.
			o, ok = s.oids[shortenedOid]

			// Then try it as a text OID.
			if !ok {
				o, ok = s.names[shortenedOid]
			}
			if ok {
				// If we've removed indexes, this should be either a column or a scalar.
				if i < len(ss) {
					if o.Kind != KindColumn && o.Kind != KindScalar {
						return nil
					} else if o.Parent == nil {
						return nil
					}
				}
				s.updateKnown(origOid, o)
				return o
			}
		}
		return nil
	}

	o, ok := s.names[oid]
	if ok {
		s.updateKnown(origOid, o)
	}
	return o
}

// GetOids returns the map of numerical oid to corresponding parsed Object.
func (s *store) GetOids() map[string]*Object {
	return s.oids
}

// GetNames returns the map of oid name to corresponding parsed Object.
func (s *store) GetNames() map[string]*Object {
	return s.names
}

// GetModules returns the map of module-name to the corresponding parsed Module.
func (s *store) GetModules() map[string]*Module {
	return s.modules
}

func resolveOID(po *parseObject, store *store) error {
	// If we've already resolved this one, don't do it again unless
	// this is a newer version of an object we've already resolved.
	if o, ok := store.names[po.object.Name]; ok {
		if po.object.Module == o.Module ||
			moduleUpgrade(po.object.Module, po.object.Name) == o.Module ||
			moduleUpgrade(po.object.Module, po.object.Name) != po.object.Module {
			return nil
		}
	}

	newOid := []string{}
	for _, subid := range strings.Split(po.object.Oid, ".") {
		if _, err := strconv.Atoi(subid); err != nil {
			if subid == "iso" {
				newOid = append(newOid, "1")
			} else {
				p, ok := store.names[subid]
				if !ok {
					return fmt.Errorf("Could not find OID for '%s', module '%s'",
						subid, po.object.Module)
				}
				newOid = append(newOid, p.Oid)
			}
		} else {
			newOid = append(newOid, subid)
		}
	}
	po.object.Oid = strings.Join(newOid, ".")

	store.names[po.object.Name] = po.object
	store.oids[po.object.Oid] = po.object
	return nil
}

// Pull in the indexes for AUGMENTS rows.
func resolveIndexes(po *parseObject, store *store) error {
	if po.object.Kind != KindRow || po.augments == "" {
		return nil
	}
	ao, ok := store.names[po.augments]
	if !ok {
		return fmt.Errorf("Could not find augmented object: %s", po.augments)
	}
	po.object.Indexes = make([]string, len(ao.Indexes))
	copy(po.object.Indexes, ao.Indexes)
	return nil
}

func resolveTree(po *parseObject, store *store, keepgoing bool) error {
	if err := resolveOID(po, store); err != nil && !keepgoing {
		return err
	}
	if err := resolveIndexes(po, store); err != nil && !keepgoing {
		return err
	}

	for _, child := range po.children {
		if err := resolveTree(child, store, keepgoing); err != nil &&
			!keepgoing {
			return err
		}
	}

	return nil
}

// Once the parsing is finished we have OIDs that look like
// {hrSystem 1}. To get the full numeric OIDs we have to resolve
// the text part to its numeric value.
func resolveModule(moduleName string, store *store,
	parseModules map[string]*parseModule,
	resolvedModules map[string]bool) error {
	if _, ok := resolvedModules[moduleName]; ok {
		return nil
	}

	pm, ok := parseModules[moduleName]
	if !ok {
		return fmt.Errorf("Can't resolve unparsed module '%s'", moduleName)
	}

	// Resolve any modules this module imports.
	for _, imp := range pm.imports {
		mr := moduleUpgrade(imp.Module, imp.Object)
		if err := resolveModule(mr, store, parseModules,
			resolvedModules); err != nil {
			return err
		}
	}

	// Link orphans to parent objects.
	for _, orphan := range pm.orphans {
		if len(strings.Split(orphan.object.Oid, ".")) > 0 {
			parentName := strings.Split(orphan.object.Oid, ".")[0]
			if o, ok := store.names[parentName]; ok {
				o.Children = append(o.Children, orphan.object)
				orphan.object.Parent = o
			}
		}
	}

	// Try twice to resolve each OID, in case they're declared out of
	// order in the MIB.
	for pass := 1; pass <= 2; pass++ {
		for _, obj := range pm.objectTree {
			if err := resolveTree(obj, store, pass != 2); err != nil &&
				pass == 2 {
				return err
			}
		}
	}

	resolvedModules[moduleName] = true

	return nil
}

func createModule(pm *parseModule) *Module {
	m := &Module{
		Name:       pm.name,
		ObjectTree: []*Object{},
		Imports:    pm.imports,
	}
	for _, o := range pm.objectTree {
		m.ObjectTree = append(m.ObjectTree, o.object)
	}
	return m
}
