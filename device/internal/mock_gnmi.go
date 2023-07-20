// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package internal

import (
	"context"
	"fmt"
	"io"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

// MockClientStream is used as a mock gnmi stream
type MockClientStream struct {
	SubReq  chan *gnmi.SubscribeRequest
	SubResp chan *gnmi.SubscribeResponse
	ErrC    chan error
	grpc.ClientStream
}

// Send sends a subscribe request
func (mcs *MockClientStream) Send(req *gnmi.SubscribeRequest) error {
	mcs.SubReq <- req
	return nil
}

// Recv reads from the stream
func (mcs *MockClientStream) Recv() (*gnmi.SubscribeResponse, error) {
	select {
	case err, ok := <-mcs.ErrC:
		if !ok {
			return nil, io.EOF
		}
		return nil, err
	case r, ok := <-mcs.SubResp:
		if !ok {
			return nil, io.EOF
		}
		return r, nil
	}
}

var _ gnmi.GNMI_SubscribeClient = (*MockClientStream)(nil)

// MockClient is a mock gnmi client
type MockClient struct {
	SubscribeStream chan *MockClientStream

	SetReq  chan *gnmi.SetRequest
	SetResp chan *gnmi.SetResponse
}

var _ gnmi.GNMIClient = (*MockClient)(nil)

// Subscribe opens a subscription
func (mc *MockClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return <-mc.SubscribeStream, nil
}

// Set pushes a SetRequest to the server
func (mc *MockClient) Set(ctx context.Context, req *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case mc.SetReq <- req:
	}

	setResp := <-mc.SetResp
	return setResp, nil
}

// Capabilities retrieves capabilities
func (mc *MockClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

// Get reads from the server
func (mc *MockClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
