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
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/rulex_api_server/common"
	"github.com/hootrhino/rulex/plugin/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 内部事件
*
 */
type InternalNotifyVo struct {
	UUID    string `json:"uuid"`           // UUID
	Type    string `json:"type"`           // INFO | ERROR | WARNING
	Status  int    `json:"status"`         // 1 未读 2 已读
	Event   string `json:"event"`          // 字符串
	Ts      uint64 `json:"ts"`             // 时间戳
	Summary string `json:"summary"`        // 概览，为了节省流量，在消息列表只显示这个字段，Info值为“”
	Info    string `json:"info,omitempty"` // 消息内容，是个文本，详情显示
}

/*
*
* 站内消息
*
 */
func InternalNotifiesHeader(c *gin.Context, ruleEngine typex.RuleX) {
	data := []InternalNotifyVo{}
	models := service.AllInternalNotifiesHeader()
	for _, model := range models {
		data = append(data, InternalNotifyVo{
			UUID:    model.UUID,
			Type:    model.Type,
			Event:   model.Event,
			Ts:      model.Ts,
			Summary: model.Summary,
			Status:  model.Status,
		})

	}
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

/*
*
* 站内消息
*
 */
func InternalNotifies(c *gin.Context, ruleEngine typex.RuleX) {
	data := []InternalNotifyVo{}
	models := service.AllInternalNotifies()
	for _, model := range models {
		data = append(data, InternalNotifyVo{
			UUID:    model.UUID,
			Type:    model.Type,
			Event:   model.Event,
			Ts:      model.Ts,
			Summary: model.Summary,
			Info:    model.Info,
			Status:  model.Status,
		})

	}
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

/*
*
* 清空
*
 */
func ClearInternalNotifies(c *gin.Context, ruleEngine typex.RuleX) {
	if err := service.ClearInternalNotifies(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 阅读
*
 */
func ReadInternalNotifies(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if err := service.ReadInternalNotifies(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
