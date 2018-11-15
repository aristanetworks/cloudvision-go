// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"arista/device"
	"arista/kernel"
	"arista/provider"

	"github.com/aristanetworks/glog"
)

func init() {
	// Set options
	options := map[string]device.Option{
		"procfsPeriod": device.Option{
			Description: "Interval at which to collect various stats " +
				"from /proc (0 to disable)",
			Default:  "15s",
			Required: false,
		},
		"systemID": device.Option{
			Description: "Device system ID override",
			Default:     "",
			Required:    false,
		},
	}

	// Register
	device.RegisterDevice("cvp", NewCvp, options)
}

type cvpDevice struct {
	procfsPeriod   time.Duration
	systemID       string
	kernelProvider provider.EOSProvider
}

func (c *cvpDevice) CheckAlive() (bool, error) {
	return true, nil
}

func (c *cvpDevice) Type() device.Type {
	return device.ManagementSystem
}

func (c *cvpDevice) DeviceID() (string, error) {
	return c.systemID, nil
}

func (c *cvpDevice) Providers() []provider.Provider {
	provs := []provider.Provider{
		c.kernelProvider,
	}
	return provs
}

// getSerialNumber returns the UUID of the CVP node
// where the process is currently running
func getSerialNumber() (string, error) {
	output, err := exec.Command("dmidecode").CombinedOutput()
	if err == nil {
		uuidLine := regexp.MustCompile("UUID: [a-zA-Z0-9-]+").Find(output)
		if uuidLine != nil {
			return string(uuidLine[6:]), nil
		}
	}
	glog.Infof("Failed to fetch serial number from host UUID: %s\n%s", err, output)

	if hostname, hostnameErr := os.Hostname(); hostnameErr == nil {
		hostname = strings.Replace(hostname, ".", "-", -1)
		return hostname, nil
	}

	return "", err
}

func getProcfsPeriod(opt map[string]string) (time.Duration, error) {
	pp, ok := opt["procfsPeriod"]
	if !ok {
		return time.Duration(0), errors.New("No option procfsPeriod")
	}
	ppd, err := time.ParseDuration(pp)
	if err != nil {
		return time.Duration(0), err
	}
	return ppd, nil
}

func getSystemID(opt map[string]string) (string, error) {
	sys, ok := opt["systemID"]
	if !ok {
		return "", errors.New("No option systemID")
	}
	if sys == "" {
		sn, err := getSerialNumber()
		if err != nil {
			return "", err
		}
		sys = sn
	}
	return sys, nil
}

// NewCvp returns a cvp device.
func NewCvp(opt map[string]string) (device.Device, error) {
	// Create pidfile for cvpi
	procPid := os.Getpid()
	file, err := os.Create("/var/run/cvpi/monitor.pid")
	if err != nil {
		return nil, errors.New("Unable to create process pid file under /var/run/cvpi/")
	}
	fmt.Fprintf(file, strconv.Itoa(procPid))
	file.Close()

	cvp := &cvpDevice{}

	// Extract config
	pp, err := getProcfsPeriod(opt)
	if err != nil {
		return nil, err
	}
	sys, err := getSystemID(opt)
	glog.Infof("SystemID set to serial number: %s", sys)
	if err != nil {
		return nil, err
	}

	cvp.procfsPeriod = pp
	cvp.systemID = sys
	cvp.kernelProvider = kernel.New(pp)

	return cvp, nil
}
