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
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/core"
	common "github.com/hootrhino/rulex/plugin/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* 获取一机一密
*
 */

func GetVendorKey(c *gin.Context, ruleEngine typex.RuleX) {
	cfg, _ := ini.ShadowLoad(core.INIPath)
	sections := cfg.ChildSections("plugin")
	license := ""
	for _, section := range sections {
		name := strings.TrimPrefix(section.Name(), "plugin.")
		if name == "license_manager" {
			license_path, err1 := section.GetKey("license_path")
			if err1 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err1))
				return
			}
			readBytes, err2 := os.ReadFile(license_path.String())
			if err2 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err2))
				return
			}
			license += string(readBytes)
			break
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData(license))
}
