// Copyright (C) 2024 wwhai
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
	"github.com/hootrhino/rulex/component/rtspserver"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type FlvStreamVo struct {
	Type       string           `json:"type"` // push | pull
	LiveId     string           `json:"liveId"`
	Pulled     bool             `json:"pulled"`
	Resolution utils.Resolution `json:"resolution"`
}

/*
*
*  FlvStream列表
*
 */
func GetFlvStreamList(c *gin.Context, ruleEngine typex.RuleX) {
	FlvStreamVos := []FlvStreamVo{}
	FlvStreams := rtspserver.FlvStreamSourceList()
	for _, v := range FlvStreams {
		FlvStreamVos = append(FlvStreamVos, FlvStreamVo{
			Type:       v.Type,
			LiveId:     v.LiveId,
			Pulled:     v.Pulled,
			Resolution: v.Resolution,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(FlvStreamVos))

}

/*
*
*  FlvStream详情
*
 */
func GetFlvStreamDetail(c *gin.Context, ruleEngine typex.RuleX) {
	liveId, _ := c.GetQuery("liveId")
	FlvStream, err := rtspserver.GetFlvStreamSource(liveId)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(FlvStreamVo{
		Type:       FlvStream.Type,
		LiveId:     FlvStream.LiveId,
		Pulled:     FlvStream.Pulled,
		Resolution: FlvStream.Resolution,
	}))

}
