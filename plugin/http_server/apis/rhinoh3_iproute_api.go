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

package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/ossupport"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* DHCP 配置
*
 */
type DHCPVo struct {
	Iface       string `json:"iface" validate:"required"`         // 用来做子网的那个网卡的网卡名
	Ip          string `json:"ip" validate:"required"`            // 用来做子网的那个网卡的IP地址
	Gateway     string `json:"gateway" validate:"required"`       // 用来做子网的那个网卡的网关
	Network     string `json:"network" validate:"required"`       // 用来做子网的那个网卡的网段
	Netmask     string `json:"netmask" validate:"required"`       // 用来做子网的那个网卡子网掩码
	IpPoolBegin string `json:"ip_pool_begin" validate:"required"` // DHCP IP地址池起始
	IpPoolEnd   string `json:"ip_pool_end" validate:"required"`   // DHCP IP地址池结束
	//------------------------------------
	// IP 路由方向, 默认 ETH1 透传到 4G
	//------------------------------------
	IfaceFrom string `json:"iface_from" validate:"required"` // 流量入口,固定ETH1
	IfaceTo   string `json:"iface_to" validate:"required"`   // 流量出口,固定4G
}

func SetDHCP(c *gin.Context, ruleEngine typex.RuleX) {
	dhcpVo := DHCPVo{}
	if err := c.ShouldBindJSON(&dhcpVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	err1 := service.UpdateIpRoute(model.MIpRoute{
		Iface:       dhcpVo.Iface,
		Ip:          dhcpVo.Ip,
		Gateway:     dhcpVo.Gateway,
		Network:     dhcpVo.Network,
		Netmask:     dhcpVo.Netmask,
		IpPoolBegin: dhcpVo.IpPoolBegin,
		IpPoolEnd:   dhcpVo.IpPoolEnd,
		IfaceFrom:   dhcpVo.IfaceFrom,
		IfaceTo:     dhcpVo.IfaceTo,
	})
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(200, common.Ok())
}
func GetDHCP(c *gin.Context, ruleEngine typex.RuleX) {
	Model, err := service.GetDefaultIpRoute()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(200, common.OkWithData(DHCPVo{
		Iface:       Model.Iface,
		Ip:          Model.Ip,
		Gateway:     Model.Gateway,
		Network:     Model.Network,
		Netmask:     Model.Netmask,
		IpPoolBegin: Model.IpPoolBegin,
		IpPoolEnd:   Model.IpPoolEnd,
		IfaceFrom:   Model.IfaceFrom,
		IfaceTo:     Model.IfaceTo,
	}))
}

/*
*
  - 更新默认路由
    1 取上一次的路由
    2 配置最新的
    4 删除上一次的
    5 更新最新的路由

*
*/
type IpRouteVo struct {
	UUID  string `json:"uuid"`
	Ip    string `json:"ip" validate:"required"`
	Iface string `json:"iface" validate:"required"`
}

/*
*
* 获取上一次的路由
*
 */
func GetOldDefaultIpRoute(c *gin.Context, ruleEngine typex.RuleX) {
	MIpRoute, err := service.IpRouteDetail()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(IpRouteVo{
		Ip:    MIpRoute.Ip,
		Iface: MIpRoute.Iface,
	}))

}

/*
*
* 设置默认路由
*
 */

func SetNewDefaultIpRoute(c *gin.Context, ruleEngine typex.RuleX) {
	form := IpRouteVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ifaces, err1 := ossupport.GetBSPNetIfaces()
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	if !utils.SContains(ifaces, form.Iface) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Only support iface:%v", ifaces)))
		return
	}
	err3 := service.UpdateIpRoute(model.MIpRoute{
		Ip:    form.Ip,
		Iface: form.Iface,
	})
	if err3 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err3))
		return
	}
	err2 := service.UpdateDefaultRoute(form.Ip, form.Iface)
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取当前在线的DHCP主机列表
*
 */
type DhcpLeaseVo struct {
	MacAddress string `json:"mac_address"` // MAC地址
	IpAddress  string `json:"ip_address"`  // IP地址
	Hostname   string `json:"hostname"`    // 主机名
}

func GetDhcpClients(c *gin.Context, ruleEngine typex.RuleX) {
	// GetDhcpClients, err := ossupport.GetDhcpList()
	// if err != nil {
	// 	c.JSON(common.HTTP_OK, common.Error400(err))
	// 	return
	// }
	// 测试假数据
	Clients := []DhcpLeaseVo{
		{
			MacAddress: "a8:a1:59:2e:a2:d9",
			IpAddress:  "192.168.1.175",
			Hostname:   "rulex-h1",
		},
		{
			MacAddress: "a8:a1:59:2e:a2:d9",
			IpAddress:  "192.168.1.176",
			Hostname:   "rulex-h2",
		},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Clients))
}
