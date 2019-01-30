// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

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
