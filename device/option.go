// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"net"
	"time"
)

// Option defines a command-line option accepted by a device.
type Option struct {
	Description string
	Default     string
	Required    bool
}

// sanitizedOptions takes the map of device option keys and values
// passed in at the command line and checks it against the device
// or manager's exported list of accepted options, returning an
// error if there are inappropriate or missing options.
func sanitizedOptions(options map[string]Option,
	config map[string]string) (map[string]string, error) {
	sopt := make(map[string]string)

	// Check whether the user gave us bad options.
	for k, v := range config {
		_, ok := options[k]
		if !ok {
			return nil, fmt.Errorf("Bad option '%s'", k)
		}
		sopt[k] = v
	}

	// Check that all required options were specified, and fill in
	// any others with defaults.
	for k, v := range options {
		_, found := sopt[k]
		if v.Required && !found {
			return nil, fmt.Errorf("Required option '%s' not provided", k)
		}
		if !found {
			sopt[k] = v.Default
		}
	}

	return sopt, nil
}

// Create map of option key to description.
func helpDesc(options map[string]Option) map[string]string {
	hd := make(map[string]string)

	for k, v := range options {
		desc := v.Description
		// Add default if there's a non-empty one.
		if v.Default != "" {
			desc = desc + " (default " + v.Default + ")"
		}
		hd[k] = desc
	}
	return hd
}

// GetStringOption returns the option specified by optionName as a
// string.
func GetStringOption(optionName string,
	options map[string]string) (string, error) {
	o, ok := options[optionName]
	if !ok {
		return "", fmt.Errorf("No option '%s'", optionName)
	}
	return o, nil
}

// GetAddressOption returns the option specified by optionName as a
// validated IP address or hostname.
func GetAddressOption(optionName string,
	options map[string]string) (string, error) {
	addr, ok := options[optionName]
	if !ok {
		return "", fmt.Errorf("No option '%s'", optionName)
	}

	// Validate IP
	ip := net.ParseIP(addr)
	if ip != nil {
		return ip.String(), nil
	}

	// Try for hostname if it's not an IP
	addrs, err := net.LookupIP(addr)
	if err != nil {
		return "", err
	}
	return addrs[0].String(), nil
}

// GetDurationOption returns the option specified by optionName as a
// time.Duration.
func GetDurationOption(optionName string,
	options map[string]string) (time.Duration, error) {
	o, ok := options[optionName]
	if !ok {
		return 0, fmt.Errorf("No option '%s'", optionName)
	}
	return time.ParseDuration(o)
}
