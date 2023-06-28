// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aristanetworks/cloudvision-go/provider"
	"google.golang.org/grpc"
)

// ManagedDeviceStatus contains the status of a managed device.
type ManagedDeviceStatus string

var (
	// StatusActive indicates that a device is active.
	StatusActive ManagedDeviceStatus = "DEVICE_STATUS_ACTIVE"
	// StatusInactive indicates that a device is inactive, should
	// still be tracked by CloudVision.
	StatusInactive ManagedDeviceStatus = " DEVICE_STATUS_INACTIVE"
	// StatusRemoved indicates that a device should no longer
	// be tracked by CloudVision.
	StatusRemoved ManagedDeviceStatus = "DEVICE_STATUS_REMOVED"
)

// badConfigError is an error that is impossible to fix without user intervention.
// This error is use to prevent retrying runs of the device that will never work.
type badConfigError struct {
	error
}

// NewBadConfigError returns a BadConfigError
func NewBadConfigError(err error) error {
	return badConfigError{err}
}

// NewBadConfigErrorf returns a BadConfigError
func NewBadConfigErrorf(msg string, args ...interface{}) error {
	return badConfigError{fmt.Errorf(msg, args...)}
}

// IsBadConfigError check if given errors is a bad configuration error.
func IsBadConfigError(err error) bool {
	var b badConfigError
	return errors.As(err, &b)
}

// A Device knows how to interact with a specific device.
type Device interface {
	Alive(ctx context.Context) (bool, error)
	DeviceID(ctx context.Context) (string, error)
	Providers() ([]provider.Provider, error)
	// Type should return the type of the device. The returned
	// values should be one of the constants defined for the purpose
	// in the cvclient/v2 package such as VirtualSwitch etc.
	// If this method returns an empty string, a default value
	// (NetworkElement) is used.
	Type() string
	// IPAddr should return the management IP address of the device.
	// Return "" if this is not known.
	IPAddr(ctx context.Context) (string, error)
}

// A Manager manages a device inventory, adding and deleting
// devices as appropriate.
type Manager interface {
	Device
	Manage(ctx context.Context, inventory Inventory) error
}

// Creator returns a new instance of a Device.
type Creator = func(context.Context, map[string]string, provider.Monitor) (Device, error)

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
func newDevice(ctx context.Context, config *Config, monitor provider.Monitor) (Device, error) {
	registrationInfo, ok := deviceMap[config.Device]
	if !ok {
		return nil, NewBadConfigErrorf("Device '%v' not found", config.Device)
	}
	sanitizedConfig, err := SanitizedOptions(registrationInfo.options, config.Options)
	if err != nil {
		return nil, err
	}
	return registrationInfo.creator(ctx, sanitizedConfig, monitor)
}

// NewDeviceInfo takes a device config, creates the device, and returns an device Info.
func NewDeviceInfo(ctx context.Context, config *Config, monitor provider.Monitor) (*Info, error) {
	d, err := newDevice(ctx, config, monitor)
	if err != nil {
		return nil, fmt.Errorf("Failed creating device '%v': %w", config.Device, err)
	}
	did, err := d.DeviceID(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting device ID from Device %s: %w", config.Device, err)
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
	Status  ManagedDeviceStatus
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
	DeviceID   string
	Standalone bool
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

	// Receives the grpcServer connection and returns a CredentialResolver.
	CredResolverCreator func(conn *grpc.ClientConn) (CredentialResolver, error)
}
