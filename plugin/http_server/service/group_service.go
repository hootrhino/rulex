// Copyright (C) 2023 wangwenhai
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
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

// 获取GenericGroup列表
func AllGenericGroup() []model.MGenericGroup {
	m := []model.MGenericGroup{}
	sqlitedao.Sqlite.DB().Find(&m)
	return m
}

/*
*
  - 根据分组类型查询:DEVICE, VISUAL

*
*/
func FindByType(uuid, t string) ([]model.MVisual, []model.MDevice) {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	  WHERE type = ? AND gid = ?
);`
	if t == "VISUAL" {
		m := []model.MVisual{}
		sqlitedao.Sqlite.DB().Raw(`SELECT * FROM m_visuals `+sql, t, uuid).Find(&m)
		return m, nil
	}
	if t == "DEVICE" {
		m := []model.MDevice{}
		sqlitedao.Sqlite.DB().Raw(`SELECT * FROM m_devices `+sql, t, uuid).Find(&m)
		return nil, m
	}
	return nil, nil
}

func GetGenericGroupWithUUID(uuid string) (*model.MGenericGroup, error) {
	m := model.MGenericGroup{}
	if err := sqlitedao.Sqlite.DB().
		Where("uuid=?", uuid).
		First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除GenericGroup
func DeleteGenericGroup(uuid string) error {
	return sqlitedao.Sqlite.DB().
		Where("uuid=?", uuid).
		Delete(&model.MGenericGroup{}).Error
}

// 创建GenericGroup
func InsertGenericGroup(GenericGroup *model.MGenericGroup) error {
	return sqlitedao.Sqlite.DB().Create(GenericGroup).Error
}

// 更新GenericGroup
func UpdateGenericGroup(GenericGroup *model.MGenericGroup) error {
	m := model.MGenericGroup{}
	if err := sqlitedao.Sqlite.DB().
		Where("uuid=?", GenericGroup.UUID).
		First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*GenericGroup)
		return nil
	}
}

/*
*
* 将别的东西加入到分组里面
*
 */
func BindResource(gid, rid string) error {
	m := model.MGenericGroup{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", gid).First(&m).Error; err != nil {
		return err
	}
	Relation := model.MGenericGroupRelation{
		Gid: m.UUID,
		Rid: rid,
	}
	if err := sqlitedao.Sqlite.DB().Save(&Relation).Error; err != nil {
		return err
	}
	return nil
}

/*
*
* 取消分组绑定
*
 */
func UnBindResource(gid, rid string) error {
	model := model.MGenericGroupRelation{
		Gid: gid,
		Rid: rid,
	}
	return sqlitedao.Sqlite.DB().Delete(&model).Error
}

/*
*
* 检查是否绑定
*
 */
func CheckBindResource(gid string) (uint, error) {
	sql := `SELECT count(*) FROM m_generic_group_relations WHERE gid = ?;`
	count := 0
	err := sqlitedao.Sqlite.DB().Raw(sql, gid).Find(&count).Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

/*
*
* 检查是否重复了
*
 */
func CheckAlreadyBinding(gid, rid string) (uint, error) {
	sql := `SELECT count(*) FROM m_generic_group_relations WHERE gid = ? and rid = ?;`
	count := 0
	err := sqlitedao.Sqlite.DB().Raw(sql, gid, rid).Find(&count).Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}
