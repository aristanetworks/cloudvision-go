// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmpoc

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"testing"
	"time"

	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/soniah/gosnmp"
)

var basicEntPhysicalTableResponse = `
.1.3.6.1.2.1.47.1.1.1.1.2.1 = STRING: DCS-7504 Chassis
.1.3.6.1.2.1.47.1.1.1.1.2.100002001 = STRING: Supervisor Slot 1
.1.3.6.1.2.1.47.1.1.1.1.2.100002101 = STRING: DCS-7500E-SUP Supervisor Module
.1.3.6.1.2.1.47.1.1.1.1.2.100601110 = STRING: Fan Tray 1 Fan 1
.1.3.6.1.2.1.47.1.1.1.1.5.1 = INTEGER: chassis(3)
.1.3.6.1.2.1.47.1.1.1.1.5.100002001 = INTEGER: container(5)
.1.3.6.1.2.1.47.1.1.1.1.5.100002101 = INTEGER: module(9)
.1.3.6.1.2.1.47.1.1.1.1.5.100601110 = INTEGER: fan(7)
.1.3.6.1.2.1.47.1.1.1.1.7.1 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100002101 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.7.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.8.1 = STRING: 02.00
.1.3.6.1.2.1.47.1.1.1.1.8.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.8.100002101 = STRING: 02.02
.1.3.6.1.2.1.47.1.1.1.1.8.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.1 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.10.100002101 = STRING: 4.21.0F
.1.3.6.1.2.1.47.1.1.1.1.10.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.1 = STRING: JSH11420017
.1.3.6.1.2.1.47.1.1.1.1.11.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.11.100002101 = STRING: JPE15200157
.1.3.6.1.2.1.47.1.1.1.1.11.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.12.1 = STRING: Arista Networks
.1.3.6.1.2.1.47.1.1.1.1.12.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.12.100002101 = STRING: Arista Networks
.1.3.6.1.2.1.47.1.1.1.1.12.100601110 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.13.1 = STRING: DCS-7504
.1.3.6.1.2.1.47.1.1.1.1.13.100002001 = STRING:
.1.3.6.1.2.1.47.1.1.1.1.13.100002101 = STRING: DCS-7500E-SUP
.1.3.6.1.2.1.47.1.1.1.1.13.100601110 = STRING:
`

// snmpwalk responses for six interfaces, four interface types.
var basicIfTableResponse = `
.1.3.6.1.2.1.2.2.1.1.3001 = INTEGER: 3001
.1.3.6.1.2.1.2.2.1.1.3002 = INTEGER: 3002
.1.3.6.1.2.1.2.2.1.1.999011 = INTEGER: 999011
.1.3.6.1.2.1.2.2.1.1.1000001 = INTEGER: 1000001
.1.3.6.1.2.1.2.2.1.1.2001610 = INTEGER: 2001610
.1.3.6.1.2.1.2.2.1.1.5000000 = INTEGER: 5000000
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.2.2.1.2.999011 = STRING: Management1/1
.1.3.6.1.2.1.2.2.1.2.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.2.2.1.2.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.2.2.1.2.5000000 = STRING: Loopback0
.1.3.6.1.2.1.2.2.1.3.3001 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.3002 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.999011 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.1000001 = INTEGER: 161
.1.3.6.1.2.1.2.2.1.3.2001610 = INTEGER: 136
.1.3.6.1.2.1.2.2.1.3.5000000 = INTEGER: 24
.1.3.6.1.2.1.2.2.1.4.3001 = INTEGER: 1000
.1.3.6.1.2.1.2.2.1.4.3002 = INTEGER: 1000
.1.3.6.1.2.1.2.2.1.4.999011 = INTEGER: 1000
.1.3.6.1.2.1.2.2.1.4.1000001 = INTEGER: 1000
.1.3.6.1.2.1.2.2.1.4.2001610 = INTEGER: 1000
.1.3.6.1.2.1.2.2.1.4.5000000 = INTEGER: 65536
.1.3.6.1.2.1.2.2.1.6.3001 = STRING: 74:83:ef:f:6b:6d
.1.3.6.1.2.1.2.2.1.6.3002 = STRING: 74:83:ef:f:6b:6e
.1.3.6.1.2.1.2.2.1.6.999011 = STRING: 0:1c:73:d6:22:c7
.1.3.6.1.2.1.2.2.1.6.1000001 = STRING: 0:1c:73:3c:a2:d3
.1.3.6.1.2.1.2.2.1.6.2001610 = STRING: 0:1c:73:3:13:36
.1.3.6.1.2.1.2.2.1.6.5000000 = STRING:
.1.3.6.1.2.1.2.2.1.10.3001 = Counter32: 1193032336
.1.3.6.1.2.1.2.2.1.10.3002 = Counter32: 3571250484
.1.3.6.1.2.1.2.2.1.10.999011 = Counter32: 102453638
.1.3.6.1.2.1.2.2.1.10.1000001 = Counter32: 3997687336
.1.3.6.1.2.1.2.2.1.10.2001610 = Counter32: 0
.1.3.6.1.2.1.2.2.1.10.5000000 = Counter32: 0
`
var basicIfXTableResponse = `
.1.3.6.1.2.1.31.1.1.1.1.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.31.1.1.1.1.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.31.1.1.1.1.999011 = STRING: Management1/1
.1.3.6.1.2.1.31.1.1.1.1.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.31.1.1.1.1.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.31.1.1.1.1.5000000 = STRING: Loopback0
.1.3.6.1.2.1.31.1.1.1.12.3001 = Counter64: 1303028
.1.3.6.1.2.1.31.1.1.1.12.3002 = Counter64: 5498034
.1.3.6.1.2.1.31.1.1.1.12.999011 = Counter64: 210209
.1.3.6.1.2.1.31.1.1.1.12.1000001 = Counter64: 142240878356
.1.3.6.1.2.1.31.1.1.1.12.2001610 = Counter64: 0
.1.3.6.1.2.1.31.1.1.1.12.5000000 = Counter64: 0
`

var basicLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.1 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.2 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.1 = STRING: Management1/1
.1.0.8802.1.1.2.1.3.7.1.3.451 = STRING: Ethernet3/1
.1.0.8802.1.1.2.1.3.7.1.3.452 = STRING: Ethernet3/2
`

var lldpLocalSystemDataResponseStringID = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = STRING: 50:87:89:a1:64:4f
`

var basicLldpRemTableResponse = `
.1.0.8802.1.1.2.1.4.1.1.4.0.1.1 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.4.0.451.3 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.451.4 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.5.0.1.1 = Hex-STRING: 01 AC 14 87 34
.1.0.8802.1.1.2.1.4.1.1.5.0.451.3 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.451.4 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.6.0.1.1 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.451.3 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.451.4 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.7.0.1.1 = STRING: Ethernet41
.1.0.8802.1.1.2.1.4.1.1.7.0.451.3 = STRING: p255p1
.1.0.8802.1.1.2.1.4.1.1.7.0.451.4 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.9.0.1.1 = STRING: r1-rack1-tor1.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.9.0.451.3 = STRING: server123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.9.0.451.4 = STRING: server123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.4.1.1.10.0.1.1 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.4.1.1.10.0.451.3 = STRING: Linux x.y.z
.1.0.8802.1.1.2.1.4.1.1.10.0.451.4 = STRING: Linux x.y.z
`

var basicLldpStatisticsResponse = `
.1.0.8802.1.1.2.1.2.6.1.2.1 = Counter32: 210277
.1.0.8802.1.1.2.1.2.6.1.2.2 = Counter32: 0
.1.0.8802.1.1.2.1.2.6.1.2.4 = Counter32: 0
.1.0.8802.1.1.2.1.2.6.1.2.451 = Counter32: 210214
.1.0.8802.1.1.2.1.2.6.1.2.452 = Counter32: 207597
.1.0.8802.1.1.2.1.2.7.1.2.1 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.2 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.4 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.451 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.452 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.453 = Counter32: 0
.1.0.8802.1.1.2.1.2.7.1.2.454 = Counter32: 0
`

var lldpLocalSystemDataNoChassisSubtypeResponse = `
.1.0.8802.1.1.2.1.3.2.0 = STRING: 50:87:89:a1:64:4f
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.1 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.2 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.1 = STRING: Management1/1
.1.0.8802.1.1.2.1.3.7.1.3.451 = STRING: Ethernet3/1
.1.0.8802.1.1.2.1.3.7.1.3.452 = STRING: Ethernet3/2
`

// Include another interface, Ethernet3/3, that's inactive.
var inactiveIntfLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.1 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.2 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.1 = STRING: Management1/1
.1.0.8802.1.1.2.1.3.7.1.3.451 = STRING: Ethernet3/1
.1.0.8802.1.1.2.1.3.7.1.3.452 = STRING: Ethernet3/2
.1.0.8802.1.1.2.1.3.7.1.3.453 = STRING: Ethernet3/3
`

var twoIntfIfTableResponse = `
.1.3.6.1.2.1.2.2.1.1.3001 = INTEGER: 3001
.1.3.6.1.2.1.2.2.1.1.3002 = INTEGER: 3002
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Management1/2
`

var twoIntfLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.3001 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3002 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.3001 = STRING: Eth3/1
.1.0.8802.1.1.2.1.3.7.1.3.3002 = STRING: Mgmt1/2
`

var twoIntfLldpRemTableResponse = `
.1.0.8802.1.1.2.1.4.1.1.4.0.3001.3 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3002.4 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.5.0.3001.3 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3002.4 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.6.0.3001.3 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3002.4 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.7.0.3001.3 = STRING: p255p1
.1.0.8802.1.1.2.1.4.1.1.7.0.3002.4 = STRING: macvlan-bond0
`

var basicLldpV2IntfSetupResponse = `
.1.3.6.1.2.1.2.2.1.1.6 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.1.18 = INTEGER: 18
.1.3.6.1.2.1.2.2.1.1.19 = INTEGER: 19
.1.3.6.1.2.1.2.2.1.2.6 = STRING: ethernet1/1
.1.3.6.1.2.1.2.2.1.2.18 = STRING: ethernet1/13
.1.3.6.1.2.1.2.2.1.2.19 = STRING: ethernet1/14
`

var basicLldpV2LocalSystemDataResponse = `
.1.3.111.2.802.1.1.13.1.3.1.0 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.3.2.0 = STRING: 24:0b:0a:00:70:98
.1.3.111.2.802.1.1.13.1.3.3.0 = STRING: firewall337-PAN3060
.1.3.111.2.802.1.1.13.1.3.4.0 = STRING: Palo Alto Networks 3000 series firewall
.1.3.111.2.802.1.1.13.1.3.5.0 = Hex-STRING: E8 00
.1.3.111.2.802.1.1.13.1.3.6.0 = Hex-STRING: C0 00
.1.3.111.2.802.1.1.13.1.3.7.1.2.6 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.2.18 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.2.19 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.3.7.1.3.6 = STRING: ethernet1/1
.1.3.111.2.802.1.1.13.1.3.7.1.3.18 = STRING: ethernet1/13
.1.3.111.2.802.1.1.13.1.3.7.1.3.19 = STRING: ethernet1/14
`

