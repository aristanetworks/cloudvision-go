// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmi

import (
	"arista/test/notiftest"
	"arista/types"
	"arista/util"
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func makeUpdate(path []string, val string) *gnmi.Update {
	ret := &gnmi.Update{}
	ret.Path = makePath(path)
	ret.Val = &gnmi.TypedValue{
		Value: &gnmi.TypedValue_StringVal{StringVal: val}}

	return ret
}

// Make the vals to strings to make things easier. We have other test for marshaling
func makeGNMINotif(prefix []string,
	updates, deletes [][]string, updateVals []string) *gnmi.Notification {

	ret := &gnmi.Notification{Update: []*gnmi.Update{}, Delete: []*gnmi.Path{}}
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
	anyBytes, err := proto.Marshal(&gnmi.ModelData{Name: "foobar"})
	if err != nil {
		t.Fatal(err)
	}
	anyMessage := &any.Any{TypeUrl: "gnmi/ModelData", Value: anyBytes}
	anyString := proto.CompactTextString(anyMessage)

	for name, tc := range map[string]struct {
		val *gnmi.TypedValue
		exp interface{}
	}{
		"StringVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_StringVal{StringVal: "foobar"}},
			exp: "foobar",
		},
		"IntVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_IntVal{IntVal: -42}},
			exp: int64(-42),
		},
		"UintVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_UintVal{UintVal: 42}},
			exp: uint64(42),
		},
		"BoolVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_BoolVal{BoolVal: true}},
			exp: true,
		},
		"BytesVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_BytesVal{BytesVal: []byte{0xde, 0xad}}},
			exp: "3q0=",
		},
		"FloatVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_FloatVal{FloatVal: 3.14}},
			exp: float32(3.14),
		},
		"DecimalVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_DecimalVal{
					DecimalVal: &gnmi.Decimal64{Digits: 314, Precision: 2},
				}},
			exp: float64(3.14),
		},
		"LeafListVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_LeaflistVal{
					LeaflistVal: &gnmi.ScalarArray{Element: []*gnmi.TypedValue{
						&gnmi.TypedValue{Value: &gnmi.TypedValue_BoolVal{BoolVal: true}},
						&gnmi.TypedValue{Value: &gnmi.TypedValue_AsciiVal{AsciiVal: "foobar"}},
					}},
				}},
			exp: []interface{}{true, "foobar"},
		},
		"AnyVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_AnyVal{AnyVal: anyMessage}},
			exp: anyString,
		},
		"JsonVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_JsonVal{JsonVal: []byte(`{"foo":"bar"}`)}},
			exp: []byte(`{"foo":"bar"}`),
		},
		"JsonIetfVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"foo":"bar"}`)}},
			exp: []byte(`{"foo":"bar"}`),
		},
		"AsciiVal": {
			val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_AsciiVal{AsciiVal: "foobar"}},
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
		gnmiPath  []*gnmi.PathElem
		aerisPath key.Path
		updateKey key.Key
	}{
		"Empty path": {
			gnmiPath:  []*gnmi.PathElem{},
			aerisPath: key.Path{},
			updateKey: nil,
		},
		"Simple one path": {
			gnmiPath:  []*gnmi.PathElem{&gnmi.PathElem{Name: "simple"}},
			aerisPath: key.Path{},
			updateKey: key.New("simple"),
		},
		"Simple two path": {
			gnmiPath: []*gnmi.PathElem{&gnmi.PathElem{Name: "simple"},
				&gnmi.PathElem{Name: "update"}},
			aerisPath: util.StringsToPath([]string{"simple"}),
			updateKey: key.New("update"),
		},
		"Path with update at the end": {
			gnmiPath: []*gnmi.PathElem{&gnmi.PathElem{Name: "simple"},
				&gnmi.PathElem{Name: "update", Key: map[string]string{"a": "x", "b": "y"}}},
			aerisPath: util.StringsToPath([]string{"simple", "update"}),
			updateKey: key.New(map[string]interface{}{"a": "x", "b": "y"}),
		},
		"Path with update at the middle": {
			gnmiPath: []*gnmi.PathElem{&gnmi.PathElem{Name: "simple"},
				&gnmi.PathElem{Name: "update", Key: map[string]string{"a": "x", "b": "y"}},
				&gnmi.PathElem{Name: "end"}},
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

func makePath(path []string) *gnmi.Path {
	ret := &gnmi.Path{Elem: []*gnmi.PathElem{}}
	for _, comp := range path {
		pathElem := &gnmi.PathElem{Name: comp}
		ret.Elem = append(ret.Elem, pathElem)
	}
	return ret
}

func TestConvertNotif(t *testing.T) {
	for name, tc := range map[string]struct {
		gnmiNotif   *gnmi.Notification
		aerisNotifs []types.Notification
	}{
		"Empty notification": {
			gnmiNotif:   &gnmi.Notification{},
			aerisNotifs: nil,
		},
		"Notification with Update": {
			gnmiNotif: makeGNMINotif(
				[]string{"prefix"},
				[][]string{[]string{"simple", "update"}},
				nil,
				[]string{"val"}),
			aerisNotifs: []types.Notification{types.NewNotification(
				ts42,
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
				ts42,
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

var ts42 = time.Unix(0, 42)

// The code that translates gNMI notifs to types.Notifications doesn't
// have a way of communication that it got a SyncResponse. So this just
// issues a SetRequest and then waits until 0.5s has passed without
// receiving another update before deciding we've synced.
func waitForSync(ctx context.Context, t *testing.T, client gnmi.GNMIClient,
	ch chan types.Notification) {
	_, _ = client.Set(ctx,
		&gnmi.SetRequest{
			Delete: []*gnmi.Path{
				Path("interfaces", "interface"),
			},
		})
	var to <-chan time.Time
	for {
		select {
		case <-ch:
			to = time.After(500 * time.Millisecond)
		case <-to:
			return
			//t.Fatal("Timed out waiting for sync")
		}
	}
}

func isPtr(notif types.Notification) bool {
	for _, v := range notif.Updates() {
		if _, ok := v.(key.Pointer); ok {
			return true
		}
	}
	return false
}

func checkNotifs(t *testing.T, ch chan types.Notification,
	expected []types.Notification, timeout time.Duration) {
	to := time.After(timeout)
	for len(expected) > 0 {
		select {
		case got := <-ch:
			if isPtr(got) {
				continue
			}
			le := len(expected)
			for i, want := range expected {
				if notiftest.Diff(want, got) == "" {
					expected = append(expected[:i], expected[i+1:]...)
					break
				}
				if i == le-1 {
					t.Fatalf("Notif didn't match any expected: %v (expected=%v)",
						got, expected)
				}
			}
		case <-to:
			t.Fatalf("Timed out waiting for expected notifs: %v", expected)
		}
	}
}

// Check that SetRequests handed to the Client produce the
// expected types.Notifications.
func TestSetRequestNotifications(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan types.Notification)
	errc := make(chan error)
	ctx, srv, err := Server(ctx, ch, errc, []string{"../../gopenconfig/yang/github.com"})
	if err != nil {
		t.Fatal(err)
	}
	client := Client(srv)

	waitForSync(ctx, t, client, ch)

	for _, tc := range []struct {
		desc       string
		setReq     *gnmi.SetRequest
		notifs     []types.Notification
		shouldFail bool
	}{
		{
			desc: "hostname",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(Path("system", "state", "hostname"), Strval("xyz")),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "system", "state"),
					nil,
					map[key.Key]interface{}{key.New("hostname"): "xyz"}),
			},
		},
		{
			desc: "interface",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(IntfConfigPath("intf123", "name"), Strval("intf123")),
					Update(IntfPath("intf123", "name"),
						Strval("intf123")),
					Update(IntfConfigPath("intf456", "name"),
						Strval("intf456")),
					Update(IntfPath("intf456", "name"),
						Strval("intf456")),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf123"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}, "config"),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf123"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf456"}),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf456"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf456"}, "config"),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf456"}),
			},
		},
		{
			desc: "interface counters",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(IntfStateCountersPath("intf123", "in-octets"),
						Uintval(uint64(1234))),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}, "state", "counters"),
					nil,
					map[key.Key]interface{}{key.New("in-octets"): uint64(1234)}),
			},
		},
		{
			desc: "lldp local interface",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(LldpIntfPath("intf123", "name"),
						Strval("intf123")),
					Update(LldpIntfConfigPath("intf123", "name"),
						Strval("intf123")),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf123"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}, "config"),
					nil,
					map[key.Key]interface{}{key.New("name"): "intf123"}),
			},
		},
		{
			desc: "lldp neighbor",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(LldpNeighborStatePath("intf123", "1",
						"id"), Strval("1")),
					Update(LldpNeighborStatePath("intf123", "1",
						"chassis-id"), Strval("whatever")),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp", "interfaces",
						"interface", map[string]interface{}{"name": "intf123"},
						"neighbors", "neighbor", map[string]interface{}{"id": "1"}),
					nil,
					map[key.Key]interface{}{key.New("id"): "1"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp", "interfaces",
						"interface", map[string]interface{}{"name": "intf123"},
						"neighbors", "neighbor", map[string]interface{}{"id": "1"}, "state"),
					nil,
					map[key.Key]interface{}{key.New("id"): "1"}),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp", "interfaces",
						"interface", map[string]interface{}{"name": "intf123"},
						"neighbors", "neighbor", map[string]interface{}{"id": "1"},
						"state"),
					nil,
					map[key.Key]interface{}{key.New("chassis-id"): "whatever"}),
			},
		},
		{
			desc: "interface delete",
			setReq: &gnmi.SetRequest{
				Delete: []*gnmi.Path{
					Path("interfaces", "interface"),
				},
				Replace: []*gnmi.Update{
					Update(IntfConfigPath("intf123", "name"),
						Strval("intf123")),
					Update(IntfPath("intf123", "name"),
						Strval("intf123")),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface"),
					path.New(map[string]interface{}{"name": "intf456"}),
					nil),
			},
		},
		{
			// check that we don't get updates in the above case
			desc: "another interface delete",
			setReq: &gnmi.SetRequest{
				Delete: []*gnmi.Path{
					Path("interfaces", "interface"),
					Path("lldp", "interfaces"),
				},
			},
			notifs: []types.Notification{
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface",
						map[string]interface{}{"name": "intf123"}, "state"),
					path.New("counters"),
					nil),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "interfaces", "interface"),
					path.New(map[string]interface{}{"name": "intf123"}),
					nil),
				types.NewNotification(
					ts42,
					path.New("OpenConfig", "lldp"),
					path.New("interfaces"),
					nil),
			},
		},
		{
			desc: "bogus path",
			setReq: &gnmi.SetRequest{
				Replace: []*gnmi.Update{
					Update(IntfStatePath("intf123", "bogus"),
						Uintval(uint64(12))),
				},
			},
			notifs:     nil,
			shouldFail: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := client.Set(ctx, tc.setReq)
			if err != nil && !tc.shouldFail {
				t.Fatal(err)
			}
			if err == nil && tc.shouldFail {
				t.Fatalf("Expected failure in test case %v", tc.desc)
			}
			if !tc.shouldFail {
				checkNotifs(t, ch, tc.notifs, time.Second*5)
			}
		})
	}
}

