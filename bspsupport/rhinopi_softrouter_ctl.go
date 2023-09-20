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

package archsupport

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
var __IP_TABLE_TEMPLATE = `
iptables -A FORWARD -i eth1 -o %s -j ACCEPT
iptables -A FORWARD -i %s -o eth1 -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -t nat -A POSTROUTING -o %s -j MASQUERADE
`

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
func InitSoftRouter() {

}
func __fillIpTables(iface string) string {
	return fmt.Sprintf(__IP_TABLE_TEMPLATE, iface, iface, iface)
}
func Set4GOutLink() error {
	return __execIpTables(__fillIpTables("usb0"))
}
func SetWIFIOutLink() error {
	return __execIpTables(__fillIpTables("wlan0"))
}
func SetEth1OutLink() error {
	return __execIpTables(__fillIpTables("eth1"))
}

func __execIpTables(shell string) error {
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Exec IpTables error:%s,%s", string(output), err)
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
func FlushForwardRule() error {
	// sudo iptables -F FORWARD
	cmd := exec.Command("sh", "-c", "iptables -F FORWARD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Flush IpTables error:%s,%s", string(output), err)
	}
	return nil
}
