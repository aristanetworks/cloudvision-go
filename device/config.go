// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"os"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

// Config represents a single device configuration.
type Config struct {
	Name     string            `yaml:"Name,omitempty"`
	Device   string            `yaml:"Device,omitempty"`
	NoStream bool              `yaml:"NoStream,omitempty"`
	Options  map[string]string `yaml:"Options,omitempty"`
	LogLevel string            `yaml:"LogLevel,omitempty"`

	Credentials map[string]string `yaml:"Credentials,omitempty"`
	Enabled     bool              `yaml:"Enabled,omitempty"`

	deleted bool
	syncEnd bool
}

// NewConfig creates a new Config
func NewConfig(name string) *Config {
	return &Config{
		Name:        name,
		Options:     make(map[string]string),
		Credentials: make(map[string]string),
	}
}

// NewDeletedConfig creates a config that indicates a deleted device config.
func NewDeletedConfig(name string) *Config {
	return &Config{
		Name:    name,
		deleted: true,
	}
}

// NewSyncEndConfig creates a new config that represents the end of the sync phase
func NewSyncEndConfig() *Config {
	return &Config{
		syncEnd: true,
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

// Equal returns true if the two configs are config equal
func (c *Config) Equal(o *Config) bool {
	if c == nil || o == nil {
		return false
	}

	// TODO: Switch to using https://pkg.go.dev/maps@go1.21.5 when the go version is upgraded
	return reflect.DeepEqual(c.Options, o.Options) &&
		reflect.DeepEqual(c.Credentials, o.Credentials) &&
		c.Name == o.Name &&
		c.Device == o.Device &&
		c.NoStream == o.NoStream &&
		c.LogLevel == o.LogLevel &&
		c.Enabled == o.Enabled &&
		c.deleted == o.deleted &&
		c.syncEnd == o.syncEnd
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
