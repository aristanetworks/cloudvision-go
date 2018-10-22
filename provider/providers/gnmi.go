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
	"fmt"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/gnmi"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

type gnmiProvider struct {
	provider.ReadOnly
	// Closed when we're done initialization
	ready chan struct{}

	client pb.GNMIClient
	cfg    *gnmi.Config
	paths  []string

	channel chan<- types.Notification
	isInit  bool
}

func (p *gnmiProvider) WaitForNotification() {
	<-p.ready
}

func (p *gnmiProvider) Init(s *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	p.channel = ch
	p.isInit = true
}

func (p *gnmiProvider) Run(ctx context.Context) error {
	if !p.isInit {
		return fmt.Errorf("provider is uninitialized")
	}
	respChan := make(chan *pb.SubscribeResponse)
	errChan := make(chan error)
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
		case <-ctx.Done():
			return nil
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
				GNMIEmitNotif(resp.Update, p.channel)
			}
		case err := <-errChan:
			return fmt.Errorf("Error from gNMI connection: %v", err)
		}
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client pb.GNMIClient, cfg *gnmi.Config, paths []string) provider.EOSProvider {
	return &gnmiProvider{
		ready:  make(chan struct{}),
		client: client,
		cfg:    cfg,
		paths:  paths,
	}
}
