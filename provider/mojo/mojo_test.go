// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package mojo

import (
	pgnmi "arista/provider/gnmi"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// Samples from http://apidocs.mojonetworks.com/manage-devices.
var deviceUpdate1 = `{
	"boxId": 1,
	"name": "Mojo_A3:A9:FF",
	"vendorName": "Mojo",
	"firstDetectedTime": 1394019224180,
	"macaddress": "00:11:74:A3:A9:FF",
	"locationId": {
		"type": "locallocationid",
		"id": -1
	},
	"deviceMode": "SENSOR",
	"templateId": -2,
	"templateName": "System Template",
	"model": "SS-300-AT-C-60",
	"placed": false,
	"softwareVersion": "7.0.30",
	"capability": 33,
	"monitoredVlanCount": 0,
	"totalVlanCount": 0,
	"platformId": 11,
	"meshEnabled": false,
	"failSafeMode": false,
	"deviceNote": null,
	"ipAddress": "192.168.9.63",
	"active": false,
	"upgradeStatus": 0,
	"troubleshootingStatus": 0,
	"radios": null,
	"spectrumAnalysisState": "STOPPED",
	"quarantineStatus": "QUARANTINE_STATUS_OFF",
	"clientConnectivityState": "OFF",
	"targetAPClientConnectivityStatus": "OFF",
	"networkTag": null,
	"powerSource": "NOT_AVAILABLE",
	"bleStatus": "NOT_AVAILABLE",
	"upSince": 1394213971490
}`
var deviceUpdate2 = `{
	"boxId": 2,
	"name": "Mojo_9C:BC:9F",
	"vendorName": "Mojo",
	"firstDetectedTime": 1513576892217,
	"macaddress": "00:11:74:9C:BC:9F",
	"locationId": {
		"type": "locallocationid",
		"id": -1
	},
	"deviceMode": "AP",
	"templateId": -2,
	"templateName": "System Template",
	"model": "C-55",
	"placed": false,
	"softwareVersion": "8.2.1-902.34",
	"capability": 33,
	"monitoredVlanCount": 0,
	"totalVlanCount": 1,
	"platformId": 13,
	"meshEnabled": false,
	"failSafeMode": false,
	"deviceNote": null,
	"ipAddress": "172.17.30.96",
	"active": true,
	"upgradeStatus": 0,
	"troubleshootingStatus": 0,
	"radios": [
		{
			"radioId": 2,
			"active": false,
			"noiseFloor": 0,
			"worstClientRssi": 0,
			"txPower": 0,
			"rfUtilization": 0,
			"capability": 12,
			"apMode": true,
			"operatingBand": "BAND_2_4_GHZ",
			"retryRate": 0,
			"upstreamUsage": 0,
			"downstreamUsage": 0,
			"wirelessInterfaces": [
				{
					"bssid": "00:11:74:9C:BC:92",
					"ssid": "sched_always_up",
					"active": false,
					"numAssocClients": 0,
					"security": 1,
					"secAuth": 0,
					"secGCS": 0,
					"secPWCS": 0
				},
				{
					"bssid": "00:11:74:9C:BC:91",
					"ssid": "m-sch-g",
					"active": false,
					"numAssocClients": 0,
					"security": 8,
					"secAuth": 1,
					"secGCS": 8,
					"secPWCS": 8
				},
				{
					"bssid": "00:11:74:9C:BC:90",
					"ssid": "Documentation",
					"active": true,
					"numAssocClients": 0,
					"security": 8,
					"secAuth": 1,
					"secGCS": 8,
					"secPWCS": 8
				}
			],
			"channel": 0,
			"channelWidth": 0,
			"channelOffset": 0,
			"mcsSet11ac": 0,
			"beaconInterval": 100,
			"supportedRates": -2144683266
		},
		{
			"radioId": 1,
			"active": true,
			"noiseFloor": -105,
			"worstClientRssi": 0,
			"txPower": 21,
			"rfUtilization": 4,
			"capability": 3,
			"apMode": true,
			"operatingBand": "BAND_2_4_GHZ",
			"retryRate": 0,
			"upstreamUsage": 0,
			"downstreamUsage": 0,
			"wirelessInterfaces": [
				{
					"bssid": "00:11:74:9C:BC:92",
					"ssid": "sched_always_up",
					"active": false,
					"numAssocClients": 0,
					"security": 1,
					"secAuth": 0,
					"secGCS": 0,
					"secPWCS": 0
				},
				{
					"bssid": "00:11:74:9C:BC:91",
					"ssid": "m-sch-g",
					"active": false,
					"numAssocClients": 0,
					"security": 8,
					"secAuth": 1,
					"secGCS": 8,
					"secPWCS": 8
				},
				{
					"bssid": "00:11:74:9C:BC:90",
					"ssid": "Documentation",
					"active": true,
					"numAssocClients": 0,
					"security": 8,
					"secAuth": 1,
					"secGCS": 8,
					"secPWCS": 8
				}
			],
			"channel": 6,
			"channelWidth": 0,
			"channelOffset": 0,
			"mcsSet11ac": 0,
			"beaconInterval": 100,
			"supportedRates": -2144683266
		}
	],
	"spectrumAnalysisState": "STOPPED",
	"quarantineStatus": "QUARANTINE_STATUS_OFF",
	"clientConnectivityState": "OFF",
	"targetAPClientConnectivityStatus": "OFF",
	"networkTag": "192.168.55.0/24",
	"powerSource": "NOT_AVAILABLE",
	"bleStatus": "NOT_AVAILABLE",
	"upSince": 1513749001743
}`

type testCase struct {
	name         string
	deviceUpdate string
	expected     *gnmi.SetRequest
}

func testUpdate(t *testing.T, mj *mojo, gc *pgnmi.TestClient, tc testCase) {
	u := ManagedDevice{}
	if err := json.Unmarshal([]byte(tc.deviceUpdate), &u); err != nil {
		t.Fatalf("Error in Unmarshal: %s", err)
	}

	mj.deviceUpdateChan <- &u
	got := <-gc.Out

	if !reflect.DeepEqual(got, tc.expected) {
		t.Fatalf("SetRequests not equal. Expected %v\nGot: %v",
			tc.expected, got)
	}
}

func TestMojo(t *testing.T) {
	gNMIClient := &pgnmi.TestClient{
		Out: make(chan *gnmi.SetRequest),
	}
	mj := NewMojoProvider(make(chan *ManagedDevice)).(*mojo)
	mj.InitGNMIOpenConfig(gNMIClient)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = mj.Run(ctx)
	}()

	for _, tc := range []testCase{
		{
			name:         "deviceUpdate1",
			deviceUpdate: deviceUpdate1,
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{
					pgnmi.Path("system", "state"),
					pgnmi.Path("components")},
				Replace: []*gnmi.Update{
					pgnmi.Update(platformComponentConfigPath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentPath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentStatePath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentStatePath("hardware-version"),
						pgnmi.Strval("SS-300-AT-C-60")),
					pgnmi.Update(platformComponentStatePath("software-version"),
						pgnmi.Strval("7.0.30")),
					pgnmi.Update(pgnmi.Path("system", "state", "hostname"),
						pgnmi.Strval("001174A3A9FF")),
					pgnmi.Update(pgnmi.Path("system", "state", "boot-time"),
						pgnmi.Uintval(uint64(139421397149000))),
				},
			},
		},
		{
			name:         "deviceUpdate2",
			deviceUpdate: deviceUpdate2,
			expected: &gnmi.SetRequest{
				Delete: []*gnmi.Path{
					pgnmi.Path("system", "state"),
					pgnmi.Path("components")},
				Replace: []*gnmi.Update{
					pgnmi.Update(platformComponentConfigPath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentPath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentStatePath("name"), pgnmi.Strval("chassis")),
					pgnmi.Update(platformComponentStatePath("hardware-version"),
						pgnmi.Strval("C-55")),
					pgnmi.Update(platformComponentStatePath("software-version"),
						pgnmi.Strval("8.2.1-902.34")),
					pgnmi.Update(pgnmi.Path("system", "state", "hostname"),
						pgnmi.Strval("0011749CBC9F")),
					pgnmi.Update(pgnmi.Path("system", "state", "boot-time"),
						pgnmi.Uintval(uint64(151374900174300))),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testUpdate(t, mj, gNMIClient, tc)
		})
	}
}
