// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package devices

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/device/internal"
	pg "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestOpenConfigDeviceID(t *testing.T) {

	type expectedSubs = []struct {
		req         *gnmi.SubscribeRequest
		responses   []*gnmi.SubscribeResponse
		responseErr error
	}

	subscribeOnce := func(path string) *gnmi.SubscribeRequest {
		gpath := pg.PathFromString(path)
		// Element is set for backwards compatibility, so we need to include it too.
		gpath.Element = strings.Split(path, "/") // nolint: staticcheck
		return &gnmi.SubscribeRequest{
			Request: &gnmi.SubscribeRequest_Subscribe{
				Subscribe: &gnmi.SubscriptionList{
					Prefix: &gnmi.Path{},
					Mode:   gnmi.SubscriptionList_ONCE,
					Subscription: []*gnmi.Subscription{
						{
							Path: gpath,
						},
					},
				},
			},
		}
	}

	subscribeUpdates := func(ups ...*gnmi.Update) *gnmi.SubscribeResponse {
		return &gnmi.SubscribeResponse{
			Response: &gnmi.SubscribeResponse_Update{
				Update: &gnmi.Notification{
					Update: ups,
				},
			},
		}
	}

	for _, tc := range []struct {
		name         string
		expectedID   string
		expectedSubs expectedSubs
	}{
		{
			name:       "from components state",
			expectedID: "serial-123",
			expectedSubs: expectedSubs{
				{
					req: subscribeOnce("components/component/state"),
					responses: []*gnmi.SubscribeResponse{
						subscribeUpdates(
							pg.Update(
								pg.PathFromString("components/component[name=x]/state/serial-no"),
								agnmi.TypedValue("serial-123")),
							pg.Update(
								pg.PathFromString("components/component[name=x]/state/type"),
								agnmi.TypedValue("openconfig-platform-types:CHASSIS")),
						),
					},
				},
			},
		},
		{
			name:       "from lldp state",
			expectedID: "serial-1234",
			expectedSubs: expectedSubs{
				{
					req:       subscribeOnce("components/component/state"),
					responses: []*gnmi.SubscribeResponse{},
				},
				{
					req: subscribeOnce("lldp/state/chassis-id"),
					responses: []*gnmi.SubscribeResponse{
						subscribeUpdates(
							pg.Update(
								pg.PathFromString("lldp/state/chassis-id"),
								agnmi.TypedValue("serial-1234")),
						),
					},
				},
			},
		},
		{
			name:       "from address",
			expectedID: "the-address",
			expectedSubs: expectedSubs{
				{
					req:       subscribeOnce("components/component/state"),
					responses: []*gnmi.SubscribeResponse{},
				},
				{
					req:       subscribeOnce("lldp/state/chassis-id"),
					responses: []*gnmi.SubscribeResponse{},
				},
			},
		},
		{
			name:       "from lldp state when components had error",
			expectedID: "serial-1234",
			expectedSubs: expectedSubs{
				{
					req:         subscribeOnce("components/component/state"),
					responses:   nil,
					responseErr: errors.New("components not supported"),
				},
				{
					req: subscribeOnce("lldp/state/chassis-id"),
					responses: []*gnmi.SubscribeResponse{
						subscribeUpdates(
							pg.Update(
								pg.PathFromString("lldp/state/chassis-id"),
								agnmi.TypedValue("serial-1234")),
						),
					},
				},
			},
		},
		{
			name:       "from address when both components and lldp had error",
			expectedID: "the-address",
			expectedSubs: expectedSubs{
				{
					req:         subscribeOnce("components/component/state"),
					responses:   nil,
					responseErr: errors.New("components not supported"),
				},
				{
					req:         subscribeOnce("lldp/state/chassis-id"),
					responses:   nil,
					responseErr: errors.New("lldp not supported"),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ocClient := &internal.MockClient{
				SubscribeStream: make(chan *internal.MockClientStream),
			}
			oc := openconfigDevice{
				config: &agnmi.Config{
					Addr: "the-address",
				},
				gNMIClient: ocClient,
			}

			wait := sync.WaitGroup{}
			wait.Add(1)
			go func() {
				defer wait.Done()
				d, err := oc.DeviceID(context.Background())
				if err != nil {
					t.Error(err)
				}
				if d != tc.expectedID {
					t.Errorf("Expected %v but got %v", tc.expectedID, d)
				}
			}()

			for len(tc.expectedSubs) > 0 {
				stream := &internal.MockClientStream{
					SubReq:  make(chan *gnmi.SubscribeRequest),
					SubResp: make(chan *gnmi.SubscribeResponse),
					ErrC:    make(chan error),
				}
				ocClient.SubscribeStream <- stream

				// Check that the subscription matches
				subReq := <-stream.SubReq
				theSub := tc.expectedSubs[0]
				if !proto.Equal(subReq, theSub.req) {
					t.Fatalf("Expected\n%v\nbut got\n%v", theSub.req, subReq)
				}
				// Push test responses or error
				for _, r := range theSub.responses {
					stream.SubResp <- r
				}
				if theSub.responseErr != nil {
					stream.ErrC <- theSub.responseErr
				}

				// go to next sub
				tc.expectedSubs = tc.expectedSubs[1:]
				close(stream.SubResp)
				close(stream.SubReq)
				close(stream.ErrC)
			}
			close(ocClient.SubscribeStream)

			// Wait for DeviceID goroutine to finish
			wait.Wait()
		})
	}
}

