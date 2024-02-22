// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package snmpoc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"sync"
	"testing"
	"time"

	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/pdu"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/gosnmp/gosnmp"
	"github.com/openconfig/gnmi/proto/gnmi"
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
.1.3.6.1.2.1.2.2.1.5.3001 = Gauge32: 1000000000
.1.3.6.1.2.1.2.2.1.5.3002 = Gauge32: 100000000
.1.3.6.1.2.1.2.2.1.5.999011 = Gauge32: 1000000000
.1.3.6.1.2.1.2.2.1.5.1000001 = Gauge32: 2500000000
.1.3.6.1.2.1.2.2.1.5.2001610 = Gauge32: 10000000
.1.3.6.1.2.1.2.2.1.5.5000000 = Gauge32: 0
`

// snmpwalk responses for one interface
var basicIfTableCiscoDevResponse = `
.1.3.6.1.2.1.2.2.1.1.438132736 = INTEGER: 438132736
.1.3.6.1.2.1.2.2.1.2.438132736 = STRING: Ethernet4/6/3
`

var ifTableHighSpeedResponse = `
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.2.2.1.2.999011 = STRING: Management1/1
.1.3.6.1.2.1.2.2.1.2.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.2.2.1.2.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.2.2.1.2.5000000 = STRING: Loopback0
.1.3.6.1.2.1.2.2.1.5.3001 = Gauge32: 4294967295
.1.3.6.1.2.1.2.2.1.5.3002 = Gauge32: 100000000
.1.3.6.1.2.1.2.2.1.5.999011 = Gauge32: 1000000000
.1.3.6.1.2.1.2.2.1.5.1000001 = Gauge32: 2500000000
.1.3.6.1.2.1.2.2.1.5.2001610 = Gauge32: 4294967295
.1.3.6.1.2.1.2.2.1.5.5000000 = Gauge32: 0
.1.3.6.1.2.1.31.1.1.1.15.3001 = Gauge32: 100000
.1.3.6.1.2.1.31.1.1.1.15.3002 = Gauge32: 100
.1.3.6.1.2.1.31.1.1.1.15.999011 = Gauge32: 1000
.1.3.6.1.2.1.31.1.1.1.15.1000001 = Gauge32: 2500
.1.3.6.1.2.1.31.1.1.1.15.2001610 = Gauge32: 10000
.1.3.6.1.2.1.31.1.1.1.15.5000000 = Gauge32: 0
`

var ifTableOnlyIfSpeedResponse = `
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Ethernet3/2
.1.3.6.1.2.1.2.2.1.2.999011 = STRING: Management1/1
.1.3.6.1.2.1.2.2.1.2.1000001 = STRING: Port-Channel1
.1.3.6.1.2.1.2.2.1.2.2001610 = STRING: Vlan1610
.1.3.6.1.2.1.2.2.1.2.5000000 = STRING: Loopback0
.1.3.6.1.2.1.2.2.1.5.3001 = Gauge32: 1000000000
.1.3.6.1.2.1.2.2.1.5.3002 = Gauge32: 100000000
.1.3.6.1.2.1.2.2.1.5.999011 = Gauge32: 1000000000
.1.3.6.1.2.1.2.2.1.5.1000001 = Gauge32: 2500000000
.1.3.6.1.2.1.2.2.1.5.2001610 = Gauge32: 10000000
.1.3.6.1.2.1.2.2.1.5.5000000 = Gauge32: 0
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

