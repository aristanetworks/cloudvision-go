// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

// A GNMIProvider has access to a gnmi client which can be used to either generate gnmi
// notifications or access written states.
type GNMIProvider interface {
	Provider

	// InitGNMI initializes the provider by a given gnmi client.
	InitGNMI(gnmiClient gnmi.GNMIClient)
}
