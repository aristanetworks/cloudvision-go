// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmi

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/openconfig/gnmi/proto/gnmi"
)

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

var inOctets uint64
var inMcastPkts = uint64(42)
var wg sync.WaitGroup

const npoll uint64 = 3

var testIfName = "intf123"

func expectedSetRequest(inOctets uint64) *gnmi.SetRequest {
	return &gnmi.SetRequest{
		Delete: []*gnmi.Path{
			Path("interfaces", "interface"),
		},
		Replace: []*gnmi.Update{
			Update(IntfPath(testIfName, "name"), Strval(testIfName)),
			Update(IntfConfigPath(testIfName, "name"), Strval(testIfName)),
			Update(IntfStatePath(testIfName, "name"), Strval(testIfName)),
			Update(IntfStateCountersPath(testIfName, "in-octets"),
				Uintval(inOctets)),
			Update(IntfStateCountersPath(testIfName, "in-multicast-pkts"),
				Uintval(inMcastPkts)),
		},
	}
}

// Toy poller that just increments the in-octets interface counter
// and leaves the in-multicast-pkts counter. It will poll three times
// and then give up.
func testPoller() (*gnmi.SetRequest, error) {
	inOctets++
	if inOctets > npoll {
		return nil, nil
	}
	if inOctets == npoll {
		wg.Done()
	}
	return expectedSetRequest(inOctets), nil
}

// Test the poller API.
func TestPollForever(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &TestClient{
		Out: make(chan *gnmi.SetRequest, 3),
	}

	errc := make(chan error)

	wg.Add(1)
	// Run poller 3 times.
	go PollForever(ctx, client, 100*time.Millisecond, testPoller, errc)

	wg.Wait() // Wait for the poller to poll 3x.

	var i uint64

	for i = 1; i <= npoll; i++ {
		got := <-client.Out
		exp := expectedSetRequest(i)
		if !reflect.DeepEqual(got, exp) {
			t.Fatalf("SetRequests not equal. \nExpected %v\nGot: %v", exp, got)
		}

	}
}
