// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	cvmock "github.com/aristanetworks/cloudvision-go/device/cvclient/mock"
	"github.com/aristanetworks/cloudvision-go/device/internal"
	dmock "github.com/aristanetworks/cloudvision-go/device/mock"
	gmock "github.com/aristanetworks/cloudvision-go/mock"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/golang/mock/gomock"
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
	err error
}

func (c *crashProvider) InitGRPC(conn *grpc.ClientConn) {}

func (c *crashProvider) Run(ctx context.Context) error {
	if c.err != nil {
		return c.err
	}
	panic("Crash!")
}

var _ provider.GRPCProvider = (*crashProvider)(nil)

type mockDevice struct {
	id             string
	config         map[string]string
	isAlive        bool
	shutDownReason error
}

var _ Device = (*mockDevice)(nil)
var _ Manager = (*mockDevice)(nil)

func (m *mockDevice) Alive(ctx context.Context) (bool, error) {
	if m.isAlive || m.shutDownReason == nil {
		return m.isAlive, nil
	}
	return m.isAlive, m.shutDownReason
}

func (m *mockDevice) DeviceID(ctx context.Context) (string, error) {
	return m.id, nil
}

func (m *mockDevice) Providers() ([]provider.Provider, error) {
	v, ok := m.config["crash"]
	if ok {
		switch v {
		case "": // not set
		case "manager":
			return nil, nil
		case "provider":
			return []provider.Provider{&crashProvider{}}, nil
		case "no-retry":
			return []provider.Provider{&crashProvider{err: badConfigError{errors.New(v)}}}, nil
		default:
			return []provider.Provider{&crashProvider{err: errors.New(v)}}, nil
		}
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

func newMockDevice(ctx context.Context, opt map[string]string,
	monitor provider.Monitor) (Device, error) {
	deviceID, err := GetStringOption("id", opt)
	if err != nil {
		return nil, err
	}
	return &mockDevice{
		id:      deviceID,
		config:  opt,
		isAlive: true,
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
	"input2": {
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
	"cred2": {
		Description: "somedata",
		Required:    false,
	},
	"log-level": {
		Description: "log level for mock datasource",
		Required:    false,
	},
}

type sensorTestCase struct {
	name     string
	substeps []*sensorTestCase

	stateSubResps              []*gnmi.SubscribeResponse
	configSubResps             []*gnmi.SubscribeResponse
	waitForMetadataPreSync     []string
	waitForMetadataPostSync    []string
	expectSet                  []*gnmi.SetRequest
	dynamicConfigs             []*Config
	ignoreDatasourceHeartbeats bool
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
	opts map[string]string, creds map[string]string, logLevel string) []*gnmi.Update {
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
	upds = append(upds, datasourceOptUpdate(id, name, "option", "log-level", logLevel)...)
	return upds
}

func waitMetas(ctx context.Context, t *testing.T,
	metadataCh chan string, expectation []string) error {
	expectMetas := map[string]int{}
	for _, m := range expectation {
		expectMetas[m]++
	}
	t.Logf("Waiting for meta: %v", expectMetas)
	for len(expectMetas) > 0 {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Meta canceled, but still waiting for meta: %v", expectMetas)
		case gotMeta, ok := <-metadataCh:
			if !ok {
				return fmt.Errorf("Meta chan closed, but was waiting for: %v", expectMetas)
			}
			count, ok := expectMetas[gotMeta]
			if !ok {
				return fmt.Errorf("Unexpected meta: %v", gotMeta)
			}
			count--
			t.Logf("got metadata %v", gotMeta)
			if count == 0 {
				delete(expectMetas, gotMeta)
			} else {
				expectMetas[gotMeta] = count
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
	sensor := NewSensor("abc", 100.0,
		WithSensorGNMIClient(gnmic),
		WithSensorHeartbeatInterval(50*time.Millisecond),
		WithSensorFailureRetryBackoffBase(1*time.Second),
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

						// Always send the set response
						gnmic.SetResp <- &gnmi.SetResponse{
							Prefix: pgnmi.PathFromString("ok-dont-care")}

						if !ok {
							if len(tc.expectSet) > 0 {
								return fmt.Errorf("Set channel closed but expected more sets: %v",
									tc.expectSet)
							}
							return nil
						}

						// change date in set request so we can easily match it
						sensorHeartBeat, datasourceHeartBeat := false, false
						for _, u := range setReq.Update {
							if pgnmi.PathMatch(u.Path, pgnmi.Path("streaming-start")) {
								u.Val = agnmi.TypedValue(42)
							}
							if pgnmi.PathMatch(u.Path, pgnmi.Path("last-seen")) {
								u.Val = agnmi.TypedValue(43)
								if len(setReq.Update) == 1 &&
									pgnmi.PathMatch(sensor.statePrefix, setReq.Prefix) {
									sensorHeartBeat = true
								} else if (len(setReq.Update) == 1 || len(setReq.Update) == 2) &&
									setReq.Prefix.Elem[0].Name == "datasource" {
									datasourceHeartBeat = true
								}
							}
						}

						if tc.ignoreDatasourceHeartbeats && datasourceHeartBeat {
							t.Log("skipping datasource heartbeat")
							continue
						}

						if sensorHeartBeat {
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
							lastIndex := len(tc.expectSet) - 1
							tc.expectSet[found] = tc.expectSet[lastIndex]
							tc.expectSet = tc.expectSet[:lastIndex]
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
				t.Error(err)
			}
		})
	}
	close(configStream.SubResp)
	close(configStream.SubReq)
	close(configStream.ErrC)
	cancel()
	if err := errg.Wait(); err != nil && err != context.Canceled {
		t.Error(err)
	}
}

func TestSensor(t *testing.T) {
	testCases := []sensorTestCase{
		{
			name: "Pre-existing datasource",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value1"},
						nil, "LOG_LEVEL_INFO")...),
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
			waitForMetadataPostSync: []string{"123|" +
				"map[id:123 input1:value1 log-level:LOG_LEVEL_INFO]"},
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "Update existing datasource",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value1"},
						map[string]string{"cred1": "credv1"}, "LOG_LEVEL_INFO")...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value2"},
						map[string]string{"cred1": "credv2"}, "LOG_LEVEL_DEBUG")...),
			},
			waitForMetadataPreSync: []string{"123|" +
				"map[cred1:credv1 id:123 input1:value1 log-level:LOG_LEVEL_INFO]"},
			waitForMetadataPostSync: []string{"123|" +
				"map[cred1:credv2 id:123 input1:value2 log-level:LOG_LEVEL_DEBUG]"},
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
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "Update existing datasource 2",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{
							"id":     "123",
							"input1": "value1",
							"input2": "value2"},
						map[string]string{"cred1": "credv1"},
						"LOG_LEVEL_INFO")...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				{
					Response: &gnmi.SubscribeResponse_Update{
						Update: &gnmi.Notification{
							Delete: []*gnmi.Path{
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "option", "input2"), "key"),
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "option", "input2"), "value"),
							},
							Update: []*gnmi.Update{
								pgnmi.Update(pgnmi.PathAppend(
									datasourceOptPath("abc", "xyz", "option", "input1"), "key"),
									agnmi.TypedValue("input1")),
								pgnmi.Update(pgnmi.PathAppend(
									datasourceOptPath("abc", "xyz", "option", "input1"), "value"),
									agnmi.TypedValue("value2")),
							},
						},
					},
				},
			},
			waitForMetadataPreSync: []string{"123|" +
				"map[cred1:credv1 id:123 input1:value1 input2:value2 log-level:LOG_LEVEL_INFO]"},
			waitForMetadataPostSync: []string{"123|" +
				"map[cred1:credv1 id:123 input1:value2 log-level:LOG_LEVEL_INFO]"},
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
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "delete existing datasource config items",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{
							"id":     "123",
							"input1": "value1",
							"input2": "value2"},
						map[string]string{"cred1": "credv1", "cred2": "credv2"},
						"LOG_LEVEL_INFO")...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				{
					Response: &gnmi.SubscribeResponse_Update{
						Update: &gnmi.Notification{
							Delete: []*gnmi.Path{
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "option", "input2"), "key"),
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "option", "input2"), "value"),
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "credential", "cred2"), "key"),
								pgnmi.PathAppend(datasourceOptPath(
									"abc", "xyz", "credential", "cred2"), "value"),
							},
						},
					},
				},
			},
			waitForMetadataPreSync: []string{
				"123|map[cred1:credv1 cred2:credv2 id:123 " +
					"input1:value1 input2:value2 log-level:LOG_LEVEL_INFO]"},
			waitForMetadataPostSync: []string{
				"123|map[cred1:credv1 id:123 input1:value1 log-level:LOG_LEVEL_INFO]"},
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
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "Disable existing datasource",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "input1": "value1"}, nil,
						"LOG_LEVEL_INFO")...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						false, map[string]string{"id": "123", "input1": "value2"}, nil,
						"LOG_LEVEL_INFO")...),
			},
			waitForMetadataPreSync: []string{"123|" +
				"map[id:123 input1:value1 log-level:LOG_LEVEL_INFO]"},
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
						pgnmi.Update(pgnmi.Path("last-error"),
							agnmi.TypedValue("Data source disabled")),
					},
				},
			},
			ignoreDatasourceHeartbeats: true,
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
						true, map[string]string{"id": "123"}, nil,
						"LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123 log-level:LOG_LEVEL_INFO]"},
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
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "Datasource with invalid config should keep others going (Pre-sync)",
			configSubResps: []*gnmi.SubscribeResponse{
				subscribeUpdates(
					datasourceUpdates("config", "abc", "bad1", "invalidtype",
						true, map[string]string{"id": "111"}, nil, "")...),
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123"}, nil, "LOG_LEVEL_INFO")...),
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
			},
			waitForMetadataPreSync: []string{"123|map[id:123 log-level:LOG_LEVEL_INFO]"},
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
							"Data source stopped: "+
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
			ignoreDatasourceHeartbeats: true,
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
						true, map[string]string{"id": "111"}, nil, "")...),
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123"}, nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123 log-level:LOG_LEVEL_INFO]"},
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
							"Data source stopped: "+
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
			ignoreDatasourceHeartbeats: true,
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
						true, map[string]string{"id": "123", "crash": "provider"},
						nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{
				"123|map[crash:provider id:123 log-level:LOG_LEVEL_INFO]",
				"123|map[crash:provider id:123 log-level:LOG_LEVEL_INFO]"},
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
								"fatal error in *device.crashProvider.Run: Crash!. "+
								"Retrying in 1s")),
					},
				},
				// Should eventually be restarted, and we want to see the retry time increase
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
								"fatal error in *device.crashProvider.Run: Crash!. "+
								"Retrying in 2s")),
					},
				},
			},
			ignoreDatasourceHeartbeats: true,
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
						true, map[string]string{"id": "123", "crash": "manager"},
						nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|" +
				"map[crash:manager id:123 log-level:LOG_LEVEL_INFO]"},
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
								"fatal error in Manage: Crash manager!. Retrying in 1s")),
					},
				},
			},
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "Datasource with Manager Provider crash",
			configSubResps: []*gnmi.SubscribeResponse{
				{
					Response: &gnmi.SubscribeResponse_SyncResponse{
						SyncResponse: true,
					},
				},
				subscribeUpdates(
					datasourceUpdates("config", "abc", "xyz", "mock",
						true, map[string]string{"id": "123", "crash": "manager-provider"},
						nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|" +
				"map[crash:manager-provider id:123 log-level:LOG_LEVEL_INFO]"},
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
								"provider *device.crashProvider exiting with error: "+
								"manager-provider. Retrying in 1s")),
					},
				},
			},
			ignoreDatasourceHeartbeats: true,
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
								true, map[string]string{"id": "123"}, nil, "LOG_LEVEL_INFO")...),
					},
					waitForMetadataPostSync: []string{"123|map[id:123 log-level:LOG_LEVEL_INFO]"},
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
					ignoreDatasourceHeartbeats: true,
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
			ignoreDatasourceHeartbeats: true,
		},
		{
			name: "No configs, no sets",
			configSubResps: []*gnmi.SubscribeResponse{
				{Response: &gnmi.SubscribeResponse_SyncResponse{SyncResponse: true}},
			},
			expectSet:                  []*gnmi.SetRequest{},
			ignoreDatasourceHeartbeats: true,
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
						true, map[string]string{"id": "123"}, nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|map[id:123 log-level:LOG_LEVEL_INFO]"},
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
			ignoreDatasourceHeartbeats: false,
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
							"managed": "m1,m2"}, nil, "LOG_LEVEL_INFO")...),
			},
			waitForMetadataPostSync: []string{"123|" +
				"map[id:123 log-level:LOG_LEVEL_INFO managed:m1,m2]"},
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
			ignoreDatasourceHeartbeats: false,
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
						{
							Prefix: datasourcePath("state", "abc", "device-1", ""),
							Update: []*gnmi.Update{
								pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43)),
								pgnmi.Update(pgnmi.Path("streaming-start"), agnmi.TypedValue(42)),
							},
						},
					},
					waitForMetadataPostSync:    []string{"123|map[id:123 input1:value1]"},
					ignoreDatasourceHeartbeats: false,
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
								pgnmi.Update(pgnmi.Path("last-seen"), agnmi.TypedValue(43))},
						},
					},
					ignoreDatasourceHeartbeats: false,
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
					ignoreDatasourceHeartbeats: true,
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

	Register(deviceType, func(ctx context.Context, m map[string]string,
		monitor provider.Monitor) (Device, error) {
		createCalls <- fmt.Sprintf("%v", m)
		return newMockDevice(ctx, m, nil)
	}, mockDeviceOptions)

	// Setup one datasource that we will use for all test cases
	gnmic := &internal.MockClient{
		SubscribeStream: make(chan *internal.MockClientStream),
		SetReq:          make(chan *gnmi.SetRequest, 100),
		SetResp:         make(chan *gnmi.SetResponse),
	}
	close(gnmic.SetResp) // we don't care about these responses, so make it always return nil

	metadataCh := make(chan string, 100)
	sensor := NewSensor("default", 100.0, WithSensorClientFactory(
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
		name                string
		config              datasourceConfig
		expectDeviceCreate  int
		expectRedeploy      int
		redeployBaseBackoff time.Duration
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
				loglevel: logrus.InfoLevel,
			},
			expectDeviceCreate: 1,
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
			expectDeviceCreate: 0,
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
			expectDeviceCreate: 1,
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
				loglevel: logrus.InfoLevel,
			},
			expectDeviceCreate: 1,
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
				loglevel: logrus.InfoLevel,
			},
			expectDeviceCreate: 1,
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
				loglevel: logrus.InfoLevel,
			},
			expectDeviceCreate: 0,
		},
		{
			name: "modify loglevel will not redeploy",
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
				loglevel: logrus.DebugLevel,
			},
			expectDeviceCreate: 0,
		},
		{
			name:                "device run failure should trigger restart",
			redeployBaseBackoff: 10 * time.Millisecond,
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id":    "124",
					"crash": "fail message",
				},
				credential: map[string]string{
					"cred1": "abc",
				},
			},
			expectDeviceCreate: 2,
			// we expect to see 2 redeploys but we stop redeploying after the first,
			// as we only need to check that it repeats once.
			expectRedeploy: 2,
		},
		{
			name:                "device run failure of type no retry should not trigger restart",
			redeployBaseBackoff: 10 * time.Millisecond,
			config: datasourceConfig{
				name:    deviceName,
				typ:     deviceType,
				enabled: true,
				option: map[string]string{
					"id":    "124",
					"crash": "no-retry",
				},
				credential: map[string]string{
					"cred1": "abc",
				},
			},
			expectDeviceCreate: 1,
			expectRedeploy:     0,
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

			if tc.redeployBaseBackoff != 0 {
				ds := sensor.getDatasource(ctx, deviceName)
				ds.failureRetryTimer.BackoffBase = tc.redeployBaseBackoff
				ds.failureRetryTimer.Reset()
			}

			err := sensor.runDatasourceConfig(ctx, deviceName)
			if err != nil {
				t.Fatal(err)
			}

			// Check if we got a create call.
			// Create is synchronous so we should either have it in the channel or not.
			createSeen := 0
			redeploySeen := 0
			lastCheck := time.After(40 * time.Millisecond)
			for {
				select {
				case <-createCalls:
					createSeen++
				case r := <-sensor.redeployDatasource:
					redeploySeen++
					if redeploySeen < tc.expectRedeploy {
						if err := sensor.runDatasourceConfig(ctx, r); err != nil {
							t.Fatal(err)
						}
					}
				case <-lastCheck:
					if createSeen == tc.expectDeviceCreate && redeploySeen == tc.expectRedeploy {
						return // test pass!
					}
					t.Fatalf("Got: createSeen: %d, redeploySeen: %d", createSeen, redeploySeen)
				}
				if createSeen == tc.expectDeviceCreate && redeploySeen == tc.expectRedeploy {
					t.Log("Last check...")
					lastCheck = time.After(40 * time.Millisecond)
				}
				t.Logf("createSeen: %d, redeploySeen: %d", createSeen, redeploySeen)
			}
		})
	}
}