var basicLldpV2RemTableResponse = `
.1.3.111.2.802.1.1.13.1.4.1.1.5.0.18.1.1 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.4.1.1.5.0.19.1.2 = INTEGER: 4
.1.3.111.2.802.1.1.13.1.4.1.1.6.0.18.1.1 = STRING: 28:99:3a:bf:26:46
.1.3.111.2.802.1.1.13.1.4.1.1.6.0.19.1.2 = STRING: 28:99:3a:bf:23:f6
.1.3.111.2.802.1.1.13.1.4.1.1.7.0.18.1.1 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.4.1.1.7.0.19.1.2 = INTEGER: 5
.1.3.111.2.802.1.1.13.1.4.1.1.8.0.18.1.1 = STRING: Ethernet46
.1.3.111.2.802.1.1.13.1.4.1.1.8.0.19.1.2 = STRING: Ethernet46
.1.3.111.2.802.1.1.13.1.4.1.1.10.0.18.1.1 = STRING: switch123.sjc.aristanetworks.com
.1.3.111.2.802.1.1.13.1.4.1.1.10.0.19.1.2 = STRING: switch124.sjc.aristanetworks.com
.1.3.111.2.802.1.1.13.1.4.1.1.11.0.18.1.1 = STRING: Arista Networks EOS version x.y.z
.1.3.111.2.802.1.1.13.1.4.1.1.11.0.19.1.2 = STRING: Arista Networks EOS version x.y.z
.1.3.111.2.802.1.1.13.1.4.1.1.12.0.18.1.1 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.12.0.19.1.2 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.13.0.18.1.1 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.13.0.19.1.2 = Hex-STRING: 28 00
.1.3.111.2.802.1.1.13.1.4.1.1.14.0.18.1.1 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.14.0.19.1.2 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.15.0.18.1.1 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.1.1.15.0.19.1.2 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.2.1.3.1.9.49.53.48.46.48.46.48.46.49 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.2.1.3.1.9.49.53.48.46.48.46.48.46.50 = INTEGER: 2
.1.3.111.2.802.1.1.13.1.4.3.1.2.127 = STRING: 0x00120f0427d8
`

var basicLldpV2StatisticsResponse = `
.1.3.111.2.802.1.1.13.1.2.6.1.3.18.1 = Counter32: 118331
.1.3.111.2.802.1.1.13.1.2.6.1.3.19.1 = Counter32: 118329
.1.3.111.2.802.1.1.13.1.2.7.1.5.18.2 = Counter32: 118219
.1.3.111.2.802.1.1.13.1.2.7.1.5.19.3 = Counter32: 118194
`

var ifTable64BitResponse = `
.1.3.6.1.2.1.2.2.1.1.3001 = INTEGER: 3001
.1.3.6.1.2.1.2.2.1.1.3002 = INTEGER: 3002
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.2.2.1.3.3001 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.3.3002 = INTEGER: 6
.1.3.6.1.2.1.2.2.1.10.3001 = Counter32: 103001
.1.3.6.1.2.1.2.2.1.10.3002 = Counter32: 103002
`

var ifXTable64BitResponse = `
.1.3.6.1.2.1.31.1.1.1.1.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.31.1.1.1.1.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.31.1.1.1.4.3001 = Counter32: 43001
.1.3.6.1.2.1.31.1.1.1.4.3002 = Counter32: 43002
.1.3.6.1.2.1.31.1.1.1.6.3001 = Counter64: 1030011
.1.3.6.1.2.1.31.1.1.1.6.3002 = Counter64: 1030022
.1.3.6.1.2.1.31.1.1.1.8.3001 = Counter64: 83001
.1.3.6.1.2.1.31.1.1.1.8.3002 = Counter64: 83002
.1.3.6.1.2.1.31.1.1.1.10.3001 = Counter64: 103001
.1.3.6.1.2.1.31.1.1.1.10.3002 = Counter64: 103002
`

// mockget and mockwalk are the SNMP get and walk routines used for
// injecting mocked SNMP data into the polling routines.
func mockget(oids []string, responses map[string][]*gosnmp.SnmpPDU,
	mibStore smi.Store) (*gosnmp.SnmpPacket, error) {
	pkt := &gosnmp.SnmpPacket{}
	for _, oid := range oids {
		obj := mibStore.GetObject(oid)
		if obj != nil {
			oid = obj.Name
		}
		r, ok := responses[oid]
		if !ok || len(r) == 0 {
			pkt.Variables = append(pkt.Variables, *PDU(oid, gosnmp.NoSuchObject, nil))
			continue
		}
		if len(r) > 1 {
			return nil, fmt.Errorf("too many PDUs for OID %s", oid)
		}
		pkt.Variables = append(pkt.Variables, *(r[0]))
	}
	return pkt, nil
}

func mockwalk(oid string, walker gosnmp.WalkFunc,
	responses map[string][]*gosnmp.SnmpPDU,
	mibStore smi.Store) error {
	obj := mibStore.GetObject(oid)
	if obj != nil {
		oid = obj.Name
	}
	pdus, ok := responses[oid]
	if !ok {
		return nil
	}
	for _, p := range pdus {
		if err := walker(*p); err != nil {
			return err
		}
	}
	return nil
}

type translatorTestCase struct {
	name                string
	responses           map[string][]*gosnmp.SnmpPDU
	mappings            map[string][]Mapper
	updatePaths         []string
	expectedSetRequests []*gnmi.SetRequest
	setRequestMatchAll  bool
	expectedErr         error
}

