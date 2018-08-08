// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"io/ioutil"
	"path"
	"plugin"
)

func LoadPlugins(pluginDir string) error {
	if pluginDir == "" {
		return nil
	}
	pluginFiles, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return err
	}
	for _, file := range pluginFiles {
		if _, err = plugin.Open(path.Join(pluginDir, file.Name())); err != nil {
			return err
		}
	}
	return nil
}
