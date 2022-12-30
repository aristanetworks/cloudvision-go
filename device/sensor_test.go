// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	"github.com/aristanetworks/cloudvision-go/device/internal"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
	exit := m.Run()
	os.Exit(exit)
}

type mockCVClient struct {
	id         string
	gnmic      gnmi.GNMIClient
	metadataCh chan string
	config     map[string]string
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

func (m *mockCVClient) SetManagedDevices(a []string) {
	_, err := m.gnmic.Set(context.Background(), &gnmi.SetRequest{
		Prefix: pgnmi.Path("managed-devices"),
		Update: []*gnmi.Update{
			pgnmi.Update(pgnmi.Path("ids"), agnmi.TypedValue(fmt.Sprintf("%v", a))),
		},
	})
	if err != nil {
		panic(err)
	}
}

func newMockCVClient(gc gnmi.GNMIClient, info *Info, metadata chan string) cvclient.CVClient {
	return &mockCVClient{
		id:         fmt.Sprintf("%v|%v", info.ID, info.Config.Options),
		gnmic:      gc,
		metadataCh: metadata,
		config:     info.Config.Options,
	}
}

type crashProvider struct {
}

func (c *crashProvider) InitGRPC(conn *grpc.ClientConn) {}

func (c *crashProvider) Run(ctx context.Context) error {
	panic("Crash!")
}

var _ provider.GRPCProvider = (*crashProvider)(nil)

type mockDevice struct {
	id     string
	config map[string]string
}

var _ Device = (*mockDevice)(nil)
var _ Manager = (*mockDevice)(nil)

func (m *mockDevice) Alive(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *mockDevice) DeviceID(ctx context.Context) (string, error) {
	return m.id, nil
}

func (m *mockDevice) Providers() ([]provider.Provider, error) {
	if v := m.config["crash"]; v == "provider" {
		return []provider.Provider{&crashProvider{}}, nil
	}
	return nil, nil
}

func (m *mockDevice) Type() string { return "" }

func (m *mockDevice) IPAddr(ctx context.Context) (string, error) {
	return "", nil
}

func (m *mockDevice) Manage(ctx context.Context, inventory Inventory) error {
	if v := m.config["crash"]; v == "manager" {
		panic("Crash manager!")
	}
	if v, ok := m.config["managed"]; ok && len(v) > 0 {
		ids := strings.Split(v, ",")
		for _, id := range ids {
			if err := inventory.Add(&Info{
				ID:      id,
				Context: ctx,
				Device:  &mockDevice{id: id},
				Config:  &Config{},
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func newMockDevice(ctx context.Context, opt map[string]string) (Device, error) {
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
	"crash": {
		Description: "makes device crash for test purposes",
		Required:    false,
	},
	"input1": {
		Description: "somedata",
		Required:    false,
	},
	"managed": {
		Description: "managed devices comma separated",
		Required:    false,
	},
	"cred1": {
		Description: "somedata",
		Required:    false,
	},
}

type sensorTestCase struct {
	name     string
	substeps []*sensorTestCase

	stateSubResps           []*gnmi.SubscribeResponse
	configSubResps          []*gnmi.SubscribeResponse
	waitForMetadataPreSync  []string
	waitForMetadataPostSync []string
	expectSet               []*gnmi.SetRequest
	dynamicConfigs          []*Config
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

func sensorPath(configOrState, id string) *gnmi.Path {
	path := fmt.Sprintf("/datasource/%s/sensor[id=%s]", configOrState, id)
	p := pgnmi.PathFromString(path)
	p.Origin = "arista"
	p.Target = "cv"
	return p
}

func sensorUpdate(configOrState, id, leaf string, val interface{}) *gnmi.Update {
	return pgnmi.Update(pgnmi.PathAppend(sensorPath(configOrState, id), leaf),
		agnmi.TypedValue(val))
}

func datasourcePath(configOrState, id, name, leaf string) *gnmi.Path {
	path := fmt.Sprintf("/datasource/%s/sensor[id=%s]/source[name=%s]/%s",
		configOrState, id, name, leaf)
	p := pgnmi.PathFromString(path)
	p.Origin = "arista"
	p.Target = "cv"
	return p
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

func waitMetas(ctx context.Context, t *testing.T,
	metadataCh chan string, expectation []string) error {
	expectMetas := map[string]struct{}{}
	for _, m := range expectation {
		expectMetas[m] = struct{}{}
	}
	t.Logf("Waiting for meta: %v", expectMetas)
	for len(expectMetas) > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case gotMeta := <-metadataCh:
			if _, ok := expectMetas[gotMeta]; ok {
				delete(expectMetas, gotMeta)
				t.Logf("got metadata %v", gotMeta)
			} else {
				return fmt.Errorf("Unexpected meta: %v", gotMeta)
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("Failed to see metadata: %v", expectMetas)
		}
	}
	return nil
}

func initialSetReq(sensor string) *gnmi.SetRequest {
	prefix := pgnmi.PathFromString("datasource/state/sensor[id=" + sensor + "]/")
	prefix.Origin = "arista"
	prefix.Target = "cv"
	return &gnmi.SetRequest{
		Prefix: prefix,
		Update: []*gnmi.Update{
			pgnmi.Update(pgnmi.Path("version"), agnmi.TypedValue("dev")),
			pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(42)),
			pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
			pgnmi.Update(pgnmi.Path("last-error"), agnmi.TypedValue("Sensor started")),
		},
	}
}

func runSensorTest(t *testing.T, tc sensorTestCase) {
	// Always run as if it was substeps.
	if len(tc.substeps) > 0 && len(tc.configSubResps) > 0 {
		t.Fatal("Should use either substeps or assume one step only")
	}
	if len(tc.substeps) == 0 {
		tc.substeps = []*sensorTestCase{&tc}
	}

	// Set up mock gNMI client
	gnmic := &internal.MockClient{
		SubscribeStream: make(chan *internal.MockClientStream),
		SetReq:          make(chan *gnmi.SetRequest),
		SetResp:         make(chan *gnmi.SetResponse),
	}
	defer close(gnmic.SubscribeStream)

	metadataCh := make(chan string)

	configCh := make(chan *Config)
	defer close(configCh)

	// Start sensor.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sensor := NewSensor("abc",
		WithSensorGNMIClient(gnmic),
		WithSensorHeartbeatInterval(50*time.Millisecond),
		WithSensorGRPCConn(nil),
		WithSensorConfigChan(configCh),
		WithSensorClientFactory(
			func(gc gnmi.GNMIClient, info *Info) cvclient.CVClient {
				return newMockCVClient(gc, info, metadataCh)
			}))
	sensor.log = sensor.log.WithField("test", tc.name)
	sensor.deviceRedeployTimer = 10 * time.Millisecond
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return sensor.Run(ctx)
	})

	configStream := &internal.MockClientStream{
		SubReq:  make(chan *gnmi.SubscribeRequest),
		SubResp: make(chan *gnmi.SubscribeResponse),
		ErrC:    make(chan error),
	}

	errg.Go(func() error {
		// Create the state stream.
		stream := &internal.MockClientStream{
			SubReq:  make(chan *gnmi.SubscribeRequest),
			SubResp: make(chan *gnmi.SubscribeResponse),
			ErrC:    make(chan error),
		}
		// Serve state stream.
		gnmic.SubscribeStream <- stream
		t.Logf("Got state sub request: %v", <-stream.SubReq)
		for _, resp := range tc.stateSubResps {
			stream.SubResp <- resp
		}
		close(stream.SubResp)
		close(stream.SubReq)
		close(stream.ErrC)
		// Serve config stream.
		gnmic.SubscribeStream <- configStream
		t.Logf("Got config sub request: %v", <-configStream.SubReq)
		return nil
	})

	for _, tc := range tc.substeps {
		// Start a sequence of checks for each substep.
		// This allow us to wait and see a set of SetRequests based on some input
		// before pushing more configs.
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			errgSub, ctx := errgroup.WithContext(ctx)

			errgSub.Go(func() error {
				// config responses
				for _, resp := range tc.configSubResps {
					t.Log("Pushing config", resp)
					select {
					case configStream.SubResp <- resp:
						if resp.GetSyncResponse() {
							if err := waitMetas(ctx, t,
								metadataCh, tc.waitForMetadataPreSync); err != nil {
								return err
							}
						}
					case <-ctx.Done():
						t.Errorf("Context canceled before pushing config resp: %v", resp)
						return ctx.Err()
					}
				}
				if err := waitMetas(ctx, t, metadataCh, tc.waitForMetadataPostSync); err != nil {
					return err
				}
				return nil
			})
			errgSub.Go(func() error {
				for i, c := range tc.dynamicConfigs {
					select {
					case configCh <- c:
					case <-ctx.Done():
						return fmt.Errorf("Failed to push dynamic config index %d: %v", i, c)
					}
				}
				return nil
			})
			errgSub.Go(func() error {
				defer cancel()

				setsIdx := 0
				extraReads := time.NewTicker(time.Hour)
				if len(tc.expectSet) == 0 {
					extraReads.Reset(5 * time.Millisecond)
				}
				defer extraReads.Stop()
				for timeout := time.After(5 * time.Second); ; {
					select {
					case <-timeout:
						if len(tc.expectSet) > 0 {
							return fmt.Errorf("Timed out reading sets, expecting: %v", tc.expectSet)
						}
						return nil
					case <-extraReads.C:
						select {
						case r := <-gnmic.SetReq:
							return fmt.Errorf("Unexpected set: %v", r)
						default:
							return nil
						}
					case setReq, ok := <-gnmic.SetReq:
						setsIdx++
						t.Logf("Got set #%d: %v", setsIdx, setReq)
						if !ok {
							if len(tc.expectSet) > 0 {
								return fmt.Errorf("Set channel closed but expected more sets: %v",
									tc.expectSet)
							}
							return nil
						}

						// Always send the set response
						select {
						case gnmic.SetResp <- &gnmi.SetResponse{
							Prefix: pgnmi.PathFromString("ok-dont-care")}:
						case <-ctx.Done():
							return nil
						}

						// change date in set request so we can easily match it
						heartbeat := false
						for _, u := range setReq.Update {
							if pgnmi.PathMatch(u.Path, pgnmi.Path("streaming-start")) {
								u.Val = agnmi.TypedValue(42)
							}
							if pgnmi.PathMatch(u.Path, pgnmi.Path("last-seen")) {
								u.Val = agnmi.TypedValue(43)
								heartbeat = len(setReq.Update) == 1
							}
						}

						// Skip sensor heartbeats
						if heartbeat && pgnmi.PathMatch(sensor.statePrefix, setReq.Prefix) {
							t.Log("skipping sensor heartbeat")
							continue
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
							copy(tc.expectSet[found:], tc.expectSet[found+1:])
							tc.expectSet = tc.expectSet[:len(tc.expectSet)-1]
							if len(tc.expectSet) == 0 {
								extraReads.Reset(5 * time.Millisecond) // schedule last check
							}
						} else {
							return fmt.Errorf("Set %d unexpected:\n%s\nexpecting:\n%v",
								setsIdx, setReq, tc.expectSet)
						}
					case <-ctx.Done():
						if len(tc.expectSet) > 0 {
							t.Errorf("Context cancel but did not match all sets: %v", tc.expectSet)
						}
						return ctx.Err()
					}
				}
			})
			if err := errgSub.Wait(); err != nil && err != context.Canceled {
				t.Fatal(err)
			}
		})
	}
	close(configStream.SubResp)
	close(configStream.SubReq)
	close(configStream.ErrC)
	cancel()
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
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
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
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
			},
		},
		{
			name: "Disable existing datasource",
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
						false, map[string]string{"id": "123", "input1": "value2"}, nil)...),
			},
			waitForMetadataPreSync: []string{"123|map[id:123 input1:value1]"},
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(false)),
					},
				},
			},
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
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
			},
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
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "bad1", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("invalidtype")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "bad1", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.PathFromString("last-error"), agnmi.TypedValue(
							"Datasource stopped unexpectedly: "+
								"Failed creating device 'invalidtype': "+
								"Device 'invalidtype' not found")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
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
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "bad1", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("invalidtype")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "bad1", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.PathFromString("last-error"), agnmi.TypedValue(
							"Datasource stopped unexpectedly: "+
								"Failed creating device 'invalidtype': "+
								"Device 'invalidtype' not found")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
			},
		},
		{
			name: "Datasource with provider crash",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "crash": "provider"}, nil)...),
			},
			waitForMetadataPostSync: []string{"123|map[crash:provider id:123]"},
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("last-error"), agnmi.TypedValue(
							"Datasource stopped unexpectedly: "+
								"error starting providers for device \"123\" (mock): "+
								"fatal error in *device.crashProvider.Run: Crash!")),
					},
				},
			},
		},
		{
			name: "Datasource with Manager crash",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "crash": "manager"}, nil)...),
			},
			waitForMetadataPostSync: []string{"123|map[crash:manager id:123]"},
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("last-error"), agnmi.TypedValue(
							"Datasource stopped unexpectedly: "+
								"fatal error in Manage: Crash manager!")),
					},
				},
			},
		},
		{
			name: "Delete sensor config and start new configs again",
			substeps: []*sensorTestCase{
				{
					name: "First onboard device and see it streaming",
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
					expectSet: []*gnmi.SetRequest{
						initialSetReq("abc"),
						{
							Prefix: datasourcePath("state", "abc", "xyz", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
								pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
								pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
							},
						},
						{
							Prefix: datasourcePath("state", "abc", "xyz", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
							},
						},
					},
				},
				{
					name: "Delete sensor config and see delete happening",
					configSubResps: []*gnmi.SubscribeResponse{
						{
							Response: &gnmi.SubscribeResponse_Update{
								Update: &gnmi.Notification{
									Delete: []*gnmi.Path{
										sensorPath("config", "abc"),
									},
								},
							},
						},
					},
					expectSet: []*gnmi.SetRequest{
						{
							Delete: []*gnmi.Path{
								datasourcePath("state", "abc", "xyz", ""),
							},
						},
						{
							Delete: []*gnmi.Path{
								sensorPath("state", "abc"),
							},
						},
					},
				},
			},
		},
		{
			name: "No configs, no sets",
			configSubResps: []*gnmi.SubscribeResponse{
				{Response: &gnmi.SubscribeResponse_SyncResponse{SyncResponse: true}},
			},
			expectSet: []*gnmi.SetRequest{},
		},
		{
			name: "Device heartbeat and streaming-start",
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
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
						pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(42)),
					},
				},
			},
		},
		{
			name: "Device heartbeat with managed devices",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{
							"id":      "123",
							"managed": "m1,m2"}, nil)...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123 managed:m1,m2]"},
			expectSet: []*gnmi.SetRequest{
				initialSetReq("abc"),
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
						pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
						pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
						pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(42)),
					},
				},
				{
					Prefix: pgnmi.Path("managed-devices"),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("ids"), agnmi.TypedValue("[m1]")),
					},
				},
				{
					Prefix: pgnmi.Path("managed-devices"),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("ids"), agnmi.TypedValue("[m1 m2]")),
					},
				},
				{
					Prefix: datasourcePath("state", "abc", "xyz", ""),
					Update: []*gnmi.Update{
						pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
					},
				},
			},
		},
		{
			name: "Custom config add and remove",
			substeps: []*sensorTestCase{
				{ // Add device with custom config
					configSubResps: []*gnmi.SubscribeResponse{
						// create config to enable sensor
						subscribeUpdates(sensorUpdate("config", "abc", "id", "abc")),
						{Response: &gnmi.SubscribeResponse_SyncResponse{SyncResponse: true}},
					},
					dynamicConfigs: []*Config{
						{
							Name:   "device-1",
							Device: "mock",
							Options: map[string]string{
								"id":     "123",
								"input1": "value1"},
						},
					},
					expectSet: []*gnmi.SetRequest{
						initialSetReq("abc"),
						{
							Prefix: datasourcePath("state", "abc", "device-1", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(lastErrorKey, agnmi.TypedValue("Datasource started")),
								pgnmi.Update(pgnmi.Path("type"), agnmi.TypedValue("mock")),
								pgnmi.Update(pgnmi.Path("enabled"), agnmi.TypedValue(true)),
							},
						},
						{
							Prefix: datasourcePath("state", "abc", "device-1", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(pgnmi.Path("source-id"), agnmi.TypedValue("123")),
							},
						},
					},
					waitForMetadataPostSync: []string{"123|map[id:123 input1:value1]"},
				},
				{ // Same config will give no updates
					dynamicConfigs: []*Config{
						{
							Name:   "device-1",
							Device: "mock",
							Options: map[string]string{
								"id":     "123",
								"input1": "value1"},
						},
					},
					expectSet: []*gnmi.SetRequest{ // wait to see status update
						{
							Prefix: datasourcePath("state", "abc", "device-1", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
								pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(42)),
							},
						},
					},
				},
				{ // Delete config
					dynamicConfigs: []*Config{
						NewDeletedConfig("device-1"),
					},
					expectSet: []*gnmi.SetRequest{
						{
							Delete: []*gnmi.Path{
								datasourcePath("state", "abc", "device-1", ""),
							},
						},
					},
					waitForMetadataPostSync: []string{"123|map[id:123 input1:value1]"},
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

	Register(deviceType, func(ctx context.Context, m map[string]string) (Device, error) {
		createCalls <- fmt.Sprintf("%v", m)
		return newMockDevice(ctx, m)
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
			if tc.expectDeviceCreate {
				<-createCalls
			}
		})
	}
}
