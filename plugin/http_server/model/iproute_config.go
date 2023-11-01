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

package model

/*
*
* Linux默认路由，该配置主要用来处理软路由相关功能
// 每次启动的时候换成最新配置的路由, 默认是ETH1 192.168.64.0
// 这个初始化的目的是为了配合软路由使用, 和isc-dhcp-server、dnsmasq 两个DHCP服务有关
*/
type MIpRoute struct {
	RulexModel
	UUID        string `gorm:"not null"`
	Iface       string `gorm:"not null"` // 用来做子网的那个网卡的网卡名
	Ip          string `gorm:"not null"` // 用来做子网的那个网卡的IP地址
	Gateway     string `gorm:"not null"` // 用来做子网的那个网卡的网关
	Network     string `gorm:"not null"` // 用来做子网的那个网卡的网段
	Netmask     string `gorm:"not null"` // 用来做子网的那个网卡子网掩码
	IpPoolBegin string `gorm:"not null"` // DHCP IP地址池起始
	IpPoolEnd   string `gorm:"not null"` // DHCP IP地址池结束
	//------------------------------------
	// IP 路由方向, 默认 ETH1 透传到 4G
	//------------------------------------
	IfaceFrom    string `gorm:"not null"` // 流量入口
	IfaceTo      string `gorm:"not null"` // 流量出口
}
