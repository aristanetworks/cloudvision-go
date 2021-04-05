// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/protobuf/proto"
)

func TestGNMIPathJoin(t *testing.T) {
	simpleOriginTargetPrefix := PathFromString("/a/b/c")
	simpleOriginTargetPrefix.Origin = "openconfig"
	simpleOriginTargetPrefix.Target = "ABC123"
	simpleOriginTargetPath := PathFromString("/d/e/f")
	simpleOriginTargetPath.Origin = "openconfig"
	simpleOriginTargetPath.Target = "ABC123"
	simpleJoinedPath := PathFromString("a/b/c/d/e/f")
	simpleJoinedPath.Origin = "openconfig"
	simpleJoinedPath.Target = "ABC123"

	for _, tc := range []struct {
		name   string
		p1     *gnmi.Path
		p2     *gnmi.Path
		result *gnmi.Path
	}{
		{
			name:   "basic",
			p1:     PathFromString("/a/b/c"),
			p2:     PathFromString("/d/e/f"),
			result: PathFromString("/a/b/c/d/e/f"),
		},
		{
			name:   "key in prefix",
			p1:     PathFromString("/a/b[x=foo][y=bar]/c"),
			p2:     PathFromString("/d/e/f"),
			result: PathFromString("/a/b[x=foo][y=bar]/c/d/e/f"),
		},
		{
			name:   "key in path",
			p1:     PathFromString("/a/b[x=foo][y=bar]/c"),
			p2:     PathFromString("/d/e[k=baz]/f"),
			result: PathFromString("/a/b[x=foo][y=bar]/c/d/e[k=baz]/f"),
		},
		{
			name:   "origin, target in prefix",
			p1:     simpleOriginTargetPrefix,
			p2:     PathFromString("/d/e/f"),
			result: simpleJoinedPath,
		},
		{
			name:   "origin, target in path",
			p1:     PathFromString("a/b/c"),
			p2:     simpleOriginTargetPath,
			result: simpleJoinedPath,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := PathJoin(tc.p1, tc.p2)
			if !reflect.DeepEqual(tc.result, res) {
				t.Fatalf("Expected %v, got %v. p1: %s, p2: %s",
					tc.result, res, tc.p1, tc.p2)
			}
		})
	}
}

func TestGNMIPatchMatch(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		pattern string
		result  bool
	}{
		{
			name:    "simple exact",
			path:    "/a/b/c",
			pattern: "/a/b/c",
			result:  true,
		},
		{
			name:    "simple exact no match",
			path:    "/a/b/c",
			pattern: "/a/b/z",
			result:  false,
		},
		{
			name:    "simple pattern too short",
			path:    "/a/b/c",
			pattern: "/a/b",
			result:  false,
		},
		{
			name:    "simple pattern too long",
			path:    "/a/b/c",
			pattern: "/a/b/c/d",
			result:  false,
		},
		{
			name:    "wildcard elem",
			path:    "/a/b/c",
			pattern: "/a/*/c",
			result:  true,
		},
		{
			name:    "wildcard final elem",
			path:    "/a/b/c",
			pattern: "/a/b/*",
			result:  true,
		},
		{
			name:    "list exact",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[foo=X][bar=Y]/c",
			result:  true,
		},
		{
			name:    "list exact no match",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[foo=A][bar=Y]/c",
			result:  false,
		},
		{
			name:    "list implicit key wildcard",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b/c",
			result:  true,
		},
		{
			name:    "list explicit full key wildcard",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[foo=*][bar=*]/c",
			result:  true,
		},
		{
			name:    "list explicit partial key wildcard",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[foo=*][bar=Y]/c",
			result:  true,
		},
		{
			name:    "list explicit partial key implicit wildcard",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[bar=Y]/c",
			result:  true,
		},
		{
			name:    "list explicit partial key wildcard no match",
			path:    "/a/b[foo=X][bar=Y]/c",
			pattern: "/a/b[foo=*][bar=B]/c",
			result:  false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := PathMatch(PathFromString(tc.path), PathFromString(tc.pattern))
			if res != tc.result {
				t.Fatalf("Expected %v, got %v. path: %s, pattern: %s",
					tc.result, res, tc.path, tc.pattern)
			}
		})
	}
}

func TestGNMIPathCopy(t *testing.T) {
	for name, tc := range map[string]struct {
		path *gnmi.Path
	}{
		"nil path": {
			path: nil,
		},
		"nil elems": {
			path: &gnmi.Path{
				Origin: "foo",
				Elem:   nil,
				Target: "bar",
			},
		},
		"multiple elements": {
			path: &gnmi.Path{
				Origin: "foo",
				Elem: []*gnmi.PathElem{
					&gnmi.PathElem{
						Name: "elem1",
						Key:  map[string]string{"key1": "val1", "key2": "val2"},
					},
					&gnmi.PathElem{
						Name: "elem2",
					},
				},
				Target: "bar",
			},
		},
	} {
		oldPath := tc.path
		t.Run(name, func(t *testing.T) {
			newPath := PathCopy(oldPath)
			if !proto.Equal(oldPath, newPath) {
				t.Fatalf("old path %v != new path %v", oldPath.String(), newPath.String())
			}
			// Now make sure that even though values are the same, pointers are not.
			if &oldPath == &newPath {
				t.Fatal("old path address == new path address")
			}
			if oldPath != nil {
				for i, pe := range oldPath.Elem {
					if &pe == &newPath.Elem[i] {
						t.Fatal("old path elem address == new path elem address")
					}
				}
			}

		})
	}
}

func TestUnmarshal(t *testing.T) {
	// requires gnmi.ModelData to be updated to a protov2 model
	// "since the v1 lib was deprecated
	/*
			anyBytes, err := proto.Marshal(&gnmi.ModelData{Name: "foobar"})
			if err != nil {
				t.Fatal(err)
			}
		anyMessage := &any.Any{TypeUrl: "gnmi/ModelData", Value: anyBytes}
		anyString := anyMessage.String()
	*/

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
		/*
			"AnyVal": {
				val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_AnyVal{AnyVal: anyMessage}},
				exp: anyString,
			},
		*/
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

func expectedSetRequest(inOctets uint64) []*gnmi.SetRequest {
	return []*gnmi.SetRequest{
		&gnmi.SetRequest{
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
		},
	}
}

// Toy poller that just increments the in-octets interface counter
// and leaves the in-multicast-pkts counter. It will poll three times
// and then give up.
func testPoller() ([]*gnmi.SetRequest, error) {
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
	out := make(chan *gnmi.SetRequest, 3)
	setFunc := func(ctx context.Context, in *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		out <- in
		return nil, nil
	}

	client := NewSimpleGNMIClient(setFunc)

	errc := make(chan error)

	wg.Add(1)
	// Run poller 3 times.
	go PollForever(ctx, client, time.Millisecond, testPoller, errc)

	wg.Wait() // Wait for the poller to poll 3x.

	var i uint64

	for i = 1; i <= npoll; i++ {
		got := <-out
		exp := expectedSetRequest(i)
		if len(exp) != 1 {
			t.Fatalf("Too many SetRequests")
		}
		if !reflect.DeepEqual(got, exp[0]) {
			t.Fatalf("SetRequests not equal. \nExpected %v\nGot: %v", exp, got)
		}

	}
}