func sortUpdates(upd []*gnmi.Update) {
	sort.Slice(upd, func(i, j int) bool {
		return upd[i].Path.String() < upd[j].Path.String()
	})
}

func sortPaths(paths []*gnmi.Path) {
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].String() < paths[j].String()
	})
}

// Check that two SetRequests have the same deletes, replaces, and
// updates, even if they're ordered differently.
func setRequestsEqual(sr1, sr2 *gnmi.SetRequest) bool {
	sortUpdates(sr1.Replace)
	sortUpdates(sr2.Replace)
	sortUpdates(sr1.Update)
	sortUpdates(sr2.Update)
	sortPaths(sr1.Delete)
	sortPaths(sr2.Delete)

	return reflect.DeepEqual(sr1, sr2)
}

func (tc *translatorTestCase) checkSetRequests(t *testing.T,
	setRequests []*gnmi.SetRequest) {
	for _, esr := range tc.expectedSetRequests {
		found := false
		for _, sr := range setRequests {
			if setRequestsEqual(esr, sr) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("No SetRequest matching %v\ngot: %v", esr, setRequests)
		}
	}

	if tc.setRequestMatchAll {
		if len(tc.expectedSetRequests) != len(setRequests) {
			t.Fatalf("Expected %d SetRequests; got %d", len(tc.expectedSetRequests),
				len(setRequests))
		}
	}
}

func runTranslatorTest(t *testing.T, mibStore smi.Store, tc translatorTestCase) {
	trans, err := NewTranslator(mibStore, &gosnmp.GoSNMP{})
	if err != nil {
		t.Fatal(err)
	}

	// Set up mock SNMP connection.
	trans.Mock = true
	if len(tc.mappings) > 0 {
		trans.Mappings = tc.mappings
	}
	trans.Getter = func(oids []string) (*gosnmp.SnmpPacket, error) {
		return mockget(oids, tc.responses, mibStore)
	}
	trans.Walker = func(oid string, walker gosnmp.WalkFunc) error {
		return mockwalk(oid, walker, tc.responses, mibStore)
	}
	now = func() time.Time {
		return time.Unix(1554954972, 0)
	}

	setReqs := []*gnmi.SetRequest{}
	client := pgnmi.NewSimpleGNMIClient(func(ctx context.Context,
		req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		setReqs = append(setReqs, req)
		return nil, nil
	})
	// Call translator.Poll and check the translator's output.
	err = trans.Poll(context.Background(), client, tc.updatePaths)
	if err != nil {
		t.Fatalf("Failure in translator.Poll: %v", err)
	}
	(&tc).checkSetRequests(t, setReqs)
}

