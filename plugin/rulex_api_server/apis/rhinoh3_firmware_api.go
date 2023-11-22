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

// RhinoH3 固件相关操作
package apis

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/ossupport"
	common "github.com/hootrhino/rulex/plugin/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 重启固件
*
 */
func ReStartRulex(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 重启操作系统
*
 */
func Reboot(c *gin.Context, ruleEngine typex.RuleX) {
	err := ossupport.Reboot()
	if err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 回复出厂
*
 */
func RecoverNew(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取本地升级日志
*
 */
func GetUpGradeLog(c *gin.Context, ruleEngine typex.RuleX) {
	byteS, _ := os.ReadFile("local-upgrade-log.txt")
	c.JSON(common.HTTP_OK, common.OkWithData(string(byteS)))
}
