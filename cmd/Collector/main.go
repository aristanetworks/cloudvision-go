// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package main

import (
	"github.com/aristanetworks/cloudvision-go/cmd/Collector/libmain"
	"github.com/aristanetworks/cloudvision-go/device"
)

func main() {
	sc := device.SensorConfig{Connector: device.NewDefaultGRPCConnector()}
	libmain.Main(sc)
}
