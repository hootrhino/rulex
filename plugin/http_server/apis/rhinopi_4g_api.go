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
	"strings"

	"github.com/gin-gonic/gin"
	archsupport "github.com/hootrhino/rulex/bspsupport"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 重启4G
*
 */
func RhinoPiRestart4G(c *gin.Context, ruleEngine typex.RuleX) {
	_, err := archsupport.RhinoPiRestart4G()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取4G基本信息
*
 */
func Get4GBaseInfo(c *gin.Context, ruleEngine typex.RuleX) {
	csq := archsupport.RhinoPiGet4GCSQ()
	cops, err1 := archsupport.RhinoPiGetCOPS()
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	cm := "UNKNOWN"
	if strings.Contains(cops, "CMCC") {
		cm = "中国移动"
	}
	if strings.Contains(cops, "MOBILE") {
		cm = "中国移动"
	}
	if strings.Contains(cops, "UNICOM") {
		cm = "中国联通"
	}
	iccid, err2 := archsupport.RhinoPiGetICCID()
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(
		map[string]interface{}{
			"csq":   csq,
			"cops":  cm,
			"iccid": strings.TrimLeft(iccid, "+QCCID: "),
		},
	))

}

/*
*
* 信号强度
*
 */
func Get4GCSQ(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(archsupport.RhinoPiGet4GCSQ()))
}

// (1,"CHINA MOBILE","CMCC","46000",0),
// (3,"CHN-UNICOM","UNICOM","46001",7),
// +COPS: 0,0,\"CHINA MOBILE\",7
func Get4GCOPS(c *gin.Context, ruleEngine typex.RuleX) {
	result, err := archsupport.RhinoPiGetCOPS()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		cm := "UNKNOWN"
		if strings.Contains(result, "CMCC") {
			cm = "中国移动"
		}
		if strings.Contains(result, "MOBILE") {
			cm = "中国移动"
		}
		if strings.Contains(result, "UNICOM") {
			cm = "中国联通"
		}
		c.JSON(common.HTTP_OK, common.OkWithData(cm))
	}
}

/*
*
* 获取电话卡ICCID
*
 */
func Get4GICCID(c *gin.Context, ruleEngine typex.RuleX) {
	result, err := archsupport.RhinoPiGetICCID()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		// +QCCID: 89860426102180397625
		iccid := strings.TrimLeft(result, "+QCCID: ")
		c.JSON(common.HTTP_OK, common.OkWithData(iccid))
	}
}
