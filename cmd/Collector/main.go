// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/aristanetworks/cloudvision-go/device"
	_ "github.com/aristanetworks/cloudvision-go/device/devices"
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/aristanetworks/glog"
	aflag "github.com/aristanetworks/goarista/flag"
	"golang.org/x/sync/errgroup"
)

var (
	v    = flag.Bool("version", false, "Print the version number")
	help = flag.Bool("help", false, "Print program options")

	// Device config
	deviceName = flag.String("device", "",
		"Device type (available devices: "+deviceList()+")")
	deviceOptions = aflag.Map{}
	managerName   = flag.String("manager", "",
		"Manager type (available managers: "+managerList()+")")
	managerOptions = aflag.Map{}

	// gNMI server config
	gnmiServerAddr = flag.String("gnmiServerAddr", "localhost:6030",
		"Address of gNMI server")
)

func init() {
	flag.Var(deviceOptions, "deviceoption", "<key>=<value> option for the Device. "+
		"May be repeated to set multiple Device options.")
	flag.Var(managerOptions, "manageroption", "<key>=<value> option for the Manager. "+
		"May be repeated to set multiple Manager options.")
	flag.BoolVar(help, "h", false, "Print program options")
}

// Return a formatted list of available devices.
func deviceList() string {
	dl := device.Registered()
	if len(dl) > 0 {
		return strings.Join(dl, ", ")
	}
	return "none"
}

// Return a formatted list of available managers.
func managerList() string {
	ml := device.RegisteredManagers()
	if len(ml) > 0 {
		return strings.Join(ml, ", ")
	}
	return "none"
}

func addHelp() error {
	var oh map[string]string
	var name string
	var optionType string
	var err error

	if *managerName != "" {
		name = *managerName
		optionType = "manager"
		oh, err = device.ManagerOptionHelp(name)
	} else {
		name = *deviceName
		optionType = "device"
		oh, err = device.OptionHelp(name)
	}
	if err != nil {
		return fmt.Errorf("addHelp: %v", err)
	}

	var formattedOptions string
	if len(oh) > 0 {
		b := new(bytes.Buffer)
		aflag.FormatOptions(b, "Help options for "+optionType+" '"+name+"':", oh)
		formattedOptions = b.String()
	}

	aflag.AddHelp("", formattedOptions)
	return nil
}

func validateConfig() {
	// A device or a device manager must be specified unless we're running with -h
	if *deviceName == "" && *managerName == "" {
		glog.Fatal("-device or -manager must be specified.")
	}

	if *deviceName != "" && *managerName != "" {
		glog.Fatal("-device and -manager should not be both specified.")
	}
}

func createDevices() (devices []device.Device, err error) {
	deviceConfigs := []device.Config{}

	// Single configured device
	if *deviceName != "" {
		deviceConfigs = []device.Config{
			device.Config{
				Device:  *deviceName,
				Options: deviceOptions,
			},
		}
	}

	// Config file
	if *deviceName == "" {
		// XXX TODO
		return nil, errors.New("Config file reading not yet implemented")
	}

	// Create devices from configs.
	for _, dc := range deviceConfigs {
		d, err := device.Create(dc.Device, dc.Options)
		if err != nil {
			return nil, fmt.Errorf("Failed creating device '%v'", dc.Device)
		}
		devices = append(devices, d)
	}
	return
}

func main() {
	flag.Parse()

	// Print version.
	if *v {
		fmt.Println(version.Version, runtime.Version())
		return
	}

	// Print help, including device/manager-specific help,
	// if requested.
	if *help {
		if *deviceName != "" || *managerName != "" {
			addHelp()
		}
		flag.Usage()
		return
	}

	// We're running for real at this point. Check that the config
	// is sane.
	validateConfig()

	// Create inventory.
	group, ctx := errgroup.WithContext(context.Background())
	inventory := device.NewInventory(ctx, group, *gnmiServerAddr)

	// Populate inventory with manager or from configured devices.
	if *managerName != "" {
		manager, err := device.CreateManager(*managerName, managerOptions)
		if err != nil {
			glog.Fatal(err)
		}
		go func() {
			err = manager.Manage(inventory)
			if err != nil {
				glog.Fatal(err)
			}
		}()
	} else {
		devices, err := createDevices()
		if err != nil {
			glog.Fatal(err)
		}
		for _, d := range devices {
			did, err := d.DeviceID()
			if err != nil {
				glog.Fatalf("Error getting device ID: %v", err)
			}
			if err = inventory.Add(did, d); err != nil {
				glog.Fatal(err)
			}
		}
	}

	// Watch for errors.
	err := group.Wait()
	if err == nil {
		err = errors.New("device routines returned unexpectedly")
	}
	glog.Fatal(err)
}
