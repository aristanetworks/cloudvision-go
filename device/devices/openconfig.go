// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package devices

import (
	"context"
	"net"
	"strings"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/goarista/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func init() {
	// Set options
	options := map[string]device.Option{
		"address": device.Option{
			Description: "gNMI server host/port",
			Required:    true,
		},
		"paths": device.Option{
			Description: "gNMI subscription path (comma-separated if multiple)",
			Default:     "/",
			Required:    false,
		},
		"username": device.Option{
			Description: "gNMI subscription username",
			Default:     "",
			Required:    false,
		},
		"password": device.Option{
			Description: "gNMI subscription password",
			Default:     "",
			Required:    false,
		},
		"cafile": device.Option{
			Description: "Path to server TLS certificate file",
			Default:     "",
			Required:    false,
		},
		"certfile": device.Option{
			Description: "Path to client TLS certificate file",
			Default:     "",
			Required:    false,
		},
		"keyfile": device.Option{
			Description: "Path to client TLS private key file",
			Default:     "",
			Required:    false,
		},
		"compression": device.Option{
			Description: "Compression method (Supported options: \"\" and \"gzip\")",
			Default:     "",
			Required:    false,
		},
		"tls": device.Option{
			Description: "Enable TLS",
			Default:     "false",
			Required:    false,
		},
		"device_id": device.Option{
			Description: "device ID",
			Default:     "",
			Required:    false,
		},
	}

	// Register
	device.Register("openconfig", newOpenConfig, options)
}

type openconfigDevice struct {
	gNMIProvider provider.GNMIProvider
	gNMIClient   pb.GNMIClient
	config       *gnmi.Config
	deviceID     string
	mgmtIP       string
}

func (o *openconfigDevice) Alive() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, o.config)
	livenessPath := "/system/processes/process/state"
	req, err := gnmi.NewGetRequest(gnmi.SplitPaths([]string{livenessPath}), "")
	if err != nil {
		return false, err
	}
	resp, err := o.gNMIClient.Get(ctx, req)
	return err == nil && resp != nil && len(resp.Notification) > 0, nil
}

func (o *openconfigDevice) Providers() ([]provider.Provider, error) {
	return []provider.Provider{o.gNMIProvider}, nil
}

type ocStringGetter func(chan *pb.SubscribeResponse) (string, error)

func (o *openconfigDevice) getStringFromSubscription(path string,
	f ocStringGetter) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = gnmi.NewContext(ctx, o.config)
	opt := &gnmi.SubscribeOptions{
		Mode:  "once",
		Paths: gnmi.SplitPaths([]string{path}),
	}
	respCh := make(chan *pb.SubscribeResponse, 1)
	result := ""
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		err := gnmi.SubscribeErr(ctx, o.gNMIClient, opt, respCh)
		// gNMI server sometimes returns unexpected EOF even if
		// we've seen all the data. Can't just check for io.EOF
		// because gRPC wraps it.
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			return err
		}
		return nil
	})
	errg.Go(func() error {
		r, err := f(respCh)
		cancel()
		result = r
		return err
	})
	if err := errg.Wait(); err != nil {
		return "", err
	}
	return result, nil
}

// getSerial assumes a subscription to /components/component/state, which
// looks like this:
// module: openconfig-platform
//   +--rw components
//      +--rw component* [name]
//         +--ro state
//         |  +--ro name?               string
//         |  +--ro type?               union
//         ...
//         |  +--ro serial-no?          string
//         ...
func getSerial(respCh chan *pb.SubscribeResponse) (string, error) {
	serials := map[string]string{}
	chassisName := ""
	// Iterate until the channel is closed because SubscribErr doesn't
	// handle context cancelation.
	for resp := range respCh {
		notif := resp.GetUpdate()
		if notif == nil {
			continue
		}
		for _, upd := range notif.Update {
			if len(upd.Path.Elem) == 0 {
				continue
			}
			leafName := upd.Path.Elem[len(upd.Path.Elem)-1].Name
			fullPath := upd.Path
			if notif.Prefix != nil {
				fullPath = gnmi.JoinPaths(notif.Prefix, upd.Path)
			}
			// Throw out anything that's not serial-no or type. If it's the
			// serial, but we don't yet know whether this is the CHASSIS
			// component, save it in a map[component]serial; if it's type, and
			// type is CHASSIS, this is the component we want.
			if leafName == "serial-no" {
				for _, elm := range fullPath.Elem {
					if elm.Name == "component" {
						serial := upd.Val.GetStringVal()
						if serial != "" {
							serials[elm.Key["name"]] = serial
						}
					}
				}
			} else if leafName == "type" {
				for _, elm := range fullPath.Elem {
					if elm.Name == "component" {
						typ := upd.Val.GetStringVal()
						if typ == "openconfig-platform-types:CHASSIS" {
							name := elm.Key["name"]
							chassisName = name
						}
					}
				}
			}
		}
	}
	if serial, ok := serials[chassisName]; ok {
		return serial, nil
	}
	return "", nil
}

