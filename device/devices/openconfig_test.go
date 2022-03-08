// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package devices

import (
	"errors"
	"strings"
	"sync"
	"testing"

	pg "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
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
			ocClient := &mockClient{
				subscribeStream: make(chan *mockClientStream),
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
				d, err := oc.DeviceID()
				if err != nil {
					t.Error(err)
				}
				if d != tc.expectedID {
					t.Errorf("Expected %v but got %v", tc.expectedID, d)
				}
			}()

			for len(tc.expectedSubs) > 0 {
				stream := &mockClientStream{
					subReq:  make(chan *gnmi.SubscribeRequest),
					subResp: make(chan *gnmi.SubscribeResponse),
					errC:    make(chan error),
				}
				ocClient.subscribeStream <- stream

				// Check that the subscription matches
				subReq := <-stream.subReq
				theSub := tc.expectedSubs[0]
				if !proto.Equal(subReq, theSub.req) {
					t.Fatalf("Expected\n%v\nbut got\n%v", theSub.req, subReq)
				}
				// Push test responses or error
				for _, r := range theSub.responses {
					stream.subResp <- r
				}
				if theSub.responseErr != nil {
					stream.errC <- theSub.responseErr
				}

				// go to next sub
				tc.expectedSubs = tc.expectedSubs[1:]
				close(stream.subResp)
				close(stream.subReq)
				close(stream.errC)
			}
			close(ocClient.subscribeStream)

			// Wait for DeviceID goroutine to finish
			wait.Wait()
		})
	}
}
