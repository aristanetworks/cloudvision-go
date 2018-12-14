// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package gnmi

import (
	"context"
	"fmt"
	"math"
	"time"

	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

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