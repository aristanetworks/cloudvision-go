// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/types"
	"arista/util"
	"reflect"
	"testing"
	"time"

	"github.com/aristanetworks/goarista/key"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

func makeUpdate(path []string, val string) *pb.Update {
	ret := &pb.Update{}
	ret.Path = makePath(path)
	ret.Val = &pb.TypedValue{
		Value: &pb.TypedValue_StringVal{StringVal: val}}

	return ret
}

// Make the vals to strings to make things easier. We have other test for marshaling
func makeGNMINotif(prefix []string,
	updates, deletes [][]string, updateVals []string) *pb.Notification {

	ret := &pb.Notification{Update: []*pb.Update{}, Delete: []*pb.Path{}}
	ret.Prefix = makePath(prefix)
	for i, update := range updates {
		ret.Update = append(ret.Update, makeUpdate(update, updateVals[i]))
	}
	for _, delete := range deletes {
		ret.Delete = append(ret.Delete, makePath(delete))
	}
	return ret
}

func TestUnmarshal(t *testing.T) {
	anyBytes, err := proto.Marshal(&pb.ModelData{Name: "foobar"})
	if err != nil {
		t.Fatal(err)
	}
	anyMessage := &any.Any{TypeUrl: "gnmi/ModelData", Value: anyBytes}
	anyString := proto.CompactTextString(anyMessage)

	for name, tc := range map[string]struct {
		val *pb.TypedValue
		exp interface{}
	}{
		"StringVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_StringVal{StringVal: "foobar"}},
			exp: "foobar",
		},
		"IntVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_IntVal{IntVal: -42}},
			exp: int64(-42),
		},
		"UintVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_UintVal{UintVal: 42}},
			exp: uint64(42),
		},
		"BoolVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_BoolVal{BoolVal: true}},
			exp: true,
		},
		"BytesVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_BytesVal{BytesVal: []byte{0xde, 0xad}}},
			exp: "3q0=",
		},
		"FloatVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_FloatVal{FloatVal: 3.14}},
			exp: float32(3.14),
		},
		"DecimalVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_DecimalVal{
					DecimalVal: &pb.Decimal64{Digits: 314, Precision: 2},
				}},
			exp: float64(3.14),
		},
		"LeafListVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_LeaflistVal{
					LeaflistVal: &pb.ScalarArray{Element: []*pb.TypedValue{
						&pb.TypedValue{Value: &pb.TypedValue_BoolVal{BoolVal: true}},
						&pb.TypedValue{Value: &pb.TypedValue_AsciiVal{AsciiVal: "foobar"}},
					}},
				}},
			exp: []interface{}{true, "foobar"},
		},
		"AnyVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_AnyVal{AnyVal: anyMessage}},
			exp: anyString,
		},
		"JsonVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_JsonVal{JsonVal: []byte(`{"foo":"bar"}`)}},
			exp: []byte(`{"foo":"bar"}`),
		},
		"JsonIetfVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"foo":"bar"}`)}},
			exp: []byte(`{"foo":"bar"}`),
		},
		"AsciiVal": {
			val: &pb.TypedValue{
				Value: &pb.TypedValue_AsciiVal{AsciiVal: "foobar"}},
			exp: "foobar",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := Unmarshal(tc.val)
			if !reflect.DeepEqual(got, tc.exp) {
				t.Errorf("Expected: %v:%T Got: %v:%T", tc.exp, tc.exp, got, got)
			}
		})
	}
}