func TestTranslator(t *testing.T) {
	mibStore, err := smi.NewStore("../smi/mibs")
	if err != nil {
		t.Fatalf("Error in smi.NewStore: %s", err)
	}

	for _, tc := range []translatorTestCase{
		{
			name:        "updateSystemStateBasic",
			updatePaths: []string{"^/system/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"sysName": []*gosnmp.SnmpPDU{
					PDU("sysName", octstr, []byte("device123.sjc.aristanetworks.com")),
				},
				"hrSystemUptime": []*gosnmp.SnmpPDU{
					PDU("hrSystemUptime", timeticks, 162275519),
				},
				"sysUpTimeInstance": []*gosnmp.SnmpPDU{
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("device123")),
						update(pgnmi.Path("system", "state", "domain-name"),
							strval("sjc.aristanetworks.com")),
						update(pgnmi.Path("system", "state", "boot-time"), intval(1553332217)),
					},
				},
			},
		},
		{
			name:        "updatePlatformBasic",
			updatePaths: []string{"^/components/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"entPhysicalEntry": PDUsFromString(basicEntPhysicalTableResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("components")},
					Replace: []*gnmi.Update{
						update(pgnmi.PlatformComponentConfigPath("1", "name"),
							strval("1")),
						update(pgnmi.PlatformComponentPath("1", "name"),
							strval("1")),
						update(pgnmi.PlatformComponentStatePath("1", "name"),
							strval("1")),
						update(pgnmi.PlatformComponentStatePath("1", "id"),
							strval("1")),
						update(pgnmi.PlatformComponentStatePath("1", "description"),
							strval("DCS-7504 Chassis")),
						update(pgnmi.PlatformComponentConfigPath("100002001", "name"),
							strval("100002001")),
						update(pgnmi.PlatformComponentPath("100002001", "name"),
							strval("100002001")),
						update(pgnmi.PlatformComponentStatePath("100002001", "name"),
							strval("100002001")),
						update(pgnmi.PlatformComponentStatePath("100002001", "id"),
							strval("100002001")),
						update(pgnmi.PlatformComponentStatePath("100002001", "description"),
							strval("Supervisor Slot 1")),
						update(pgnmi.PlatformComponentConfigPath("100002101", "name"),
							strval("100002101")),
						update(pgnmi.PlatformComponentPath("100002101", "name"),
							strval("100002101")),
						update(pgnmi.PlatformComponentStatePath("100002101", "name"),
							strval("100002101")),
						update(pgnmi.PlatformComponentStatePath("100002101", "id"),
							strval("100002101")),
						update(pgnmi.PlatformComponentStatePath("100002101", "description"),
							strval("DCS-7500E-SUP Supervisor Module")),
						update(pgnmi.PlatformComponentConfigPath("100601110", "name"),
							strval("100601110")),
						update(pgnmi.PlatformComponentPath("100601110", "name"),
							strval("100601110")),
						update(pgnmi.PlatformComponentStatePath("100601110", "name"),
							strval("100601110")),
						update(pgnmi.PlatformComponentStatePath("100601110", "id"),
							strval("100601110")),
						update(pgnmi.PlatformComponentStatePath("100601110", "description"),
							strval("Fan Tray 1 Fan 1")),
						update(pgnmi.PlatformComponentStatePath("1", "type"),
							strval("openconfig-platform-types:CHASSIS")),
						update(pgnmi.PlatformComponentStatePath("100601110", "type"),
							strval("openconfig-platform-types:FAN")),
						update(pgnmi.PlatformComponentStatePath("100002101", "software-version"),
							strval("4.21.0F")),
						update(pgnmi.PlatformComponentStatePath("1", "serial-no"),
							strval("JSH11420017")),
						update(pgnmi.PlatformComponentStatePath("100002101", "serial-no"),
							strval("JPE15200157")),
						update(pgnmi.PlatformComponentStatePath("1", "mfg-name"),
							strval("Arista Networks")),
						update(pgnmi.PlatformComponentStatePath("100002101", "mfg-name"),
							strval("Arista Networks")),
						update(pgnmi.PlatformComponentStatePath("1", "hardware-version"),
							strval("DCS-7504")),
						update(pgnmi.PlatformComponentStatePath("100002101", "hardware-version"),
							strval("DCS-7500E-SUP")),
					},
				},
			},
		},
		{
			name:        "updateInterfacesBasic",
			updatePaths: []string{"^/interfaces/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":  PDUsFromString(basicIfTableResponse),
				"ifXTable": PDUsFromString(basicIfXTableResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfStatePath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfConfigPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfStatePath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfConfigPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfStatePath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfConfigPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfStatePath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfConfigPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Management1/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Port-Channel1", "type"),
							strval("iana-if-type:ieee8023adLag")),
						update(pgnmi.IntfStatePath("Vlan1610", "type"),
							strval("iana-if-type:l3ipvlan")),
						update(pgnmi.IntfStatePath("Loopback0", "type"),
							strval("iana-if-type:softwareLoopback")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Management1/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Port-Channel1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Vlan1610", "mtu"), uintval(1000)),
						// 65536 should be converted to 65535
						update(pgnmi.IntfStatePath("Loopback0", "mtu"), uintval(65535)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-octets"),
							uintval(uint(1193032336))),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-octets"),
							uintval(uint(3571250484))),
						update(pgnmi.IntfStateCountersPath("Management1/1", "in-octets"),
							uintval(uint(102453638))),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "in-octets"),
							uintval(uint(3997687336))),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "out-multicast-pkts"),
							uintval(1303028)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "out-multicast-pkts"),
							uintval(5498034)),
						update(pgnmi.IntfStateCountersPath("Management1/1", "out-multicast-pkts"),
							uintval(210209)),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "out-multicast-pkts"),
							uintval(142240878356)),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "out-multicast-pkts"),
							uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "out-multicast-pkts"),
							uintval(0)),
					},
				},
			},
		},
		{
			name:        "updateLldpBasic",
			updatePaths: []string{"^/interfaces/", "^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":             PDUsFromString(basicIfTableResponse),
				"ifXTable":            PDUsFromString(basicIfXTableResponse),
				"lldpLocalSystemData": PDUsFromString(basicLldpLocalSystemDataResponse),
				"lldpRemTable":        PDUsFromString(basicLldpRemTableResponse),
				"lldpStatistics":      PDUsFromString(basicLldpStatisticsResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfStatePath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfConfigPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfStatePath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfConfigPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfStatePath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfConfigPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfStatePath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfConfigPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Management1/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Port-Channel1", "type"),
							strval("iana-if-type:ieee8023adLag")),
						update(pgnmi.IntfStatePath("Vlan1610", "type"),
							strval("iana-if-type:l3ipvlan")),
						update(pgnmi.IntfStatePath("Loopback0", "type"),
							strval("iana-if-type:softwareLoopback")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Management1/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Port-Channel1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Vlan1610", "mtu"), uintval(1000)),
						// 65536 should be converted to 65535
						update(pgnmi.IntfStatePath("Loopback0", "mtu"), uintval(65535)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-octets"),
							uintval(uint(1193032336))),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-octets"),
							uintval(uint(3571250484))),
						update(pgnmi.IntfStateCountersPath("Management1/1", "in-octets"),
							uintval(uint(102453638))),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "in-octets"),
							uintval(uint(3997687336))),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "out-multicast-pkts"),
							uintval(1303028)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "out-multicast-pkts"),
							uintval(5498034)),
						update(pgnmi.IntfStateCountersPath("Management1/1", "out-multicast-pkts"),
							uintval(210209)),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "out-multicast-pkts"),
							uintval(142240878356)),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "out-multicast-pkts"),
							uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "out-multicast-pkts"),
							uintval(0)),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
						update(pgnmi.LldpStatePath("chassis-id"), strval("00:1c:73:03:13:36")),
						update(pgnmi.LldpStatePath("system-name"),
							strval("device123.sjc.aristanetworks.com")),
						update(pgnmi.LldpStatePath("system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpIntfConfigPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfStatePath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfConfigPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfStatePath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfConfigPath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
						update(pgnmi.LldpIntfPath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
						update(pgnmi.LldpIntfStatePath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "id"),
							strval("1")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "chassis-id-type"),
							strval("NETWORK_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "id"), strval("3")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "id"), strval("4")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "chassis-id"),
							strval("172.20.135.52")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "port-id"),
							strval("Ethernet41")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id"),
							strval("p255p1")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "port-id"),
							strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1", "system-name"),
							strval("r1-rack1-tor1.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "system-name"),
							strval("server123.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4", "system-name"),
							strval("server123.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("Management1/1", "1",
							"system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3",
							"system-description"),
							strval("Linux x.y.z")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "4",
							"system-description"),
							strval("Linux x.y.z")),
						update(pgnmi.LldpIntfCountersPath("Management1/1", "frame-out"),
							uintval(210277)),
						update(pgnmi.LldpIntfCountersPath("Ethernet3/1", "frame-out"),
							uintval(210214)),
						update(pgnmi.LldpIntfCountersPath("Ethernet3/2", "frame-out"),
							uintval(207597)),
						update(pgnmi.LldpIntfCountersPath("Management1/1", "frame-discard"),
							uintval(0)),
						update(pgnmi.LldpIntfCountersPath("Ethernet3/1", "frame-discard"),
							uintval(0)),
						update(pgnmi.LldpIntfCountersPath("Ethernet3/2", "frame-discard"),
							uintval(0)),
					},
				},
			},
		},
		{
			name:        "updateLldpOmitInactiveIntfs",
			updatePaths: []string{"^/interfaces/", "^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":  PDUsFromString(basicIfTableResponse),
				"ifXTable": PDUsFromString(basicIfXTableResponse),
				"lldpLocalSystemData": PDUsFromString(basicLldpLocalSystemDataResponse +
					inactiveIntfLldpLocalSystemDataResponse),
				"lldpRemTable":   []*gosnmp.SnmpPDU{},
				"lldpStatistics": []*gosnmp.SnmpPDU{},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfStatePath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfConfigPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.IntfStatePath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfConfigPath("Port-Channel1", "name"),
							strval("Port-Channel1")),
						update(pgnmi.IntfStatePath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfConfigPath("Vlan1610", "name"), strval("Vlan1610")),
						update(pgnmi.IntfStatePath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfConfigPath("Loopback0", "name"), strval("Loopback0")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Management1/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Port-Channel1", "type"),
							strval("iana-if-type:ieee8023adLag")),
						update(pgnmi.IntfStatePath("Vlan1610", "type"),
							strval("iana-if-type:l3ipvlan")),
						update(pgnmi.IntfStatePath("Loopback0", "type"),
							strval("iana-if-type:softwareLoopback")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Management1/1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Port-Channel1", "mtu"), uintval(1000)),
						update(pgnmi.IntfStatePath("Vlan1610", "mtu"), uintval(1000)),
						// 65536 should be converted to 65535
						update(pgnmi.IntfStatePath("Loopback0", "mtu"), uintval(65535)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-octets"),
							uintval(uint(1193032336))),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-octets"),
							uintval(uint(3571250484))),
						update(pgnmi.IntfStateCountersPath("Management1/1", "in-octets"),
							uintval(uint(102453638))),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "in-octets"),
							uintval(uint(3997687336))),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "in-octets"), uintval(0)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "out-multicast-pkts"),
							uintval(1303028)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "out-multicast-pkts"),
							uintval(5498034)),
						update(pgnmi.IntfStateCountersPath("Management1/1", "out-multicast-pkts"),
							uintval(210209)),
						update(pgnmi.IntfStateCountersPath("Port-Channel1", "out-multicast-pkts"),
							uintval(142240878356)),
						update(pgnmi.IntfStateCountersPath("Vlan1610", "out-multicast-pkts"),
							uintval(0)),
						update(pgnmi.IntfStateCountersPath("Loopback0", "out-multicast-pkts"),
							uintval(0)),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
						update(pgnmi.LldpStatePath("chassis-id"), strval("00:1c:73:03:13:36")),
						update(pgnmi.LldpStatePath("system-name"),
							strval("device123.sjc.aristanetworks.com")),
						update(pgnmi.LldpStatePath("system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpIntfConfigPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfPath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfStatePath("Management1/1", "name"),
							strval("Management1/1")),
						update(pgnmi.LldpIntfConfigPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfStatePath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfConfigPath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
						update(pgnmi.LldpIntfPath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
						update(pgnmi.LldpIntfStatePath("Ethernet3/2", "name"),
							strval("Ethernet3/2")),
					},
				},
			},
		},
		{
			name:        "updateLldpDifferentIntfName",
			updatePaths: []string{"^/interfaces/", "^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":             PDUsFromString(twoIntfIfTableResponse),
				"lldpLocalSystemData": PDUsFromString(twoIntfLldpLocalSystemDataResponse),
				"lldpRemTable":        PDUsFromString(twoIntfLldpRemTableResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.IntfPath("Management1/2", "name"), strval("Management1/2")),
						update(pgnmi.IntfConfigPath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
						update(pgnmi.LldpStatePath("chassis-id"), strval("00:1c:73:03:13:36")),
						update(pgnmi.LldpStatePath("system-name"),
							strval("device123.sjc.aristanetworks.com")),
						update(pgnmi.LldpStatePath("system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpIntfConfigPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfPath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfStatePath("Ethernet3/1", "name"),
							strval("Ethernet3/1")),
						update(pgnmi.LldpIntfConfigPath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.LldpIntfPath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.LldpIntfStatePath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "id"), strval("3")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "id"),
							strval("4")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id"),
							strval("p255p1")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "port-id"),
							strval("macvlan-bond0")),
					},
				},
			},
		},
		{
			// chassis ID as string rather than hex-STRING
			name:        "updateLldpStringChassisID",
			updatePaths: []string{"^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":             PDUsFromString(basicIfTableResponse),
				"ifXTable":            PDUsFromString(basicIfXTableResponse),
				"lldpLocalSystemData": PDUsFromString(lldpLocalSystemDataResponseStringID),
				"lldpRemTable":        []*gosnmp.SnmpPDU{},
				"lldpStatistics":      []*gosnmp.SnmpPDU{},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.LldpStatePath("chassis-id"), strval("50:87:89:a1:64:4f")),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
					},
				},
			},
		},
		{
			// sysName is hostname only--no domain name
			name:        "updateSystemStateHostnameOnly",
			updatePaths: []string{"^/system/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"sysName": []*gosnmp.SnmpPDU{
					PDU("sysName", octstr, []byte("deviceABC")),
				},
				"sysUpTimeInstance": []*gosnmp.SnmpPDU{
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
				"hrSystemUptime": []*gosnmp.SnmpPDU{
					PDU("hrSystemUptime", timeticks, 162275519),
				},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("deviceABC")),
						update(pgnmi.Path("system", "state", "boot-time"), intval(1553332217)),
					},
				},
			},
		},
		{
			name:        "updateSystemStateUpTimeOnly",
			updatePaths: []string{"^/system/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"sysName": []*gosnmp.SnmpPDU{
					PDU("sysName", octstr, []byte("deviceABC")),
				},
				"sysUpTimeInstance": []*gosnmp.SnmpPDU{
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
				"hrSystemUptime": []*gosnmp.SnmpPDU{},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("deviceABC")),
						update(pgnmi.Path("system", "state", "boot-time"), intval(1553332356)),
					},
				},
			},
		},
		{
			name:        "lldpV2IntfSetup",
			updatePaths: []string{"^/interfaces/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":  PDUsFromString(basicLldpV2IntfSetupResponse),
				"ifXTable": []*gosnmp.SnmpPDU{},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfPath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfConfigPath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfPath("ethernet1/13", "name"), strval("ethernet1/13")),
						update(pgnmi.IntfConfigPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfPath("ethernet1/14", "name"), strval("ethernet1/14")),
						update(pgnmi.IntfConfigPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
					},
				},
			},
		},
		{
			name:        "updateLldpV2Basic",
			updatePaths: []string{"^/interfaces/", "^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":               PDUsFromString(basicLldpV2IntfSetupResponse),
				"ifXTable":              []*gosnmp.SnmpPDU{},
				"lldpV2LocalSystemData": PDUsFromString(basicLldpV2LocalSystemDataResponse),
				"lldpV2RemTable":        PDUsFromString(basicLldpV2RemTableResponse),
				"lldpV2Statistics":      PDUsFromString(basicLldpV2StatisticsResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfConfigPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfConfigPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfConfigPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
						update(pgnmi.LldpStatePath("chassis-id"), strval("24:0b:0a:00:70:98")),
						update(pgnmi.LldpStatePath("system-name"),
							strval("firewall337-PAN3060")),
						update(pgnmi.LldpStatePath("system-description"),
							strval("Palo Alto Networks 3000 series firewall")),
						update(pgnmi.LldpIntfConfigPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.LldpIntfPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.LldpIntfStatePath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.LldpIntfConfigPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.LldpIntfPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.LldpIntfStatePath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.LldpIntfConfigPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.LldpIntfPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.LldpIntfStatePath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "id"),
							strval("1")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "id"), strval("2")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "chassis-id"),
							strval("28:99:3a:bf:26:46")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "chassis-id"),
							strval("28:99:3a:bf:23:f6")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "port-id"),
							strval("Ethernet46")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "port-id"),
							strval("Ethernet46")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1", "system-name"),
							strval("switch123.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2", "system-name"),
							strval("switch124.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/13", "1",
							"system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpNeighborStatePath("ethernet1/14", "2",
							"system-description"),
							strval("Arista Networks EOS version x.y.z")),
						update(pgnmi.LldpIntfCountersPath("ethernet1/13", "frame-out"),
							uintval(118331)),
						update(pgnmi.LldpIntfCountersPath("ethernet1/14", "frame-out"),
							uintval(118329)),
						update(pgnmi.LldpIntfCountersPath("ethernet1/13", "frame-in"),
							uintval(118219)),
						update(pgnmi.LldpIntfCountersPath("ethernet1/14", "frame-in"),
							uintval(118194)),
					},
				},
			},
		},
		{
			name:        "updateInterfaces64BitFirstPoll",
			updatePaths: []string{"^/interfaces/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":  PDUsFromString(ifTable64BitResponse),
				"ifXTable": PDUsFromString(ifXTable64BitResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfConfigPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "type"),
							strval("iana-if-type:ethernetCsmacd")),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-octets"),
							uintval(1030011)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-octets"),
							uintval(1030022)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "in-multicast-pkts"),
							uintval(83001)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "in-multicast-pkts"),
							uintval(83002)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/1", "out-octets"),
							uintval(103001)),
						update(pgnmi.IntfStateCountersPath("Ethernet3/2", "out-octets"),
							uintval(103002)),
					},
				},
			},
		},
		{
			// chassis ID has no subtype
			name:        "lldpStringChassisIDNoSubtype",
			updatePaths: []string{"^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"lldpLocalSystemData": PDUsFromString(lldpLocalSystemDataNoChassisSubtypeResponse),
				"lldpRemTable":        []*gosnmp.SnmpPDU{},
				"lldpStatistics":      []*gosnmp.SnmpPDU{},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				&gnmi.SetRequest{
					Delete: []*gnmi.Path{pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.LldpStatePath("chassis-id"), strval("50:87:89:a1:64:4f")),
						update(pgnmi.LldpStatePath("system-name"),
							strval("device123.sjc.aristanetworks.com")),
						update(pgnmi.LldpStatePath("system-description"),
							strval("Arista Networks EOS version x.y.z")),
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runTranslatorTest(t, mibStore, tc)
		})
	}
}

