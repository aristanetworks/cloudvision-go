// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/device/internal"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockCVClient struct {
	id         string
	gnmic      gnmi.GNMIClient
	metadataCh chan string
}

func (m *mockCVClient) Capabilities(ctx context.Context,
	in *gnmi.CapabilityRequest, opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return m.gnmic.Capabilities(ctx, in, opts...)
}

func (m *mockCVClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return m.gnmic.Get(ctx, in, opts...)
}

func (m *mockCVClient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	return m.gnmic.Set(ctx, in, opts...)
}

func (m *mockCVClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return m.gnmic.Subscribe(ctx, opts...)
}

func (m *mockCVClient) SendDeviceMetadata(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.metadataCh <- m.id:
	}
	return nil
}

func (m *mockCVClient) SendHeartbeat(ctx context.Context, alive bool) error {
	return nil
}

func (m *mockCVClient) ForProvider(p provider.GNMIProvider) cvclient.CVClient {
	return m
}

func newMockCVClient(gc gnmi.GNMIClient, info *Info, metadata chan string) cvclient.CVClient {
	return &mockCVClient{
		id:         fmt.Sprintf("%v|%v", info.ID, info.Config.Options),
		gnmic:      gc,
		metadataCh: metadata,
	}
}

type mockDevice struct {
	id     string
	config map[string]string
}

var _ Device = (*mockDevice)(nil)

func (m *mockDevice) Alive() (bool, error) {
	return true, nil
}

func (m *mockDevice) DeviceID() (string, error) {
	return m.id, nil
}

