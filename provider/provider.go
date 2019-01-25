// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"context"
)

// A Provider "owns" some states on a target device streams out notifications on any
// changes to those states.
type Provider interface {

	// Run() kicks off the provider.  This method does not return until ctx
	// is cancelled or an error is encountered,
	// and is thus usually invoked by doing `go provider.Run()'.
	Run(ctx context.Context) error
}
