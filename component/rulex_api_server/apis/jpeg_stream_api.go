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
	"github.com/hootrhino/rulex/component/jpegstream"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type JpegStreamVo struct {
	Type       string           `json:"type"` // push | pull
	LiveId     string           `json:"liveId"`
	Pulled     bool             `json:"pulled"`
	Resolution utils.Resolution `json:"resolution"`
}

/*
*
*  JpegStream列表
*
 */
func GetJpegStreamList(c *gin.Context, ruleEngine typex.RuleX) {
	JpegStreamVos := []JpegStreamVo{}
	JpegStreams := jpegstream.JpegStreamSourceList()
	for _, v := range JpegStreams {
		JpegStreamVos = append(JpegStreamVos, JpegStreamVo{
			Type:       v.Type,
			LiveId:     v.LiveId,
			Pulled:     v.Pulled,
			Resolution: v.Resolution,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(JpegStreamVos))

}

/*
*
*  JpegStream详情
*
 */
func GetJpegStreamDetail(c *gin.Context, ruleEngine typex.RuleX) {
	liveId, _ := c.GetQuery("liveId")
	JpegStream, err := jpegstream.GetJpegStreamSource(liveId)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(JpegStreamVo{
		Type:       JpegStream.Type,
		LiveId:     JpegStream.LiveId,
		Pulled:     JpegStream.Pulled,
		Resolution: JpegStream.Resolution,
	}))

}
