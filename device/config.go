// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config represents a single device configuration.
type Config struct {
	Device  string            `yaml:"Device,omitempty"`
	Options map[string]string `yaml:"Options,omitempty"`
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
	}
	return configs, nil
}
