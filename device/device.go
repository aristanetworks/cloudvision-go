// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"
	"strings"

	"github.com/aristanetworks/cloudvision-go/provider"
	"google.golang.org/grpc"
)

// A Device knows how to interact with a specific device.
type Device interface {
	Alive() (bool, error)
	DeviceID() (string, error)
	Providers() ([]provider.Provider, error)
	// Type should return the type of the device. The returned
	// values should be one of the constants defined for the purpose
	// in the cvclient/v2 package such as VirtualSwitch etc.
	// If this method returns an empty string, a default value
	// (NetworkElement) is used.
	Type() string
	// IPAddr should return the management IP address of the device.
	// Return "" if this is not known.
	IPAddr() string
}

// A Manager manages a device inventory, adding and deleting
// devices as appropriate.
type Manager interface {
	Device
	Manage(ctx context.Context, inventory Inventory) error
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

// newDevice takes a device config and returns a Device.
func newDevice(config *Config) (Device, error) {
	registrationInfo, ok := deviceMap[config.Device]
	if !ok {
		return nil, fmt.Errorf("Device '%v' not found", config.Device)
	}
	sanitizedConfig, err := SanitizedOptions(registrationInfo.options, config.Options)
	if err != nil {
		return nil, err
	}
	return registrationInfo.creator(sanitizedConfig)
}

// NewDeviceInfo takes a device config, creates the device, and returns an device Info.
func NewDeviceInfo(config *Config) (*Info, error) {
	d, err := newDevice(config)
	if err != nil {
		return nil, fmt.Errorf("Failed creating device '%v': %v", config.Device, err)
	}
	did, err := d.DeviceID()
	if err != nil {
		return nil, fmt.Errorf("Error getting device ID from Device %s: %w",
			config.Device, err)
	}
	return &Info{Device: d, ID: did, Config: config}, nil
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
	ID      string
	Context context.Context
	Device  Device
	Config  *Config
}

func (i *Info) String() string {
	template := "Device %s config:{%s}"
	if i.Config == nil {
		return fmt.Sprintf(template, i.ID, "")
	}
	var options []string
	for k, v := range i.Config.Options {
		options = append(options, fmt.Sprintf("deviceoption: %s=%s", k, v))
	}
	optStr := strings.Join(options, ", ")
	configStr := fmt.Sprintf("type: %s, %s", i.Config.Device, optStr)
	return fmt.Sprintf(template, i.ID, configStr)
}

// GRPCConnectorConfig used to pass configuration parameters to GRPCConnector
// interface
type GRPCConnectorConfig struct {
	DeviceID string
}

// GRPCConnector allows callers to supply one gRPC connection and
// to create another to be used by a device implementation
type GRPCConnector interface {
	Connect(ctx context.Context, conn *grpc.ClientConn,
		addr string, config GRPCConnectorConfig) (*grpc.ClientConn, error)
}

// defaultGRPCConnector default implementation of GRPCConnector interface
type defaultGRPCConnector struct {
}

// NewDefaultGRPCConnector return empty object
func NewDefaultGRPCConnector() GRPCConnector {
	return &defaultGRPCConnector{}
}

// Connect returns grpc connection
func (dgc *defaultGRPCConnector) Connect(ctx context.Context,
	conn *grpc.ClientConn, addr string, config GRPCConnectorConfig) (*grpc.ClientConn, error) {
	return conn, nil
}

// SensorConfig to store GRPCConnector config
type SensorConfig struct {
	Connector GRPCConnector
}
