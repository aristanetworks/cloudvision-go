// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package managers

import (
	"arista/device"
	"arista/device/devices"
	pmojo "arista/provider/mojo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	options := map[string]device.Option{
		"apiurl": device.Option{
			Description: "MWM API endpoint",
			Required:    true,
		},
		"pollInterval": device.Option{
			Description: "Polling interval, in seconds",
			Default:     "30",
		},
		"jsessionid": device.Option{
			Description: "JSESSIONID cookie (temporary; for testing until we have an API key)",
			Required:    true,
		},
	}
	device.RegisterManager("mwm", newMWM, options)
}

type mwm struct {
	jsessionid       string
	apiURL           string
	pollInterval     time.Duration
	client           *http.Client
	inventory        device.Inventory
	deviceUpdateChan map[string]chan *pmojo.ManagedDevice
}

type managedDevicesResponse struct {
	TotalCount   int                   `json:"totalCount"`
	NextLink     string                `json:"nextLink"`
	PreviousLink string                `json:"previousLink"`
	Devices      []pmojo.ManagedDevice `json:"managedDevices"`
}

func (m *mwm) addDevice(d pmojo.ManagedDevice) (device.Device, error) {
	ch := make(chan *pmojo.ManagedDevice)
	did := pmojo.DeviceIDFromMac(d.Macaddress)
	md := devices.NewMojo(did, ch)
	if err := m.inventory.Add(did, md); err != nil {
		return nil, err
	}
	m.deviceUpdateChan[did] = ch
	return md, nil
}

func (m *mwm) removeDevice(key string) error {
	if err := m.inventory.Delete(key); err != nil {
		return err
	}
	delete(m.deviceUpdateChan, key)
	return nil
}

func (m *mwm) handleDeviceUpdate(d pmojo.ManagedDevice) error {
	did := pmojo.DeviceIDFromMac(d.Macaddress)
	_, ok := m.inventory.Get(did)
	if !ok {
		if !d.Active {
			return nil
		}
		_, err := m.addDevice(d)
		if err != nil {
			return err
		}
	}

	if !d.Active {
		if err := m.removeDevice(did); err != nil {
			return err
		}
		return nil
	}

	m.deviceUpdateChan[did] <- &d
	return nil
}

func (m *mwm) getManagedDevices(url string) (*managedDevicesResponse, error) {
	if url == "" {
		url = m.apiURL + "/new/webservice/v6/devices/manageddevices"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: m.jsessionid})
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting manageddevices: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	mdr := &managedDevicesResponse{}
	if err := json.Unmarshal(body, mdr); err != nil {
		return nil, err
	}
	return mdr, nil
}

func (m *mwm) pollMWM() error {
	// Keep requesting more devices until we've got them all.
	url := ""
	for {
		resp, err := m.getManagedDevices(url)
		if err != nil {
			return err
		}

		for _, d := range resp.Devices {
			if err := m.handleDeviceUpdate(d); err != nil {
				return err
			}
		}

		// XXX_jcr: Is 100 always the limit per request?
		if len(resp.Devices) == 100 {
			url = resp.NextLink
		} else {
			break
		}
	}
	return nil
}

func (m *mwm) Manage(inventory device.Inventory) error {
	m.inventory = inventory

	m.client = &http.Client{
		Timeout: time.Second * 20,
	}

	if err := m.pollMWM(); err != nil {
		return err
	}

	tick := time.NewTicker(m.pollInterval)
	defer tick.Stop()
	for range tick.C {
		if err := m.pollMWM(); err != nil {
			return err
		}
	}

	return nil
}

func newMWM(options map[string]string) (device.Manager, error) {
	m := &mwm{}
	var err error

	m.apiURL, err = device.GetStringOption("apiurl", options)
	if err != nil {
		return nil, err
	}

	m.jsessionid, err = device.GetStringOption("jsessionid", options)
	if err != nil {
		return nil, err
	}

	m.pollInterval, err = device.GetDurationOption("pollInterval", options)
	if err != nil {
		return nil, err
	}
	m.deviceUpdateChan = make(map[string]chan *pmojo.ManagedDevice)

	return m, nil
}
