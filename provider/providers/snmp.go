// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/gopenconfig/eos"
	"arista/gopenconfig/eos/converter"
	"arista/gopenconfig/model/node"
	"arista/provider"
	"arista/schema"
	"arista/types"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aristanetworks/glog"
	"github.com/soniah/gosnmp"
)

type snmp struct {
	provider.ReadOnly
	ready          chan struct{}
	done           chan struct{}
	errc           chan error
	ctx            context.Context
	cancel         context.CancelFunc
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

// Given an incoming PDU, update the appropriate interface state.
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
		err = OpenConfigUpdateLeaf(s.ctx, intfPath(intfName, "state"),
			"name", string(pdu.Value.([]byte)))
	case ifAdminStatus:
		err = OpenConfigUpdateLeaf(s.ctx, intfPath(intfName, "state"),
			"admin-status", ifStatus(pdu.Value.(int)))
	case ifOperStatus:
		err = OpenConfigUpdateLeaf(s.ctx, intfPath(intfName, "state"),
			"oper-status", ifStatus(pdu.Value.(int)))
	case ifInOctets:
		err = OpenConfigUpdateLeaf(s.ctx, intfPath(intfName, "state", "counters"),
			"in-octets", uint64(pdu.Value.(uint)))
	case ifOutOctets:
		err = OpenConfigUpdateLeaf(s.ctx, intfPath(intfName, "state", "counters"),
			"out-octets", uint64(pdu.Value.(uint)))
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

	return OpenConfigUpdateLeaf(s.ctx, node.NewPath("system", "config"),
		"hostname", hostname)
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
			err = s.updateSystemConfig()
			if err != nil {
				glog.Errorf("Failure in updateSystemConfig: %v", err)
				return
			}
			err = s.updateInterfaces()
			if err != nil {
				glog.Infof("Failure in updateInterfaces: %s", err)
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
	pollInt = pollInterval
	return &snmp{
		ready:          make(chan struct{}),
		done:           make(chan struct{}),
		errc:           make(chan error),
		interfaceIndex: make(map[string]string),
		address:        address,
		community:      community,
	}
}
