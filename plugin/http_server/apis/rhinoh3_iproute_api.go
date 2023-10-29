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
	archsupport "github.com/hootrhino/rulex/bspsupport"
	"github.com/hootrhino/rulex/ossupport"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

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
	UUID  string `json:"uuid,omitempty"`
	Ip    string `json:"ip" validate:"required"`
	Iface string `json:"iface" validate:"required"`
}

/*
*
* 获取上一次的路由
*
 */
func GetOldDefaultIpRoute(c *gin.Context, ruleEngine typex.RuleX) {
	MIpRoute, err := ossupport.IpRouteDetail()
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
	ifaces, err1 := archsupport.GetBSPNetIfaces()
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	if !utils.SContains(ifaces, form.Iface) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Only support iface:%v", ifaces)))
		return
	}
	err3 := ossupport.UpdateIpRoute(model.MIpRoute{
		Ip:    form.Ip,
		Iface: form.Iface,
	})
	if err3 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err3))
		return
	}
	err2 := ossupport.UpdateDefaultRoute(form.Ip, form.Iface)
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
func GetDhcpClients(c *gin.Context, ruleEngine typex.RuleX) {
	GetDhcpClients, err := ossupport.GetDhcpList()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(GetDhcpClients))
}
