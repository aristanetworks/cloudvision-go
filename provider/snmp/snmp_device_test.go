// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmp

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
	"google.golang.org/grpc"
)

// These tests run the SNMP provider against dumps generated from
// real devices we care about. Checking that the provider generates
// the exact SetRequests we want would be extremely tedious, so
// we check instead that it's updating a few interesting paths.

// deviceTestCase describes a test of an SNMP provider against a real
// device dump.
type deviceTestCase struct {
	deviceName     string
	deviceDumpFile string
	expectedPaths  []*gnmi.Path
	polls          int
}

type snmpReqType string

const (
	get      string = "get"
	bulkwalk string = "bulkwalk"
)

// walkMap represents a walk of an SNMP tree as a map from OID
// to the set of PDUs returned by the walk that are rooted at the
// specified OID.
type walkMap map[string][]*gosnmp.SnmpPDU

// Generate a slice of walkMaps from a dump file.
func walkMapsFromDump(filename string) ([]walkMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open dump file: %v", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("Failed to unzip dump: %v", err)
	}

	walkMaps := []walkMap{}
	wm := make(walkMap)

	scanner := bufio.NewScanner(gz)
	firstLine := ""
	for scanner.Scan() {
		line := scanner.Text()

		pdu := pduFromString(scanner.Text())
		if pdu == nil {
			continue
		}

		// If this is the start of a new poll, store the existing walkMap
		// and start another.
		if firstLine == "" {
			firstLine = line
		} else if line == firstLine {
			walkMaps = append(walkMaps, wm)
			wm = make(walkMap)
		}

		// Store a pointer to this PDU for each of the valid prefixes
		// of its OID. So, e.g., a PDU with OID .1.2.3.4 would be
		// stored in the maps of ".1", ".1.2", ".1.2.3", and
		// ".1.2.3.4".
		sections := strings.Split(pdu.Name, ".")[1:]
		prefix := ""
		for _, section := range sections {
			prefix += "." + section
			wm[prefix] = append(wm[prefix], pdu)
		}
	}
	walkMaps = append(walkMaps, wm)

	return walkMaps, nil
}

// testGNMIClient implements gnmi.GNMIClient. All it does is check
// the provider's SetRequests for the specified paths. When it finds
// one of the paths it's looking for, it removes that path from
// pathsRemaining. Once pathsRemaining is empty, it decrements
// pollsRemaining. Once pollsRemaining is zero, it considers the test
// a success.
type testGNMIClient struct {
	ctx            context.Context
	cancel         context.CancelFunc
	polls          int
	pollsRemaining int
	paths          []*gnmi.Path
	pathsRemaining []*gnmi.Path
	lock           *sync.Mutex
}

func (m *testGNMIClient) Capabilities(ctx context.Context,
	in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	panic("not implemented")
}

func (m *testGNMIClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	panic("not implemented")
}

func pathElemMatches(e1, e2 *gnmi.PathElem) bool {
	if e1.Name != e2.Name {
		return false
	}
	if len(e1.Key) != len(e2.Key) {
		return false
	}
	for k, v := range e1.Key {
		v2, ok := e2.Key[k]
		if !ok {
			return false
		}
		if v != v2 && v2 != "*" {
			return false
		}
	}
	return true
}

// Check whether the path p1 matches the path p2, which may have
// wildcard ("*") elements.
func pathMatches(p1, p2 *gnmi.Path) bool {
	if len(p1.Elem) != len(p2.Elem) {
		return false
	}
	for i, e := range p1.Elem {
		if !pathElemMatches(e, p2.Elem[i]) {
			return false
		}
	}
	return true
}

func (m *testGNMIClient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// For each update in in.Replace, if it matches one of the paths
	// we want to see for each poll, mark that path as found. Once all
	// paths are found we can mark the poll complete.
	for _, update := range in.Replace {
		for i, expectedPath := range m.pathsRemaining {
			if pathMatches(update.Path, expectedPath) {
				m.pathsRemaining = append(m.pathsRemaining[:i], m.pathsRemaining[(i+1):]...)
				if len(m.pathsRemaining) == 0 {
					m.pollsRemaining--
					m.pathsRemaining = m.paths
				}
				if m.pollsRemaining == 0 {
					m.cancel()
				}
				return nil, nil
			}
		}
	}
	return nil, nil
}

func (m *testGNMIClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	panic("not implemented")
}

