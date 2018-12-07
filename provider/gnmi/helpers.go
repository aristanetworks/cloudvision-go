// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmi

import (
	aaagrpc "arista/aaa/grpc"
	"arista/aaa/provider/rote"
	apiserver "arista/aeris/apiserver/client"
	"arista/gopenconfig/event/router"
	"arista/gopenconfig/gnmi/server"
	"arista/gopenconfig/model"
	"arista/gopenconfig/model/notifier"
	"arista/gopenconfig/modules"
	"arista/gopenconfig/yang"
	"arista/types"
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/aristanetworks/glog"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// convertPath returns all but the last element of the gNMI path as aeris
// path, and the last element as the update key.
func convertPath(gnmiPath []*gnmi.PathElem) (key.Path, key.Key) {
	aerisPath := key.Path{}
	for _, elm := range gnmiPath {
		aerisPath = path.Append(aerisPath, key.New(elm.Name))
		if len(elm.Key) > 0 {
			keyMap := map[string]interface{}{}
			for k, v := range elm.Key {
				keyMap[k] = v
			}
			aerisPath = path.Append(aerisPath, keyMap)
		}
	}
	if len(aerisPath) == 0 {
		return aerisPath, nil
	}
	return aerisPath[:len(aerisPath)-1], aerisPath[len(aerisPath)-1]
}

// Unmarshal will return an interface representing the supplied value.
func Unmarshal(val *gnmi.TypedValue) interface{} {
	switch v := val.GetValue().(type) {
	case *gnmi.TypedValue_StringVal:
		return v.StringVal
	case *gnmi.TypedValue_JsonIetfVal:
		return v.JsonIetfVal
	case *gnmi.TypedValue_JsonVal:
		return v.JsonVal
	case *gnmi.TypedValue_IntVal:
		return v.IntVal
	case *gnmi.TypedValue_UintVal:
		return v.UintVal
	case *gnmi.TypedValue_BoolVal:
		return v.BoolVal
	case *gnmi.TypedValue_BytesVal:
		return agnmi.StrVal(val)
	case *gnmi.TypedValue_DecimalVal:
		d := v.DecimalVal
		return float64(d.Digits) / math.Pow(10, float64(d.Precision))
	case *gnmi.TypedValue_FloatVal:
		return v.FloatVal
	case *gnmi.TypedValue_LeaflistVal:
		ret := []interface{}{}
		for _, val := range v.LeaflistVal.Element {
			ret = append(ret, Unmarshal(val))
		}
		return ret
	case *gnmi.TypedValue_AsciiVal:
		return v.AsciiVal
	case *gnmi.TypedValue_AnyVal:
		return v.AnyVal.String()
	default:
		panic(v)
	}
}

func convertNotif(notif *gnmi.Notification) []types.Notification {

	var ret []types.Notification

	gnmiPath := []*gnmi.PathElem{&gnmi.PathElem{Name: "OpenConfig"}}
	if notif.Prefix != nil {
		gnmiPath = append(gnmiPath, notif.Prefix.Elem...)
	}
	for _, update := range notif.Update {
		gnmiUpdatePath := append(gnmiPath, update.Path.Elem...)
		aerisUpdatePath, updateKey := convertPath(gnmiUpdatePath)
		ret = append(ret, types.NewNotification(time.Now(), aerisUpdatePath, nil,
			map[key.Key]interface{}{updateKey: Unmarshal(update.Val)}))
	}
	for _, delete := range notif.Delete {
		gnmiDeletePath := append(gnmiPath, delete.Elem...)
		aerisDeletePath, deleteKey := convertPath(gnmiDeletePath)
		ret = append(ret,
			types.NewNotification(time.Now(), aerisDeletePath, []key.Key{deleteKey}, nil))
	}
	return ret
}

// EmitNotif converts a gNMI notification into a series of
// types.Notifications and puts those on the provided notif channel,
// adding any necessary pointer path notifications.
func EmitNotif(notif *gnmi.Notification, ch chan<- types.Notification) {
	transformer := apiserver.NewPathPointerCreator(true)
	for _, update := range convertNotif(notif) {
		for _, notif := range transformer.Transform(update) {
			ch <- types.NewNotification(
				notif.Timestamp(), notif.Path(), notif.Deletes(), notif.Updates())
		}
	}
}

func doYANGModuleSetup(ctx context.Context, yangPaths []string) (context.Context, error) {
	yangModules := modules.New()

	if err := yangModules.AddRootPath(yangPaths...); err != nil {
		return ctx, fmt.Errorf("error adding YANG directories: %s", err)
	}

	if err := yangModules.Import(yang.DefaultModules...); err != nil {
		return nil, fmt.Errorf("error importing YANG modules: %v", err)
	}

	return modules.NewContext(ctx, yangModules), nil
}

func doDatastoresSetup(ctx context.Context) (context.Context,
	model.Datastores, error) {
	ms, _ := modules.FromContext(ctx)
	runningConfig := model.New(model.DSRunning)
	n := notifier.New()
	go n.Run(ctx)
	runningConfig.RootNode().SetNotifier(n)
	runningConfig.RootNode().Metadata().AddField("ReadOnlyWritable", true)
	datastores := model.NewDatastores()
	ctx = context.WithValue(ctx, model.DatastoresKey, datastores)
	yangErrs := model.PopulateDataModel(ms, runningConfig, yang.DefaultModules...)
	if len(yangErrs) > 0 {
		return ctx, nil, fmt.Errorf("YANG import errors: %v", yangErrs)
	}
	err := datastores.SetDatastore(runningConfig)
	return ctx, datastores, err
}

// openConfigSetUpDatastores sets up an OpenConfig tree using the
// default modules.
func openConfigSetUpDatastores(ctx context.Context, yangPaths []string) (context.Context,
	model.Datastores, error) {
	ctx, err := doYANGModuleSetup(ctx, yangPaths)
	if err != nil {
		return ctx, nil, err
	}

	return doDatastoresSetup(ctx)
}

// subscribeStream implements the stream interface needed by the gNMI
// server. All it needs to be able to do is convert the server's
// outgoing updates to types.Notifications.
type subscribeStream struct {
	req  chan *gnmi.SubscribeRequest
	resp chan<- types.Notification
	ctx  context.Context
}

func (f *subscribeStream) Send(resp *gnmi.SubscribeResponse) error {
	select {
	case <-f.ctx.Done():
		return f.ctx.Err()
	default:
		switch r := resp.Response.(type) {
		case *gnmi.SubscribeResponse_Error:
			glog.Errorf("gNMI SubscribeResponse error: %v", r.Error.Message)
		case *gnmi.SubscribeResponse_Update:
			EmitNotif(r.Update, f.resp)
		case *gnmi.SubscribeResponse_SyncResponse:
			if !r.SyncResponse {
				glog.Errorf("gNMI sync failed")
			}
		default:
			glog.Errorf("Unexpected gNMI SubscribeResponse type: %v", r)
		}
	}
	return nil
}

func (f *subscribeStream) Recv() (*gnmi.SubscribeRequest, error) {
	select {
	case req := <-f.req:
		if req == nil {
			return nil, io.EOF
		}
		return req, nil
	case <-f.ctx.Done():
		return nil, f.ctx.Err()
	}
}

func (f *subscribeStream) SetHeader(metadata.MD) error {
	panic("not implemented")
}
func (f *subscribeStream) SendHeader(metadata.MD) error {
	panic("not implemented")
}
func (f *subscribeStream) SetTrailer(metadata.MD) {
	panic("not implemented")
}
func (f *subscribeStream) Context() context.Context {
	return f.ctx
}
func (f *subscribeStream) SendMsg(m interface{}) error {
	panic("not implemented")
}
func (f *subscribeStream) RecvMsg(m interface{}) error {
	panic("not implemented")
}

// Create a gNMI server with an existing datastores object.
func newServerWithStore(stores model.Datastores,
	ms *modules.Modules) (gnmi.GNMIServer, *aaagrpc.Handler) {
	p := rote.New(rote.WithAuthenPass(), rote.WithAuthorPass())
	yangRouter := router.New(nil)
	ctx := router.NewContext(context.Background(), yangRouter)
	ctx = model.NewLocksContext(ctx, model.NewLocks())
	if ms != nil {
		ctx = modules.NewContext(ctx, ms)
	}

	options := []server.Option{
		server.WithContext(ctx),
		server.WithAcctProvider(p),
		server.WithAuthenProvider(p),
		server.WithAuthzProvider(p),
		server.WithDatastores(stores),
	}
	return server.New(options...)
}

// Given a subscription stream, subscribe to all updates.
func subscribeAll(stream *subscribeStream) {
	stream.req <- &gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: &gnmi.SubscriptionList{
				Subscription: []*gnmi.Subscription{&gnmi.Subscription{
					Path: &gnmi.Path{Element: nil},
				}},
			},
		},
	}
}

