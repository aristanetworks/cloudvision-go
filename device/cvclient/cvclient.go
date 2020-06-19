// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package cvclient defines an interface for connecting to and
// communicating with CloudVision.
package cvclient

import (
	"context"

	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// CVClient defines the interface that needs to be implemented for communicating
// with CloudVision.
type CVClient interface {
	// Clients will usually communicate with CloudVision using gNMI
	// but can use some other protocol as well, all they need to do
	// is to implement the GNMIClient interface.
	gnmi.GNMIClient
	// SendDeviceMetadata sends metadata for the device to CV.
	SendDeviceMetadata(ctx context.Context) error
	// SendPeriodicUpdate sends information about device liveness to CV.
	SendHeartbeat(ctx context.Context, alive bool) error
	// ForProvider returns a new client instance for the same device but
	// customized for a specific provider.
	ForProvider(p provider.GNMIProvider) CVClient
}