var inOctets uint64
var inMcastPkts = uint64(42)
var wg sync.WaitGroup

const npoll uint64 = 3

// Toy poller that just increments the in-octets interface counter
// and leaves the in-multicast-pkts counter. It will poll three times
// and then give up.
func testPoller() (*gnmi.SetRequest, error) {
	ifName := "intf123"
	inOctets++
	if inOctets > npoll {
		return nil, nil
	}
	setreq := &gnmi.SetRequest{
		Delete: []*gnmi.Path{
			Path("interfaces", "interface"),
		},
		Replace: []*gnmi.Update{
			Update(IntfPath(ifName, "name"), Strval(ifName)),
			Update(IntfConfigPath(ifName, "name"), Strval(ifName)),
			Update(IntfStatePath(ifName, "name"), Strval(ifName)),
			Update(IntfStateCountersPath(ifName, "in-octets"),
				Uintval(inOctets)),
			Update(IntfStateCountersPath(ifName, "in-multicast-pkts"),
				Uintval(inMcastPkts)),
		},
	}
	if inOctets == npoll {
		wg.Done()
	}
	return setreq, nil
}

// Test the poller API.
func TestPollForever(t *testing.T) {
	// Set up.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan types.Notification)
	errc := make(chan error)
	ctx, srv, err := Server(ctx, ch, errc, []string{"../../gopenconfig/yang/github.com"})
	if err != nil {
		t.Fatal(err)
	}
	client := Client(srv)

	// Clear queued notifs.
	waitForSync(ctx, t, client, ch)

	wg.Add(1)
	// Run poller 3 times.
	go PollForever(ctx, client, 100*time.Millisecond, testPoller, errc)

	wg.Wait() // Wait for the poller to poll 3x.

	// Check that we see the notifications we expect to see, and
	// nothing else.
	expected := []types.Notification{
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}),
			nil,
			map[key.Key]interface{}{key.New("name"): "intf123"}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "config"),
			nil,
			map[key.Key]interface{}{key.New("name"): "intf123"}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "state"),
			nil,
			map[key.Key]interface{}{key.New("name"): "intf123"}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "state", "counters"),
			nil,
			map[key.Key]interface{}{key.New("in-octets"): uint64(1)}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "state", "counters"),
			nil,
			map[key.Key]interface{}{key.New("in-multicast-pkts"): uint64(42)}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "state", "counters"),
			nil,
			map[key.Key]interface{}{key.New("in-octets"): uint64(2)}),
		types.NewNotification(
			ts42,
			path.New("OpenConfig", "interfaces", "interface",
				map[string]interface{}{"name": "intf123"}, "state", "counters"),
			nil,
			map[key.Key]interface{}{key.New("in-octets"): uint64(3)}),
	}

	checkNotifs(t, ch, expected, 2*time.Second)
}
