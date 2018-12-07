// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	psnmp "arista/provider/snmp"
	"time"
)

func init() {
	options := map[string]device.Option{
		"address": device.Option{
			Description: "Hostname or address of device",
			Required:    true,
		},
		"port": device.Option{
			Description: "Device SNMP port to use",
			Default:     "161",
		},
		"community": device.Option{
			Description: "SNMP community string",
			Required:    true,
		},
		"pollInterval": device.Option{
			Description: "Polling interval, in seconds",
			Default:     "20",
		},
	}

	device.Register("snmp", newSnmp, options)
}

type snmp struct {
	address      string
	community    string
	pollInterval time.Duration
	systemID     string
	snmpProvider provider.GNMIOpenConfigProvider
}

func (s *snmp) Type() device.Type {
	return device.Target
}

// XXX_jcr: For now, we return an error rather than just returning false. We
// may want to rethink that in the future.
func (s *snmp) CheckAlive() (bool, error) {
	_, err := s.snmpProvider.(*psnmp.Snmp).CheckAlive()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *snmp) DeviceID() (string, error) {
	if s.systemID != "" {
		return s.systemID, nil
	}
	systemID, err := s.snmpProvider.(*psnmp.Snmp).DeviceID()
	if err != nil {
		return "", err
	}
	s.systemID = systemID
	return s.systemID, nil
}

func (s *snmp) Providers() ([]provider.Provider, error) {
	return []provider.Provider{s.snmpProvider}, nil
}

// XXX_jcr: The network operations here could fail on startup, and if
// they do, the error will be passed back to Collector and it will fail.
// Are we OK with this or should we be doing retries?
func newSnmp(options map[string]string) (device.Device, error) {
	s := &snmp{}
	var err error

	s.address, err = device.GetAddressOption("address", options)
	if err != nil {
		return nil, err
	}

	s.community, err = device.GetStringOption("community", options)
	if err != nil {
		return nil, err
	}

	s.pollInterval, err = device.GetDurationOption("pollInterval", options)
	if err != nil {
		return nil, err
	}

	s.snmpProvider = psnmp.NewSNMPProvider(s.address, s.community, s.pollInterval)

	return s, nil
}
