// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmieos

import (
	"arista/provider"
	pgnmi "arista/provider/gnmi"
	"arista/schema"
	"arista/types"
	"context"
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type gnmieos struct {
	provider.ReadOnly
	prov        provider.GNMIProvider
	notifChan   chan<- types.Notification
	gnmiChan    chan *gnmi.Notification
	initialized bool
	ready       chan struct{}
}

func (g *gnmieos) Init(s *schema.Schema, root types.Entity,
	notification chan<- types.Notification) {
	g.notifChan = notification
	g.gnmiChan = make(chan *gnmi.Notification)
	g.initialized = true
}

func (g *gnmieos) Run(ctx context.Context) error {
	if !g.initialized {
		return fmt.Errorf("Provider is uninitialized")
	}

	g.prov.InitGNMI(g.gnmiChan)
	close(g.ready)
	go func() {
		for {
			notif := <-g.gnmiChan
			pgnmi.EmitNotif(notif, g.notifChan)
		}
	}()
	err := g.prov.Run(ctx)
	return err
}

func (g *gnmieos) WaitForNotification() {
	<-g.ready
}

// NewGNMIEOSProvider takes in a GNMIProvider and returns the same
// provider, converted to an EOSProvider
func NewGNMIEOSProvider(gp provider.GNMIProvider) provider.EOSProvider {
	return &gnmieos{
		prov:  gp,
		ready: make(chan struct{}),
	}
}
