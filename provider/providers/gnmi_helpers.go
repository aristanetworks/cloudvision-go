// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	aaagrpc "arista/aaa/grpc"
	"arista/aaa/provider/rote"
	apiserver "arista/aeris/apiserver/client"
	"arista/gopenconfig/event/router"
	"arista/gopenconfig/gnmi/server"
	"arista/gopenconfig/model"
	"arista/gopenconfig/model/base"
	"arista/gopenconfig/model/node"
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

// GNMIEmitNotif converts a gNMI notification into a series of
// types.Notifications and puts those on the provided notif channel,
// adding any necessary pointer path notifications.
func GNMIEmitNotif(notif *gnmi.Notification, ch chan<- types.Notification) {
	transformer := apiserver.NewPathPointerCreator(true)
	for _, update := range convertNotif(notif) {
		for _, notif := range transformer.Transform(update) {
			ch <- types.NewNotification(
				notif.Timestamp(), notif.Path(), notif.Deletes(), notif.Updates())
		}
	}
}

// XXX_jcr: Pass in a module list?
func doYANGModuleSetup(ctx context.Context) (context.Context, error) {
	yangModules := modules.New()
	yangRoot := []string{"../../gopenconfig/yang/github.com"}

	if err := yangModules.AddRootPath(yangRoot...); err != nil {
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
func openConfigSetUpDatastores(ctx context.Context) (context.Context,
	model.Datastores, error) {
	ctx, err := doYANGModuleSetup(ctx)
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
			GNMIEmitNotif(r.Update, f.resp)
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

// gNMIStreamUpdates takes a data tree and channel of types.Notifications
// and streams all changes to the tree into the channel.
func gNMIStreamUpdates(ctx context.Context, datastores model.Datastores,
	ch chan<- types.Notification, errc chan error) (context.Context, error) {
	ms, _ := modules.FromContext(ctx)
	server, aaa := newServerWithStore(datastores, ms)
	ctx = aaa.TagConn(ctx, nil)
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

// Return the tree root from ctx.
func treeRoot(ctx context.Context) node.Node {
	return model.GetRootYANGNode(ctx)
}

// Lock the YANG tree.
func treeLock(ctx context.Context) {
	models := model.DatastoresFromCtx(ctx)
	ds := models.Datastore(model.DSRunning)
	model.LockDatastore(ds)
}

// Unlock the YANG tree.
func treeUnlock(ctx context.Context) {
	models := model.DatastoresFromCtx(ctx)
	ds := models.Datastore(model.DSRunning)
	model.UnlockDatastore(ds)
}

// OpenConfigNotifyingTree sets up an OpenConfig data tree and a
// gNMI server for streaming out updates to the tree (converted to
// types.Notifications).
func OpenConfigNotifyingTree(ctx context.Context, ch chan<- types.Notification,
	errc chan error) (context.Context, error) {
	ctx, datastores, err := openConfigSetUpDatastores(ctx)
	if err != nil {
		return ctx, fmt.Errorf("Error setting up data tree: %v", err)
	}

	ctx, err = gNMIStreamUpdates(ctx, datastores, ch, errc)
	if err != nil {
		return ctx, fmt.Errorf("Errorf setting up notif stream: %v", err)
	}
	return ctx, nil
}

// Given a container, setLeaf sets the leaf specified by attr.
func setLeaf(ctr base.Container, attr string, val interface{}) error {
	if ctr == nil {
		return fmt.Errorf("container does not exist")
	}
	leaf := ctr.Leaf(attr)
	if leaf == nil {
		return fmt.Errorf("failed creating leaf %v", attr)
	}
	return leaf.Set(val)
}

// XXX_jcr: Need to add delete APIs as well.

// OpenConfigUpdateLeaf creates a container at the specified path and
// updates the indicated leaf.
func OpenConfigUpdateLeaf(ctx context.Context, path node.Path, leafName string,
	val interface{}) error {
	root := treeRoot(ctx)
	treeLock(ctx)
	defer treeUnlock(ctx)
	n, err := model.CreateNodes(root, path)
	if err != nil {
		return err
	}
	c, ok := n.(base.Container)
	if !ok {
		return fmt.Errorf("Provided path %v is not a container", path)
	}
	return setLeaf(c, leafName, val)
}