// implements gnmi.GNMIServer
type mockGnmiServer struct {
	serveCapability func(ctx context.Context,
		r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error)
}

func (g *mockGnmiServer) Capabilities(ctx context.Context, r *gnmi.CapabilityRequest) (
	*gnmi.CapabilityResponse, error) {
	return g.serveCapability(ctx, r, g)
}

func (g *mockGnmiServer) Get(_ context.Context, _ *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	panic("Not implemented")
}

func (g *mockGnmiServer) Set(_ context.Context, _ *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	panic("Not implemented")
}

func (g *mockGnmiServer) Subscribe(_ gnmi.GNMI_SubscribeServer) error {
	panic("not implemented")
}

func TestNewOpenConfig(t *testing.T) {

	//
	nonCompliantServerLst, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to create gRPC listener: %s", err)
	}
	defer nonCompliantServerLst.Close()
	nonCompliantGrpcServer := grpc.NewServer()
	defer nonCompliantGrpcServer.Stop()
	go func() {
		_ = nonCompliantGrpcServer.Serve(nonCompliantServerLst)
	}()
	t.Logf("nonCompliantAddr: %v", nonCompliantServerLst.Addr().String())

	nonGRPCServLst, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to create gRPC listener: %s", err)
	}
	defer nonGRPCServLst.Close()
	t.Logf("nonGRPCServer: %v", nonGRPCServLst.Addr().String())

	logrus.SetLevel(logrus.DebugLevel)

	for _, tcase := range []struct {
		name            string
		opts            map[string]string
		serveCapability func(ctx context.Context,
			r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error)
		expectErr string
	}{
		{
			name: "Bad address, can't connect",
			opts: map[string]string{
				"address": "wrong:12345",
				"timeout": "400ms", // so test is quicker
			},
			expectErr: "failed to dial: context deadline exceeded",
		},
		{
			name: "Good connection, but auth fails",
			opts: map[string]string{
				"address": "ACTUAL_SERVER",
			},
			serveCapability: func(ctx context.Context,
				r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error) {
				return nil, status.Error(codes.Unauthenticated, "auth err")
			},
			expectErr: "failed to reach device: rpc error: code = Unauthenticated desc = auth err",
		},
		{
			name: "Can reach address, but it is not compliant to gNMI",
			opts: map[string]string{
				"address": nonCompliantServerLst.Addr().String(),
			},
			expectErr: "failed to reach device: " +
				"rpc error: code = Unimplemented desc = unknown service gnmi.gNMI",
		},
		{
			name: "Can reach address, but it is not gRPC",
			opts: map[string]string{
				"address": nonGRPCServLst.Addr().String(),
				"timeout": "400ms",
			},
			expectErr: "failed to dial: context deadline exceeded",
		},
		{
			name: "Good connection, capabilities returns non implemented should work",
			opts: map[string]string{
				"address": "ACTUAL_SERVER",
			},
			serveCapability: func(ctx context.Context,
				r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error) {
				return nil, status.Error(codes.Unimplemented, "not implemented")
			},
		},
		{
			name: "Good connection but gets connection refused with code Unavailable, fails",
			opts: map[string]string{
				"address": "ACTUAL_SERVER",
				"timeout": "400ms",
			},
			serveCapability: func(ctx context.Context,
				r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error) {
				return nil, status.Error(codes.Unavailable, "something connect: connection refused")
			},
			expectErr: "failed to reach device: " +
				"rpc error: code = Unavailable desc = something connect: connection refused",
		},
		{
			name: "All good",
			opts: map[string]string{
				"address": "ACTUAL_SERVER",
			},
			serveCapability: func(ctx context.Context,
				r *gnmi.CapabilityRequest, s *mockGnmiServer) (*gnmi.CapabilityResponse, error) {
				return &gnmi.CapabilityResponse{}, nil
			},
		},
	} {
		t.Run(tcase.name, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			errg, _ := errgroup.WithContext(ctx)

			// Setup mock gnmi server
			server := &mockGnmiServer{
				serveCapability: tcase.serveCapability,
			}

			// Start grpc server
			grpcServer := grpc.NewServer()
			gRPCLis, err := net.Listen("tcp", "localhost:0")
			if err != nil {
				t.Fatalf("failed to create gRPC listener: %s", err)
			}
			gnmi.RegisterGNMIServer(grpcServer, server)
			defer gRPCLis.Close()
			defer grpcServer.Stop()
			errg.Go(func() error {
				return grpcServer.Serve(gRPCLis)
			})

			// Set actual server address if needed
			if tcase.opts["address"] == "ACTUAL_SERVER" {
				tcase.opts["address"] = gRPCLis.Addr().String()
			}

			// Sanitized options so we don't need to specify what is not required
			opts, err := device.SanitizedOptions(openConfigOptions(), tcase.opts)
			if err != nil {
				t.Fatalf("Bad test options passed: %v", err)
			}

			dev, err := newOpenConfig(ctx, opts, nil)
			if err != nil {
				if err.Error() != tcase.expectErr {
					t.Fatalf("Expected err: %v, got %v", tcase.expectErr, err)
				}
			}
			if err == nil && len(tcase.expectErr) > 0 {
				t.Fatalf("Expected error %v, but got none", tcase.expectErr)
			}
			if len(tcase.expectErr) == 0 && dev == nil {
				t.Fatalf("Expected a valid device object, got nil")
			}
			grpcServer.Stop()
			cancel()

			if err := errg.Wait(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