var basicIPAddressTableResponse = `
.1.3.6.1.2.1.4.34.1.3.1.4.172.20.193.1 = INTEGER: 3001
.1.3.6.1.2.1.4.34.1.3.1.4.172.20.193.2 = INTEGER: 3002
.1.3.6.1.2.1.4.34.1.3.1.4.172.20.253.81 = INTEGER: 1000001
.1.3.6.1.2.1.4.34.1.3.1.4.172.20.252.98 = INTEGER: 5000000
.1.3.6.1.2.1.4.34.1.3.2.16.253.122.98.159.82.164.32.193.0.0.0.0.0.0.0.1 = INTEGER: 3001
.1.3.6.1.2.1.4.34.1.3.2.16.253.122.98.159.82.164.47.32.0.0.0.0.0.0.0.17 = INTEGER: 3002
.1.3.6.1.2.1.4.34.1.3.2.16.253.122.98.159.82.164.47.128.0.0.0.0.0.0.0.152 = INTEGER: 5000000
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

var basicLldpLocalSystemDataCiscoDevResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = STRING: "34:f8:e7:a5:fa:41"
.1.0.8802.1.1.2.1.3.3.0 = STRING: "dut373"
.1.0.8802.1.1.2.1.3.7.1.2.535 = INTEGER: 7
.1.0.8802.1.1.2.1.3.7.1.3.535 = STRING: "Ethernet4/6/3"
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

// snmpwalk responses for one interface with lldpRemChassisId oid having a
// string value with NULL bytes.
var basicLldpRemTableCiscoDevResponse = `
.1.0.8802.1.1.2.1.4.1.1.4.0.535.3 = INTEGER: 6
.1.0.8802.1.1.2.1.4.1.1.5.0.535.3 = Hex-STRING: 76 6D 6E 69 63 35 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
.1.0.8802.1.1.2.1.4.1.1.6.0.535.3 = INTEGER: 3
.1.0.8802.1.1.2.1.4.1.1.7.0.535.3 = Hex-STRING: 48 DF 37 12 60 99
.1.0.8802.1.1.2.1.4.1.1.9.0.535.3 = STRING: "tst-esx-93.sjc.aristanetworks.com"
.1.0.8802.1.1.2.1.4.1.1.10.0.535.3 = STRING: "VMware ESX Releasebuild-13006603"
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
.1.3.6.1.2.1.2.2.1.1.3003 = INTEGER: 3003
.1.3.6.1.2.1.2.2.1.1.3004 = INTEGER: 3004
.1.3.6.1.2.1.2.2.1.1.3005 = INTEGER: 3005
.1.3.6.1.2.1.2.2.1.1.3006 = INTEGER: 3006
.1.3.6.1.2.1.2.2.1.1.3007 = INTEGER: 3007
.1.3.6.1.2.1.2.2.1.1.3008 = INTEGER: 3008
.1.3.6.1.2.1.2.2.1.2.3001 = STRING: Ethernet3/1
.1.3.6.1.2.1.2.2.1.2.3002 = STRING: Management1/2
.1.3.6.1.2.1.2.2.1.2.3003 = STRING: GigabitEthernet0/1
.1.3.6.1.2.1.2.2.1.2.3004 = STRING: TwoGigabitEthernet0/1
.1.3.6.1.2.1.2.2.1.2.3005 = STRING: FiveGigabitEthernet0/1
.1.3.6.1.2.1.2.2.1.2.3006 = STRING: TenGigabitEthernet0/1
.1.3.6.1.2.1.2.2.1.2.3007 = STRING: TwentyFiveGigE0/1
.1.3.6.1.2.1.2.2.1.2.3008 = STRING: FortyGigabitEthernet0/1
`

var twoIntfLldpLocalSystemDataResponse = `
.1.0.8802.1.1.2.1.3.1.0 = INTEGER: 4
.1.0.8802.1.1.2.1.3.2.0 = Hex-STRING: 00 1C 73 03 13 36
.1.0.8802.1.1.2.1.3.3.0 = STRING: device123.sjc.aristanetworks.com
.1.0.8802.1.1.2.1.3.4.0 = STRING: Arista Networks EOS version x.y.z
.1.0.8802.1.1.2.1.3.7.1.2.3001 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3002 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3003 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3004 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3005 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3006 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3007 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.2.3008 = INTEGER: 5
.1.0.8802.1.1.2.1.3.7.1.3.3001 = STRING: Eth3/1
.1.0.8802.1.1.2.1.3.7.1.3.3002 = STRING: Mgmt1/2
.1.0.8802.1.1.2.1.3.7.1.3.3003 = STRING: Gi0/1
.1.0.8802.1.1.2.1.3.7.1.3.3004 = STRING: Tw0/1
.1.0.8802.1.1.2.1.3.7.1.3.3005 = STRING: Fi0/1
.1.0.8802.1.1.2.1.3.7.1.3.3006 = STRING: Te0/1
.1.0.8802.1.1.2.1.3.7.1.3.3007 = STRING: Twe0/1
.1.0.8802.1.1.2.1.3.7.1.3.3008 = STRING: Fo0/1
`

var twoIntfLldpRemTableResponse = `
.1.0.8802.1.1.2.1.4.1.1.4.0.3001.3 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3002.4 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3003.5 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3004.6 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3005.7 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3006.8 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3007.9 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.4.0.3008.10 = INTEGER: 4
.1.0.8802.1.1.2.1.4.1.1.5.0.3001.3 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3002.4 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3003.5 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3004.6 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3005.7 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3006.8 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3007.9 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.5.0.3008.10 = Hex-STRING: 02 82 9B 3E E5 FA
.1.0.8802.1.1.2.1.4.1.1.6.0.3001.3 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3002.4 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3003.5 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3004.6 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3005.7 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3006.8 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3007.9 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.6.0.3008.10 = INTEGER: 5
.1.0.8802.1.1.2.1.4.1.1.7.0.3001.3 = STRING: p255p1
.1.0.8802.1.1.2.1.4.1.1.7.0.3002.4 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3003.5 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3004.6 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3005.7 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3006.8 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3007.9 = STRING: macvlan-bond0
.1.0.8802.1.1.2.1.4.1.1.7.0.3008.10 = STRING: macvlan-bond0
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
	tc.checkSetRequests(t, setReqs)
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
				"sysName": {
					PDU("sysName", octstr, []byte("device123.sjc.aristanetworks.com")),
				},
				"hrSystemUptime": {
					PDU("hrSystemUptime", timeticks, 162275519),
				},
				"sysUpTimeInstance": {
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("device123")),
						update(pgnmi.Path("system", "state", "domain-name"),
							strval("sjc.aristanetworks.com")),
						update(pgnmi.Path("system", "state", "boot-time"),
							intval(1553332216810000000)),
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
				{
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
						update(pgnmi.PlatformComponentStatePath("100002101", "type"),
							strval("openconfig-platform-types:MODULE")),
						update(pgnmi.PlatformComponentStatePath("100002001", "type"),
							strval("openconfig-platform-types:CONTAINER")),
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
						update(pgnmi.PlatformComponentStatePath("1", "part-no"),
							strval("DCS-7504")),
						update(pgnmi.PlatformComponentStatePath("100002101", "part-no"),
							strval("DCS-7500E-SUP")),
						update(pgnmi.PlatformComponentStatePath("1", "hardware-version"),
							strval("02.00")),
						update(pgnmi.PlatformComponentStatePath("100002101", "hardware-version"),
							strval("02.02")),
					},
				},
			},
		},
		{
			name:        "updateInterfacesBasic",
			updatePaths: []string{"^/interfaces/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":        PDUsFromString(basicIfTableResponse),
				"ifXTable":       PDUsFromString(basicIfXTableResponse),
				"ipAddressTable": PDUsFromString(basicIPAddressTableResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
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
						update(pgnmi.IntfStatePath("Loopback0", "ifindex"), uintval(5000000)),
						update(pgnmi.IntfStatePath("Ethernet3/1", "ifindex"), uintval(3001)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "ifindex"), uintval(3002)),
						update(pgnmi.IntfStatePath("Management1/1", "ifindex"), uintval(999011)),
						update(pgnmi.IntfStatePath("Port-Channel1", "ifindex"), uintval(1000001)),
						update(pgnmi.IntfStatePath("Vlan1610", "ifindex"), uintval(2001610)),
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
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "mac-address"),
							strval("74:83:ef:0f:6b:6d")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "mac-address"),
							strval("74:83:ef:0f:6b:6e")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "mac-address"),
							strval("00:1c:73:d6:22:c7")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "mac-address"),
							strval("00:1c:73:3c:a2:d3")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "mac-address"),
							strval("00:1c:73:03:13:36")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "port-speed"),
							strval("SPEED_100MB")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "port-speed"),
							strval("SPEED_2500MB")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "port-speed"),
							strval("SPEED_10MB")),
						update(pgnmi.IntfEthernetStatePath("Loopback0", "port-speed"),
							strval("SPEED_UNKNOWN")),
						update(pgnmi.IntfSubIntfIPPath("Ethernet3/1", "ip", "ipv4",
							"172.20.193.1"),
							strval("172.20.193.1")),
						update(pgnmi.IntfSubIntfIPPath("Ethernet3/2", "ip", "ipv4",
							"172.20.193.2"),
							strval("172.20.193.2")),
						update(pgnmi.IntfSubIntfIPPath("Port-Channel1", "ip", "ipv4",
							"172.20.253.81"),
							strval("172.20.253.81")),
						update(pgnmi.IntfSubIntfIPPath("Loopback0", "ip", "ipv4",
							"172.20.252.98"),
							strval("172.20.252.98")),
						update(pgnmi.IntfSubIntfIPPath("Ethernet3/1", "ip", "ipv6",
							"fd7a:629f:52a4:20c1:0000:0000:0000:0001"),
							strval("fd7a:629f:52a4:20c1:0000:0000:0000:0001")),
						update(pgnmi.IntfSubIntfIPPath("Ethernet3/2", "ip", "ipv6",
							"fd7a:629f:52a4:2f20:0000:0000:0000:0011"),
							strval("fd7a:629f:52a4:2f20:0000:0000:0000:0011")),
						update(pgnmi.IntfSubIntfIPPath("Loopback0", "ip", "ipv6",
							"fd7a:629f:52a4:2f80:0000:0000:0000:0098"),
							strval("fd7a:629f:52a4:2f80:0000:0000:0000:0098")),
					},
				},
			},
		},
		{
			name:        "update port-speed from ifHighSpeed",
			updatePaths: []string{"^/interfaces/.*/port-speed"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable": PDUsFromString(ifTableHighSpeedResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "port-speed"),
							strval("SPEED_100GB")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "port-speed"),
							strval("SPEED_100MB")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "port-speed"),
							strval("SPEED_2500MB")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "port-speed"),
							strval("SPEED_10GB")),
						update(pgnmi.IntfEthernetStatePath("Loopback0", "port-speed"),
							strval("SPEED_UNKNOWN")),
					},
				},
			},
		},
		{
			name:        "update port-speed from ifSpeed",
			updatePaths: []string{"^/interfaces/.*/port-speed"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable": PDUsFromString(ifTableOnlyIfSpeedResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "port-speed"),
							strval("SPEED_100MB")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "port-speed"),
							strval("SPEED_2500MB")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "port-speed"),
							strval("SPEED_10MB")),
						update(pgnmi.IntfEthernetStatePath("Loopback0", "port-speed"),
							strval("SPEED_UNKNOWN")),
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
				{
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
						update(pgnmi.IntfStatePath("Loopback0", "ifindex"), uintval(5000000)),
						update(pgnmi.IntfStatePath("Ethernet3/1", "ifindex"), uintval(3001)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "ifindex"), uintval(3002)),
						update(pgnmi.IntfStatePath("Management1/1", "ifindex"), uintval(999011)),
						update(pgnmi.IntfStatePath("Port-Channel1", "ifindex"), uintval(1000001)),
						update(pgnmi.IntfStatePath("Vlan1610", "ifindex"), uintval(2001610)),
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
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "mac-address"),
							strval("74:83:ef:0f:6b:6d")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "mac-address"),
							strval("74:83:ef:0f:6b:6e")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "mac-address"),
							strval("00:1c:73:d6:22:c7")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "mac-address"),
							strval("00:1c:73:3c:a2:d3")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "mac-address"),
							strval("00:1c:73:03:13:36")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "port-speed"),
							strval("SPEED_100MB")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "port-speed"),
							strval("SPEED_2500MB")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "port-speed"),
							strval("SPEED_10MB")),
						update(pgnmi.IntfEthernetStatePath("Loopback0", "port-speed"),
							strval("SPEED_UNKNOWN")),
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
			name:        "updateLldpForCiscoDevice with NULL lldpRemChassisId bytes in value",
			updatePaths: []string{"^/interfaces/", "^/lldp/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":             PDUsFromString(basicIfTableCiscoDevResponse),
				"ifXTable":            {},
				"lldpLocalSystemData": PDUsFromString(basicLldpLocalSystemDataCiscoDevResponse),
				"lldpRemTable":        PDUsFromString(basicLldpRemTableCiscoDevResponse),
				"lldpStatistics":      {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet4/6/3", "name"), strval(
							"Ethernet4/6/3")),
						update(pgnmi.IntfPath("Ethernet4/6/3", "name"), strval("Ethernet4/6/3")),
						update(pgnmi.IntfConfigPath("Ethernet4/6/3", "name"), strval(
							"Ethernet4/6/3")),
						update(pgnmi.IntfStatePath("Ethernet4/6/3", "ifindex"), uintval(438132736)),
						update(pgnmi.LldpStatePath("chassis-id-type"),
							strval(openconfig.LLDPChassisIDType(4))),
						update(pgnmi.LldpStatePath("chassis-id"), strval("34:f8:e7:a5:fa:41")),
						update(pgnmi.LldpStatePath("system-name"), strval("dut373")),
						update(pgnmi.LldpIntfConfigPath("Ethernet4/6/3", "name"),
							strval("Ethernet4/6/3")),
						update(pgnmi.LldpIntfPath("Ethernet4/6/3", "name"),
							strval("Ethernet4/6/3")),
						update(pgnmi.LldpIntfStatePath("Ethernet4/6/3", "name"),
							strval("Ethernet4/6/3")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "id"),
							strval("3")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "chassis-id-type"),
							strval("INTERFACE_NAME")),
						// After timming the NULL bytes from lldpRemChassisId value,
						// chassis-id value is set correctly.
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "chassis-id"),
							strval("vmnic5")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "port-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "port-id"),
							strval("48:df:37:12:60:99")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3", "system-name"),
							strval("tst-esx-93.sjc.aristanetworks.com")),
						update(pgnmi.LldpNeighborStatePath("Ethernet4/6/3", "3",
							"system-description"), strval("VMware ESX Releasebuild-13006603")),
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
				"lldpRemTable":   {},
				"lldpStatistics": {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
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
						update(pgnmi.IntfStatePath("Loopback0", "ifindex"), uintval(5000000)),
						update(pgnmi.IntfStatePath("Ethernet3/1", "ifindex"), uintval(3001)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "ifindex"), uintval(3002)),
						update(pgnmi.IntfStatePath("Management1/1", "ifindex"), uintval(999011)),
						update(pgnmi.IntfStatePath("Port-Channel1", "ifindex"), uintval(1000001)),
						update(pgnmi.IntfStatePath("Vlan1610", "ifindex"), uintval(2001610)),
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
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "mac-address"),
							strval("74:83:ef:0f:6b:6d")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "mac-address"),
							strval("74:83:ef:0f:6b:6e")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "mac-address"),
							strval("00:1c:73:d6:22:c7")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "mac-address"),
							strval("00:1c:73:3c:a2:d3")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "mac-address"),
							strval("00:1c:73:03:13:36")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Ethernet3/2", "port-speed"),
							strval("SPEED_100MB")),
						update(pgnmi.IntfEthernetStatePath("Management1/1", "port-speed"),
							strval("SPEED_1GB")),
						update(pgnmi.IntfEthernetStatePath("Port-Channel1", "port-speed"),
							strval("SPEED_2500MB")),
						update(pgnmi.IntfEthernetStatePath("Vlan1610", "port-speed"),
							strval("SPEED_10MB")),
						update(pgnmi.IntfEthernetStatePath("Loopback0", "port-speed"),
							strval("SPEED_UNKNOWN")),
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
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "ifindex"), uintval(3001)),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.IntfStatePath("Management1/2", "ifindex"), uintval(3002)),
						update(pgnmi.IntfPath("Management1/2", "name"), strval("Management1/2")),
						update(pgnmi.IntfConfigPath("Management1/2", "name"),
							strval("Management1/2")),
						update(pgnmi.IntfStatePath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("GigabitEthernet0/1", "ifindex"), uintval(3003)),
						update(pgnmi.IntfPath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.IntfConfigPath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("TwoGigabitEthernet0/1", "ifindex"),
							uintval(3004)),
						update(pgnmi.IntfPath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.IntfConfigPath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("FiveGigabitEthernet0/1", "ifindex"),
							uintval(3005)),
						update(pgnmi.IntfPath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.IntfConfigPath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("TenGigabitEthernet0/1", "ifindex"),
							uintval(3006)),
						update(pgnmi.IntfPath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.IntfConfigPath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.IntfStatePath("TwentyFiveGigE0/1", "ifindex"), uintval(3007)),
						update(pgnmi.IntfPath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.IntfConfigPath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.IntfStatePath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
						update(pgnmi.IntfStatePath("FortyGigabitEthernet0/1", "ifindex"),
							uintval(3008)),
						update(pgnmi.IntfPath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
						update(pgnmi.IntfConfigPath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
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
						update(pgnmi.LldpIntfConfigPath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.LldpIntfPath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.LldpIntfStatePath("GigabitEthernet0/1", "name"),
							strval("GigabitEthernet0/1")),
						update(pgnmi.LldpIntfConfigPath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.LldpIntfPath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.LldpIntfStatePath("TwoGigabitEthernet0/1", "name"),
							strval("TwoGigabitEthernet0/1")),
						update(pgnmi.LldpIntfConfigPath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.LldpIntfPath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.LldpIntfStatePath("FiveGigabitEthernet0/1", "name"),
							strval("FiveGigabitEthernet0/1")),
						update(pgnmi.LldpIntfConfigPath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.LldpIntfPath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.LldpIntfStatePath("TenGigabitEthernet0/1", "name"),
							strval("TenGigabitEthernet0/1")),
						update(pgnmi.LldpIntfConfigPath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.LldpIntfPath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.LldpIntfStatePath("TwentyFiveGigE0/1", "name"),
							strval("TwentyFiveGigE0/1")),
						update(pgnmi.LldpIntfConfigPath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
						update(pgnmi.LldpIntfPath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
						update(pgnmi.LldpIntfStatePath("FortyGigabitEthernet0/1", "name"),
							strval("FortyGigabitEthernet0/1")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "id"),
							strval("3")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "id"),
							strval("4")),
						update(pgnmi.LldpNeighborStatePath("GigabitEthernet0/1", "5", "id"),
							strval("5")),
						update(pgnmi.LldpNeighborStatePath("TwoGigabitEthernet0/1", "6", "id"),
							strval("6")),
						update(pgnmi.LldpNeighborStatePath("FiveGigabitEthernet0/1", "7", "id"),
							strval("7")),
						update(pgnmi.LldpNeighborStatePath("TenGigabitEthernet0/1", "8", "id"),
							strval("8")),
						update(pgnmi.LldpNeighborStatePath("TwentyFiveGigE0/1", "9", "id"),
							strval("9")),
						update(pgnmi.LldpNeighborStatePath("FortyGigabitEthernet0/1", "10", "id"),
							strval("10")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "chassis-id-type"),
							strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("GigabitEthernet0/1", "5",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("TwoGigabitEthernet0/1", "6",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("FiveGigabitEthernet0/1", "7",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("TenGigabitEthernet0/1", "8",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("TwentyFiveGigE0/1", "9",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("FortyGigabitEthernet0/1", "10",
							"chassis-id-type"), strval("MAC_ADDRESS")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "chassis-id"),
							strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("GigabitEthernet0/1", "5",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("TwoGigabitEthernet0/1", "6",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("FiveGigabitEthernet0/1", "7",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("TenGigabitEthernet0/1", "8",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("TwentyFiveGigE0/1", "9",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("FortyGigabitEthernet0/1", "10",
							"chassis-id"), strval("02:82:9b:3e:e5:fa")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "port-id-type"),
							strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("GigabitEthernet0/1", "5",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("TwoGigabitEthernet0/1", "6",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("FiveGigabitEthernet0/1", "7",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("TenGigabitEthernet0/1", "8",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("TwentyFiveGigE0/1", "9",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("FortyGigabitEthernet0/1", "10",
							"port-id-type"), strval("INTERFACE_NAME")),
						update(pgnmi.LldpNeighborStatePath("Ethernet3/1", "3", "port-id"),
							strval("p255p1")),
						update(pgnmi.LldpNeighborStatePath("Management1/2", "4", "port-id"),
							strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("GigabitEthernet0/1", "5",
							"port-id"), strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("TwoGigabitEthernet0/1", "6",
							"port-id"), strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("FiveGigabitEthernet0/1", "7",
							"port-id"), strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("TenGigabitEthernet0/1", "8",
							"port-id"), strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("TwentyFiveGigE0/1", "9",
							"port-id"), strval("macvlan-bond0")),
						update(pgnmi.LldpNeighborStatePath("FortyGigabitEthernet0/1", "10",
							"port-id"), strval("macvlan-bond0")),
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
				"lldpRemTable":        {},
				"lldpStatistics":      {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
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
				"sysName": {
					PDU("sysName", octstr, []byte("deviceABC")),
				},
				"sysUpTimeInstance": {
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
				"hrSystemUptime": {
					PDU("hrSystemUptime", timeticks, 162275519),
				},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("deviceABC")),
						update(pgnmi.Path("system", "state", "boot-time"),
							intval(1553332216810000000)),
					},
				},
			},
		},
		{
			name:        "updateSystemStateUpTimeOnly",
			updatePaths: []string{"^/system/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"sysName": {
					PDU("sysName", octstr, []byte("deviceABC")),
				},
				"sysUpTimeInstance": {
					PDU("sysUpTimeInstance", timeticks, 162261667),
				},
				"hrSystemUptime": {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("system")},
					Replace: []*gnmi.Update{
						update(pgnmi.Path("system", "state", "hostname"), strval("deviceABC")),
						update(pgnmi.Path("system", "state", "boot-time"),
							intval(1553332355330000000)),
					},
				},
			},
		},
		{
			name:        "lldpV2IntfSetup",
			updatePaths: []string{"^/interfaces/"},
			responses: map[string][]*gosnmp.SnmpPDU{
				"ifTable":  PDUsFromString(basicLldpV2IntfSetupResponse),
				"ifXTable": {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/1", "ifindex"), uintval(6)),
						update(pgnmi.IntfPath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfConfigPath("ethernet1/1", "name"), strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/13", "ifindex"), uintval(18)),
						update(pgnmi.IntfPath("ethernet1/13", "name"), strval("ethernet1/13")),
						update(pgnmi.IntfConfigPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfStatePath("ethernet1/14", "ifindex"), uintval(19)),
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
				"ifXTable":              {},
				"lldpV2LocalSystemData": PDUsFromString(basicLldpV2LocalSystemDataResponse),
				"lldpV2RemTable":        PDUsFromString(basicLldpV2RemTableResponse),
				"lldpV2Statistics":      PDUsFromString(basicLldpV2StatisticsResponse),
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces"), pgnmi.Path("lldp")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/1", "ifindex"),
							uintval(6)),
						update(pgnmi.IntfConfigPath("ethernet1/1", "name"),
							strval("ethernet1/1")),
						update(pgnmi.IntfStatePath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/13", "ifindex"),
							uintval(18)),
						update(pgnmi.IntfConfigPath("ethernet1/13", "name"),
							strval("ethernet1/13")),
						update(pgnmi.IntfStatePath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfPath("ethernet1/14", "name"),
							strval("ethernet1/14")),
						update(pgnmi.IntfStatePath("ethernet1/14", "ifindex"),
							uintval(19)),
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
				{
					Delete: []*gnmi.Path{pgnmi.Path("interfaces")},
					Replace: []*gnmi.Update{
						update(pgnmi.IntfStatePath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfConfigPath("Ethernet3/1", "name"), strval("Ethernet3/1")),
						update(pgnmi.IntfStatePath("Ethernet3/1", "ifindex"), uintval(3001)),
						update(pgnmi.IntfStatePath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfPath("Ethernet3/2", "name"), strval("Ethernet3/2")),
						update(pgnmi.IntfStatePath("Ethernet3/2", "ifindex"), uintval(3002)),
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
				"lldpRemTable":        {},
				"lldpStatistics":      {},
			},
			expectedSetRequests: []*gnmi.SetRequest{
				{
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
				"interfaces-lldp": {
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
				"system": {
					name: "system",
					models: map[string]*model{
						"system": supportedModels["system"],
					},
					updatePaths: map[string][]string{
						"system": {"/system/state/hostname",
							"/system/state/boot-time"},
					},
				},
			},
		},
		{
			name:  "none",
			paths: []string{},
			expectedMappingGroups: map[string]*mappingGroup{
				"interfaces-lldp": {
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
				"system": {
					name: "system",
					models: map[string]*model{
						"system": supportedModels["system"],
					},
					updatePaths: map[string][]string{
						"system": allSystemPaths,
					},
				},
				"platform": {
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

func TestMacAddrStrVal(t *testing.T) {

	for _, tc := range []valTestCase{
		{
			name: "valid mac",
			in:   "00:45:75:6f:a0:c0",
			out:  pgnmi.Strval("00:45:75:6f:a0:c0"),
		},
		{
			name: "valid mac with invalid group size",
			in:   "0:45:75:6f:a0:c",
			out:  pgnmi.Strval("00:45:75:6f:a0:0c"),
		},
		{
			name: "invalid mac 1",
			in:   "0:45",
			out:  nil,
		},
		{
			name: "invalid mac 2",
			in:   "00:45:75:6f:a0:c0:0f:dd:ff",
			out:  nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			inStr, ok := tc.in.(string)
			if !ok {
				t.Fatalf("input must be of type of string")
			}
			got := macAddrStrVal(macAddrToByteHWAddr(inStr))
			if !reflect.DeepEqual(got, tc.out) {
				t.Fatalf("got: %v, expected: %v", got, tc.out)
			}
		})
	}

	// below tests need to run without calling the macAddrToByteHWAddr() helper func
	for _, tc := range []valTestCase{
		{
			name: "data type not []byte",
			in:   "00:45:75:6f:a0:0c",
			out:  nil,
		},
		{
			name: "nil value",
			in:   nil,
			out:  nil,
		},
		{
			name: "empty byte slice",
			in:   []byte{},
			out:  nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := macAddrStrVal(tc.in)
			if !reflect.DeepEqual(got, tc.out) {
				t.Fatalf("got: %v, expected: %v", got, tc.out)
			}
		})
	}
}

func TestIfHighSpeedStrVal(t *testing.T) {
	unknownSpeed := pgnmi.Strval("SPEED_UNKNOWN")
	for _, tc := range []valTestCase{
		{
			name: "SPEED_10MB",
			in:   uint32(10),
			out:  pgnmi.Strval("SPEED_10MB"),
		},
		{
			name: "SPEED_100MB",
			in:   uint32(100),
			out:  pgnmi.Strval("SPEED_100MB"),
		},
		{
			name: "SPEED_1GB",
			in:   uint32(1000),
			out:  pgnmi.Strval("SPEED_1GB"),
		},
		{
			name: "SPEED_2500MB",
			in:   uint32(2500),
			out:  pgnmi.Strval("SPEED_2500MB"),
		},
		{
			name: "SPEED_5GB",
			in:   uint32(5000),
			out:  pgnmi.Strval("SPEED_5GB"),
		},
		{
			name: "SPEED_10GB",
			in:   uint32(1e4),
			out:  pgnmi.Strval("SPEED_10GB"),
		},
		{
			name: "SPEED_25GB",
			in:   uint32(2.5e4),
			out:  pgnmi.Strval("SPEED_25GB"),
		},
		{
			name: "SPEED_40GB",
			in:   uint32(4e4),
			out:  pgnmi.Strval("SPEED_40GB"),
		},
		{
			name: "SPEED_50GB",
			in:   uint32(5e4),
			out:  pgnmi.Strval("SPEED_50GB"),
		},
		{
			name: "SPEED_100GB",
			in:   uint32(1e5),
			out:  pgnmi.Strval("SPEED_100GB"),
		},
		{
			name: "SPEED_200GB",
			in:   uint32(2e5),
			out:  pgnmi.Strval("SPEED_200GB"),
		},
		{
			name: "SPEED_400GB",
			in:   uint32(4e5),
			out:  pgnmi.Strval("SPEED_400GB"),
		},
		{
			name: "SPEED_600GB",
			in:   uint32(6e5),
			out:  pgnmi.Strval("SPEED_600GB"),
		},
		{
			name: "SPEED_800GB",
			in:   uint32(8e5),
			out:  pgnmi.Strval("SPEED_800GB"),
		},
		{
			name: "zero bps",
			in:   uint32(0),
			out:  unknownSpeed,
		},
		{
			name: "non-zero invalid speed",
			in:   uint32(86700000),
			out:  unknownSpeed,
		},
		{
			name: "invalid value type",
			in:   "100000000",
			out:  unknownSpeed,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ifHighSpeedStrVal(tc.in)
			if !reflect.DeepEqual(got, tc.out) {
				t.Errorf("got: %v, expected: %v", got, tc.out)
			}
		})
	}
}

func TestIfSpeedStrVal(t *testing.T) {
	unknownSpeed := pgnmi.Strval("SPEED_UNKNOWN")
	for _, tc := range []valTestCase{
		{
			name: "SPEED_10MB",
			in:   uint32(1e7),
			out:  pgnmi.Strval("SPEED_10MB"),
		},
		{
			name: "SPEED_100MB",
			in:   uint32(1e8),
			out:  pgnmi.Strval("SPEED_100MB"),
		},
		{
			name: "SPEED_1GB",
			in:   uint32(1e9),
			out:  pgnmi.Strval("SPEED_1GB"),
		},
		{
			name: "SPEED_2500MB",
			in:   uint32(2.5e9),
			out:  pgnmi.Strval("SPEED_2500MB"),
		},
		{
			name: "high speed",
			in:   uint32(4294967295),
			out:  unknownSpeed,
		},
		{
			name: "zero bps",
			in:   uint32(0),
			out:  unknownSpeed,
		},
		{
			name: "invalid value type",
			in:   "100000000",
			out:  unknownSpeed,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ifSpeedStrVal(tc.in)
			if !reflect.DeepEqual(got, tc.out) {
				t.Errorf("got: %v, expected: %v", got, tc.out)
			}
		})
	}
}

func runTranslatorErrTest(t *testing.T, mibStore smi.Store,
	mockGetFunc func(oids []string) (*gosnmp.SnmpPacket, error),
	mockWalkFunc func(oid string, walker gosnmp.WalkFunc) error) error {
	trans, err := NewTranslator(mibStore, &gosnmp.GoSNMP{})
	if err != nil {
		t.Fatal(err)
	}

	// Set up mock SNMP connection.
	trans.Mock = true

	trans.Getter = mockGetFunc
	trans.Walker = mockWalkFunc

	setReqs := []*gnmi.SetRequest{}
	client := pgnmi.NewSimpleGNMIClient(func(ctx context.Context,
		req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
		setReqs = append(setReqs, req)
		return nil, nil
	})

	return trans.Poll(context.Background(), client, []string{})
}

func TestTranslatorErr(t *testing.T) {
	mibStore, err := smi.NewStore("../smi/mibs")
	if err != nil {
		t.Fatalf("Error in smi.NewStore: %s", err)
	}

	for _, errmsg := range []string{
		"request timeout",
		"closed network connection",
		"error reading from socket"} {
		err = runTranslatorErrTest(t, mibStore,
			func(oids []string) (*gosnmp.SnmpPacket, error) {
				return nil, errors.New(errmsg)
			},
			func(oid string, walker gosnmp.WalkFunc) error {
				return mockwalk(oid, walker, map[string][]*gosnmp.SnmpPDU{}, mibStore)
			})
		if err.Error() != fmt.Sprintf("SNMP Getter failed: %s", errmsg) {
			t.Fatalf("Expected err: %s, but got %s", errmsg, err.Error())
		}

		err = runTranslatorErrTest(t, mibStore,
			func(oids []string) (*gosnmp.SnmpPacket, error) {
				return mockget(oids, map[string][]*gosnmp.SnmpPDU{}, mibStore)
			},
			func(oid string, walker gosnmp.WalkFunc) error {
				return errors.New(errmsg)
			})
		if err.Error() != fmt.Sprintf("SNMP Walker failed: %s", errmsg) {
			t.Fatalf("Expected err: %s, but got %s", errmsg, err.Error())
		}
	}
}

func TestBuildIPAddrMap(t *testing.T) {

	for _, tc := range []struct {
		name           string
		expectedMapper map[string]string
		pdus           []*gosnmp.SnmpPDU
	}{
		{
			name: "Correct IP to interface name mapping with puds from both ifTable and" +
				" ipAddressTable",
			pdus: append(PDUsFromString(basicIfTableResponse),
				PDUsFromString(basicIPAddressTableResponse)...),
			expectedMapper: map[string]string{
				"172.20.193.1":  "Ethernet3/1",
				"172.20.193.2":  "Ethernet3/2",
				"172.20.253.81": "Port-Channel1",
				"172.20.252.98": "Loopback0",
				"fd7a:629f:52a4:20c1:0000:0000:0000:0001": "Ethernet3/1",
				"fd7a:629f:52a4:2f20:0000:0000:0000:0011": "Ethernet3/2",
				"fd7a:629f:52a4:2f80:0000:0000:0000:0098": "Loopback0",
			},
		},
		{
			name: "No interface name mapped for ip when there are no pdus from ifTable",
			pdus: PDUsFromString(`.1.3.6.1.2.1.4.34.1.3.1.4.172.20.193.1 = INTEGER: 3001`),
			expectedMapper: map[string]string{
				"172.20.193.1": "",
			},
		},
		{
			name:           "No mapping when there are no pdus from ipAddressTable",
			pdus:           PDUsFromString(`1.3.6.1.2.1.2.2.1.2.301 = STRING: Ethernet3/1`),
			expectedMapper: map[string]string{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mibStore, _ := smi.NewStore("../smi/mibs")
			mapper := sync.Map{}

			ps, _ := pdu.NewStore(mibStore)
			for _, p := range tc.pdus {
				err := ps.Add(p)
				if err != nil {
					t.Fatalf("error adding pdu to store: %v", err)
				}
			}

			err := buildIPAddrMap(mibStore, ps, &mapper, nil)
			if err != nil {
				t.Fatalf("expected no error, found: %v", err)
			}

			i, _ := mapper.Load("ipAddrIntfName")
			if !reflect.DeepEqual(i.(map[string]string), tc.expectedMapper) {
				t.Errorf("got: %v, \nexpected: %v", i.(map[string]string), tc.expectedMapper)
			}
		})
	}
}
