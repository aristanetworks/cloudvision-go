// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"context"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/glog"
	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
)

type dumpInfo struct {
	writePath string
	lock      sync.Mutex
	startTime time.Time
	timeout   time.Duration
	doneGroup sync.WaitGroup
	done      bool
	file      *os.File
}

func (d *dumpInfo) processRequest(ctx context.Context,
	req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.done || reflect.DeepEqual(req, &gnmi.SetRequest{}) {
		return nil, nil
	}
	if d.file == nil {
		var err error
		d.file, err = os.Create(d.writePath)
		if err != nil {
			return nil, err
		}
	}
	err := proto.CompactText(d.file, req)
	if err != nil {
		return nil, err
	}
	_, err = d.file.WriteString("\n")
	if err != nil {
		return nil, err
	}
	if time.Since(d.startTime) < d.timeout {
		return nil, nil
	}
	d.done = true
	d.doneGroup.Done()
	d.file.Close()
	return nil, nil
}

func newDumpInfo() *dumpInfo {
	return &dumpInfo{
		writePath: *dumpFile,
		startTime: time.Now(),
		timeout:   *dumpTimeout,
	}
}

func runDump(ctx context.Context) {
	dumpInfo := newDumpInfo()
	dumpInfo.doneGroup.Add(1)
	inventory := device.NewInventory(ctx, pgnmi.NewSimpleGNMIClient(dumpInfo.processRequest))
	devices, err := device.CreateDevices(*deviceName, *deviceConfigFile, deviceOptions)
	if err != nil {
		glog.Fatal(err)
	}
	for _, info := range devices {
		err := inventory.Add(info.ID, info.Device)
		if err != nil {
			glog.Fatalf("Error in inventory.Add(): %v", err)
		}
	}
	glog.V(2).Info("Dump Collector is running")
	dumpInfo.doneGroup.Wait()
}
