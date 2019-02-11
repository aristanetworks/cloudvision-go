// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"fmt"

	"github.com/aristanetworks/cloudvision-go/provider"

	"github.com/aristanetworks/glog"
	agnmi "github.com/aristanetworks/goarista/gnmi"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type gnmiProvider struct {
	cfg         *agnmi.Config
	paths       []string
	inClient    gnmi.GNMIClient
	outClient   gnmi.GNMIClient
	initialized bool
}

func (p *gnmiProvider) InitGNMI(client gnmi.GNMIClient) {
	p.outClient = client
	p.initialized = true
}

func (p *gnmiProvider) OpenConfig() bool {
	return false
}

func (p *gnmiProvider) Run(ctx context.Context) error {
	if !p.initialized {
		return fmt.Errorf("provider is uninitialized")
	}
	respChan := make(chan *gnmi.SubscribeResponse)
	errChan := make(chan error)
	ctx = agnmi.NewContext(ctx, p.cfg)

	subscribeOptions := &agnmi.SubscribeOptions{
		Mode:       "stream",
		StreamMode: "target_defined",
		Paths:      agnmi.SplitPaths(p.paths),
	}
	go agnmi.Subscribe(ctx, p.inClient, subscribeOptions, respChan, errChan)
	for {
		select {
		case <-ctx.Done():
			return nil
		case response := <-respChan:
			switch resp := response.Response.(type) {
			case *gnmi.SubscribeResponse_Error:
				// Not sure if this is recoverable so it doesn't return and hope things get better
				glog.Errorf("gNMI SubscribeResponse Error: %v", resp.Error.Message)
			case *gnmi.SubscribeResponse_SyncResponse:
				if !resp.SyncResponse {
					glog.Errorf("gNMI sync failed")
				}
			case *gnmi.SubscribeResponse_Update:
				// One SetRequest per update:
				p.outClient.Set(ctx, &gnmi.SetRequest{Update: resp.Update.Update})
			}
		case err := <-errChan:
			return fmt.Errorf("Error from gNMI connection: %v", err)
		}
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client gnmi.GNMIClient, cfg *agnmi.Config,
	paths []string) provider.GNMIProvider {
	return &gnmiProvider{
		inClient: client,
		cfg:      cfg,
		paths:    paths,
	}
}