func (m *mockDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (m *mockDevice) Type() string { return "" }

func (m *mockDevice) IPAddr() string {
	return ""
}

func newMockDevice(opt map[string]string) (Device, error) {
	fmt.Println("newMockDevice:", opt)
	deviceID, err := GetStringOption("id", opt)
	if err != nil {
		return nil, err
	}
	return &mockDevice{
		id:     deviceID,
		config: opt,
	}, nil
}

var mockDeviceOptions = map[string]Option{
	"id": {
		Description: "whatever",
		Required:    true,
	},
	"input1": {
		Description: "somedata",
		Required:    false,
	},
	"cred1": {
		Description: "somedata",
		Required:    false,
	},
}

type sensorTestCase struct {
	name                    string
	stateSubResps           []*gnmi.SubscribeResponse
	configSubResps          []*gnmi.SubscribeResponse
	waitForMetadataPreSync  []string
	waitForMetadataPostSync []string
	expectSet               []*gnmi.SetRequest
}

func subscribeUpdates(ups ...*gnmi.Update) *gnmi.SubscribeResponse {
	return &gnmi.SubscribeResponse{
		Response: &gnmi.SubscribeResponse_Update{
			Update: &gnmi.Notification{
				Update: ups,
			},
		},
	}
}

func datasourcePath(configOrState, id, name, leaf string) *gnmi.Path {
	path := fmt.Sprintf("/datasource/%s/sensor[id=%s]/source[name=%s]/%s",
		configOrState, id, name, leaf)
	return pgnmi.PathFromString(path)
}

func datasourceOptPath(id, name, optOrCred, key string) *gnmi.Path {
	path := fmt.Sprintf("/datasource/config/sensor[id=%s]/source[name=%s]/%s[key=%s]",
		id, name, optOrCred, key)
	return pgnmi.PathFromString(path)
}

func datasourceUpdate(configOrState, id, name, leaf string, val interface{}) *gnmi.Update {
	return pgnmi.Update(datasourcePath(configOrState, id, name, leaf),
		agnmi.TypedValue(val))
}

func datasourceOptUpdate(id, name, optOrCred, key, val string) []*gnmi.Update {
	return []*gnmi.Update{
		pgnmi.Update(pgnmi.PathAppend(datasourceOptPath(id, name, optOrCred, key), "value"),
			agnmi.TypedValue(val)),
		pgnmi.Update(pgnmi.PathAppend(datasourceOptPath(id, name, optOrCred, key), "key"),
			agnmi.TypedValue(key)),
	}
}

func datasourceUpdates(configOrState, id, name, typ string, enabled bool,
	opts map[string]string, creds map[string]string) []*gnmi.Update {
	upds := []*gnmi.Update{
		datasourceUpdate(configOrState, id, name, "name", name),
		datasourceUpdate(configOrState, id, name, "type", typ),
		datasourceUpdate(configOrState, id, name, "enabled", enabled),
	}
	for k, v := range opts {
		upds = append(upds, datasourceOptUpdate(id, name, "option", k, v)...)
	}
	for k, v := range creds {
		upds = append(upds, datasourceOptUpdate(id, name, "credential", k, v)...)
	}
	return upds
}

func runSensorTest(t *testing.T, tc sensorTestCase) {
	// Set up mock gNMI client
	gnmic := &internal.MockClient{
		SubscribeStream: make(chan *internal.MockClientStream),
		SetReq:          make(chan *gnmi.SetRequest),
		SetResp:         make(chan *gnmi.SetResponse),
	}
	defer close(gnmic.SubscribeStream)

	metadataCh := make(chan string)

	// Start sensor.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sensor := NewSensor("abc",
		WithSensorGNMIClient(gnmic),
		WithSensorGRPCConn(nil),
		WithSensorClientFactory(
			func(gc gnmi.GNMIClient, info *Info) cvclient.CVClient {
				return newMockCVClient(gc, info, metadataCh)
			}))
	sensor.deviceRedeployTimer = 10 * time.Millisecond
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error { return sensor.Run(ctx) })
	errg.Go(func() error {
		// state responses
		stream := &internal.MockClientStream{
			SubReq:  make(chan *gnmi.SubscribeRequest),
			SubResp: make(chan *gnmi.SubscribeResponse),
			ErrC:    make(chan error),
		}
		gnmic.SubscribeStream <- stream
		t.Logf("Got state sub request: %v", <-stream.SubReq)
		for _, resp := range tc.stateSubResps {
			stream.SubResp <- resp
		}
		close(stream.SubResp)
		close(stream.SubReq)
		close(stream.ErrC)

		// config responses
		stream = &internal.MockClientStream{
			SubReq:  make(chan *gnmi.SubscribeRequest),
			SubResp: make(chan *gnmi.SubscribeResponse),
			ErrC:    make(chan error),
		}

		waitMetas := func(expectation []string) {
			expectMetas := map[string]struct{}{}
			for _, m := range expectation {
				expectMetas[m] = struct{}{}
			}
			t.Logf("Waiting for meta: %v", expectMetas)
			for len(expectMetas) > 0 {
				select {
				case <-ctx.Done():
					return
				case gotMeta := <-metadataCh:
					if _, ok := expectMetas[gotMeta]; ok {
						delete(expectMetas, gotMeta)
						t.Logf("got metadata %v", gotMeta)
					} else {
						t.Fatalf("Unexpected meta: %v", gotMeta)
					}
				case <-time.After(10 * time.Second):
					t.Fatalf("Failed to see metadata: %v", expectMetas)
				}
			}
		}

		gnmic.SubscribeStream <- stream
		t.Logf("Got config sub request: %v", <-stream.SubReq)
		for _, resp := range tc.configSubResps {
			t.Log("Pushing config", resp)
			select {
			case stream.SubResp <- resp:
				if resp.GetSyncResponse() {
					waitMetas(tc.waitForMetadataPreSync)
				}
			case <-ctx.Done():
				t.Errorf("Context canceled before pushing config resp: %v", resp)
				return ctx.Err()
			}
		}
		waitMetas(tc.waitForMetadataPostSync)
		t.Log("Canceling context as the test is done!")
		cancel()

		close(stream.SubResp)
		close(stream.SubReq)
		close(stream.ErrC)
		return nil
	})
	errg.Go(func() error {
		for {
			select {
			case setReq, ok := <-gnmic.SetReq:
				if !ok {
					if len(tc.expectSet) > 0 {
						return fmt.Errorf("Set channel closed but expected more sets: %v",
							tc.expectSet)
					}
					return nil
				}
				found := -1
				for i, expectSet := range tc.expectSet {
					if proto.Equal(expectSet, setReq) {
						found = i
						break
					}
				}
				if found >= 0 {
					t.Log("Matched set")
					select {
					case gnmic.SetResp <- &gnmi.SetResponse{
						Prefix: pgnmi.PathFromString("ok-dont-care")}:
					case <-ctx.Done():
						return nil
					}
					copy(tc.expectSet[found:], tc.expectSet[found+1:])
					tc.expectSet = tc.expectSet[:len(tc.expectSet)-1]
				} else {
					return fmt.Errorf("Found unexpected set: %v\nexpecting: %v",
						setReq, tc.expectSet)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	if err := errg.Wait(); err != nil && err != context.Canceled {
		t.Fatal(err)
	}
}

func TestSensor(t *testing.T) {
	testCases := []sensorTestCase{
		{
			name: "Pre-existing datasource",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value1"}, nil)...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
			},
			waitForMetadataPostSync: []string{"123|map[id:123 input1:value1]"},
		},
		{
			name: "Update existing datasource",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value1"}, nil)...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value2"}, nil)...),
			},
			waitForMetadataPreSync:  []string{"123|map[id:123 input1:value1]"},
			waitForMetadataPostSync: []string{"123|map[id:123 input1:value2]"},
		},
		{
			name: "Datasource added after sync",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123"}, nil)...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123]"},
		},
		{
			name: "Datasource with invalid config should keep others going (Pre-sync)",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "bad1", "invalidtype",
						true, map[string]string{"id": "111"}, nil)...),
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123"}, nil)...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
			},
			waitForMetadataPreSync: []string{"123|map[id:123]"},
			expectSet: []*gnmi.SetRequest{
				{
					Prefix: pgnmi.PathFromString(
						fmt.Sprintf("datasource/state/sensor[id=%s]/source[name=%s]",
							"abc", "bad1")),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.PathFromString("last-error"), agnmi.TypedValue(
							"Failed creating device 'invalidtype': "+
								"Device 'invalidtype' not found")),
					},
				},
			},
		},
		{
			name: "Datasource with invalid config should keep others going (Post-sync)",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "bad1", "invalidtype",
						true, map[string]string{"id": "111"}, nil)...),
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123"}, nil)...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123]"},
			expectSet: []*gnmi.SetRequest{
				{
					Prefix: pgnmi.PathFromString(
						fmt.Sprintf("datasource/state/sensor[id=%s]/source[name=%s]",
							"abc", "bad1")),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.PathFromString("last-error"), agnmi.TypedValue(
							"Failed creating device 'invalidtype': "+
								"Device 'invalidtype' not found")),
					},
				},
			},
		},
	}

	// Register mockDevice.
	Register("mock", newMockDevice, mockDeviceOptions)

	// Run through test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runSensorTest(t, tc)
		})
	}
}