// Server creates a new OpenConfig tree and returns a
// gnmi.GNMIServer that operates on that tree.
func Server(ctx context.Context, ch chan<- types.Notification,
	errc chan error, yangPaths []string) (context.Context, gnmi.GNMIServer, error) {
	ctx, datastores, err := openConfigSetUpDatastores(ctx, yangPaths)
	if err != nil {
		return ctx, nil, fmt.Errorf("Error setting up data tree: %v", err)
	}
	ctx, server := gNMIServer(ctx, datastores)
	ctx, err = gNMIStreamUpdates(ctx, server, ch, errc)
	if err != nil {
		return ctx, nil, fmt.Errorf("Errorf setting up notif stream: %v", err)
	}
	return ctx, server, nil
}

type gnmiclient struct {
	s gnmi.GNMIServer
}

// Client takes a gnmi.GNMIServer and returns a gnmi.GNMIClient
// that will translate client Sets to server Sets without doing any
// RPC.
func Client(s gnmi.GNMIServer) gnmi.GNMIClient {
	return &gnmiclient{s: s}
}

func (g *gnmiclient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return g.s.Capabilities(ctx, in)
}

func (g *gnmiclient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return g.s.Get(ctx, in)
}

func (g *gnmiclient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	return g.s.Set(ctx, in)
}

