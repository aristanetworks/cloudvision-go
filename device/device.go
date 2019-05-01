// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aristanetworks/cloudvision-go/provider"
	yaml "gopkg.in/yaml.v2"
)

// A Device knows how to interact with a specific device.
type Device interface {
	Alive() (bool, error)
	DeviceID() (string, error)
	Providers() ([]provider.Provider, error)
}

// A Manager manages a device inventory, adding and deleting
// devices as appropriate.
type Manager interface {
	Device
	Manage(inventory Inventory) error
}

// Creator returns a new instance of a Device.
type Creator = func(map[string]string) (Device, error)

// registrationInfo contains all the information about a device that's
// knowable before it's instantiated: its name, its factory function,
// and the options it supports.
type registrationInfo struct {
	name    string
	creator Creator
	options map[string]Option
}

var (
	deviceMap = map[string]registrationInfo{}
	// DeviceIDFile is the path to output file used to associate
	// device IDs with device configuration.
	DeviceIDFile string
)

// Register registers a function that can create a new Device
// of the given name.
func Register(name string, creator Creator, options map[string]Option) {
	deviceMap[name] = registrationInfo{
		name:    name,
		creator: creator,
		options: options,
	}
}

// Unregister removes a device from the registry.
func Unregister(name string) {
	delete(deviceMap, name)
}

// Registered returns a list of registered device names.
func Registered() (keys []string) {
	for k := range deviceMap {
		keys = append(keys, k)
	}
	return
}

// Create takes a config map, sanitizes the provided config, and returns
// a Device.
func Create(name string, config map[string]string) (Device, error) {
	registrationInfo, ok := deviceMap[name]
	if !ok {
		return nil, fmt.Errorf("Device '%v' not found", name)
	}
	sanitizedConfig, err := SanitizedOptions(registrationInfo.options, config)
	if err != nil {
		return nil, err
	}
	return registrationInfo.creator(sanitizedConfig)
}

// OptionHelp returns the options and associated help strings of the
// specified device.
func OptionHelp(deviceName string) (map[string]string, error) {
	registrationInfo, ok := deviceMap[deviceName]
	if !ok {
		return nil, fmt.Errorf("Device '%v' not found", deviceName)
	}
	return helpDesc(registrationInfo.options), nil
}

// Info contains the running state of an instantiated device.
type Info struct {
	ID     string
	Config Config
	Device Device
}

func (i *Info) String() string {
	var options []string
	for k, v := range i.Config.Options {
		options = append(options, fmt.Sprintf("deviceoption: %s=%s", k, v))
	}
	optStr := strings.Join(options, ", ")
	return fmt.Sprintf("Device %s {device: %s, %s}", i.ID, i.Config.Device, optStr)
}

func createDevice(name string, options map[string]string) (*Info, error) {
	config := Config{
		Device:  name,
		Options: options,
	}
	d, err := Create(name, options)
	if err != nil {
		return nil, fmt.Errorf("Failed creating device '%v': %v",
			config.Device, err)
	}
	did, err := d.DeviceID()
	if err != nil {
		return nil, fmt.Errorf(
			"Error getting device ID from Device %s with options %v: %v", name, options, err)
	}
	return &Info{Config: config, ID: did, Device: d}, nil
}

// CreateDevices returns a list of Info from either a single target device or a config file.
func CreateDevices(name, configPath string, options map[string]string) ([]*Info, error) {
	var infos []*Info
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
		readConfigs, err := ReadConfigs(configPath)
		if err != nil {
			return nil, fmt.Errorf("Error reading config file: %v", err)
		}
		errStrs := []string{}
		for _, config := range readConfigs {
			info, err := createDevice(config.Device, config.Options)
			if err != nil {
				errStrs = append(errStrs, err.Error())
			} else {
				infos = append(infos, info)
			}
		}
		if len(errStrs) != 0 {
			return nil, fmt.Errorf(strings.Join(errStrs, "\n"))
		}
	}
	return infos, nil
}

// DumpDeviceIDs dumps devices to the output file.
func DumpDeviceIDs(devices []*Info, outputFile string) error {
	idToConfig := map[string]Config{}
	for _, info := range devices {
		idToConfig[info.ID] = info.Config
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
	err = os.Rename(f.Name(), outputFile)
	if err != nil {
		return fmt.Errorf("Error in os.Rename: %v", err)
	}
	return nil
}
