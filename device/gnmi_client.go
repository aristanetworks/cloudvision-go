// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// gNMIClientWrapper wraps a GNMIClient and updates the passed context
// to reflect whether the gNMI server should perform OpenConfig type-
// checking.
type gNMIClientWrapper struct {
	client     gnmi.GNMIClient
	deviceID   string
	openConfig bool
}

func (g *gNMIClientWrapper) updatedContext(ctx context.Context) context.Context {
	oc := "false"
	if g.openConfig {
		oc = "true"
	}
	return metadata.AppendToOutgoingContext(ctx,
		deviceIDMetadata, g.deviceID,
		openConfigMetadata, oc)
}

func (g *gNMIClientWrapper) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return g.client.Capabilities(g.updatedContext(ctx), in, opts...)
}

func (g *gNMIClientWrapper) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return g.client.Get(g.updatedContext(ctx), in, opts...)
}

func (g *gNMIClientWrapper) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	return g.client.Set(g.updatedContext(ctx), in, opts...)
}

func (g *gNMIClientWrapper) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return g.client.Subscribe(g.updatedContext(ctx), opts...)
}

func newGNMIClientWrapper(client gnmi.GNMIClient, deviceID string,
	openConfig bool) *gNMIClientWrapper {
	return &gNMIClientWrapper{
		client:     client,
		deviceID:   deviceID,
		openConfig: openConfig,
	}
}
