package snmpoc

import (
	"fmt"
	"strings"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider"
	pgnmi "github.com/aristanetworks/cloudvision-go/provider/gnmi"
	"github.com/aristanetworks/cloudvision-go/provider/openconfig"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/pdu"
	"github.com/aristanetworks/cloudvision-go/provider/snmp/smi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

type ValueProcessor func(interface{}) *gnmi.TypedValue

// Generic value processors
func strval(s interface{}) *gnmi.TypedValue {
	switch t := s.(type) {
	case string:
		if t == "" {
			return nil
		}
		return pgnmi.Strval(t)
	case []byte:
		return pgnmi.Strval(string(t))
	}
	return nil
}

func uintval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToUint64(u); err == nil {
		return pgnmi.Uintval(v)
	}
	return nil
}

func intval(u interface{}) *gnmi.TypedValue {
	if v, err := provider.ToInt64(u); err == nil {
		return pgnmi.Intval(v)
	}
	return nil
}

func update(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Update {
	return pgnmi.Update(path, val)
}

func scalarMapper(ps pdu.Store, path string, oid string,
	vp ValueProcessor) (*gnmi.Update, error) {
	pdu, err := ps.GetScalar(oid)
	if err != nil {
		return nil, err
	} else if pdu == nil {
		return nil, nil
	}

	val := vp(pdu.Value)
	if val == nil {
		return nil, nil
	}

	return update(pgnmi.PathFromString(path), val), nil
}

func scalarMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store, kvs KVStore) ([]*gnmi.Update, error) {
		u, err := scalarMapper(ps, path, oid, vp)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, nil
		}
		return []*gnmi.Update{u}, nil
	}
}

// Some implementations will return a hostname only, while others
// will return a fully qualified domain name.
func processSysName(s interface{}) *gnmi.TypedValue {
	sysName := s.(string)
	return strval(strings.SplitN(sysName, ".", 2)[0])
}

// Get boot-time by subtracting the target device's uptime from
// the Collector's current time. This isn't really correct--the
// boot time we produce is a blend of UNIX time on the Collector
// and UNIX time on the target device (which may not be in sync),
// and the target device's time may be recorded well before the
// Collector's. There doesn't seem to be a great way to get the
// device's time using SNMP. Assuming the two devices are
// roughly in sync, though, this shouldn't be disastrous.
func processBootTime(x interface{}) *gnmi.TypedValue {
	t, err := provider.ToInt64(x)
	if err != nil {
		return nil
	}
	if t == 0 {
		return nil
	}
	return intval(time.Now().Unix() - t/100)
}

// interface helpers
func setIntfName(kvs KVStore, ifIndex, ifDescr string) {
	m := kvs.Get("intfName")
	if m == nil {
		kvs.Set("intfName", make(map[string]string))
		m = kvs.Get("intfName")
	}

	mp := m.(map[string]string)
	mp[ifIndex] = ifDescr
}

func getIntfName(kvs KVStore, ifIndex string) string {
	m := kvs.Get("intfName")
	if m == nil {
		return ""
	}
	mp := m.(map[string]string)
	return mp[ifIndex]
}

// generic mapper for PDUs from ifTable
func ifTableMapper(ss smi.Store, ps pdu.Store, kvs KVStore, path string,
	oid string, vp ValueProcessor) ([]*gnmi.Update, error) {
	pdus, err := ps.GetTabular(oid)
	if err != nil {
		return nil, err
	} else if len(pdus) == 0 {
		return nil, nil
	}

	// ifDescr is our index, so handle it separately
	o := ss.GetObject(pdus[0].Name)
	if o == nil {
		return nil, fmt.Errorf("GetObject failed on OID %s", oid)
	}

	updates := []*gnmi.Update{}
	for _, p := range pdus {
		ifIndex := pdu.IndexValues(ss, p)[0]
		ifDescr := getIntfName(kvs, ifIndex)
		if ifDescr == "" && o.Name == "ifDescr" {
			ifDescr = p.Value.(string)
			setIntfName(kvs, ifIndex, ifDescr)
		} else if ifDescr == "" {
			return nil, fmt.Errorf("No ifDescr for ifIndex %s", ifIndex)
		}
		fullPath := pgnmi.PathFromString(fmt.Sprintf(path, ifDescr))
		updates = append(updates, update(fullPath, vp(p.Value)))
	}
	return updates, nil
}

func ifTableMapperFn(path, oid string, vp ValueProcessor) Mapper {
	return func(ss smi.Store, ps pdu.Store, kvs KVStore) ([]*gnmi.Update, error) {
		return ifTableMapper(ss, ps, kvs, path, oid, vp)
	}
}

