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
	"github.com/hootrhino/rulex/component/hwportmanager"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type HwPortVo struct {
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`   // 接口名称
	Type        string         `json:"type"`   // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias       string         `json:"alias"`  // 别名
	Config      any            `json:"config"` // 配置
	Busy        bool           `json:"busy"`   // 运行时数据，是否被占
	OccupyBy    HwPortOccupyVo `json:"occupyBy"`
	Description string         `json:"description"` // 额外备注

}
type HwPortOccupyVo struct {
	UUID string `json:"uuid"` // UUID
	Type string `json:"type"` // DEVICE, Other......
	Name string `json:"name"` // DEVICE, Other......
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
* 针对刚插入硬件的情况，需要及时刷新
*
 */
func RefreshPortList(c *gin.Context, ruleEngine typex.RuleX) {
	if err := service.InitHwPortConfig(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 硬件接口
*
 */
func AllHwPorts(c *gin.Context, ruleEngine typex.RuleX) {
	HwPortVos := []HwPortVo{}
	for _, port := range hwportmanager.AllHwPort() {
		HwPortVos = append(HwPortVos, HwPortVo{
			UUID:  port.UUID,
			Name:  port.Name,
			Type:  port.Type,
			Alias: port.Alias,
			Busy:  port.Busy,
			OccupyBy: HwPortOccupyVo{
				UUID: port.OccupyBy.UUID,
				Type: port.OccupyBy.Type,
				Name: port.OccupyBy.Name,
			},
			Description: port.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(HwPortVos))
}

/*
*
* 更新接口参数
*
 */
func UpdateHwPortConfig(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUID   string       `json:"uuid"`
		Config UartConfigVo `json:"config"` // 配置, 串口配置、或者网卡、USB等
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.UpdateHwPortConfig(model.MHwPort{
		UUID:   form.UUID,
		Config: form.Config.JsonString(),
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MHwPort, err1 := service.GetHwPortConfig(form.UUID)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	HwIPort := hwportmanager.RhinoH3HwPort{
		UUID:        MHwPort.UUID,
		Name:        MHwPort.Name,
		Type:        MHwPort.Type,
		Alias:       MHwPort.Alias,
		Description: MHwPort.Description,
	}
	// 串口类
	if MHwPort.Type == "UART" {
		config := hwportmanager.UartConfig{}
		utils.BindConfig(MHwPort.GetConfig(), &config)
		HwIPort.Config = config
	}
	if MHwPort.Type == "FD" {
		HwIPort.Config = nil
	}
	// 刷新接口参数
	hwportmanager.RefreshPort(HwIPort)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 获取详情
*
 */
func GetHwPortDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Port, err1 := hwportmanager.GetHwPort(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(HwPortVo{
		UUID:        Port.UUID,
		Name:        Port.Name,
		Type:        Port.Type,
		Alias:       Port.Alias,
		Config:      Port.Config,
		Description: Port.Description,
		Busy:        Port.Busy,
		OccupyBy: HwPortOccupyVo{
			Port.OccupyBy.UUID, Port.OccupyBy.Type, Port.Name,
		},
	}))

}
