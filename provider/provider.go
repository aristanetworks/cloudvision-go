// Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package provider

import (
	"arista/types"
)

// A Provider "owns" certain entities.  There are providers for entities
// coming from different sources (from Sysdb, from Smash, from /proc, etc.).
// Providers typically run in their own Goroutine(s), e.g. to read from the
// socket from Sysdb or from the shared memory files for Smash.  Providers can
// be asked to stop.  They also have a method used to write an update back to
// the source (e.g. send a message to Sysdb or update a shared-memory file for
// Smash).  Some providers can be read-only (e.g. the provider exposing data
// from /proc).
type Provider interface {
	// Run() kicks off the provider.  This method does not return until Stop()
	// is invoked, and is thus usually invoked by doing `go provider.Run()'.
	Run()

	// Stop() asks the provider to stop executing and clean up any Goroutines
	// it has started and release any resources it had acquired.
	// The provider will then stop, asynchronously.
	Stop()

	// Write asks the provider to apply the updates carried by the given
	// Notification to its data source (e.g. by sending an update to Sysdb
	// or updating a Smash table, etc.).  The error is returned asynchronsouly.
	Write(notif types.Notification, result chan<- error)
}
