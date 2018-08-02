// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package devices

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"arista/agent"
	"arista/device"
	"arista/kernel"
	"arista/provider"

	"github.com/aristanetworks/glog"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	device.RegisterDevice("cvp", NewCvp)
}

type cvpConfig struct {
	procfsPeriod time.Duration `yaml:"procfsperiod"`
	systemID     string        `yaml:"systemid"`
}

type cvpDevice struct {
	name            string
	systemID        string
	kernelProvider  provider.Provider
	versionProvider provider.Provider
}

func (c *cvpDevice) Name() string {
	return c.name
}

func (c *cvpDevice) CheckAlive() bool {
	return true
}

func (c *cvpDevice) DeviceID() string {
	return c.systemID
}

func (c *cvpDevice) Providers() []provider.Provider {
	provs := []provider.Provider{
		c.kernelProvider,
		c.versionProvider,
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

// NewCvp returns a cvp device
func NewCvp(yamlConfig io.Reader) (device.Device, error) {
	// Create pidfile for cvpi
	procPid := os.Getpid()
	file, err := os.Create("/var/run/cvpi/monitor.pid")
	if err != nil {
		glog.Fatalf("Unable to create process pid file under /var/run/cvpi/")
	}
	fmt.Fprintf(file, strconv.Itoa(procPid))
	file.Close()

	config := &cvpConfig{
		procfsPeriod: 15 * time.Second,
	}
	if yamlConfig != nil {
		parser := yaml.NewDecoder(yamlConfig)
		parser.SetStrict(true)
		err := parser.Decode(config)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse config file: %s", err)
		}
	}
	cvp := &cvpDevice{
		name:            "CvpMonitor",
		kernelProvider:  kernel.New(config.procfsPeriod),
		versionProvider: agent.NewVersionProvider(),
	}
	if config.systemID == "" {
		serialNumber, err := getSerialNumber()
		if err != nil {
			glog.Fatalf("Unable to get serial number: %s", err)
		}
		glog.Infof("SystemID set to serial number: %s", serialNumber)
		config.systemID = serialNumber
	}
	cvp.systemID = config.systemID
	return cvp, nil
}
