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

package service

import (
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
)

// 获取UserLuaTemplate列表
func AllUserLuaTemplate() []model.MUserLuaTemplate {
	m := []model.MUserLuaTemplate{}
	interdb.DB().Find(&m)
	return m

}

/*
*
* 获取分组
*
 */
func GetUserLuaTemplateGroup(rid string) model.MGenericGroup {
	sql := `
SELECT m_generic_groups.*
  FROM m_generic_group_relations
       LEFT JOIN
       m_generic_groups ON (m_generic_groups.uuid = m_generic_group_relations.gid)
 WHERE m_generic_group_relations.rid = ?;
`
	m := model.MGenericGroup{}
	interdb.DB().Raw(sql, rid).Find(&m)
	return m
}

/*
*
* ID获取
*
 */
func GetUserLuaTemplateWithUUID(uuid string) (model.MUserLuaTemplate, error) {
	m := model.MUserLuaTemplate{}
	err := interdb.DB().Where("uuid=?", uuid).First(&m).Error
	return m, err
}

// 删除UserLuaTemplate
func DeleteUserLuaTemplate(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MUserLuaTemplate{}).Error
}

// 创建UserLuaTemplate
func InsertUserLuaTemplate(UserLuaTemplate model.MUserLuaTemplate) error {
	return interdb.DB().Create(&UserLuaTemplate).Error
}

// 更新UserLuaTemplate
func UpdateUserLuaTemplate(UserLuaTemplate model.MUserLuaTemplate) error {
	return interdb.DB().
		Model(UserLuaTemplate).
		Where("uuid=?", UserLuaTemplate.UUID).
		Updates(&UserLuaTemplate).Error
}
