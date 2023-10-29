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
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
)

type SiteConfigVo struct {
	SiteName string `json:"siteName"`
	Logo     string `json:"logo"`
	AppName  string `json:"appName"`
}

func UpdateSiteConfig(c *gin.Context, ruleEngine typex.RuleX) {

	form := SiteConfigVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.UpdateSiteConfig(model.MSiteConfig{
		SiteName: form.SiteName,
		Logo:     form.Logo,
		AppName:  form.AppName,
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
func GetSiteConfig(c *gin.Context, ruleEngine typex.RuleX) {
	Model, err := service.GetSiteConfig()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(SiteConfigVo{
		SiteName: Model.SiteName,
		Logo:     Model.Logo,
		AppName:  Model.AppName,
	}))
}
