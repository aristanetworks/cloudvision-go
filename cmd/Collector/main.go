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
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/aristanetworks/cloudvision-go/device"
	_ "github.com/aristanetworks/cloudvision-go/device/devices"
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/aristanetworks/glog"
	aflag "github.com/aristanetworks/goarista/flag"
	"golang.org/x/sync/errgroup"
	yaml "gopkg.in/yaml.v2"
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
	managerOptions   = aflag.Map{}
	deviceConfigFile = flag.String("configFile", "", "Path to the config file for devices")
	deviceIDFile     = flag.String("dumpDeviceIDs", "",
		"Path to output file used to associate device IDs with device configuration")

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
	if *deviceName == "" && *managerName == "" && *deviceConfigFile == "" {
		glog.Fatal("-device, -manager, or -config must be specified.")
	}

	if *deviceName != "" && *managerName != "" {
		glog.Fatal("-device and -manager should not be both specified.")
	}

	if *deviceConfigFile != "" && *managerName != "" {
		glog.Fatal("-config and -manager should not be both specified.")
	}

	if *deviceConfigFile != "" && *deviceName != "" {
		glog.Fatal("-config and -device should not be both specified.")
	}
}

type deviceInfo struct {
	id     string
	config device.Config
	device device.Device
}

func createDevice(name string, options map[string]string) (*deviceInfo, error) {
	config := device.Config{
		Device:  name,
		Options: options,
	}
	d, err := device.Create(name, options)
	if err != nil {
		return nil, fmt.Errorf("Failed creating device '%v'", config.Device)
	}
	did, err := d.DeviceID()
	if err != nil {
		return nil, fmt.Errorf("Error getting device ID: %v", err)
	}
	return &deviceInfo{config: config, id: did, device: d}, nil
}

func createDevices(name, configPath string, options map[string]string) ([]*deviceInfo, error) {
	var infos []*deviceInfo
	// Single configured device
	if name != "" {
		info, err := createDevice(name, options)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	// Config file
	if configPath != "" {
		readConfigs, err := device.ReadConfigs(configPath)
		if err != nil {
			return nil, fmt.Errorf("Error reading config file: %v", err)
		}
		for _, config := range readConfigs {
			info, err := createDevice(config.Device, config.Options)
			if err != nil {
				return nil, err
			}
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func dumpDeviceIDs(devices []*deviceInfo) error {
	idToConfig := map[string]device.Config{}
	for _, info := range devices {
		idToConfig[info.id] = info.config
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return fmt.Errorf("Error in ioutil.TempFile: %v", err)
	}
	defer os.Remove(f.Name())
	enc := yaml.NewEncoder(f)

	err = enc.Encode(&idToConfig)
	if err != nil {
		return fmt.Errorf("Error in yaml.Decode: %v", err)
	}
	f.Close()
	err = os.Rename(f.Name(), *deviceIDFile)
	if err != nil {
		return fmt.Errorf("Error in os.Rename: %v", err)
	}
	return nil
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
		devices, err := createDevices(*deviceName, *deviceConfigFile, deviceOptions)
		if err != nil {
			glog.Fatal(err)
		}
		for _, info := range devices {
			err := inventory.Add(info.id, info.device)
			if err != nil {
				glog.Fatalf("Error in inventory.Add(): %v", err)
			}
		}
		if *deviceIDFile != "" {
			err := dumpDeviceIDs(devices)
			if err != nil {
				glog.Fatal(err)
			}
		}
	}
	glog.V(2).Info("Collector is running")
	// Watch for errors.
	err := group.Wait()
	if err == nil {
		err = errors.New("device routines returned unexpectedly")
	}
	glog.Fatal(err)
}
