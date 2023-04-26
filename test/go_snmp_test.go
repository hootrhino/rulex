// Copyright 2012 The GoSNMP Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gosnmp/gosnmp"
	"github.com/hootrhino/rulex/glogger"
)

// https://www.alvestrand.no/objectid/top.html
type SystemInfo struct {
	snmpClient *gosnmp.GoSNMP
}

func (si *SystemInfo) SystemDescription() string {
	s := ""
	si.snmpClient.Walk(".1.3.6.1.2.1.1.1.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			s = string(variable.Value.([]byte))
		}
		return nil
	})
	return s
}
func (si *SystemInfo) PCName() string {
	s := ""
	si.snmpClient.Walk(".1.3.6.1.2.1.1.5.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			s = string(variable.Value.([]byte))
		}
		return nil
	})
	return s
}
func (si *SystemInfo) TotalMemory() int {
	v := 0
	si.snmpClient.Walk(".1.3.6.1.2.1.25.2.2.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			v = int(variable.Value.(int))
		}
		return nil
	})
	return v

}
func (si *SystemInfo) CPUs() map[string]int {
	oid := ".1.3.6.1.2.1.25.3.3.1.2"
	r := map[string]int{}
	si.snmpClient.Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			k := strings.Replace(variable.Name, ".1.3.6.1.2.1.25.3.3.1.2.", "", 1)
			r[k] = variable.Value.(int)
		}
		return nil
	})
	return r
}
func (si *SystemInfo) ProcessList() []string {
	ss := []string{}
	si.snmpClient.Walk(".1.3.6.1.2.1.25.4.2.1.2", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			ss = append(ss, string(variable.Value.([]byte)))
		}
		return nil
	})

	return ss
}
func (si *SystemInfo) InterfaceIPs() []string {
	oid := "1.3.6.1.2.1.4.20.1.2"
	r := []string{}
	si.snmpClient.Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			ip := strings.Replace(variable.Name, ".1.3.6.1.2.1.4.20.1.2.", "", 1)
			if ip != "127.0.0.1" {
				r = append(r, ip)
			}
		}
		return nil
	})
	return r
}
func (si *SystemInfo) HardwareNetInterfaceName() []string {
	oid := ".1.3.6.1.2.1.2.2.1.2"
	ss := []string{}
	si.snmpClient.Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			ss = append(ss, string(variable.Value.([]byte)))
		}
		return nil
	})
	return ss
}
func (si *SystemInfo) HardwareNetInterfaceMac() []string {
	oid := ".1.3.6.1.2.1.2.2.1.6"
	result := []string{}
	macMaps := map[string]string{}
	si.snmpClient.Walk(oid, func(variable gosnmp.SnmpPDU) error {
		fmt.Println(variable)
		if variable.Type == gosnmp.OctetString {
			macHexs := variable.Value.([]uint8)

			if len(macHexs) > 0 {
				if macHexs[0] > 0 {
					hexs := []string{}
					for _, macHex := range macHexs {
						hexs = append(hexs, fmt.Sprintf("%X", macHex))
					}
					macMaps[strings.Join(hexs, ":")] = strings.Join(hexs, ":")
				}

			}
		}
		return nil
	})
	for k := range macMaps {
		result = append(result, k)
	}
	return result
}
func TestSnmp(t *testing.T) {
	gosnmp.Default.Target = "127.0.0.1"
	gosnmp.Default.Community = "public"
	err := gosnmp.Default.Connect()
	if err != nil {
		glogger.GLogger.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()

	si := &SystemInfo{snmpClient: gosnmp.Default}
	t.Log(si.SystemDescription())
	t.Log(si.TotalMemory())
	t.Log(si.PCName())
	// t.Log(si.CPUs())
	t.Log(si.HardwareNetInterfaceName())
	// t.Log(si.HardwareNetInterfaceMac())
	// t.Log(si.InterfaceIPs())

}
