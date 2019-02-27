// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	_ "github.com/aristanetworks/cloudvision-go/device/devices"  // import all registered devices
	_ "github.com/aristanetworks/cloudvision-go/device/managers" // import all registered managers
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/aristanetworks/glog"
	aflag "github.com/aristanetworks/goarista/flag"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"golang.org/x/sync/errgroup"
)

// Main is the "real" main.
func Main() {
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

		// MockCollector config
		mock          = flag.Bool("mock", false, "Run Collector in mock mode")
		mockCheckPath = aflag.Map{}
		mockTimeout   = flag.Duration("mockTimeout", 60*time.Second,
			"Timeout for checking notifications in mock mode")

		// gNMI server config
		gnmiServerAddr = flag.String("gnmiServerAddr", "localhost:6030",
			"Address of gNMI server")
	)
	flag.Var(mockCheckPath, "mockCheckPath",
		"<path>=<feature> option for mock mode, where <path> is a path that,"+
			" if present in the Collector output, signifies that the target device supports "+
			"the feature described in <feature>")
	flag.Var(deviceOptions, "deviceoption", "<key>=<value> option for the Device. "+
		"May be repeated to set multiple Device options.")
	flag.Var(managerOptions, "manageroption", "<key>=<value> option for the Manager. "+
		"May be repeated to set multiple Manager options.")
	flag.BoolVar(help, "h", false, "Print program options")

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
			addHelp(*managerName, *deviceName)
		}
		flag.Usage()
		return
	}

	// We're running for real at this point. Check that the config
	// is sane.
	validateConfig(*managerName, *deviceName, *deviceConfigFile, *mock, mockCheckPath)

	// Create inventory.
	group, ctx := errgroup.WithContext(context.Background())
	var inventory device.Inventory
	mockInfo := newMockInfo(mockCheckPath)
	if *mock {
		inventory = device.NewInventory(ctx, group,
			pgnmi.NewSimpleGNMIClient(mockInfo.processRequest))
	} else {
		gnmiClient, err := agnmi.Dial(&agnmi.Config{Addr: *gnmiServerAddr})
		if err != nil {
			glog.Fatal(err)
		}
		inventory = device.NewInventory(ctx, group, gnmiClient)
	}

	// Populate inventory with manager or from configured devices.
	if *managerName != "" {
		manager, err := device.CreateManager(*managerName, managerOptions)
		if err != nil {
			glog.Fatal(err)
		}
		group.Go(func() error {
			return manager.Manage(inventory)
		})
	} else {
		devices, err := device.CreateDevices(*deviceName, *deviceConfigFile, deviceOptions)
		if err != nil {
			glog.Fatal(err)
		}
		mockInfo.initDevices(devices)
		for _, info := range devices {
			err := inventory.Add(info.ID, info.Device)
			if err != nil {
				glog.Fatalf("Error in inventory.Add(): %v", err)
			}
		}
		if *deviceIDFile != "" {
			err := device.DumpDeviceIDs(devices, *deviceIDFile)
			if err != nil {
				glog.Fatal(err)
			}
		}
	}
	glog.V(2).Info("Collector is running")
	errChan := make(chan error)
	go func() {
		// Watch for errors.
		err := group.Wait()
		if err == nil {
			err = errors.New("device routines returned unexpectedly")
		}
		errChan <- err
	}()
	if *mock {
		err := mockInfo.waitForUpdates(errChan, *mockTimeout)
		if err != nil {
			glog.Fatal(err)
		}
	} else {
		glog.Fatal(<-errChan)
	}
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

func addHelp(managerName, deviceName string) error {
	var oh map[string]string
	var name string
	var optionType string
	var err error

	if managerName != "" {
		name = managerName
		optionType = "manager"
		oh, err = device.ManagerOptionHelp(name)
	} else {
		name = deviceName
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

func validateConfig(managerName, deviceName, deviceConfigFile string,
	mock bool, mockCheckPath map[string]string) {
	// A device or a device manager must be specified unless we're running with -h
	if deviceName == "" && managerName == "" && deviceConfigFile == "" {
		glog.Fatal("-device, -manager, or -config must be specified.")
	}

	if deviceName != "" && managerName != "" {
		glog.Fatal("-device and -manager should not be both specified.")
	}

	if deviceConfigFile != "" && managerName != "" {
		glog.Fatal("-config and -manager should not be both specified.")
	}

	if deviceConfigFile != "" && deviceName != "" {
		glog.Fatal("-config and -device should not be both specified.")
	}

	if mock && managerName != "" {
		glog.Fatal("-manager should not be specified in mock mode")
	}

	if !mock && len(mockCheckPath) > 0 {
		glog.Fatal("-mockCheckPath is only valid in mock mode")
	}
}
