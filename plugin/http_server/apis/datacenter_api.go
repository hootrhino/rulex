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
	"github.com/hootrhino/rulex/component/datacenter"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 获取仓库细节
*
 */
func GetSchemaDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	schema := datacenter.GetSchemaDetail(uuid)
	c.JSON(common.HTTP_OK, common.OkWithData(schema))
}

/*
*
* 获取仓库列表
*
 */
func GetSchemaDefineList(c *gin.Context, ruleEngine typex.RuleX) {
	Column, err := datacenter.SchemaDefineList()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Column))
}

/*
*
* 获取单个仓库的表结构定义
*
 */
func GetSchemaDefine(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Schema, err := datacenter.GetSchemaDefine(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if Schema.UUID == "" {
		c.JSON(common.HTTP_OK, common.Error("Schema not found"))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Schema))
}

/*
*
* 获取仓库结构列表
*
 */
func GetSchemaList(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(datacenter.SchemaList()))
}

/*
*
* 执行查询
*
 */
func GetQueryData(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		Uuid  string `json:"uuid,omitempty"`
		Query string `json:"query,omitempty"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Column, err := datacenter.Query(form.Uuid, form.Query)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Column))
}
