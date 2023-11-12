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
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
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
	return nil
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
func ConfigDefaultIscServer(Iface string) error {
	// isc-dhcp-server config
	if err0 := __InitDefaultDHCPListenIface(Iface); err0 != nil {
		return err0
	}
	// dnsmasq config
	if err0 := __InitDefaultDHCPd(); err0 != nil {
		return err0
	}
	return nil
}

/*
*
rhino@RH-PI1:~$ cat /var/lib/dhcp/dhcpd.leases
# The format of this file is documented in the dhcpd.leases(5) manual page.
# This lease file was written by isc-dhcp-4.3.5

# authoring-byte-order entry is generated, DO NOT DELETE

	    authoring-byte-order little-endian;
		lease 192.168.64.101 {
		  hardware ethernet a8:a1:59:2e:a2:d9;
		}
		lease 192.168.64.102 {
		  hardware ethernet a8:a1:59:2e:a2:d9;
		}

*
*/
func GetDhcpList() ([]DhcpLease, error) {
	return __ParseDhcpConfig("/var/lib/dhcp/dhcpd.leases")
}

type DhcpLease struct {
	MacAddress string `json:"mac_address"`
	IpAddress  string `json:"ip_address"`
	Hostname   string `json:"hostname"`
}

func (v DhcpLease) JsonString() string {
	b, _ := json.Marshal(v)
	return string(b)
}

/*
*
* 解析 /var/lib/dhcp/dhcpd.leases
*
 */
func __ParseDhcpConfig(filePath string) ([]DhcpLease, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var leases []DhcpLease
	var currentLease DhcpLease

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "lease") {
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				currentLease.IpAddress = strings.TrimRight(parts[1], "{")
			}
		} else if strings.HasPrefix(line, "hardware ethernet") {
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				currentLease.MacAddress = strings.TrimRight(parts[1], ";")
			}
		} else if strings.HasPrefix(line, "client-hostname") {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				currentLease.Hostname = strings.TrimRight(parts[1], ";")
			}
			leases = append(leases, currentLease)
			currentLease = DhcpLease{"", "", ""}
		}
	}
	leasesMap := map[string]DhcpLease{}
	result := []DhcpLease{}
	for _, lease := range leases {
		leasesMap[lease.Hostname] = lease
		result = append(result, leasesMap[lease.Hostname])
	}
	return result, nil
}

/*
	    IP和MAC地址绑定
		host hostname {
		    hardware ethernet 00:11:22:33:44:55;
		    fixed-address 192.168.1.200;
			client-hostname "wwhai"
		}
*/
func BindMacAndIP() {

}

/*
*

	subnet 192.168.64.0 netmask 255.255.255.0 {
	   range 192.168.64.100 192.168.64.200;          # 开放的地址池
	   option domain-name-servers 192.168.64.100;    # DNS域名服务器，如果没有就注释掉
	   #option domain-name "internal.example.org";   # 域名
	   option routers 192.168.64.100;                # 网关地址
	   option broadcast-address 192.168.64.255;      # 广播地址
	   default-lease-time 600;                       # 默认租期，单位：秒
	   max-lease-time 7200;                          # 最大租期
	}

*isc-dhcp-server 会加载 /etc/dhcp/dhcpd.conf 配置文件
这里需要注意一下，此处配置要和网卡的配置一致
*/
type IscServerDHCPConfig struct {
	Iface       string // 用来做子网的那个网卡的网卡名
	Ip          string // 用来做子网的那个网卡的IP地址
	Gateway     string // 用来做子网的那个网卡的网关
	Network     string // 用来做子网的那个网卡的网段
	Netmask     string // 用来做子网的那个网卡子网掩码
	IpPoolBegin string // DHCP IP地址池起始
	IpPoolEnd   string // DHCP IP地址池结束
	//------------------------------------
	// IP 路由方向, 默认 ETH1 透传到 4G
	//------------------------------------
	IfaceFrom string // 流量入口
	IfaceTo   string // 流量出口
}

func ConfigDefaultIscServeDhcp(IpRoute IscServerDHCPConfig) error {
	dhcpdConf := `
subnet %s netmask %s {
    range %s %s;
    option routers %s;
    default-lease-time 600;
    max-lease-time 7200;
}

`
	if IsIPv4InDHCPRange(IpRoute.Network, IpRoute.IpPoolBegin, IpRoute.IpPoolEnd) {
		return fmt.Errorf("not Valid Ip Range:%s-%s", IpRoute.IpPoolBegin, IpRoute.IpPoolEnd)
	}
	if IsIPRangeValid(IpRoute.IpPoolBegin, IpRoute.IpPoolEnd) {
		return fmt.Errorf("not Valid Ip Range:%s-%s", IpRoute.IpPoolBegin, IpRoute.IpPoolEnd)
	}
	shell := fmt.Sprintf(dhcpdConf, IpRoute.Network, IpRoute.Netmask,
		IpRoute.IpPoolBegin, IpRoute.IpPoolEnd, IpRoute.Gateway)

	if err1 := os.WriteFile("/etc/dhcp/dhcpd.conf",
		[]byte(shell), 0644); err1 != nil {
		return err1
	}
	return nil
}

// IsIPv4InDHCPRange 检查给定的 IPv4 地址是否在 DHCP 地址池范围内
func IsIPv4InDHCPRange(ip, start, end string) bool {
	// 将字符串 IP 地址解析为 net.IP 类型
	ipAddr := net.ParseIP(ip)
	startAddr := net.ParseIP(start)
	endAddr := net.ParseIP(end)

	// 检查 IP 地址是否为空或不是 IPv4 地址
	if ipAddr == nil ||
		startAddr == nil ||
		endAddr == nil ||
		ipAddr.To4() == nil ||
		startAddr.To4() == nil ||
		endAddr.To4() == nil {
		return false
	}

	// 将 IP 地址转换为 32 位整数
	ipUint32 := ipToUint32(ipAddr.To4())
	startUint32 := ipToUint32(startAddr.To4())
	endUint32 := ipToUint32(endAddr.To4())

	// 检查 IP 地址是否在地址池范围内
	return ipUint32 >= startUint32 && ipUint32 <= endUint32
}

// 将 IPv4 地址转换为 32 位整数
func ipToUint32(ip net.IP) uint32 {
	return uint32(
		ip[0])<<24 |
		uint32(ip[1])<<16 |
		uint32(ip[2])<<8 |
		uint32(ip[3])
}
func IsIPRangeValid(startIPStr, endIPStr string) bool {
	// 解析起始和结束 IP 地址
	startIP := net.ParseIP(startIPStr)
	endIP := net.ParseIP(endIPStr)

	// 检查 IP 地址是否有效
	if startIP == nil || endIP == nil {
		return false
	}

	// 检查 IP 地址的版本是否一致（IPv4 或 IPv6）
	if startIP.To4() == nil || endIP.To4() == nil {
		return false // IP 地址版本不一致
	}

	// 将 IPv4 地址转换为大整数以进行比较
	start := big.NewInt(0)
	start.SetBytes(startIP.To4())

	end := big.NewInt(0)
	end.SetBytes(endIP.To4())

	// 检查起始 IP 是否小于等于结束 IP
	return start.Cmp(end) <= 0
}
