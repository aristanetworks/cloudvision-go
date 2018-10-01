// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/provider"
	"arista/schema"
	"arista/types"
	"context"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/gnmi"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

type gnmiProvider struct {
	provider.ReadOnly
	// Closed when we're done initialization
	ready chan struct{}
	// Closed when we want to stop Run()
	done chan struct{}

	client   pb.GNMIClient
	cfg      *gnmi.Config
	paths    []string
	typeDefs *schema.Schema
}

func (p *gnmiProvider) WaitForNotification() {
	<-p.ready
}

func (p *gnmiProvider) Stop() {
	<-p.ready
	close(p.done)
}

func (p *gnmiProvider) Run(s *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	p.typeDefs = s
	respChan := make(chan *pb.SubscribeResponse)
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, p.cfg)

	subscribeOptions := &gnmi.SubscribeOptions{
		Mode:       "stream",
		StreamMode: "target_defined",
		Paths:      gnmi.SplitPaths(p.paths),
	}
	go gnmi.Subscribe(ctx, p.client, subscribeOptions, respChan, errChan)
	close(p.ready)
	for {
		select {
		case <-p.done:
			return
		case response := <-respChan:
			switch resp := response.Response.(type) {
			case *pb.SubscribeResponse_Error:
				// Not sure if this is recoverable so it doesn't return and hope things get better
				glog.Errorf("gNMI SubscribeResponse Error: %v", resp.Error.Message)
			case *pb.SubscribeResponse_SyncResponse:
				if !resp.SyncResponse {
					glog.Errorf("gNMI sync failed")
				}
			case *pb.SubscribeResponse_Update:
				GNMIEmitNotif(resp.Update, ch)
			}
		case err := <-errChan:
			glog.Errorf("Error from gNMI connection: %v", err)
			return
		}
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client pb.GNMIClient, cfg *gnmi.Config, paths []string) provider.Provider {
	return &gnmiProvider{
		ready:  make(chan struct{}),
		done:   make(chan struct{}),
		client: client,
		cfg:    cfg,
		paths:  paths,
	}
}
