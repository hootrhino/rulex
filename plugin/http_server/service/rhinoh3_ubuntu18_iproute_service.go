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
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

/*
*
* 删除原来的垃圾路由，换成最新的配置
*
 */
func UpdateDefaultRoute(newGatewayIP, newIface string) error {
	DefaultRoutes := getDefaultRoute()
	for _, route := range DefaultRoutes {
		if err := delDefaultRoute(route); err != nil {
			return err
		}
	}

	if err := addDefaultRoute(newGatewayIP, newIface); err != nil {
		return err
	}
	return nil
}

func IpRouteDetail() (model.MIpRoute, error) {
	m := model.MIpRoute{}
	if err := interdb.DB().Where("uuid=?", "0").First(&m).Error; err != nil {
		return model.MIpRoute{}, err
	} else {
		return m, nil
	}
}
func GetDefaultIpRoute() (model.MIpRoute, error) {
	return IpRouteDetail()
}

// 更新 IpRoute
func UpdateIpRoute(IpRoute model.MIpRoute) error {
	return interdb.DB().Model(IpRoute).Where("uuid=?", "0").Updates(IpRoute).Error
}

// 每次启动的时候换成最新配置的路由, 默认是ETH1 192.168.64.0
// 这个初始化的目的是为了配合软路由使用, 和isc-dhcp-server、dnsmasq 两个DHCP服务有关
func InitDefaultIpRoute() error {
	m := model.MIpRoute{
		UUID:        "0",
		Iface:       "eth1",
		Ip:          "192.168.64.100",
		Gateway:     "192.168.64.100",
		Netmask:     "255.255.255.0",
		IpPoolBegin: "192.168.64.100",
		IpPoolEnd:   "192.168.64.150",
		IfaceFrom:   "eth1", // 默认将Eth1网口和USB 4G网卡桥接起来, 这样就可以实现4G上网
		IfaceTo:     "usb1",
	}
	if err := interdb.DB().Model(m).
		Where("uuid=?", "0").
		FirstOrCreate(&m).Error; err != nil {
		return err
	}
	return nil
	// return UpdateDefaultRoute(m.Ip, m.Iface)
}

// getDefaultRoute 返回默认路由的信息作为字符串切片
func getDefaultRoute() []string {
	cmd := exec.Command("ip", "route", "show", "default")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	output := stdout.String()
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

/*
*
* 设置默认路由，只能有一条默认路由
*
 */
func addDefaultRoute(newGatewayIP, iface string) error {
	// sudo ip route add default via 192.168.1.1 dev eth0
	cmd := exec.Command("ip", "route", "add", "default", "via", newGatewayIP, "dev", iface)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing add: %s", err.Error()+":"+string(output))
	}
	return nil
}
func delDefaultRoute(route string) error {
	// sudo ip route del default via 192.168.43.1 dev usb0
	cmd := exec.Command("sh", "-c", "ip route del %s", route)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing del: %s", err.Error()+":"+string(output))
	}
	return nil
}
