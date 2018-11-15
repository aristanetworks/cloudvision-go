// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"arista/device"
	"arista/provider"
	"arista/provider/providers"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aristanetworks/goarista/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

func init() {
	// Set options
	options := map[string]device.Option{
		"gnmi_addr": device.Option{
			Description: "gNMI server host/port",
			Required:    true,
		},
		"gnmi_paths": device.Option{
			Description: "gNMI subscription path (comma-separated if multiple)",
			Default:     "/",
			Required:    false,
		},
		"gnmi_username": device.Option{
			Description: "gNMI subscription username",
			Required:    true,
		},
	}

	// Register
	device.RegisterDevice("openconfig", newOpenConfig, options)
}

type openconfigDevice struct {
	gNMIProvider provider.GNMIProvider
	gNMIClient   pb.GNMIClient
	config       *gnmi.Config
}

func (o *openconfigDevice) Type() device.Type {
	return device.Target
}

func (o *openconfigDevice) CheckAlive() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, o.config)
	livenessPath := "/system/processes/process/state"
	req, err := gnmi.NewGetRequest(gnmi.SplitPaths([]string{livenessPath}))
	if err != nil {
		return false, err
	}
	resp, err := o.gNMIClient.Get(ctx, req)
	return err == nil && resp != nil && len(resp.Notification) > 0, nil
}

func (o *openconfigDevice) Providers() []provider.Provider {
	return []provider.Provider{o.gNMIProvider}
}

func (o *openconfigDevice) DeviceID() (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, o.config)
	// TODO: Use /components/component/state/serial-no after checking state/type is Chassis.
	// Not doing that now because EOS doesn't support it so I don't even know what the gNMI
	// response looks like..
	livenessPath := "/system/config"
	req, err := gnmi.NewGetRequest(gnmi.SplitPaths([]string{livenessPath}))
	if err != nil {
		return "", err
	}
	resp, err := o.gNMIClient.Get(ctx, req)
	if err != nil || len(resp.Notification) == 0 {
		return "", fmt.Errorf("Unable to get request to %v: %v", livenessPath, err)
	}
	var config map[string]string
	val := providers.Unmarshal(resp.Notification[0].Update[0].Val)
	err = json.Unmarshal(val.([]byte), &config)
	if err != nil {
		return "", err
	}
	return config["openconfig-system:hostname"] + "." + config["openconfig-system:domain-name"], nil
}

// newOpenConfig returns an openconfig device.
func newOpenConfig(opt map[string]string) (device.Device, error) {
	gNMIAddr := opt["gnmi_addr"]
	gNMIUsername := opt["gnmi_username"]
	gNMIPaths := strings.Split(opt["gnmi_paths"], ",")
	openconfig := &openconfigDevice{}
	config := &gnmi.Config{
		Addr:     gNMIAddr,
		Username: gNMIUsername,
	}
	client, err := gnmi.Dial(config)
	if err != nil {
		return nil, err
	}
	openconfig.gNMIClient = client
	openconfig.config = config

	openconfig.gNMIProvider = providers.NewGNMIProvider(client, config, gNMIPaths)

	return openconfig, nil
}