func TestConvertPath(t *testing.T) {
	for name, tc := range map[string]struct {
		gnmiPath  []*pb.PathElem
		aerisPath key.Path
		updateKey key.Key
	}{
		"Empty path": {
			gnmiPath:  []*pb.PathElem{},
			aerisPath: key.Path{},
			updateKey: nil,
		},
		"Simple one path": {
			gnmiPath:  []*pb.PathElem{&pb.PathElem{Name: "simple"}},
			aerisPath: key.Path{},
			updateKey: key.New("simple"),
		},
		"Simple two path": {
			gnmiPath: []*pb.PathElem{&pb.PathElem{Name: "simple"},
				&pb.PathElem{Name: "update"}},
			aerisPath: util.StringsToPath([]string{"simple"}),
			updateKey: key.New("update"),
		},
		"Path with update at the end": {
			gnmiPath: []*pb.PathElem{&pb.PathElem{Name: "simple"},
				&pb.PathElem{Name: "update", Key: map[string]string{"a": "x", "b": "y"}}},
			aerisPath: util.StringsToPath([]string{"simple", "update"}),
			updateKey: key.New(map[string]interface{}{"a": "x", "b": "y"}),
		},
		"Path with update at the middle": {
			gnmiPath: []*pb.PathElem{&pb.PathElem{Name: "simple"},
				&pb.PathElem{Name: "update", Key: map[string]string{"a": "x", "b": "y"}},
				&pb.PathElem{Name: "end"}},
			aerisPath: key.Path{key.New("simple"),
				key.New("update"), key.New(map[string]interface{}{"a": "x", "b": "y"})},
			updateKey: key.New("end"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			path, key := convertPath(tc.gnmiPath)
			if !reflect.DeepEqual(key, tc.updateKey) {
				t.Errorf("Update key mismatches: Expected: %v:%T Got: %v:%T",
					tc.updateKey, tc.updateKey, key, key)
			}
			if !reflect.DeepEqual(path, tc.aerisPath) {
				t.Errorf("Aeris path mismatches: Expected: %v:%T Got: %v:%T",
					tc.aerisPath, tc.aerisPath, path, path)
			}
		})

	}
}

func makePath(path []string) *pb.Path {
	ret := &pb.Path{Elem: []*pb.PathElem{}}
	for _, comp := range path {
		pathElem := &pb.PathElem{Name: comp}
		ret.Elem = append(ret.Elem, pathElem)
	}
	return ret
}

func TestConvertNotif(t *testing.T) {
	for name, tc := range map[string]struct {
		gnmiNotif   *pb.Notification
		aerisNotifs []types.Notification
	}{
		"Empty notification": {
			gnmiNotif:   &pb.Notification{},
			aerisNotifs: nil,
		},
		"Notification with Update": {
			gnmiNotif: makeGNMINotif(
				[]string{"prefix"},
				[][]string{[]string{"simple", "update"}},
				nil,
				[]string{"val"}),
			aerisNotifs: []types.Notification{types.NewNotification(
				time.Now(),
				util.StringsToPath([]string{"OpenConfig", "prefix", "simple"}),
				nil,
				map[key.Key]interface{}{key.New("update"): "val"})},
		},
		"Notification with Delete": {
			gnmiNotif: makeGNMINotif(
				[]string{"prefix"},
				nil,
				[][]string{[]string{"simple", "delete"}},
				nil),
			aerisNotifs: []types.Notification{types.NewNotification(
				time.Now(),
				util.StringsToPath([]string{"OpenConfig", "prefix", "simple"}),
				[]key.Key{key.New("delete")},
				nil)},
		},
	} {
		t.Run(name, func(t *testing.T) {
			notifs := convertNotif(tc.gnmiNotif)
			for i, notif := range notifs {
				if !reflect.DeepEqual(notif.Path(), tc.aerisNotifs[i].Path()) {
					t.Errorf("Notification pathes mismatch: Expected %v:%T Got %v:%T",
						tc.aerisNotifs[i].Path(), tc.aerisNotifs[i].Path(),
						notif.Path(), notif.Path())
				}
				if !reflect.DeepEqual(notif.Updates(), tc.aerisNotifs[i].Updates()) {
					t.Errorf("Notification updates mismatch: Expected %v:%T Got %v:%T",
						tc.aerisNotifs[i].Updates(), tc.aerisNotifs[i].Updates(),
						notif.Updates(), notif.Updates())
				}
				if !reflect.DeepEqual(notif.Deletes(), tc.aerisNotifs[i].Deletes()) {
					t.Errorf("Notification deletes mismatch: Expected %v:%T Got %v:%T",
						tc.aerisNotifs[i].Deletes(), tc.aerisNotifs[i].Deletes(),
						notif.Deletes(), notif.Deletes())
				}
			}
		})
	}
}
