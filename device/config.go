// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config represents a single device configuration.
type Config struct {
	Address string            `yaml:"Address,omitempty"`
	Device  string            `yaml:"Device,omitempty"`
	Options map[string]string `yaml:"Options,omitempty"`
}

func (c *Config) String() string {
	var fields []string
	fields = append(fields, fmt.Sprintf("device: %s", c.Device))
	fields = append(fields, fmt.Sprintf("address: %s", c.Address))
	for k, v := range c.Options {
		fields = append(fields, fmt.Sprintf("deviceoption: %s=%s", k, v))
	}
	return strings.Join(fields, ", ")
}

// ReadConfigs reads the config file at the specified path,
// optionally extracting just the specified devices.
func ReadConfigs(configPath string) ([]Config, error) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return readConfigsFromBytes(yamlFile)
}

func readConfigsFromBytes(yamlFile []byte) ([]Config, error) {
	configs := []Config{}
	err := yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		if config.Device == "" {
			return nil, errors.Errorf("Device must be specified")
		}
		if config.Address == "" {
			return nil, errors.Errorf("Address must be specified")
		}
	}
	return configs, nil
}
