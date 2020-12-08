// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"fmt"

	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"

	agnmi "github.com/aristanetworks/goarista/gnmi"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// A Gnmi connects to a gNMI server at a target device
// and emits updates as gNMI SetRequests.
type Gnmi struct {
	cfg         *agnmi.Config
	paths       []string
	inClient    gnmi.GNMIClient
	outClient   gnmi.GNMIClient
	initialized bool
}

// InitGNMI initializes the provider with a gNMI client.
func (p *Gnmi) InitGNMI(client gnmi.GNMIClient) {
	p.outClient = client
	p.initialized = true
}

// OpenConfig indicates whether the provider wants OpenConfig type-checking.
func (p *Gnmi) OpenConfig() bool {
	return true
}

// Origin indicates that the provider streams to an OpenConfig model.
func (p *Gnmi) Origin() string {
	return "openconfig"
}

// Run kicks off the provider.
func (p *Gnmi) Run(ctx context.Context) error {
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
		case response, ok := <-respChan:
			if !ok {
				// channel closed, we received an io.EOF from the server.
				log.Log(p).Infof("gNMI target closed subscription")
				return fmt.Errorf("gNMI target closed subscription")
			}
			switch resp := response.Response.(type) {
			case *gnmi.SubscribeResponse_Error:
				// Not sure if this is recoverable so it doesn't return and hope things get better
				log.Log(p).Infof(
					"gNMI SubscribeResponse Error: %v", resp.Error.Message)
			case *gnmi.SubscribeResponse_SyncResponse:
				if resp.SyncResponse {
					log.Log(p).Debug("gNMI sync_response")
				} else {
					log.Log(p).Infof("gNMI sync failed")
				}
			case *gnmi.SubscribeResponse_Update:
				// One SetRequest per update:
				sr := &gnmi.SetRequest{
					Prefix: resp.Update.Prefix,
					Update: resp.Update.Update,
					Delete: resp.Update.Delete,
				}
				log.Log(p).Debugf("SetRequest: %+v", sr)
				if _, err := p.outClient.Set(ctx, sr); err != nil {
					log.Log(p).Infof("Error on Set: %v", err)
				}
			}
		case err := <-errChan:
			return fmt.Errorf("Error from gNMI connection: %v", err)
		}
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client gnmi.GNMIClient, cfg *agnmi.Config,
	paths []string) provider.GNMIProvider {
	return &Gnmi{
		inClient: client,
		cfg:      cfg,
		paths:    paths,
	}
}