type valTestCase struct {
	name string
	in   interface{}
	out  *gnmi.TypedValue
}

func TestStrval(t *testing.T) {
	for _, tc := range []valTestCase{
		{
			name: "normal string",
			in:   "normal",
			out:  pgnmi.Strval("normal"),
		},
		{
			name: "normal bytes",
			in:   []byte{'n', 'o', 'r', 'm', 'a', 'l'},
			out:  pgnmi.Strval("normal"),
		},
		{
			name: "bytes with null",
			in:   []byte{'a', 'b', 0x0, 'c'},
			out:  pgnmi.Strval("abc"),
		},
		{
			name: "bytes with newline",
			in:   []byte{'a', 'b', '\n', 'c'},
			out:  pgnmi.Strval("abc"),
		},
		{
			name: "square brackets",
			in:   []byte("A B [C D][]"),
			out:  pgnmi.Strval("A B (C D)()"),
		},
		{
			name: "unconverted SNMP network address",
			in:   []byte{0x01, 0xAC, 0x14, 0x87, 0x34},
			out:  pgnmi.Strval("4"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(strval(tc.in), tc.out) {
				t.Fatalf("got: %v, expected: %v", strval(tc.in), tc.out)
			}
		})
	}
}

type mappingGroupTestCase struct {
	name                  string
	paths                 []string
	expectedMappingGroups map[string]*mappingGroup
}

