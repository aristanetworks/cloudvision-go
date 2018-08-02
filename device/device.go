// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package device

import (
	"arista/provider"
)

// A Device knows how to interact with a specific device
type Device interface {
	Name() string
	CheckAlive() bool
	DeviceID() string
	Providers() []provider.Provider
}
