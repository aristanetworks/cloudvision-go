// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aristanetworks/cloudvision-go/provider"
	"golang.org/x/sync/errgroup"

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

	monitor provider.Monitor
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

func (p *Gnmi) subscribeAndSet(ctx context.Context,
	subscribeOptions *agnmi.SubscribeOptions) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = agnmi.NewContext(ctx, p.cfg)
	respCh := make(chan *gnmi.SubscribeResponse)
	errg, ctx := errgroup.WithContext(ctx)

	// producer: subscribe to target
	errg.Go(func() error {
		return agnmi.SubscribeErr(ctx, p.inClient, subscribeOptions, respCh)
	})

	// consumer: receive SubscribeResponses, send SetRequests
	errg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case response, ok := <-respCh:
				if !ok {
					return io.EOF
				}
				switch resp := response.Response.(type) {
				case *gnmi.SubscribeResponse_Error:
					// Not sure if this is recoverable; keep going and
					// hope things get better.
					p.monitor.Infof("gNMI SubscribeResponse Error: %v",
						resp.Error.Message) // nolint: staticcheck
				case *gnmi.SubscribeResponse_SyncResponse:
					if resp.SyncResponse {
						p.monitor.Infof("gNMI sync_response")
					} else {
						p.monitor.Infof("gNMI sync failed")
					}
				case *gnmi.SubscribeResponse_Update:
					// One SetRequest per update:
					sr := &gnmi.SetRequest{
						Prefix: resp.Update.Prefix,
						Update: resp.Update.Update,
						Delete: resp.Update.Delete,
					}
					p.monitor.Debugf("SetRequest: %+v", sr)
					if _, err := p.outClient.Set(ctx, sr); err != nil {
						p.monitor.Infof("Error on Set: %v", err)
					}
				}
			}
		}
	})

	return errg.Wait()
}

// Run kicks off the provider.
func (p *Gnmi) Run(ctx context.Context) error {
	if !p.initialized {
		return fmt.Errorf("provider is uninitialized")
	}

	subscribeOptions := &agnmi.SubscribeOptions{
		Mode:       "stream",
		StreamMode: "target_defined",
		Paths:      agnmi.SplitPaths(p.paths),
	}

	p.monitor.Debugf("gNMI subscribeOptions: %+v, config: "+
		"{Addr:%s, CAFile:%s, CertFile:%s, Username:%s, TLS:%t, Compression:%s}",
		subscribeOptions, p.cfg.Addr, p.cfg.CAFile, p.cfg.CertFile, p.cfg.Username,
		p.cfg.TLS, p.cfg.Compression)

	// Initialize retry timer, setting it to one nanosecond so we can
	// go straight into the retry loop.
	backoffTimer := provider.NewBackoffTimer()

	// Retry loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-backoffTimer.Wait():
			// Subscribe, hopefully forever.
			err := p.subscribeAndSet(ctx, subscribeOptions)

			// Subscribe failed, schedule retry with backoff.
			// This is done before logging the error so we can log a precise retry delay.
			curBackoff := backoffTimer.Backoff()

			if !errors.Is(err, context.Canceled) {
				p.monitor.Infof("gNMI subscription failed, retrying in %v. Err: %v",
					curBackoff, err)
			}
		}
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client gnmi.GNMIClient, cfg *agnmi.Config,
	paths []string, monitor provider.Monitor) provider.GNMIProvider {

	g := &Gnmi{
		inClient: client,
		cfg:      cfg,
		paths:    paths,
	}
	g.monitor = monitor
	return g
}
