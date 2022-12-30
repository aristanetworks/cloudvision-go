// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config represents a single device configuration.
type Config struct {
	Name     string            `yaml:"Name,omitempty"`
	Device   string            `yaml:"Device,omitempty"`
	NoStream bool              `yaml:"NoStream,omitempty"`
	Options  map[string]string `yaml:"Options,omitempty"`

	deleted bool
}

// NewDeletedConfig creates a config that indicates a deleted device config.
func NewDeletedConfig(name string) *Config {
	return &Config{
		Name:    name,
		deleted: true,
	}
}

// IsDeleted returns true if the config is marked as deleted.
func (c *Config) IsDeleted() bool {
	return c.deleted
}

func (c *Config) validate() error {
	if c.Device == "" {
		return fmt.Errorf("Device in config cannot be empty")
	}
	return nil
}

// ReadConfigs generates device configs from the config file at the specified path.
func ReadConfigs(configPath string) ([]*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return readConfigsFromBytes(data)
}

func readConfigsFromBytes(data []byte) ([]*Config, error) {
	configs := []*Config{}
	err := yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		err := config.validate()
		if err != nil {
			return nil, err
		}
	}
	return configs, nil
}

// WriteConfigs writes a list of Config to the specified path.
func WriteConfigs(configPath string, configs []*Config) error {
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(&configs)
}
