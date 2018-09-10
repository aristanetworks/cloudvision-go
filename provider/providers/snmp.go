// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/entity"
	"arista/provider"
	"arista/schema"
	"arista/types"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"
	"github.com/soniah/gosnmp"
)

type snmp struct {
	provider.ReadOnly
	ready          chan struct{}
	done           chan struct{}
	root           types.Entity
	interfaceIndex map[string]string
	address        string
	community      string
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
	return SNMPGetByOID(".1.3.6.1.2.1.47.1.1.1.1.11.1")
}

// SNMPCheckAlive checks if device is still alive if poll interval has passed.
func SNMPCheckAlive() (bool, error) {
	if time.Since(lastAlive) < pollInt {
		return true, nil
	}
	_, err := SNMPGetByOID("1.3.6.1.2.1.1.3.0")
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

// Check a path for an entity and return it if it exists.
func getEntityAtPath(root types.Entity, path key.Path) *types.Entity {
	cwd := entity.NewCwd(root, "").ChangeDirExact(path)
	if cwd == nil {
		return nil
	}
	ent := cwd.GetEntity()

	return &ent
}

// Update the given attribute inside the state entity for the specified interface.
func updateIntfStateAttr(root types.Entity, intfName string, attr string, val interface{}) error {
	// Create state and counters types
	intfStateCountersType := types.NewEntityType("::OpenConfig::interface::state::counters")
	intfStateCountersType.AddAttribute("in-octets", types.U64Type)
	intfStateCountersType.AddAttribute("out-octets", types.U64Type)
	intfStateType := types.NewEntityType("::OpenConfig::interface::state")
	intfStateType.AddAttribute("admin-status", types.S64Type) // XXX_jcr: convert to enum?
	intfStateType.AddAttribute("oper-status", types.S64Type)
	intfStateType.AddAttribute("counters", intfStateCountersType)

	// Look for the entity. If it doesn't exist, create it and instantiate
	// the counters entity inside it.
	path := path.New("OpenConfig", "interfaces", intfName, "state")
	ent := getEntityAtPath(root, path)
	if ent == nil {
		e, err := entity.MakeDirsWithAttributes(root, path, nil, intfStateType,
			map[string]interface{}{attr: val})
		if err != nil {
			return err
		}
		err = e.SetAttribute("counters", map[string]interface{}{"name": "counters"})
		return err
	}

	err := (*ent).SetAttribute(attr, val)
	return err
}

// Update the given attribute inside the counters entity for the specified interface.
func updateIntfCountersAttr(root types.Entity, intfName string, attr string,
	val interface{}) error {
	path := path.New("OpenConfig", "interfaces", intfName, "state", "counters")

	ent := getEntityAtPath(root, path)
	if ent == nil {
		return fmt.Errorf("Unable to create entity")
	}

	err := (*ent).SetAttribute(attr, val)
	return err
}

func (s *snmp) handleInterfacePDU(pdu gosnmp.SnmpPDU) error {
	ifDescr := ".1.3.6.1.2.1.2.2.1.2"
	ifAdminStatus := ".1.3.6.1.2.1.2.2.1.7"
	ifOperStatus := ".1.3.6.1.2.1.2.2.1.8"
	ifInOctets := ".1.3.6.1.2.1.2.2.1.10"
	ifOutOctets := ".1.3.6.1.2.1.2.2.1.16"

	// Get/set interface name from index. If there's no mapping, just return and
	// wait for the mapping to show up.
	baseOid, index, err := oidIndex(pdu.Name)
	if err != nil {
		return err
	}
	intfName, ok := s.interfaceIndex[index]
	if !ok && baseOid != ifDescr {
		return nil
	} else if !ok && baseOid == ifDescr {
		intfName = string(pdu.Value.([]byte))
		s.interfaceIndex[index] = intfName
	}

	err = nil
	switch baseOid {
	case ifDescr:
		err = updateIntfStateAttr(s.root, intfName, "name", string(pdu.Value.([]byte)))
	case ifAdminStatus:
		err = updateIntfStateAttr(s.root, intfName, "admin-status", int64(pdu.Value.(int)))
	case ifOperStatus:
		err = updateIntfStateAttr(s.root, intfName, "oper-status", int64(pdu.Value.(int)))
	case ifInOctets:
		err = updateIntfCountersAttr(s.root, intfName, "in-octets", uint64(pdu.Value.(uint)))
	case ifOutOctets:
		err = updateIntfCountersAttr(s.root, intfName, "out-octets", uint64(pdu.Value.(uint)))
	}
	// default: ignore update
	return err
}

func (s *snmp) updateInterfaces() error {
	// XXX_jcr: We still need to add code for understanding deletes.
	return gosnmp.Default.Walk(".1.3.6.1.2.1.2.2",
		func(data gosnmp.SnmpPDU) error {
			return s.handleInterfacePDU(data)
		})
}

func (s *snmp) updateSystemConfig() error {
	hostname, err := SNMPGetByOID("1.3.6.1.2.1.1.5.0")
	if err != nil {
		return err
	}

	systemConfigType := types.NewEntityType("::OpenConfig::system::config")
	systemConfigType.AddAttribute("hostname", types.StringType)
	data := map[string]interface{}{"hostname": hostname}
	_, err = entity.MakeDirsWithAttributes(s.root,
		path.New("OpenConfig", "system", "config"), nil, systemConfigType, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *snmp) Run(schema *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	s.root = root

	err := snmpNetworkInit()
	if err != nil {
		glog.Infof("Failed to connect to device: %s", err)
		return
	}
	close(s.ready)

	// Do periodic state updates
	tick := time.NewTicker(pollInt)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err = s.updateSystemConfig()
			if err != nil {
				glog.Infof("Failure in updateSystemConfig: %s", err)
				return
			}
			err = s.updateInterfaces()
			if err != nil {
				glog.Infof("Failure in updateInterfaces: %s", err)
			}
		case <-s.done:
			return
		}
	}
}

// NewSNMPProvider returns a new SNMP provider for device at address using community value for
// authentication and pollInterval for rate limiting requests such as keepalive.
func NewSNMPProvider(address string, community string,
	pollInterval time.Duration) provider.Provider {
	gosnmp.Default.Target = address
	gosnmp.Default.Community = community
	pollInt = pollInterval
	return &snmp{
		ready:          make(chan struct{}),
		done:           make(chan struct{}),
		interfaceIndex: make(map[string]string),
		address:        address,
		community:      community,
	}
}
