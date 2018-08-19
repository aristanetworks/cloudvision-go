// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/provider"
	"arista/schema"
	"arista/types"
	"context"
	"math"
	"time"

	apiserver "arista/aeris/apiserver/client"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/gnmi"
	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

type gnmiProvider struct {
	provider.ReadOnly
	// Closed when we're done initialization
	ready chan struct{}
	// Closed when we want to stop Run()
	done chan struct{}

	client   pb.GNMIClient
	cfg      *gnmi.Config
	paths    []string
	typeDefs *schema.Schema
}

func (p *gnmiProvider) WaitForNotification() {
	<-p.ready
}

func (p *gnmiProvider) Stop() {
	<-p.ready
	close(p.done)
}

func (p *gnmiProvider) Run(s *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	p.typeDefs = s
	respChan := make(chan *pb.SubscribeResponse)
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, p.cfg)
	go gnmi.Subscribe(ctx, p.client, gnmi.SplitPaths(p.paths), respChan, errChan)
	close(p.ready)
	for {
		select {
		case <-p.done:
			return
		case response := <-respChan:
			switch resp := response.Response.(type) {
			case *pb.SubscribeResponse_Error:
				// Not sure if this is recoverable so it doesn't return and hope things get better
				glog.Errorf("gNMI SubscribeResponse Error: %v", resp.Error.Message)
			case *pb.SubscribeResponse_SyncResponse:
				if !resp.SyncResponse {
					glog.Errorf("gNMI sync failed")
				}
			case *pb.SubscribeResponse_Update:
				emitNotif(resp.Update, ch)
			}
		case err := <-errChan:
			glog.Errorf("Error from gNMI connection: %v", err)
			return
		}
	}
}

func emitNotif(notif *pb.Notification, ch chan<- types.Notification) {
	transformer := apiserver.NewPathPointerCreator(true)
	for _, update := range convertNotif(notif) {
		for _, notif := range transformer.Transform(update) {
			ch <- types.NewNotification(
				notif.Timestamp(), notif.Path(), notif.Deletes(), notif.Updates())
		}
	}
}

// convertPath returns all but the last element of the gNMI path as aeris path, and
// the last element as the update key
func convertPath(gnmiPath []*pb.PathElem) (key.Path, key.Key) {
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

func convertNotif(notif *pb.Notification) []types.Notification {

	var ret []types.Notification

	gnmiPath := []*pb.PathElem{&pb.PathElem{Name: "OpenConfig"}}
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

// NewGNMIProvider returns a read-only gNMI provider.
func NewGNMIProvider(client pb.GNMIClient, cfg *gnmi.Config, paths []string) provider.Provider {
	return &gnmiProvider{
		ready:  make(chan struct{}),
		done:   make(chan struct{}),
		client: client,
		cfg:    cfg,
		paths:  paths,
	}
}

// Unmarshal will return an interface representing the supplied value.
func Unmarshal(val *pb.TypedValue) interface{} {
	switch v := val.GetValue().(type) {
	case *pb.TypedValue_StringVal:
		return v.StringVal
	case *pb.TypedValue_JsonIetfVal:
		return v.JsonIetfVal
	case *pb.TypedValue_JsonVal:
		return v.JsonVal
	case *pb.TypedValue_IntVal:
		return v.IntVal
	case *pb.TypedValue_UintVal:
		return v.UintVal
	case *pb.TypedValue_BoolVal:
		return v.BoolVal
	case *pb.TypedValue_BytesVal:
		return gnmi.StrVal(val)
	case *pb.TypedValue_DecimalVal:
		d := v.DecimalVal
		return float64(d.Digits) / math.Pow(10, float64(d.Precision))
	case *pb.TypedValue_FloatVal:
		return v.FloatVal
	case *pb.TypedValue_LeaflistVal:
		ret := []interface{}{}
		for _, val := range v.LeaflistVal.Element {
			ret = append(ret, Unmarshal(val))
		}
		return ret
	case *pb.TypedValue_AsciiVal:
		return v.AsciiVal
	case *pb.TypedValue_AnyVal:
		return v.AnyVal.String()
	default:
		panic(v)
	}
}
