// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package provider

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

// A GNMIProvider emits gNMI notifications.
type GNMIProvider interface {
	Provider

	// InitGNMI initializes the provider by a given gnmi notification channel.
	InitGNMI(chan<- *gnmi.Notification)
}
