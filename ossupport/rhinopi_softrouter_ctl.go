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

package ossupport

import (
	"fmt"
	"os/exec"
)

/*
*
* RhinoPi 的软路由配置，该配置主要基于Ubuntu Ip table实现，理论上来说只要有多个网卡就适用
* 但是当前该功能仅适配于Rhino系列的产品，如果需要移植请注意网卡参数。
*
 */

/*
* iptables -A FORWARD -i eth0 -o wlan0 -j ACCEPT
* iptables -A FORWARD -i wlan0 -o eth0 -m state --state ESTABLISHED,RELATED -j ACCEPT
* iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE
 */
var __FLUSH_IP_TABLE_TEMPLATE = `
sudo iptables -F
sudo iptables -X
sudo iptables -t nat -F
sudo iptables -t nat -X
sudo iptables -t mangle -F
sudo iptables -t mangle -X
sudo iptables -P INPUT DROP
sudo iptables -P FORWARD DROP
sudo iptables -P OUTPUT DROP
sudo iptables -P INPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
sudo iptables -P OUTPUT ACCEPT

`
var __IP_TABLE_TEMPLATE = `
iptables -A FORWARD -i %s -o %s -j ACCEPT
iptables -A FORWARD -i %s -o %s -j ACCEPT
iptables -t nat -A POSTROUTING -o %s -j MASQUERADE

`

/*
*
* 重构ip table, 目前默认以Eth1 <--> 4G
*
 */
func ReInitForwardRule(ifaceFrom, ifaceTo string) error {
	if err := __FlushForwardRule(); err != nil {
		return err
	}
	cmd := exec.Command("sh", "-c", __fillIpTables(ifaceFrom, ifaceTo))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Flush IpTables error:%s,%s", string(output), err)
	}
	return nil
}

/*
*
* 保存IP Table
*
 */
func SaveIpTablesConfig() error {
	cmd := exec.Command("sh", "-c", "iptables-save > /etc/rhino_iptables.conf")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Save IpTables error:%s,%s", string(output), err)
	}
	return nil
}

/*
*
* 恢复IP Table
*
 */
func RestoreIpTablesConfig() error {
	cmd := exec.Command("sh", "-c", "iptables-restore  < /etc/rhino_iptables.conf")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Restore IpTables error:%s,%s", string(output), err)
	}
	return nil
}

/*
*
rhino@RH-PI1:~$ sudo iptables -L
Chain INPUT (policy ACCEPT)
target     prot opt source               destination

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination
ACCEPT     all  --  anywhere             anywhere             state RELATED,ESTABLISHED
ACCEPT     all  --  anywhere             anywhere             state RELATED,ESTABLISHED
ACCEPT     all  --  anywhere             anywhere             state RELATED,ESTABLISHED

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
*
*/

/*
*
* 清空ip table
*
*
 */
func __FlushForwardRule() error {
	cmd := exec.Command("sh", "-c", __FLUSH_IP_TABLE_TEMPLATE)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Flush IpTables error:%s,%s", string(output), err)
	}
	return nil
}

/*
*

  - 初始化路由表，默认选择WAN作为出口
    在此之前会将LAN的IP设置为静态IP：
    auto eth1
    iface eth1 inet static
    address 192.168.64.100
    gateway 192.168.64.1
    netmask 255.255.255.0
    dns-nameservers 8.8.8.8

*
*/

func __fillIpTables(from, to string) string {
	return fmt.Sprintf(__IP_TABLE_TEMPLATE, from, to, to, from)
}

/*
*
* 开启软路由
*
 */
func StartSoftRoute() error {
	shell := `
service dnsmasq start
service isc-dhcp-server start
`
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Stop dnsmasq error:%s,%s", string(output), err)
	}
	return nil
}

/*
*
* 关闭软路由
service dnsmasq stop
service isc-dhcp-server stop
*
*/
func StopSoftRoute() error {
	shell := `
service dnsmasq stop
service isc-dhcp-server stop
`
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Start dnsmasq error:%s,%s", string(output), err)
	}
	return nil
}

/*

  - 每次初始化软路由配置表
  - 1 数据库查上一次配置的网卡参数
    2 清除当前配置
    3 应用最新的
*/

func ConfigDefaultIpTable(Iface string) error {
	return ReInitForwardRule(Iface, "eth1")
}