// mappers
var (
	// /interfaces
	interfacePath        = "/interfaces/interface[name=%s]/"
	interfaceStatePath   = interfacePath + "state/"
	interfaceConfigPath  = interfacePath + "config/"
	interfaceCounterPath = interfaceStatePath + "counters/"
	interfaceName        = ifTableMapperFn(interfacePath+"name",
		"ifDescr", strval)
	interfaceStateName = ifTableMapperFn(interfaceStatePath+"name",
		"ifDescr", strval)
	interfaceConfigName = ifTableMapperFn(interfaceConfigPath+"name",
		"ifDescr", strval)
	interfaceMtu = ifTableMapperFn(interfaceStatePath+"mtu",
		"ifMtu", uintval)
	interfaceAdminStatus = ifTableMapperFn(interfaceStatePath+"admin-status",
		"ifAdminStatus", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.IntfAdminStatus(x.(int)))
		})
	interfaceOperStatus = ifTableMapperFn(interfaceStatePath+"oper-status",
		"ifOperStatus", func(x interface{}) *gnmi.TypedValue {
			return strval(openconfig.IntfOperStatus(x.(int)))
		})
	interfaceInOctets64 = ifTableMapperFn(interfaceCounterPath+"in-octets",
		"ifHCInOctets", uintval)
	interfaceInOctets32 = ifTableMapperFn(interfaceCounterPath+"in-octets",
		"ifInOctets", uintval)
	interfaceInUnicastPkts64 = ifTableMapperFn(interfaceCounterPath+"in-unicast-pkts",
		"ifHCInUcastPkts", uintval)
	interfaceInUnicastPkts32 = ifTableMapperFn(interfaceCounterPath+"in-unicast-pkts",
		"ifInUcastPkts", uintval)
	interfaceInMulticastPkts = ifTableMapperFn(interfaceCounterPath+"in-multicast-pkts",
		"ifHCInMulticastPkts", uintval)
	interfaceInBroadcastPkts = ifTableMapperFn(interfaceCounterPath+"in-broadcast-pkts",
		"ifHCInBroadcastPkts", uintval)
	interfaceOutMulticastPkts = ifTableMapperFn(interfaceCounterPath+"out-multicaast-pkts",
		"ifHCOutMulticastPkts", uintval)
	interfaceOutBroadcastPkts = ifTableMapperFn(interfaceCounterPath+"out-broadcast-pkts",
		"ifHCOutBroadcastPkts", uintval)
	interfaceInDiscards = ifTableMapperFn(interfaceCounterPath+"in-discards",
		"ifInDiscards", uintval)
	interfaceInErrors = ifTableMapperFn(interfaceCounterPath+"in-errors",
		"ifInErrors", uintval)
	interfaceInUnknownProtos = ifTableMapperFn(interfaceCounterPath+"in-unknown-protos",
		"ifInUnknownProtos", uintval)
	interfaceOutOctets64 = ifTableMapperFn(interfaceCounterPath+"out-octets",
		"ifHCOutOctets", uintval)
	interfaceOutOctets32 = ifTableMapperFn(interfaceCounterPath+"out-octets",
		"ifOutOctets", uintval)
	interfaceOutUnicastPkts64 = ifTableMapperFn(interfaceCounterPath+"out-unicast-pkts",
		"ifHCOutUcastPkts", uintval)
	interfaceOutUnicastPkts32 = ifTableMapperFn(interfaceCounterPath+"out-unicast-pkts",
		"ifOutUcastPkts", uintval)
	interfaceOutDiscards = ifTableMapperFn(interfaceCounterPath+"out-discards",
		"ifOutDiscards", uintval)
	interfaceOutErrors = ifTableMapperFn(interfaceCounterPath+"out-errors",
		"ifOutErrors", uintval)

	// /system/state
	systemStatePath     = "/system/state/"
	systemStateHostname = scalarMapperFn(systemStatePath+"hostname",
		"sysName", processSysName)
	systemStateHostnameLldp = scalarMapperFn(systemStatePath+"hostname",
		"lldpLocSysName", processSysName)
	systemStateBootTime64 = scalarMapperFn(systemStatePath+"boot-time",
		"hrSystemUptime", processBootTime)
	systemStateBootTime32 = scalarMapperFn(systemStatePath+"boot-time",
		"sysUpTimeInstance", processBootTime)
)

// DefaultMappings defines an ordered set of mappings per supported path.
var DefaultMappings = map[string][]Mapper{
	// interface
	"/interfaces/interface[name=name]/name":               []Mapper{interfaceName},
	"/interfaces/interface[name=name]/state/name":         []Mapper{interfaceStateName},
	"/interfaces/interface[name=name]/config/name":        []Mapper{interfaceConfigName},
	"/interfaces/interface[name=name]/state/mtu":          []Mapper{interfaceMtu},
	"/interfaces/interface[name=name]/state/admin-status": []Mapper{interfaceAdminStatus},
	"/interfaces/interface[name=name]/state/oper-status":  []Mapper{interfaceOperStatus},
	"/interfaces/interface[name=name]/state/in-octets": []Mapper{interfaceInOctets64,
		interfaceInOctets32},
	"/interfaces/interface[name=name]/state/in-unicast-pkts": []Mapper{interfaceInUnicastPkts64,
		interfaceInUnicastPkts32},
	"/interfaces/interface[name=name]/state/in-multicast-pkts": []Mapper{interfaceInMulticastPkts},
	"/interfaces/interface[name=name]/state/in-broadcast-pkts": []Mapper{interfaceInBroadcastPkts},
	"/interfaces/interface[name=name]/state/out-multicast-pkts": []Mapper{
		interfaceOutMulticastPkts},
	"/interfaces/interface[name=name]/state/out-broadcast-pkts": []Mapper{
		interfaceOutBroadcastPkts},
	"/interfaces/interface[name=name]/state/in-discards":       []Mapper{interfaceInDiscards},
	"/interfaces/interface[name=name]/state/in-errors":         []Mapper{interfaceInErrors},
	"/interfaces/interface[name=name]/state/in-unknown-protos": []Mapper{interfaceInUnknownProtos},
	"/interfaces/interface[name=name]/state/out-octets": []Mapper{interfaceOutOctets64,
		interfaceOutOctets32},
	"/interfaces/interface[name=name]/state/out-unicast-pkts": []Mapper{interfaceOutUnicastPkts64,
		interfaceOutUnicastPkts32},
	"/interfaces/interface[name=name]/state/out-discards": []Mapper{interfaceOutDiscards},
	"/interfaces/interface[name=name]/state/out-errors":   []Mapper{interfaceOutErrors},

	// system
	"/system/state/hostname":  []Mapper{systemStateHostname, systemStateHostnameLldp},
	"/system/state/boot-time": []Mapper{systemStateBootTime64, systemStateBootTime32},
}
