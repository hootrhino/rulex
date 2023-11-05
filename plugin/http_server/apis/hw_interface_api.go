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
	"encoding/json"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type hwInterfaceVo struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`             // 接口名称
	Type        string        `json:"type"`             // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias       string        `json:"alias"`            // 别名
	Config      *UartConfigVo `json:"config,omitempty"` // 配置
	Description string        `json:"description"`      // 额外备注

}
type UartConfigVo struct {
	Timeout  int    `json:"timeout"`
	Uart     string `json:"uart"`
	BaudRate int    `json:"baudRate"`
	DataBits int    `json:"dataBits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stopBits"`
}

func (u UartConfigVo) JsonString() string {
	if bytes, err := json.Marshal(u); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

/*
*
* 硬件接口
*
 */
func AllHwInterfaces(c *gin.Context, ruleEngine typex.RuleX) {
	MHwInterfaces, err := service.AllHwInterface()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	hwInterfaceVos := []hwInterfaceVo{}
	for _, ifce := range MHwInterfaces {
		hwInterfaceVos = append(hwInterfaceVos, hwInterfaceVo{
			UUID:        ifce.UUID,
			Name:        ifce.Name,
			Type:        ifce.Type,
			Alias:       ifce.Alias,
			Description: ifce.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(hwInterfaceVos))
}

/*
*
*
*
 */
func UpdateHwInterfaceConfig(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUID   string       `json:"uuid"`
		Config UartConfigVo `json:"config"` // 配置, 串口配置、或者网卡、USB等
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.UpdateHwInterfaceConfig(model.MHwInterface{
		UUID:   form.UUID,
		Config: form.Config.JsonString(),
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// TODO: 更新串口
	// ApplyNewestConfig()
	//
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 获取详情
*
 */
func GetHwInterfaceDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Iface, err := service.GetHwInterfaceConfig(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 将 Model String 转成结构
	config := UartConfigVo{}
	if err := utils.BindConfig(Iface.GetConfig(), &config); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(hwInterfaceVo{
		UUID:        Iface.UUID,
		Name:        Iface.Name,
		Type:        Iface.Type,
		Alias:       Iface.Alias,
		Config:      &config,
		Description: Iface.Description,
	}))

}
