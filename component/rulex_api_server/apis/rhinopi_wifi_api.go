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
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/ossupport"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 扫描WIFI
*
 */
func ScanWIFIWithNmcli(c *gin.Context, ruleEngine typex.RuleX) {
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	SupportWifi := false
	for _, IFace := range interfaces {
		if strings.Contains(IFace.Name, "wlan") {
			SupportWifi = true
			break
		}
	}
	if !SupportWifi {
		c.JSON(common.HTTP_OK, common.Error("Device not support Wifi"))
		return
	}
	Wlans, err := ossupport.ScanWIFIWithNmcli()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Wlans))
}
