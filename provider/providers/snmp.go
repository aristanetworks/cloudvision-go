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
	"time"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/path"
	"github.com/soniah/gosnmp"
)

type snmp struct {
	provider.ReadOnly
	ready  chan struct{}
	done   chan struct{}
	period time.Duration
	root   types.Entity
}

var connected bool

const updatePeriod = 10 * time.Second

func SNMPNetworkInit() error {
	if connected {
		return nil
	}
	err := gosnmp.Default.Connect()
	if err == nil {
		connected = true
	}
	return err
}

func SNMPGetByOID(oid string) (string, error) {
	oids := []string{oid}

	// Ask for object
	result, err := gosnmp.Default.Get(oids)
	if err != nil {
		return "", err
	}

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

func (s *snmp) WaitForNotification() {
	<-s.ready
}

func (s *snmp) Stop() {
	<-s.ready
	gosnmp.Default.Conn.Close()
	close(s.done)
}

func (s *snmp) setSystemConfig() error {
	hostname, err := SNMPGetByOID("1.3.6.1.2.1.1.5.0")
	if err != nil {
		return err
	}

	systemType := types.NewEntityType("::OpenConfig::system::config")
	systemType.RemoveAttribute("name")
	systemType.AddAttribute("hostname", types.StringType)
	data := map[string]interface{}{"hostname": hostname}
	_, err = entity.MakeDirsWithAttributes(s.root,
		path.New("OpenConfig", "system", "config"), nil, systemType, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *snmp) updateInterfaceCounters() {
}

func (s *snmp) Run(schema *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	s.root = root

	err := SNMPNetworkInit()
	if err != nil {
		glog.Infof("Failed to connect to device: %s", err)
		return
	}

	// Update unchanging system state
	err = s.setSystemConfig()
	if err != nil {
		glog.Fatal("Failed to set system config")
	}
	close(s.ready)

	// Do periodic state updates
	tick := time.NewTicker(s.period)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			s.updateInterfaceCounters()
		case <-s.done:
			return
		}
	}
}

func NewSNMPProvider() provider.Provider {
	return &snmp{
		ready:  make(chan struct{}),
		done:   make(chan struct{}),
		period: updatePeriod,
	}
}
