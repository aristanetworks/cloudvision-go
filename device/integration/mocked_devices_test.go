// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

// +build integration,testeddevices

package devicetest

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arista/aeris/atypes"
	_ "arista/aeris/turbine/blades/general"
	_ "arista/aeris/turbine/blades/network"
	_ "arista/aeris/turbine/blades/statistics"
	_ "arista/aeris/turbine/blades/versions"
	"arista/aeris/turbine/test/vdc"
	"arista/test/notiftest"
	"arista/types"

	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"
)

var (
	jsonDumpURLs  = atypes.FlagOptions{}
	jsonDumpPaths = atypes.FlagOptions{}
)

func init() {
	flag.Var(jsonDumpURLs, "jsondumpurl",
		"The json dump url to be downloaded in format <url>=<deviceid>."+
			"May be repeated to test multiple devices.")
	flag.Var(jsonDumpPaths, "jsondumppath",
		"The local json dump path in format <filepath>=<deviceid>."+
			"May be repeated to test multiple devices.")
}

func downloadDump(url string) (filePath string, err error) {
	filePath = filepath.Join("/tmp", filepath.Base(url))
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Bad status: %s", resp.Status)
	}
	// Writer the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return
}

func notifsFromDumpFile(filePath string) ([]types.Notification, error) {
	var notifs []types.Notification
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		notif := types.NewNotification(time.Now(), nil, nil, nil)
		err = json.Unmarshal(scanner.Bytes(), &notif)
		if err != nil {
			return nil, fmt.Errorf("Error in json.Unmarshal: %v", err)
		}
		notifs = append(notifs, notif)
	}
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error in file scanner: %v", err)
	}
	return notifs, nil
}

func testDumpFileUpdatePaths(t *testing.T, filePath, deviceID string, v vdc.Vdc,
	chanToExpectedPaths map[<-chan types.Notification][]key.Path) {

	notifs, err := notifsFromDumpFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, notif := range notifs {
		v.Publish(deviceID, notif)
	}
	for ch, paths := range chanToExpectedPaths {
		notiftest.CheckUnorderedPaths(t, ch, paths)
	}
}

func TestMockedDevice(t *testing.T) {

	if len(jsonDumpURLs) == 0 && len(jsonDumpPaths) == 0 {
		t.Fatal("Either jsondumpurl or jsondumppath should be specified")
	}

	for url, deviceid := range jsonDumpURLs {
		filePath, err := downloadDump(url)
		if err != nil {
			t.Fatalf("Error downloading file from %s: %v", url, err)
		}
		jsonDumpPaths[filePath] = deviceid
		defer os.Remove(filePath)
	}

	for filePath, deviceid := range jsonDumpPaths {
		v := vdc.New(t)

		vdc.StartTurbine(v, t, "../../aeris/cmd/turbine/configs/streaming-status.yml", nil)
		vdc.StartTurbine(v, t, "../../aeris/cmd/turbine/configs/version-lldp-state.yml", nil)
		vdc.StartTurbine(v, t, "../../aeris/cmd/turbine/configs/version-lldp-neighbors.yml", nil)

		_, versionDevicesOut := vdc.StartTurbine(v, t,
			"../../aeris/cmd/turbine/configs/version-devices.yml", nil)
		_, devicesOut := vdc.StartTurbine(v, t, "../../aeris/cmd/turbine/configs/devices.yml", nil)
		_, lldpOut := vdc.StartTurbine(v, t,
			"../../aeris/cmd/turbine/configs/network-lldp-neighbors.yml", nil)
		_, intfOut := vdc.StartTurbine(v, t,
			"../../aeris/cmd/turbine/configs/rate-openconfig-intf-counters.yml", nil)

		testDumpFileUpdatePaths(t, filePath, deviceid, v, map[<-chan types.Notification][]key.Path{
			versionDevicesOut: []key.Path{
				path.New("Devices", deviceid, "versioned-data", "Device"),
			},
			devicesOut: []key.Path{
				path.New("DatasetInfo", "Devices"),
			},
			lldpOut: []key.Path{
				path.New("network", "v1", "topology", "nodes"),
			},
			intfOut: []key.Path{
				path.New("Devices", deviceid, "versioned-data", "interfaces", "data",
					path.Wildcard, "rates"),
			},
		})
		v.Stop()
	}
}
