// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/snmpoc"
	"github.com/sirupsen/logrus"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

const (
	snmpEntPhysicalClass          = ".1.3.6.1.2.1.47.1.1.1.1.5"
	snmpEntPhysicalSerialNum      = ".1.3.6.1.2.1.47.1.1.1.1.11"
	snmpLldpLocChassisID          = ".1.0.8802.1.1.2.1.3.2.0"
	snmpLldpLocChassisIDSubtype   = ".1.0.8802.1.1.2.1.3.1.0"
	snmpLldpV2LocChassisID        = ".1.3.111.2.802.1.1.13.1.3.2.0"
	snmpLldpV2LocChassisIDSubtype = ".1.3.111.2.802.1.1.13.1.3.1.0"
	snmpSysUpTimeInstance         = ".1.3.6.1.2.1.1.3.0"
	snmpUsmStatsUnknownUserNames  = ".1.3.6.1.6.3.15.1.1.3.0"

	logFieldDeviceID = "deviceID"
)

// Split the final index off an OID and return it along with the remaining OID.
func oidSplitEnd(oid string) (string, string, error) {
	finalDotPos := strings.LastIndex(oid, ".")
	if finalDotPos < 0 {
		return "", "", fmt.Errorf("oid '%s' does not match expected format", oid)
	}
	return oid[:finalDotPos], oid[(finalDotPos + 1):], nil
}

// Snmp contains everything needed to implement an SNMP provider.
type Snmp struct {
	client gnmi.GNMIClient

	gsnmp      *gosnmp.GoSNMP // gosnmp object for Snmp's use
	tgsnmp     *gosnmp.GoSNMP // gosnmp object for translator's use
	mock       bool           // if true, don't do any network init
	translator *snmpoc.Translator

	// gosnmp can't handle parallel gets, so we also need to lock
	// access to its connection object.
	connectionLock sync.Mutex

	pollInterval time.Duration
	lastAlive    time.Time
	initialized  bool
	deviceID     string

	// List of files or directories to search for supported MIBs.
	mibs     []string
	mibStore smi.Store

	// Alternative Walk() and Get() for mock testing.
	getter func([]string) (*gosnmp.SnmpPacket, error)
	walker func(string, gosnmp.WalkFunc) error

	// Alternative time.Now() for mock testing.
	now func() time.Time

	// Logger
	logger *logrus.Entry
}

