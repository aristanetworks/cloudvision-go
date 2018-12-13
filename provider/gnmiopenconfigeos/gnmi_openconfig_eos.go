// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmiopenconfigeos

import (
	"arista/provider"
	pgnmi "arista/provider/gnmi"
	"arista/schema"
	"arista/types"
	"context"
	"fmt"

	"github.com/aristanetworks/glog"
	"github.com/openconfig/gnmi/proto/gnmi"
)

type gnmiopenconfigeos struct {
	provider.ReadOnly
	prov        provider.GNMIOpenConfigProvider
	errc        chan error
	server      gnmi.GNMIServer
	client      gnmi.GNMIClient
	ready       chan struct{}
	initialized bool
	notifChan   chan<- types.Notification
	yangPaths   []string
}

func (g *gnmiopenconfigeos) Init(s *schema.Schema, root types.Entity,
	notification chan<- types.Notification) {
	g.notifChan = notification
	g.initialized = true
}

func (g *gnmiopenconfigeos) Run(ctx context.Context) error {
	if !g.initialized {
		return fmt.Errorf("Provider is uninitialized")
	}

	var err error
	ctx, g.server, err = pgnmi.Server(ctx, g.notifChan, g.errc, g.yangPaths)
	if err != nil {
		glog.Errorf("Error in creating GNMIServer: %v", err)
	}

	g.client = pgnmi.Client(g.server)
	g.prov.InitGNMIOpenConfig(g.client)
	close(g.ready)
	err = g.prov.Run(ctx)
	return err
}

func (g *gnmiopenconfigeos) WaitForNotification() {
	<-g.ready
}

// NewGNMIOpenConfigEOSProvider takes in a GNMIOpenConfigProvider and returns the same
// provider, converted to an EOSProvider.
func NewGNMIOpenConfigEOSProvider(gp provider.GNMIOpenConfigProvider,
	yangPaths []string) provider.EOSProvider {
	return &gnmiopenconfigeos{
		prov:      gp,
		errc:      make(chan error),
		ready:     make(chan struct{}),
		yangPaths: yangPaths,
	}
}
