// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package devices

import (
	"context"
	"fmt"
	"io"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

type mockClientStream struct {
	subReq  chan *gnmi.SubscribeRequest
	subResp chan *gnmi.SubscribeResponse
	errC    chan error
	grpc.ClientStream
}

func (mcs *mockClientStream) Send(req *gnmi.SubscribeRequest) error {
	mcs.subReq <- req
	return nil
}

func (mcs *mockClientStream) Recv() (*gnmi.SubscribeResponse, error) {
	select {
	case err := <-mcs.errC:
		return nil, err
	case r, ok := <-mcs.subResp:
		if !ok {
			return nil, io.EOF
		}
		return r, nil
	}
}

type mockClient struct {
	subscribeStream chan *mockClientStream

	setReq  chan *gnmi.SetRequest
	setResp chan *gnmi.SetResponse
}

func (mc *mockClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return <-mc.subscribeStream, nil
}

func (mc *mockClient) Set(ctx context.Context, req *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case mc.setReq <- req:
	}
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case setResp := <-mc.setResp:
		return setResp, nil
	}
}

func (mc *mockClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (mc *mockClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
