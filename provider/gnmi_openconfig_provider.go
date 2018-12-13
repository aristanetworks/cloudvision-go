// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

// A GNMIOpenConfigProvider interacts with a GNMI client.
type GNMIOpenConfigProvider interface {
	Provider

	// InitGNMI initializes the provider by a given gnmi notification channel.
	InitGNMIOpenConfig(client gnmi.GNMIClient)
}
