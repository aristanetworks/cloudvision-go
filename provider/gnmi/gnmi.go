// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/sirupsen/logrus"
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

	logger *logrus.Entry
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
		p.logger.Infof("gNMI subscribeOptions: %+v, config: %+v",
			subscribeOptions, p.cfg)
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
					p.logger.Info("gNMI target closed subscription")
					return io.EOF
				}
				switch resp := response.Response.(type) {
				case *gnmi.SubscribeResponse_Error:
					// Not sure if this is recoverable; keep going and
					// hope things get better.
					p.logger.Infof("gNMI SubscribeResponse Error: %v",
						resp.Error.Message) // nolint: staticcheck
				case *gnmi.SubscribeResponse_SyncResponse:
					if resp.SyncResponse {
						p.logger.Info("gNMI sync_response")
					} else {
						p.logger.Infof("gNMI sync failed")
					}
				case *gnmi.SubscribeResponse_Update:
					// One SetRequest per update:
					sr := &gnmi.SetRequest{
						Prefix: resp.Update.Prefix,
						Update: resp.Update.Update,
						Delete: resp.Update.Delete,
					}
					p.logger.Debugf("SetRequest: %+v", sr)
					if _, err := p.outClient.Set(ctx, sr); err != nil {
						p.logger.Infof("Error on Set: %v", err)
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

	// Initialize retry timer, setting it to one nanosecond so we can
	// go straight into the retry loop.
	retryTimer := time.NewTimer(time.Nanosecond)

	// Retry loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-retryTimer.C:
			// Subscribe, hopefully forever.
			if err := p.subscribeAndSet(ctx, subscribeOptions); err != nil {
				if err != io.EOF {
					return err
				}
			}
			// The server closed the connection with no error. Schedule a retry.
			p.logger.Info("gNMI subscription disconnected, retrying...")
			retryTimer.Reset(5 * time.Second)
		}
	}
}

type gnmiProviderOptions struct {
	deviceID string
}

// ProviderOption is a function to set options for the GNMIProvider.
type ProviderOption func(ops *gnmiProviderOptions)

// WithDeviceID idenfifies the source of the data.
// This is used for logging and errors only.
func WithDeviceID(identifier string) ProviderOption {
	return func(o *gnmiProviderOptions) {
		o.deviceID = identifier
	}
}

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client gnmi.GNMIClient, cfg *agnmi.Config,
	paths []string, options ...ProviderOption) provider.GNMIProvider {
	var opts gnmiProviderOptions
	for _, o := range options {
		o(&opts)
	}

	g := &Gnmi{
		inClient: client,
		cfg:      cfg,
		paths:    paths,
	}
	g.logger = log.Log(g)
	if len(opts.deviceID) > 0 {
		g.logger = g.logger.WithField("deviceID", opts.deviceID)
	} else {
		g.logger = g.logger.WithField("deviceAddr", cfg.Addr)
	}
	return g
}
