// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package v1 implements the version v1 client for communicating with
// CloudVision.
package v1

import (
	"context"
	"strconv"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/version"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type v1client struct {
	gnmiClient   gnmi.GNMIClient // underlying raw GNMI client
	deviceID     string
	isMgmtSystem bool
	openConfig   bool
	typeCheck    bool
}

func (c *v1client) ctxWithDeviceMetadata(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		deviceIDMetadata, c.deviceID,
		openConfigMetadata, strconv.FormatBool(c.openConfig),
		typeCheckMetadata, strconv.FormatBool(c.typeCheck))
}

func (c *v1client) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	ctx = c.ctxWithDeviceMetadata(ctx)
	return c.gnmiClient.Capabilities(ctx, in, opts...)
}

func (c *v1client) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	ctx = c.ctxWithDeviceMetadata(ctx)
	return c.gnmiClient.Get(ctx, in, opts...)
}

func (c *v1client) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	ctx = c.ctxWithDeviceMetadata(ctx)
	return c.gnmiClient.Set(ctx, in, opts...)
}

func (c *v1client) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	ctx = c.ctxWithDeviceMetadata(ctx)
	return c.gnmiClient.Subscribe(ctx, opts...)
}

func (c *v1client) SendDeviceMetadata(ctx context.Context) error {
	ctx = metadata.AppendToOutgoingContext(ctx, collectorVersionMetadata, version.Version)
	if c.isMgmtSystem {
		// ManagementSystem is a system managing other devices which itself
		// shouldn't be treated as an actual streaming device in CloudVision.
		ctx = metadata.AppendToOutgoingContext(ctx,
			deviceTypeMetadata, "managementSystem")
	} else {
		// Target is an ordinary device streaming to CloudVision.
		ctx = metadata.AppendToOutgoingContext(ctx,
			deviceTypeMetadata, "target")
	}
	_, err := c.Set(ctx, &gnmi.SetRequest{})
	return err
}

func (c *v1client) SendHeartbeat(ctx context.Context, alive bool) error {
	if !alive {
		return nil
	}
	ctx = metadata.AppendToOutgoingContext(ctx, deviceLivenessMetadata, "true")
	_, err := c.Set(ctx, &gnmi.SetRequest{})
	return err
}

// NewV1Client returns a new v1client object.
func NewV1Client(gc gnmi.GNMIClient, deviceID string, isMgmtSystem bool) cvclient.CVClient {
	return &v1client{
		gnmiClient:   gc,
		deviceID:     deviceID,
		isMgmtSystem: isMgmtSystem,
	}
}

func (c *v1client) ForProvider(p provider.GNMIProvider) cvclient.CVClient {
	var openConfig, typeCheck bool
	openConfig = p.OpenConfig()
	typeCheck = openConfig
	// special case for Gnmi provider
	if _, ok := p.(*pgnmi.Gnmi); ok {
		typeCheck = false
	}
	return &v1client{
		gnmiClient:   c.gnmiClient,
		deviceID:     c.deviceID,
		isMgmtSystem: c.isMgmtSystem,
		openConfig:   openConfig,
		typeCheck:    typeCheck,
	}
}