func (s *Snmp) snmpNetworkInit() error {
	if s.initialized || s.mock {
		return nil
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()
	if err := s.tgsnmp.Connect(); err != nil {
		return err
	}

	if err := s.gsnmp.Connect(); err != nil {
		return err
	}

	// Do an initial Get. If auth fails, give up before we start.
	pkt, err := s.unsafeGet(snmpSysUpTimeInstance)
	if err != nil {
		return err
	} else if len(pkt.Variables) > 0 {
		pdu := pkt.Variables[0]
		if oidExists(pdu) && pdu.Name == snmpUsmStatsUnknownUserNames {
			return fmt.Errorf("unknown username, can't proceed: %s", pdu.Name)
		}
	}

	s.initialized = true
	return nil
}

// get, but without a lock. Don't use this.
func (s *Snmp) unsafeGet(oid string) (*gosnmp.SnmpPacket, error) {
	pkt, err := s.getter([]string{oid})
	if err != nil {
		return nil, err
	}

	// Handle packet errors.
	if pkt.Error != gosnmp.NoError {
		errstr, ok := snmpoc.SNMPErrCodes[pkt.Error]
		if !ok {
			errstr = "Unknown error"
		}
		return nil, fmt.Errorf("Error in packet (%v): %v", pkt, errstr)
	}
	return pkt, nil
}

func (s *Snmp) get(oid string) (*gosnmp.SnmpPacket, error) {
	s.logger.Debugf("get (OID = %s)", oid)
	if s.getter == nil {
		return nil, errors.New("SNMP getter not set")
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()

	pkt, err := s.unsafeGet(oid)
	if err != nil {
		return nil, err
	}

	s.lastAlive = s.now()

	return pkt, nil
}

func oidExists(pdu gosnmp.SnmpPDU) bool {
	return pdu.Type != gosnmp.NoSuchObject && pdu.Type != gosnmp.NoSuchInstance
}

func (s *Snmp) getFirstPDU(oid string) (*gosnmp.SnmpPDU, error) {
	pkt, err := s.get(oid)
	if err != nil {
		return nil, err
	}
	if len(pkt.Variables) == 0 {
		return nil, fmt.Errorf("No variables in SNMP packet for OID %s", oid)
	}
	return &pkt.Variables[0], err
}

// getString does a Get on the specified OID, an octet string, and
// returns the result as a string.
func (s *Snmp) getString(oid string) (string, error) {
	pdu, err := s.getFirstPDU(oid)

	// Accept a noSuchObject or noSuchInstance, but otherwise, if it's not
	// an octet string, something went wrong.
	if err != nil || !oidExists(*pdu) {
		return "", err
	}
	if pdu.Type != gosnmp.OctetString {
		return "", fmt.Errorf("Variable type in PDU for OID %s is not octet string", oid)
	}

	return string(pdu.Value.([]byte)), nil
}

func (s *Snmp) walk(rootOid string, walkFn gosnmp.WalkFunc) error {
	s.logger.Debugf("walk (OID = %s)", rootOid)
	if s.walker == nil {
		return errors.New("SNMP walker not set")
	}

	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()
	err := s.walker(rootOid, walkFn)
	if err != nil {
		return err
	}
	s.logger.Debugf("walk complete (OID = %s)", rootOid)
	s.lastAlive = s.now()
	return err
}

var errStopWalk = errors.New("stop walk")

func (s *Snmp) getSerialNumber() (string, error) {
	serial := ""
	var done bool
	chassisIndex := ""
	var snmpEntPhysicalClassTypeChassis = 3

	// Get the serial number corresponding to the index whose class
	// type is chassis(3).
	s.logger.Tracef("getSerialNumber")
	entPhysicalWalk := func(data gosnmp.SnmpPDU) error {
		// If we're finished, throw a pseudo-error to indicate to the
		// walker that no more walking is required.
		if done {
			return errStopWalk
		}
		baseOid, index, err := oidSplitEnd(data.Name)
		if err != nil {
			return err
		}
		// If the physical class is "chassis", this is the index we want.
		if baseOid == snmpEntPhysicalClass {
			if data.Value == snmpEntPhysicalClassTypeChassis {
				chassisIndex = index
			}
		}
		if baseOid == snmpEntPhysicalSerialNum {
			// Take the first non-empty serial number as a backup, in
			// case there isn't a chassis serial number.
			if serial == "" {
				serial = string(data.Value.([]byte))
			}
			if index == chassisIndex && string(data.Value.([]byte)) != "" {
				serial = string(data.Value.([]byte))
				done = true
			}
		}

		return nil
	}

	if err := s.walk(snmpEntPhysicalClass, entPhysicalWalk); err != nil {
		return "", err
	}
	if err := s.walk(snmpEntPhysicalSerialNum, entPhysicalWalk); err != nil {
		if err != errStopWalk {
			return "", err
		}
	}
	s.logger.Tracef("getSerialNumber complete (serial = %v)", serial)
	return serial, nil
}

func (s *Snmp) getChassisID() (string, error) {
	s.logger.Tracef("getChassisID")
	var subtype string
	for _, subtypeOID := range []string{snmpLldpLocChassisIDSubtype,
		snmpLldpV2LocChassisIDSubtype} {
		pdu, err := s.getFirstPDU(subtypeOID)
		if err != nil {
			return "", err
		}
		if oidExists(*pdu) {
			v, ok := pdu.Value.(int)
			if !ok {
				s.logger.Errorf("Unexpected type '%T' for PDU value: %v"+
					" (req OID=%v)", pdu.Value, pdu.Value, subtypeOID)
				continue
			}
			subtype = openconfig.LLDPChassisIDType(v)
			break
		}
	}
	if subtype == "" {
		return "", nil
	}

	for _, oid := range []string{snmpLldpLocChassisID, snmpLldpV2LocChassisID} {
		pdu, err := s.getFirstPDU(oid)
		if err != nil {
			return "", err
		}
		if oidExists(*pdu) {
			v, ok := pdu.Value.([]byte)
			if !ok {
				s.logger.Errorf("Unexpected type '%T' for PDU value: %v "+
					"(OID=%v)", pdu.Value, pdu.Value, oid)
				continue
			}
			s.logger.Tracef("getChassisID (chassisID = %v)",
				chassisID(v, subtype))
			return chassisID(v, subtype), nil
		}
	}
	s.logger.Traceln("getChassisID: no chassis ID")
	return "", nil
}

// DeviceID returns the device ID.
func (s *Snmp) DeviceID() (string, error) {
	s.logger.Trace("Snmp.DeviceID")
	if err := s.snmpNetworkInit(); err != nil {
		return "", fmt.Errorf("Error connecting to device %q: %w",
			s.gsnmp.Target, err)
	}

	if s.deviceID != "" {
		return s.deviceID, nil
	}

	did, err := s.getSerialNumber()
	if err != nil {
		return did, err
	} else if did != "" {
		s.deviceID = did
		s.loggerUseDeviceID()
		return did, nil
	}

	did, err = s.getChassisID()
	if err != nil {
		return did, err
	} else if did != "" {
		s.deviceID = did
		s.loggerUseDeviceID()
		return did, nil
	}

	// The device didn't give us a serial number. Use the device
	// address instead. It's not great but better than nothing.
	s.logger.Infof("Failed to retrieve serial number for device '%s'; "+
		"using address for device ID", s.gsnmp.Target)
	s.deviceID = s.gsnmp.Target
	s.loggerUseDeviceID()
	return s.gsnmp.Target, nil
}

func (s *Snmp) loggerUseDeviceID() {
	s.logger = s.logger.WithField(logFieldDeviceID, s.deviceID)
}

// Alive checks if device is still alive if poll interval has passed.
func (s *Snmp) Alive() (bool, error) {
	s.logger.Debugf("Alive")
	if err := s.snmpNetworkInit(); err != nil {
		return false, fmt.Errorf("Error connecting to device: %v", err)
	}
	if time.Since(s.lastAlive) < s.pollInterval {
		return true, nil
	}
	_, err := s.get(snmpSysUpTimeInstance)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *Snmp) stop() {
	if !s.mock {
		s.tgsnmp.Conn.Close()
		s.gsnmp.Conn.Close()
	}
}

var chassisIDSubtypeMacAddress = openconfig.LLDPChassisIDType(4)
var chassisIDSubtypeNetworkAddress = openconfig.LLDPChassisIDType(5)

func chassisID(b []byte, subtype string) string {
	if subtype == chassisIDSubtypeMacAddress {
		return snmpoc.MacFromBytes(b)
	} else if subtype == chassisIDSubtypeNetworkAddress {
		return snmpoc.IPFromBytes(b)
	}
	return snmpoc.BytesToSanitizedString(b)
}

// InitGNMI initializes the Snmp provider with a gNMI client.
func (s *Snmp) InitGNMI(client gnmi.GNMIClient) {
	s.client = client
}

// OpenConfig indicates that this provider wants OpenConfig
// type-checking.
func (s *Snmp) OpenConfig() bool {
	return true
}

// Origin indicates that this provider streams to an OpenConfig model.
func (s *Snmp) Origin() string {
	return "openconfig"
}

func (s *Snmp) sendUpdates(ctx context.Context) error {
	return s.translator.Poll(ctx, s.client, []string{".*"})
}

func ignoredError(err error) bool {
	if err == io.EOF || err == context.Canceled {
		return true
	}
	return false
}

// Run sets the Snmp provider running and returns only on error.
func (s *Snmp) Run(ctx context.Context) error {
	if s.client == nil {
		return errors.New("Run called before InitGNMI")
	}
	s.logger.Debugf("Run")

	if err := s.snmpNetworkInit(); err != nil {
		return fmt.Errorf("Error connecting to device: %v", err)
	}
	s.logger.Debugf("gosnmp.Connect complete")

	mibStore, err := smi.NewStore(s.mibs...)
	if err != nil {
		return fmt.Errorf("Error creating MIB store: %s", err)
	}
	s.mibStore = mibStore

	translator, err := snmpoc.NewTranslator(mibStore, s.tgsnmp)
	if err != nil {
		return fmt.Errorf("Failed creating Translator: %v", err)
	}

	s.translator = translator
	s.translator.Mock = s.mock
	s.translator.Walker = s.walker
	s.translator.Getter = s.getter
	s.translator.Logger = s.logger

	// Do periodic state updates forever.
	if err := s.sendUpdates(ctx); err != nil && !ignoredError(err) {
		s.logger.Infof("Error in sendUpdates: %s", err)
	}

	tick := time.NewTicker(s.pollInterval)
	for {
		select {
		case <-tick.C:
			if err := s.sendUpdates(ctx); err != nil && !ignoredError(err) {
				s.logger.Infof("Error in sendUpdates: %s", err)
			}
		case <-ctx.Done():
			goto finish
		}
	}

finish:
	s.stop()
	return nil
}

// V3Params contains options related to SNMPv3.
type V3Params struct {
	SecurityModel gosnmp.SnmpV3SecurityModel
	Level         gosnmp.SnmpV3MsgFlags
	UsmParams     *gosnmp.UsmSecurityParameters
}

// NewSNMPProvider returns a new SNMP provider for the device at 'address'
// using a community value for authentication and pollInterval for rate
// limiting requests.
func NewSNMPProvider(address string, port uint16, community string,
	pollInt time.Duration, version gosnmp.SnmpVersion,
	v3Params *V3Params, mibs []string, mock bool) provider.GNMIProvider {
	gsnmp := gosnmp.GoSNMP{
		Port:               port,
		Version:            version,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
		Target:             address,
		Community:          community,
		Timeout:            time.Duration(2) * time.Second,
		Logger:             nil,
		MaxRepetitions:     12,
	}
	if v3Params != nil {
		gsnmp.MsgFlags = v3Params.Level
		gsnmp.SecurityModel = v3Params.SecurityModel
		gsnmp.SecurityParameters = v3Params.UsmParams
	}
	translatorGoSNMP := gsnmp

	s := &Snmp{
		tgsnmp:       &translatorGoSNMP,
		gsnmp:        &gsnmp,
		pollInterval: pollInt,
		mibs:         mibs,
		mock:         mock,
		getter:       gsnmp.Get,
		walker:       gsnmp.BulkWalk,
		now:          time.Now,
	}
	s.logger = log.Log(s).WithField(logFieldDeviceID,
		fmt.Sprintf("%s:%d", address, port))

	s.logger.Debugf("NewSNMPProvider, address: %v, version: %v",
		address, version)

	return s
}