// getChassisID assumes a subcsription to /lldp/state/chassis-id.
// module: openconfig-lldp
//   +--rw lldp
//      +--ro state
//      ...
//      |  +--ro chassis-id?                   string
//      ...
func getChassisID(respCh chan *pb.SubscribeResponse) (string, error) {
	chassisID := ""
	for resp := range respCh {
		notif := resp.GetUpdate()
		if notif == nil {
			continue
		}
		for _, upd := range notif.Update {
			chassisID = upd.Val.GetStringVal()
		}
	}
	return chassisID, nil
}

func (o *openconfigDevice) DeviceID() (string, error) {
	if o.deviceID != "" {
		return o.deviceID, nil
	}

	// Try for serial first.
	did, err := o.getStringFromSubscription("/components/component/state",
		getSerial)
	if err != nil || did != "" {
		return did, err
	}

	// Then go with chassis-id (MAC), if it's there.
	did, err = o.getStringFromSubscription("/lldp/state/chassis-id",
		getChassisID)
	if err != nil || did != "" {
		return did, err
	}

	// Fall back on the configured address.
	return o.config.Addr, nil
}

func (o *openconfigDevice) Type() string {
	return ""
}

func (o *openconfigDevice) IPAddr() string {
	if o.mgmtIP == "" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", o.config.Addr)
		if err == nil {
			o.mgmtIP = tcpAddr.IP.String()
		}
	}
	return o.mgmtIP
}

func parseGNMIOptions(opt map[string]string) (*gnmi.Config, error) {
	config := &gnmi.Config{}
	var err error
	config.Addr, err = device.GetStringOption("address", opt)
	if err != nil {
		return nil, err
	}
	config.Username, err = device.GetStringOption("username", opt)
	if err != nil {
		return nil, err
	}
	config.Password, err = device.GetStringOption("password", opt)
	if err != nil {
		return nil, err
	}
	config.CAFile, err = device.GetStringOption("cafile", opt)
	if err != nil {
		return nil, err
	}
	config.CertFile, err = device.GetStringOption("certfile", opt)
	if err != nil {
		return nil, err
	}
	config.KeyFile, err = device.GetStringOption("keyfile", opt)
	if err != nil {
		return nil, err
	}
	config.Compression, err = device.GetStringOption("compression", opt)
	if err != nil {
		return nil, err
	}
	config.TLS, err = device.GetBoolOption("tls", opt)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// newOpenConfig returns an openconfig device.
func newOpenConfig(opt map[string]string) (device.Device, error) {
	deviceID, err := device.GetStringOption("device_id", opt)
	if err != nil {
		return nil, err
	}
	gNMIPaths, err := device.GetStringOption("paths", opt)
	if err != nil {
		return nil, err
	}
	openconfig := &openconfigDevice{}
	config, err := parseGNMIOptions(opt)
	if err != nil {
		return nil, err
	}
	config.DialOptions = []grpc.DialOption{grpc.WithBlock()}
	log.Log(openconfig).Infof("Dialing gNMI target device: %s", config.Addr)
	client, err := gnmi.Dial(config)
	if err != nil {
		return nil, err
	}
	log.Log(openconfig).Infof("Connected to gNMI target device: %s", config.Addr)
	openconfig.gNMIClient = client
	openconfig.config = config
	openconfig.deviceID = deviceID

	openconfig.gNMIProvider = pgnmi.NewGNMIProvider(client, config, strings.Split(gNMIPaths, ","))

	return openconfig, nil
}
