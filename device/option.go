// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"arista/flag"
	"bytes"
	"fmt"
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

// Return help string for a given set of options.
func help(options map[string]Option, optionType, name string) string {
	b := new(bytes.Buffer)
	hd := helpDesc(options)
	// Don't print out device separator if the device has no options.
	if len(hd) == 0 {
		return ""
	}
	flag.FormatOptions(b, "Help options for "+optionType+" '"+name+"':", hd)
	return b.String()
}
