// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"context"
	"sync"
)

// A Provider "owns" some states on a target device streams out notifications on any
// changes to those tates.
type Provider interface {

	// Run() kicks off the provider.  This method does not return until ctx
	// is cancelled or an error is encountered,
	// and is thus usually invoked by doing `go provider.Run()'.
	Run(ctx context.Context) error
}

// Run runs a single provider and returns a function that will stop it, wait for it to finish,
// and return any error returned by Run. Intended for use within tests.
func Run(p Provider) func() error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	go func() {
		err = p.Run(ctx)
		wg.Done()
	}()
	return func() error {
		cancel()
		wg.Wait()
		return err
	}
}
