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
	"fmt"
	"os/exec"
	"strings"

	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

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
func delDefaultRoute(ip, iface string) error {
	// sudo ip route del default via 192.168.43.1 dev usb0
	cmd := exec.Command("ip", "route", "del", "default", "via", ip, "dev", iface)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing del: %s", err.Error()+":"+string(output))
	}
	return nil
}

func UpdateDefaultRoute(oldGatewayIP, oldIface, newGatewayIP, newIface string) error {
	existsDefaultRoute, err := checkDefaultRoute(oldGatewayIP, oldIface)
	if err != nil {
		return err
	}
	if existsDefaultRoute {
		if err := delDefaultRoute(oldGatewayIP, oldIface); err != nil {
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
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", "0").First(&m).Error; err != nil {
		return model.MIpRoute{}, err
	} else {
		return m, nil
	}
}

// 更新 IpRoute
func UpdateIpRoute(IpRoute model.MIpRoute) error {
	return sqlitedao.Sqlite.DB().Model(IpRoute).Where("uuid=?", "0").Updates(IpRoute).Error
}

// Init
func InitDefaultIpRoute() error {
	m := model.MIpRoute{
		UUID:  "0",
		Ip:    "192.168.200.0",
		Iface: "eth1",
	}
	return sqlitedao.Sqlite.DB().Model(m).Where("uuid=?", "0").FirstOrCreate(&m).Error
}

// checkDefaultRoute 检查是否存在默认路由
func checkDefaultRoute(oldIp, oldIface string) (bool, error) {
	// 执行命令
	cmd := exec.Command("sh", "-c", "ip route | awk 'NR==1 {print $1 ,$2, $3, $4, $5}'")
	// 捕获命令的标准输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("checkDefaultRoute error: %s", string(output))
	}
	outputStr := string(output)
	//
	// ip route | awk 'NR==1 {print $1 ,$2, $3, $4, $5}'
	// default via 192.168.199.1 dev wlx0cc6551c5026
	// default via %s dev %s
	//
	// 将 AWK 输出转换为字符串并去除空白字符
	awkResult := strings.TrimSpace(outputStr)

	// 检查是否为 "default"
	return awkResult == fmt.Sprintf("default via %s dev %s", oldIp, oldIface), nil
}
