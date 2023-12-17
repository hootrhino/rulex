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
package model

import (
	"github.com/hootrhino/rulex/component/rulex_api_server/dto"
	"gopkg.in/square/go-jose.v2/json"
)

/*
*
* 用户自定义代码模板
*
 */
type MUserLuaTemplate struct {
	RulexModel
	UUID      string
	Gid       string // 分组
	Type      string // 类型 固定为 'function'
	Label     string //快捷代码名称
	Apply     string //快捷代码
	Variables string //变量
	Detail    string
}

/*
*
* 获取其变量表
*
 */
func (md MUserLuaTemplate) GetVariables() []dto.LuaTemplateVariables {
	result := make([]dto.LuaTemplateVariables, 0)
	err := json.Unmarshal([]byte(md.Variables), &result)
	if err != nil {
		return result
	}
	return result
}

/*
*
* 生成字符串
*
 */
func (md MUserLuaTemplate) GenVariables(V []dto.LuaTemplateVariables) (string, error) {
	B, err := json.Marshal(V)
	if err != nil {
		return "", err
	}
	return string(B), nil
}
