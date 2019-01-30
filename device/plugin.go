// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

// loadPlugins recursively loads all the plugin files with the suffix .so,
// starting at the given directory.
func loadPlugins(pluginDir string) error {
	if pluginDir == "" {
		return nil
	}
	return filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to access path %s: %s", path, err)
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			if _, err = plugin.Open(path); err != nil {
				return err
			}
		}
		return nil
	})
}