func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

func modelsEqual(m1, m2 *model) bool {
	if m1.name != m2.name || m1.rootPath != m2.rootPath {
		return false
	}
	if !stringSliceEqual(m1.dependencies, m1.dependencies) ||
		!stringSliceEqual(m1.snmpGetOIDs, m1.snmpGetOIDs) ||
		!stringSliceEqual(m1.snmpWalkOIDs, m1.snmpWalkOIDs) {
		return false
	}
	return true
}

func mappingGroupsEqual(mg1, mg2 *mappingGroup) bool {
	if mg1.name != mg2.name {
		return false
	}
	if len(mg1.models) != len(mg2.models) {
		return false
	}
	for k, m := range mg1.models {
		if !modelsEqual(m, mg2.models[k]) {
			return false
		}
	}
	for k, m := range mg1.updatePaths {
		if !stringSliceEqual(m, mg2.updatePaths[k]) {
			return false
		}
	}
	return true
}

func matchingPaths(pattern string, paths []string) []string {
	mp := []string{}
	for _, p := range paths {
		if match, _ := regexp.MatchString(pattern, p); match {
			mp = append(mp, p)
		}
	}
	return mp
}

func TestMappingGroups(t *testing.T) {
	defaultPaths := []string{}
	for k := range DefaultMappings() {
		defaultPaths = append(defaultPaths, k)
	}
	allIntfPaths := matchingPaths("^/interfaces/.*", defaultPaths)
	allPlatformPaths := matchingPaths("^/components/.*", defaultPaths)
	allSystemPaths := matchingPaths("^/system/.*", defaultPaths)
	allLldpPaths := matchingPaths("^/lldp/.*", defaultPaths)
	for _, tc := range []mappingGroupTestCase{
		{
			name:  "interfaces",
			paths: []string{"^/interfaces/*"},
			expectedMappingGroups: map[string]*mappingGroup{
				"interfaces-lldp": &mappingGroup{
					name: "interfaces-lldp",
					models: map[string]*model{
						"interfaces": supportedModels["interfaces"],
					},
					updatePaths: map[string][]string{
						"interfaces": allIntfPaths,
					},
				},
			},
		},
		{
			name:  "system",
			paths: []string{"/system/state/hostname", "/system/state/boot-time"},
			expectedMappingGroups: map[string]*mappingGroup{
				"system": &mappingGroup{
					name: "system",
					models: map[string]*model{
						"system": supportedModels["system"],
					},
					updatePaths: map[string][]string{
						"system": []string{"/system/state/hostname",
							"/system/state/boot-time"},
					},
				},
			},
		},
		{
			name:  "none",
			paths: []string{},
			expectedMappingGroups: map[string]*mappingGroup{
				"interfaces-lldp": &mappingGroup{
					name: "interfaces-lldp",
					models: map[string]*model{
						"interfaces": supportedModels["interfaces"],
						"lldp":       supportedModels["lldp"],
					},
					updatePaths: map[string][]string{
						"interfaces": allIntfPaths,
						"lldp":       allLldpPaths,
					},
				},
				"system": &mappingGroup{
					name: "system",
					models: map[string]*model{
						"system": supportedModels["system"],
					},
					updatePaths: map[string][]string{
						"system": allSystemPaths,
					},
				},
				"platform": &mappingGroup{
					name: "platform",
					models: map[string]*model{
						"platform": supportedModels["platform"],
					},
					updatePaths: map[string][]string{
						"platform": allPlatformPaths,
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mibStore, err := smi.NewStore("../smi/mibs")
			if err != nil {
				t.Fatal(err)
			}

			tr, err := NewTranslator(mibStore, &gosnmp.GoSNMP{})
			if err != nil {
				t.Fatal(err)
			}
			tr.Mock = true

			mg, err := tr.mappingGroupsFromPaths(tc.paths)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			if len(mg) != len(tc.expectedMappingGroups) {
				t.Fatalf("for: %+v, expected: %+v", mg, tc.expectedMappingGroups)
			}

			for k, m := range mg {
				if !mappingGroupsEqual(m, tc.expectedMappingGroups[k]) {
					t.Fatalf("got: %+v, expected: %+v", mg, tc.expectedMappingGroups)
				}
			}
		})
	}
}
