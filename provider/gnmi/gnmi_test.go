// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/aristanetworks/cloudvision-go/provider/mock"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

// inClient is a mock gNMI client that will stream out a pre-determined set of
// subscribe responses
type inClient struct {
	cancel    context.CancelFunc
	responses []*gnmi.SubscribeResponse
	t         *testing.T
}

// inSubClient is the GNMI_SubcribeClient returned by inClient.Subscribe.
type inSubClient struct {
	cancel    context.CancelFunc
	responses []*gnmi.SubscribeResponse
	grpc.ClientStream
}

// Send just validates the SubscribeRequest (currently a no-op).
func (sc *inSubClient) Send(*gnmi.SubscribeRequest) error {
	return nil
}

// Recv returns pre-determined Subscribe Responses one by one.
func (sc *inSubClient) Recv() (*gnmi.SubscribeResponse, error) {
	if len(sc.responses) == 0 {
		sc.cancel()
		return nil, io.EOF
	}
	r := sc.responses[0]
	sc.responses = sc.responses[1:]
	return r, nil
}

func (c *inClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	c.t.Errorf("Capabilites method should not be called on inClient")
	return nil, errors.New("capabilities call not implemented")
}

func (c *inClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	c.t.Errorf("Get method should not be called on inClient")
	return nil, errors.New("get call not implemented")
}

func (c *inClient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	c.t.Errorf("Set method should not be called on inClient")
	return nil, errors.New("set call not implemented")
}

func (c *inClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	resp := c.responses
	c.responses = nil // so we don't send the same responses again.
	return &inSubClient{cancel: c.cancel, responses: resp}, nil
}

// outClient is a mock gNMI client that verifies that a pre-determined
// set of set requests are being made by the caller.
type outClient struct {
	requests []*gnmi.SetRequest
	t        *testing.T
}

func (c *outClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	c.t.Errorf("Capabilites method should not be called on inClient")
	return nil, errors.New("capabilities call not implemented")
}

func (c *outClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	c.t.Errorf("Get method should not be called on inClient")
	return nil, errors.New("get call not implemented")
}

func (c *outClient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	if len(c.requests) == 0 {
		c.t.Errorf("unexpected Set request on outClient: %v", in)
		return nil, errors.New("unexpected set request")
	}
	req := c.requests[0]
	c.requests = c.requests[1:]
	if !reflect.DeepEqual(req, in) {
		c.t.Errorf("incorrect Seq request received.\n\tExpected: %v\n\tReceived: %v", req, in)
		return nil, nil // we don't return error, so that the caller doesn't bail out.
	}
	return nil, nil
}

func (c *outClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	c.t.Errorf("Subscribe method should not be called on inClient")
	return nil, errors.New("subscribe call not implemented")
}

func gnmiUpdate(path string, val interface{}) *gnmi.Update {
	return &gnmi.Update{
		Path: PathFromString(path),
		Val:  jsonValue(val),
	}
}

func gnmiDelete(paths []string) []*gnmi.Path {
	d := make([]*gnmi.Path, len(paths))
	for i, p := range paths {
		d[i] = PathFromString(p)
	}
	return d
}

func newSubscribeResponse(ts int64, prefix *gnmi.Path,
	updates []*gnmi.Update, deletes []*gnmi.Path) *gnmi.SubscribeResponse {
	return &gnmi.SubscribeResponse{
		Response: &gnmi.SubscribeResponse_Update{
			Update: &gnmi.Notification{
				Timestamp: ts,
				Prefix:    prefix,
				Update:    updates,
				Delete:    deletes,
			},
		},
	}
}
func TestGNMIProvider(t *testing.T) {
	cases := []struct {
		name     string
		subResps []*gnmi.SubscribeResponse
		setReqs  []*gnmi.SetRequest
		paths    []string
	}{
		{
			name: "basic update and delete",
			subResps: []*gnmi.SubscribeResponse{
				newSubscribeResponse(
					435, //timestamp
					nil,
					[]*gnmi.Update{
						gnmiUpdate("/u/l1", "x"),
						gnmiUpdate("/u/l2", "y"),
					},
					gnmiDelete([]string{"/d1", "/d2"}),
				),
			},
			setReqs: []*gnmi.SetRequest{
				{
					Update: []*gnmi.Update{
						gnmiUpdate("/u/l1", "x"),
						gnmiUpdate("/u/l2", "y"),
					},
					Delete: gnmiDelete([]string{"/d1", "/d2"}),
				},
			},
			paths: []string{"/"},
		},
		{
			name: "handle prefix",
			subResps: []*gnmi.SubscribeResponse{
				newSubscribeResponse(
					435, //timestamp
					PathFromString("/prefix"),
					[]*gnmi.Update{
						gnmiUpdate("/u/l1", "x"),
						gnmiUpdate("/u/l2", "y"),
					},
					gnmiDelete([]string{"/d1", "/d2"}),
				),
			},
			setReqs: []*gnmi.SetRequest{
				{
					Prefix: PathFromString("/prefix"),
					Update: []*gnmi.Update{
						gnmiUpdate("/u/l1", "x"),
						gnmiUpdate("/u/l2", "y"),
					},
					Delete: gnmiDelete([]string{"/d1", "/d2"}),
				},
			},
			paths: []string{"/"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			incl := &inClient{cancel: cancel, t: t, responses: tc.subResps}
			outcl := &outClient{t: t, requests: tc.setReqs}
			cfg := &agnmi.Config{}
			monitor := mock.NewMockMonitor()
			p := NewGNMIProvider(incl, cfg, tc.paths, monitor)
			p.InitGNMI(outcl)
			_ = p.Run(ctx)
		})
	}
}
