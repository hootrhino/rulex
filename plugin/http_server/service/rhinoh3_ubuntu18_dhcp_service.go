// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"os"
)

// /etc/default/isc-dhcp-server

// INTERFACESv4="eth1"
// INTERFACESv6="eth1"

func __InitDefaultDHCPListenIface(iface string) error {
	ss := `
INTERFACESv4="%s"
INTERFACESv6="%s"
`
	err := os.WriteFile("/etc/default/isc-dhcp-server", []byte(ss), 0644)
	if err != nil {
		return err
	}
	return err
}
func __InitDefaultDHCPd() error {
	// 	ss := `
	// subnet %s netmask %s {
	// 	range %s %s;
	// 	option routers%s;
	// 	option broadcast-address %s;
	// 	default-lease-time 600;
	// 	max-lease-time 7200;
	// }
	// `
	return nil
}

/*
*
# vim /etc/dhcp/dhcpd.conf

	subnet 192.168.64.0 netmask 255.255.255.0 {
	   range 192.168.64.100 192.168.64.200;      # 开放的地址池
	   option routers 192.168.64.100;            # 网关地址
	   option broadcast-address 192.168.64.255;  # 广播地址
	   default-lease-time 600;                   # 默认租期，单位：秒
	   max-lease-time 7200;                      # 最大租期
	}

* 这个初始化特殊在咬对两个软件的配置进行刷新，一个是 dnsmasq， 一个是isc-dhcp-server
*/
func InitDefaultDhcp() error {
	MIpRoute, err := GetDefaultIpRoute()
	if err != nil {
		return err
	}
	// isc-dhcp-server config
	if err0 := __InitDefaultDHCPListenIface(MIpRoute.Iface); err0 != nil {
		return err0
	}
	// dnsmasq config
	if err0 := __InitDefaultDHCPd(); err0 != nil {
		return err0
	}
	return nil
}
