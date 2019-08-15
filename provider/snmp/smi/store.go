package smi

import (
	"fmt"
	"strconv"
	"strings"
)

// A Store contains the SMI parse tree and allows users to query
// for objects.
type Store interface {
	GetObject(oid string) *Object
}

// store implements the Store interface.
type store struct {
	modules map[string]*Module
	oids    map[string]*Object
	names   map[string]*Object
}

// NewStore returns a Store.
func NewStore(files ...string) (Store, error) {
	parseModules, err := parseFiles(files...)
	if err != nil {
		return nil, err
	}

	store := &store{
		modules: make(map[string]*Module),
		oids:    make(map[string]*Object),
		names:   make(map[string]*Object),
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

// GetObject takes a text or numeric object identifier and returns
// the corresponding parsed Object, if one exists.
func (s *store) GetObject(oid string) *Object {
	// Text OID
	ss := strings.Split(oid, "::")
	if len(ss) >= 2 {
		oid = ss[1]
	}

	// Numeric OID
	if strings.Contains(oid, ".") {
		// Remove leading "." if there is one.
		if oid[0] == '.' {
			oid = oid[1:]
		}

		// Remove trailing ".0" for scalars.
		ss = strings.Split(oid, ".")
		scalar := false
		if ss[len(ss)-1] == "0" {
			oid = strings.Join(ss[:(len(ss)-1)], ".")
			scalar = true
		}
		o, ok := s.oids[oid]
		if ok {
			return o
		}
		if scalar {
			return nil
		}

		// Start removing possible index values from the OID.
		// If we find an object with a matching OID and the right
		// number of indexes, return that.
		for i := len(ss) - 1; i > 0; i-- {
			if o, ok = s.oids[strings.Join(ss[:i], ".")]; ok {
				if o.Kind != KindColumn {
					return nil
				} else if o.Parent == nil {
					return nil
				}
				return o
			}
		}
		return nil
	}
	return s.names[oid]
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
