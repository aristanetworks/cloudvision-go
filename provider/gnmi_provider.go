// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package provider

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

// A GNMIProvider emits updates as gNMI SetRequests.
type GNMIProvider interface {
	Provider

	// InitGNMI initializes the provider with a gNMI client.
	InitGNMI(client gnmi.GNMIClient)

	// OpenConfig indicates whether the provider wants OpenConfig
	// type-checking. Used only by v1 client.
	OpenConfig() bool

	// Origin of the YANG data model that this provider streams,
	// should be one of "arista", "fmp", "openconfig". This method
	// is used only by v2 client.
	Origin() string
}