func TestSendPeriodicUpdates(t *testing.T) {
	type mocks struct {
		end      func()
		device   *dmock.MockDevice
		gnmic    *gmock.MockGNMIClient
		cvclient *cvmock.MockCVClient
	}

	verifyLastErrorMatches := func(expect string, set *gnmi.SetRequest) {
		t.Helper()
		// look for last-error update
		for _, u := range set.Update {
			if pgnmi.PathMatch(u.Path, pgnmi.PathFromString("last-error")) {
				if got := u.Val.GetStringVal(); got != expect {
					t.Fatalf("Expected %q but got %q", expect, got)
				}
			}
		}
	}

	expectAll := func(m mocks, calls ...*gomock.Call) {
		calls = append(calls,
			// last Alive call will call end() to cancel the context and finish the test
			m.device.EXPECT().Alive(gomock.Any()).DoAndReturn(
				func(_ context.Context) (bool, error) {
					m.end()
					return false, nil
				}),
			// ignore further gnmic set calls originated after the end()
			m.gnmic.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes())
		gomock.InOrder(calls...)
	}

	for _, tc := range []struct {
		name         string
		expectations func(ctx context.Context, m mocks)
	}{
		{
			name: "shut down without err",
			expectations: func(ctx context.Context, m mocks) {
				expectAll(m,
					m.device.EXPECT().Alive(ctx).Return(false, nil),
					m.gnmic.EXPECT().Set(ctx, gomock.Any()).DoAndReturn(
						func(_ context.Context, set *gnmi.SetRequest,
							_ ...grpc.CallOption) (*gnmi.SetResponse, error) {
							verifyLastErrorMatches("Device not alive", set)
							return nil, nil
						}),
				)
			},
		},
		{
			name: "shut down with err",
			expectations: func(ctx context.Context, m mocks) {
				expectAll(m,
					m.device.EXPECT().Alive(ctx).Return(false, fmt.Errorf("some reason")),
					m.gnmic.EXPECT().Set(ctx, gomock.Any()).DoAndReturn(
						func(_ context.Context, set *gnmi.SetRequest,
							_ ...grpc.CallOption) (*gnmi.SetResponse, error) {
							verifyLastErrorMatches("Device not alive: some reason", set)
							return nil, nil
						}),
				)
			},
		},
		{
			name: "recover from failure should write a back alive message",
			expectations: func(ctx context.Context, m mocks) {
				expectAll(m,
					m.device.EXPECT().Alive(ctx).Return(false, nil),
					m.gnmic.EXPECT().Set(ctx, gomock.Any()).DoAndReturn(
						func(_ context.Context, set *gnmi.SetRequest,
							_ ...grpc.CallOption) (*gnmi.SetResponse, error) {
							verifyLastErrorMatches("Device not alive", set)
							return nil, nil
						}),
					m.device.EXPECT().Alive(ctx).Return(true, nil),
					m.gnmic.EXPECT().Set(ctx, gomock.Any()).DoAndReturn(
						func(_ context.Context, set *gnmi.SetRequest,
							_ ...grpc.CallOption) (*gnmi.SetResponse, error) {
							verifyLastErrorMatches("Device is back alive", set)
							return nil, nil
						}),
					m.cvclient.EXPECT().SendHeartbeat(ctx, true),
				)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			fakeDevice := dmock.NewMockDevice(ctrl)

			datasource := &datasource{
				heartbeatInterval: 1 * time.Millisecond,
				log:               logrus.WithField("test", t.Name()),
			}
			datasource.info = &Info{
				ID:      "test",
				Context: ctx,
				Device:  fakeDevice,
				Config:  &Config{},
			}

			gnmic := gmock.NewMockGNMIClient(ctrl)
			datasource.gnmic = gnmic
			mockCvClient := cvmock.NewMockCVClient(ctrl)
			datasource.cvClient = mockCvClient

			tc.expectations(ctx, mocks{
				device:   fakeDevice,
				gnmic:    gnmic,
				cvclient: mockCvClient,
				end:      cancel,
			})

			if err := datasource.sendPeriodicUpdates(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}