func (g *gnmiclient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	panic("not implemented")
}

func gNMIServer(ctx context.Context, datastores model.Datastores) (context.Context,
	gnmi.GNMIServer) {
	ms, _ := modules.FromContext(ctx)
	server, aaa := newServerWithStore(datastores, ms)
	ctx = aaa.TagConn(ctx, nil)
	return ctx, server
}

// gNMIStreamUpdates takes a data tree and channel of types.Notifications
// and streams all changes to the tree into the channel.
func gNMIStreamUpdates(ctx context.Context, server gnmi.GNMIServer,
	ch chan<- types.Notification, errc chan error) (context.Context, error) {
	// Set up stream
	stream := &subscribeStream{
		req:  make(chan *gnmi.SubscribeRequest),
		resp: ch,
		ctx:  ctx,
	}

	// Subscribe to /
	go func() {
		err := server.Subscribe(stream)
		if err != nil {
			errc <- err
		}
	}()
	subscribeAll(stream)
	return ctx, nil
}

// TestClient is a pass-through gNMI client for testing.
type TestClient struct {
	Out chan *gnmi.SetRequest
}

// Capabilities is not implemented.
func (g *TestClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest,
	opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	panic("not implemented")
}

// Get is not implemented.
func (g *TestClient) Get(ctx context.Context, in *gnmi.GetRequest,
	opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	panic("not implemented")
}

// Set pipes the specified SetRequest out to the TestClient's
// SetRequest channel.
func (g *TestClient) Set(ctx context.Context, in *gnmi.SetRequest,
	opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	g.Out <- in
	return nil, nil
}

// Subscribe is not implemented.
func (g *TestClient) Subscribe(ctx context.Context,
	opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	panic("not implemented")
}

// Path returns a gnmi.Path given a set of elements.
func Path(element ...string) *gnmi.Path {
	p, err := agnmi.ParseGNMIElements(element)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse GNMI elements: %v", element))
	}
	p.Element = nil
	return p
}

// gNMI TypedValues: Everything is converted to a JsonVal for now
// because those code paths are more mature in the gopenconfig code.
func jsonValue(v interface{}) *gnmi.TypedValue {
	vb := []byte(fmt.Sprintf(`"%v"`, v))
	return &gnmi.TypedValue{Value: &gnmi.TypedValue_JsonVal{JsonVal: vb}}
}

// Strval returns a gnmi.TypedValue from a string.
func Strval(s string) *gnmi.TypedValue {
	return jsonValue(s)
}

