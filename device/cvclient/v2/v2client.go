// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package v2 implements the v2 protocol for communicating with CloudVision.
package v2

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/version"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

const (
	// NetworkElement is a generic network device.
	NetworkElement = "DEVICE_TYPE_NETWORK_ELEMENT"
	// DeviceManager is a manager of network devices.
	DeviceManager = "DEVICE_TYPE_DEVICE_MANAGER"
	// WirelessAP is a wireless access point device.
	WirelessAP = "DEVICE_TYPE_WIRELESS_AP"
)

// versionString describes the version of the collector, the SDK, and the go runtime.
// This is sent as device metadata to CV and can be used for debugging.
var versionString = fmt.Sprintf("Collector version: %s, SDK version: %s, Go version: %s",
	version.CollectorVersion, version.Version, runtime.Version())

type v2Client struct {
	gnmiClient gnmi.GNMIClient // underlying raw GNMI client
	deviceID   string
	deviceType string
	origin     string
	device     device.Device
}

// setTargetAndOrigin sets target and origin fields in a GNMI path based on values in c.
func (c *v2Client) setTargetAndOrigin(p *gnmi.Path) *gnmi.Path {
	if p == nil {
		p = &gnmi.Path{}
	}
	p.Target = c.deviceID
	if p.Origin == "" {
		// set default origin if not set by provider.
		p.Origin = c.origin
	}
	return p
}

func (c *v2Client) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return c.gnmiClient.Capabilities(ctx, in, opts...)
}

func (c *v2Client) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	in.Prefix = c.setTargetAndOrigin(in.Prefix)
	return c.gnmiClient.Get(ctx, in, opts...)
}

func (c *v2Client) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	in.Prefix = c.setTargetAndOrigin(in.Prefix)
	return c.gnmiClient.Set(ctx, in, opts...)
}

func (c *v2Client) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	// TODO: intercept subscribe requests and add target.
	return c.gnmiClient.Subscribe(ctx, opts...)
}

func metadataPrefix() *gnmi.Path {
	prefix := pgnmi.Path("device-metadata", "state", "metadata")
	return prefix
}

func (c *v2Client) metadataRequest() *gnmi.SetRequest {
	u := []*gnmi.Update{
		pgnmi.Update(pgnmi.Path("type"), pgnmi.Strval(c.deviceType)),
		pgnmi.Update(pgnmi.Path("collector-version"), pgnmi.Strval(versionString)),
	}
	if ip := c.device.IPAddr(); ip != "" {
		u = append(u,
			pgnmi.Update(pgnmi.Path("ip-addr"), pgnmi.Strval(ip)),
		)
	}
	// TODO: list of managed devices.
	return &gnmi.SetRequest{
		Prefix: metadataPrefix(),
		Update: u,
	}
}

func (c *v2Client) SendDeviceMetadata(ctx context.Context) error {
	req := c.metadataRequest()
	_, err := c.Set(ctx, req)
	return err
}

func (c *v2Client) heartbeatRequest() *gnmi.SetRequest {
	now := time.Now()
	nanos := now.Unix()*int64(time.Second) + now.UnixNano()
	u := []*gnmi.Update{pgnmi.Update(pgnmi.Path("last-seen"), pgnmi.Intval(nanos))}
	return &gnmi.SetRequest{
		Prefix: metadataPrefix(),
		Update: u,
	}
}

func (c *v2Client) SendHeartbeat(ctx context.Context, alive bool) error {
	if !alive {
		return nil
	}
	req := c.heartbeatRequest()
	_, err := c.Set(ctx, req)
	return err
}

// NewV2Client returns a new client object for communication
// with CV using the v2 protocol.
func NewV2Client(gc gnmi.GNMIClient, dev device.Device) cvclient.CVClient {
	deviceType := NetworkElement
	if _, ok := dev.(device.Manager); ok {
		deviceType = DeviceManager
	} else {
		if dev.Type() != "" {
			deviceType = dev.Type()
		}
	}
	deviceID, _ := dev.DeviceID()
	return &v2Client{
		gnmiClient: gc,
		deviceID:   deviceID,
		deviceType: deviceType,
		origin:     "arista",
		device:     dev,
	}
}

func (c *v2Client) ForProvider(p provider.GNMIProvider) cvclient.CVClient {
	return &v2Client{
		gnmiClient: c.gnmiClient,
		deviceID:   c.deviceID,
		deviceType: c.deviceType,
		origin:     p.Origin(),
	}
}
