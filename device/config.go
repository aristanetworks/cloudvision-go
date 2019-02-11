// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

// Config represents a single device configuration.
type Config struct {
	Device  string            `yaml:"Device,omitempty"`
	Options map[string]string `yaml:"Options,omitempty"`
}
