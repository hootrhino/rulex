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
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/utils"
	"gorm.io/gorm"
)

// 获取GenericGroup列表
func AllGenericGroup() []model.MGenericGroup {
	m := []model.MGenericGroup{}
	interdb.DB().Find(&m)
	return m
}
func ListByGroupType(t string) []model.MGenericGroup {
	m := []model.MGenericGroup{}
	interdb.DB().Where("type=?", t).Find(&m)
	return m
}

/*
*
* 查询分组下的设备
*
 */
func FindDeviceByGroup(gid string) []model.MDevice {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC;`

	m := []model.MDevice{}
	interdb.DB().Raw(`SELECT * FROM m_devices `+sql, gid).Find(&m)
	return m

}

/*
*
* 新增的分页获取
*
 */
func PageDeviceByGroup(current, size int, gid string) (int64, []model.MDevice) {
	sql := `
SELECT * FROM m_devices WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN m_generic_group_relations ON
		(m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC limit ? offset ?;`
	MDevices := []model.MDevice{}
	offset := (current - 1) * size
	interdb.DB().Raw(sql, gid, size, offset).Find(&MDevices)
	var count int64
	interdb.DB().Model(&model.MDevice{}).Count(&count)
	return count, MDevices
}

/*
*
* 查询分组吓得大屏
*
 */
func FindVisualByGroup(uuid string) []model.MVisual {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	  WHERE type = 'VISUAL' AND gid = ?
) ORDER BY created_at DESC;`

	m := []model.MVisual{}
	interdb.DB().Raw(`SELECT * FROM m_visuals `+sql, uuid).Find(&m)
	return m

}

/*
*
  - 根据分组类型查询:代码模板

*~
*/
func FindUserTemplateByGroup(uuid string) []model.MUserLuaTemplate {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	  WHERE type = 'USER_LUA_TEMPLATE' AND gid = ?
) ORDER BY created_at DESC;`
	m := []model.MUserLuaTemplate{}
	interdb.DB().Raw(`SELECT * FROM m_user_lua_templates `+sql, uuid).Find(&m)
	return m

}

func GetGenericGroupWithUUID(uuid string) (*model.MGenericGroup, error) {
	m := model.MGenericGroup{}
	if err := interdb.DB().
		Where("uuid=?", uuid).
		First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除GenericGroup
func DeleteGenericGroup(uuid string) error {
	return interdb.DB().
		Where("uuid=?", uuid).
		Delete(&model.MGenericGroup{}).Error
}

// 创建GenericGroup
func InsertGenericGroup(GenericGroup *model.MGenericGroup) error {
	return interdb.DB().Create(GenericGroup).Error
}

// 创建GenericGroup
func InitGenericGroup(GenericGroup *model.MGenericGroup) error {
	return interdb.DB().Model(GenericGroup).
		Where("type=?", GenericGroup.Type).
		FirstOrCreate(GenericGroup).Error
}

// 更新GenericGroup
func UpdateGenericGroup(GenericGroup *model.MGenericGroup) error {
	m := model.MGenericGroup{}
	if err := interdb.DB().
		Where("uuid=?", GenericGroup.UUID).
		First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*GenericGroup)
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
	if err := interdb.DB().Where("uuid=?", gid).First(&m).Error; err != nil {
		return err
	}
	Relation := model.MGenericGroupRelation{
		UUID: utils.MakeUUID("GR"),
		Gid:  m.UUID,
		Rid:  rid,
	}
	if err := interdb.DB().Save(&Relation).Error; err != nil {
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
	return interdb.DB().
		Where("gid=? and rid =?", gid, rid).
		Delete(&model.MGenericGroupRelation{}).Error
}

/*
*
* 检查是否绑定
*
 */
func CheckBindResource(gid string) (uint, error) {
	sql := `SELECT count(*) FROM m_generic_group_relations WHERE gid = ?;`
	count := 0
	err := interdb.DB().Raw(sql, gid).Find(&count).Error
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
	err := interdb.DB().Raw(sql, gid, rid).Find(&count).Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

/*
*
* 重新绑定分组需要事务支持
*
 */
func ReBindResource(action func(tx *gorm.DB) error, Rid, Gid string) error {
	return interdb.DB().Transaction(func(tx *gorm.DB) error {
		// 1 执行的操作
		if err0 := action(tx); err0 != nil {
			return err0
		}
		// 2 解除分组关联
		sql := `
SELECT m_generic_groups.*
FROM m_generic_group_relations
LEFT JOIN
m_generic_groups ON (m_generic_groups.uuid = m_generic_group_relations.gid)
WHERE m_generic_group_relations.rid = ?;
		`
		OldGroup := model.MGenericGroup{}
		if errA := tx.Raw(sql, Rid).Find(&OldGroup).Error; errA != nil {
			return errA
		}
		if err1 := tx.Model(model.MGenericGroupRelation{}).
			Where("gid=? and rid =?", OldGroup.UUID, Rid).
			Delete(&model.MGenericGroupRelation{}).Error; err1 != nil {
			return err1
		}
		// 3 重新绑定分组,首先确定分组是否存在
		MGroup := model.MGenericGroup{}
		if err2 := tx.Model(MGroup).
			Where("uuid=?", Gid).
			First(&MGroup).Error; err2 != nil {
			return err2
		}
		// 4 重新绑定分组
		err3 := tx.Save(&model.MGenericGroupRelation{
			UUID: utils.MakeUUID("GRLT"),
			Gid:  Gid,
			Rid:  Rid,
		}).Error
		if err3 != nil {
			return err3
		}
		return nil
	})
}