func TestDatasourceDeployLoop(t *testing.T) {

	createCalls := make(chan string, 1)

	const deviceName = "dev1"
	const deviceType = "mock"

	Register(deviceType, func(m map[string]string) (Device, error) {
		createCalls <- fmt.Sprintf("%v", m)
		return newMockDevice(m)
	}, mockDeviceOptions)

	// Setup one datasource that we will use for all test cases
	gnmic := &internal.MockClient{
		SubscribeStream: make(chan *internal.MockClientStream),
		SetReq:          make(chan *gnmi.SetRequest, 100),
		SetResp:         make(chan *gnmi.SetResponse),
	}
	close(gnmic.SetResp) // we don't care about these responses, so make it always return nil

	metadataCh := make(chan string, 100)
	closedCh := make(chan struct{})
	close(closedCh)
	sensor := NewSensor("default", WithSensorClientFactory(
		func(gc gnmi.GNMIClient, info *Info) cvclient.CVClient {
			return newMockCVClient(gnmic, info, metadataCh)
		},
	), WithSensorGNMIClient(gnmic))
	sensor.datasourceConfig[deviceName] = &datasourceConfig{
		name:       deviceName,
		typ:        deviceType,
		option:     map[string]string{},
		credential: map[string]string{},
	}

	for _, tc := range []struct {
		name               string
		config             datasourceConfig
		expectDeviceCreate bool
	}{
		{
			name: "create valid config",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id": "123",
				},
			},
			expectDeviceCreate: true,
		},
		{
			name: "disable config",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: false,
				option: map[string]string{
					"id": "123",
				},
			},
			expectDeviceCreate: false,
		},
		{
			name: "re-enable config",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id": "123",
				},
			},
			expectDeviceCreate: true,
		},
		{
			name: "modify config will restart device",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id": "124", // change should trigger restart
				},
			},
			expectDeviceCreate: true,
		},
		{
			name: "modify credential will restart device",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id": "124",
				},
				credential: map[string]string{
					"cred1": "abc", // change should trigger restart
				},
			},
			expectDeviceCreate: true,
		},
		{
			name: "same config will not redeploy",
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id": "124",
				},
				credential: map[string]string{
					"cred1": "abc",
				},
			},
			expectDeviceCreate: false,
		},
	} {
		ctx := context.Background()
		t.Run(tc.name, func(t *testing.T) {
			cfg := sensor.datasourceConfig[deviceName]

			cfg.enabled = tc.config.enabled
			// simulate update behavior on maps
			for k := range cfg.option {
				delete(cfg.option, k)
			}
			for k, v := range tc.config.option {
				cfg.option[k] = v
			}
			for k := range cfg.credential {
				delete(cfg.credential, k)
			}
			for k, v := range tc.config.credential {
				cfg.credential[k] = v
			}

			err := sensor.runDatasourceConfig(ctx, deviceName)
			if err != nil {
				t.Fatal(err)
			}

			// Check if we got a create call.
			// Create is synchronous so we should either have it in the channel or not.
			select {
			case x := <-createCalls:
				if !tc.expectDeviceCreate {
					t.Fatalf("Dit not expect device being created! Got: %v", x)
				}
			default:
				if tc.expectDeviceCreate {
					t.Fatal("Dit not see device being created after 10s!")
				}
			}
		})
	}
}