// Uintval returns a gnmi.TypedValue from a uint64.
func Uintval(u uint64) *gnmi.TypedValue {
	return jsonValue(u)
}

// Intval returns a gnmi.TypedValue from an int64.
func Intval(i int64) *gnmi.TypedValue {
	return jsonValue(i)
}

// Boolval returns a gnmi.TypedValue from a bool.
func Boolval(b bool) *gnmi.TypedValue {
	return jsonValue(b)
}

// Update creates a gNMI.Update.
func Update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return &gnmi.Update{
		Path: path,
		Val:  val,
	}
}

// A PollFn polls a target device and returns a gNMI SetRequest.
type PollFn func() (*gnmi.SetRequest, error)

func pollOnce(ctx context.Context, client gnmi.GNMIClient,
	poller PollFn) error {
	setreq, err := poller()
	if err != nil {
		return err
	}
	if setreq != nil {
		_, err = client.Set(ctx, setreq)
	}
	return err
}

// PollForever takes a polling function that performs a
// complete update of some part of the OpenConfig tree and calls it
// at the specified interval.
func PollForever(ctx context.Context, client gnmi.GNMIClient,
	interval time.Duration, poller PollFn, errc chan error) {

	// Poll immediately.
	if err := pollOnce(ctx, client, poller); err != nil {
		errc <- err
	}

	// Poll at intervals forever.
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if err := pollOnce(ctx, client, poller); err != nil {
				errc <- err
			}
		case <-ctx.Done():
			return
		}
	}
}

// Helpers for creating gNMI paths for places of interest in the
// OpenConfig tree.

// ListWithKey formats a gNMI keyed list and key as a string.
func ListWithKey(listName, keyName, key string) string {
	return fmt.Sprintf("%s[%s=%s]", listName, keyName, key)
}

// Interface paths of interest:

// IntfPath returns an interface path.
func IntfPath(intfName, leafName string) *gnmi.Path {
	return Path("interfaces", ListWithKey("interface", "name", intfName),
		leafName)
}

// IntfConfigPath returns an interface config path.
func IntfConfigPath(intfName, leafName string) *gnmi.Path {
	return Path("interfaces", ListWithKey("interface", "name", intfName),
		"config", leafName)
}

// IntfStatePath returns an interface state path.
func IntfStatePath(intfName, leafName string) *gnmi.Path {
	return Path("interfaces", ListWithKey("interface", "name", intfName),
		"state", leafName)
}

// IntfStateCountersPath returns an interface state counters path.
func IntfStateCountersPath(intfName, leafName string) *gnmi.Path {
	return Path("interfaces", ListWithKey("interface", "name", intfName),
		"state", "counters", leafName)
}

// LLDP paths of interest:

// LldpStatePath returns an LLDP state path.
func LldpStatePath(leafName string) *gnmi.Path {
	return Path("lldp", "state", leafName)
}

// LldpIntfPath returns an LLDP interface path.
func LldpIntfPath(intfName, leafName string) *gnmi.Path {
	return Path("lldp", "interfaces", ListWithKey("interface", "name",
		intfName), leafName)
}

// LldpIntfConfigPath returns an LLDP interface config path.
func LldpIntfConfigPath(intfName, leafName string) *gnmi.Path {
	return Path("lldp", "interfaces", ListWithKey("interface", "name",
		intfName), "config", leafName)
}

// LldpIntfStatePath returns an LLDP interface state path.
func LldpIntfStatePath(intfName, leafName string) *gnmi.Path {
	return Path("lldp", "interfaces", ListWithKey("interface", "name",
		intfName), "state", leafName)
}

// LldpIntfCountersPath returns an LLDP interface counters path.
func LldpIntfCountersPath(intfName, leafName string) *gnmi.Path {
	return Path("lldp", "interfaces", ListWithKey("interface", "name",
		intfName), "state", "counters", leafName)
}

// LldpNeighborStatePath returns an LLDP neighbor state path.
func LldpNeighborStatePath(intfName, id, leafName string) *gnmi.Path {
	return Path("lldp", "interfaces", ListWithKey("interface", "name",
		intfName), "neighbors", ListWithKey("neighbor", "id", id),
		"state", leafName)
}
