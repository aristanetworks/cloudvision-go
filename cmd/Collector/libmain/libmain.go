// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	_ "github.com/aristanetworks/cloudvision-go/device/devices" // import all registered devices
	"github.com/aristanetworks/cloudvision-go/version"
	"github.com/aristanetworks/glog"
	aflag "github.com/aristanetworks/goarista/flag"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"golang.org/x/sync/errgroup"
)

var (
	v    = flag.Bool("version", false, "Print the version number")
	help = flag.Bool("help", false, "Print program options")

	// Device config
	deviceName = flag.String("device", "",
		"Device type (available devices: "+deviceList()+")")
	deviceOptions    = aflag.Map{}
	deviceConfigFile = flag.String("configFile", "", "Path to the config file for devices")
	deviceIDFile     = flag.String("dumpDeviceIDFile", "",
		"Path to output file used to associate device IDs with device configuration")

	// MockCollector config
	mock        = flag.Bool("mock", false, "Run Collector in mock mode")
	mockFeature = aflag.Map{}
	mockTimeout = flag.Duration("mockTimeout", 60*time.Second,
		"Timeout for checking notifications in mock mode")

	// Dump Collector config
	dump        = flag.Bool("dump", false, "Run Collector in dump mode")
	dumpFile    = flag.String("dumpFile", "", "Path to output file used to dump gNMI SetRequests")
	dumpTimeout = flag.Duration("dumpTimeout", 20*time.Second,
		"Timeout for dumping gNMI SetRequests")

	// gNMI server config
	gnmiServerAddr = flag.String("gnmiServerAddr", "localhost:6030",
		"Address of gNMI server")
)

// Main is the "real" main.
func Main() {
	flag.Var(mockFeature, "mockFeature",
		"<feature>=<path> option for mock mode, where <path> is a path that, "+
			"if present in the Collector output, signifies that the target device supports "+
			"the feature described in <feature>")
	flag.Var(deviceOptions, "deviceoption", "<key>=<value> option for the Device. "+
		"May be repeated to set multiple Device options.")
	flag.BoolVar(help, "h", false, "Print program options")

	flag.Parse()

	// Print version.
	if *v {
		fmt.Println(version.Version, runtime.Version())
		return
	}

	// Print help, including device-specific help,
	// if requested.
	if *help {
		if *deviceName != "" {
			addHelp()
		}
		flag.Usage()
		return
	}

	// We're running for real at this point. Check that the config
	// is sane.
	validateConfig()

	if *mock {
		runMock(context.Background())
		return
	}
	if *dump {
		runDump(context.Background())
		return
	}
	runMain(context.Background())
}

func runMain(ctx context.Context) {
	gnmiClient, err := agnmi.Dial(&agnmi.Config{Addr: *gnmiServerAddr})
	if err != nil {
		glog.Fatal(err)
	}
	// Create inventory.
	inventory := device.NewInventory(ctx, gnmiClient)
	devices, err := device.CreateDevices(*deviceName, *deviceConfigFile, deviceOptions)
	if err != nil {
		glog.Fatal(err)
	}

	group, ctx := errgroup.WithContext(ctx)
	for _, info := range devices {
		err := inventory.Add(info.ID, info.Device)
		if err != nil {
			glog.Fatalf("Error in inventory.Add(): %v", err)
		}
		if manager, ok := info.Device.(device.Manager); ok {
			group.Go(func() error {
				return manager.Manage(inventory)
			})
		} else {
			group.Go(func() error {
				select {
				case <-ctx.Done():
					return nil
				}
			})
		}
	}
	if *deviceIDFile != "" {
		err := device.DumpDeviceIDs(devices, *deviceIDFile)
		if err != nil {
			glog.Fatal(err)
		}
	}
	glog.V(2).Info("Collector is running")
	if err := group.Wait(); err != nil {
		glog.Fatal(err)
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

func addHelp() error {
	oh, err := device.OptionHelp(*deviceName)
	if err != nil {
		return fmt.Errorf("addHelp: %v", err)
	}

	var formattedOptions string
	if len(oh) > 0 {
		b := new(bytes.Buffer)
		aflag.FormatOptions(b, "Help options for device '"+*deviceName+"':", oh)
		formattedOptions = b.String()
	}

	aflag.AddHelp("", formattedOptions)
	return nil
}

func validateConfig() {
	// A device or a device config must be specified unless we're running with -h
	if *deviceName == "" && *deviceConfigFile == "" {
		glog.Fatal("-device or -config must be specified.")
	}

	if *deviceConfigFile != "" && *deviceName != "" {
		glog.Fatal("-config and -device should not be both specified.")
	}

	if !*mock && len(mockFeature) > 0 {
		glog.Fatal("-mockFeature is only valid in mock mode")
	}

	if *mock && *dump {
		glog.Fatal("-mock and -dump should not be both specified")
	}

	if *dump && *dumpFile == "" {
		glog.Fatal("-dumpFile must be specified in dump mode")
	}
}
