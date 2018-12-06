// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	psnmp "arista/provider/snmp"
	"errors"
	"net"
	"strconv"
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
	snmpProvider provider.GNMIProvider
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

func (s *snmp) Providers() []provider.Provider {
	return []provider.Provider{s.snmpProvider}
}

func getAddress(options map[string]string) (string, error) {
	addr, ok := options["address"]
	if !ok {
		return "", errors.New("No option 'address'")
	}

	// Validate IP
	ip := net.ParseIP(addr)
	if ip != nil {
		return ip.String(), nil
	}

	// Try for hostname if it's not an IP
	addrs, err := net.LookupIP(addr)
	if err != nil {
		return "", err
	}
	return addrs[0].String(), nil
}

func getCommunity(options map[string]string) (string, error) {
	comm, ok := options["community"]
	if !ok {
		return "", errors.New("No option 'community'")
	}
	return comm, nil
}

func getPollInterval(options map[string]string) (time.Duration, error) {
	interval, ok := options["pollInterval"]
	if !ok {
		return 0, errors.New("No option 'pollInterval'")
	}
	intv, err := strconv.ParseInt(interval, 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(intv) * time.Second, nil
}

// XXX_jcr: The network operations here could fail on startup, and if
// they do, the error will be passed back to Collector and it will fail.
// Are we OK with this or should we be doing retries?
func newSnmp(options map[string]string) (device.Device, error) {
	s := &snmp{}
	var err error

	s.address, err = getAddress(options)
	if err != nil {
		return nil, err
	}

	s.community, err = getCommunity(options)
	if err != nil {
		return nil, err
	}

	s.pollInterval, err = getPollInterval(options)
	if err != nil {
		return nil, err
	}

	s.snmpProvider = psnmp.NewSNMPProvider(s.address, s.community, s.pollInterval)

	return s, nil
}