// testget and testwalk implement a primitive/wrong SNMP
// pseudo-agent. They aren't correct but given the simple nature of
// the request types we care about and the provider error handling,
// they should be correct enough for the provider to behave the same
// way it would to a real response.
func testget(oids []string, wm walkMap) (*gosnmp.SnmpPacket, error) {
	pkt := &gosnmp.SnmpPacket{}
	if len(oids) > 1 {
		panic("testget doesn't support multiple OIDs")
	}
	oid := oids[0]
	pdus, ok := wm[oid]
	if !ok {
		pkt.Variables = []gosnmp.SnmpPDU{
			gosnmp.SnmpPDU{
				Name:  oid,
				Type:  gosnmp.NoSuchObject,
				Value: nil,
			},
		}
		return pkt, nil
	}
	for _, p := range pdus {
		pkt.Variables = append(pkt.Variables, *p)
	}
	return pkt, nil
}

func testwalk(oid string, walker gosnmp.WalkFunc, wm walkMap) error {
	pdus, ok := wm[oid]
	if !ok {

		return nil
	}
	for _, pdu := range pdus {
		if err := walker(*pdu); err != nil {
			return err
		}
	}
	return nil
}

func newTestGNMIClient(cancel context.CancelFunc,
	tc deviceTestCase) *testGNMIClient {
	client := &testGNMIClient{
		cancel: cancel,
		polls:  tc.polls,
		paths:  tc.expectedPaths,
		lock:   &sync.Mutex{},
	}
	client.pollsRemaining = client.polls
	client.pathsRemaining = client.paths
	return client
}

func newSNMPProvider(client *testGNMIClient,
	walkMaps []walkMap) provider.GNMIProvider {
	p := NewSNMPProvider("whatever", "stuff", 10*time.Millisecond, false, true)

	// Set up provider with special getter + walker, keeping track of
	// which poll we're on.
	p.(*Snmp).getter = func(oids []string) (*gosnmp.SnmpPacket, error) {
		client.lock.Lock()
		defer client.lock.Unlock()
		poll := client.polls - client.pollsRemaining
		if poll >= client.polls {
			poll = client.polls - 1
		}
		return testget(oids, walkMaps[poll])
	}
	p.(*Snmp).walker = func(oid string, walker gosnmp.WalkFunc) error {
		client.lock.Lock()
		defer client.lock.Unlock()
		poll := client.polls - client.pollsRemaining
		if poll >= client.polls {
			poll = client.polls - 1
		}
		return testwalk(oid, walker, walkMaps[poll])
	}
	return p
}

func runDeviceTest(t *testing.T, tc deviceTestCase) {
	walkMaps, err := walkMapsFromDump(tc.deviceDumpFile)
	if err != nil {
		t.Fatalf("Failed processing dump file: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := newTestGNMIClient(cancel, tc)
	prov := newSNMPProvider(client, walkMaps)
	prov.InitGNMI(client)

	if err := prov.Run(ctx); err != nil {
		t.Fatalf("runDeviceTest failed in provider.Run: %v", err)
	}

	client.lock.Lock()
	defer client.lock.Unlock()
	if client.pollsRemaining != 0 {
		t.Fatal("runDeviceTest did not finish polling")
	}
}

func basicPaths() []*gnmi.Path {
	return []*gnmi.Path{
		pgnmi.Path("system", "state", "hostname"),
		pgnmi.PlatformComponentStatePath("*", "description"),
		pgnmi.IntfStateCountersPath("*", "in-octets"),
		pgnmi.LldpNeighborStatePath("*", "*", "chassis-id"),
	}
}

// Run the SNMP provider against a set of dumps generated by the SNMP
// provider running against real target devices. Check that the
// provider generates the sorts of SetRequests we expect.
func TestDevices(t *testing.T) {
	for _, tc := range []deviceTestCase{
		{
			deviceName:     "Arista_DCS-7150S-24_4.21.3F-2GB-INT",
			deviceDumpFile: "dumps/Arista_DCS-7150S-24_4.21.3F-2GB-INT_20190301.gz",
			expectedPaths:  basicPaths(),
			polls:          2,
		},
		{
			deviceName:     "Arista_DCS-7508_4.21.3F-INT",
			deviceDumpFile: "dumps/Arista_DCS-7508_4.21.3F-INT_20190301.gz",
			expectedPaths:  basicPaths(),
			polls:          2,
		},
		{
			deviceName:     "PaloAlto_PA-3060_8.1.4",
			deviceDumpFile: "dumps/PaloAlto_PA-3060_8.1.4_20190301.gz",
			expectedPaths:  basicPaths(),
			polls:          2,
		},
		{
			deviceName:     "Cisco_N9K-C93128TX_7.0(3)I7(1)",
			deviceDumpFile: "dumps/Cisco_N9K-C93128TX_7.03I71_20190301.gz",
			expectedPaths:  basicPaths(),
			polls:          2,
		},
	} {
		t.Run(tc.deviceName, func(t *testing.T) {
			runDeviceTest(t, tc)
		})
	}
}
