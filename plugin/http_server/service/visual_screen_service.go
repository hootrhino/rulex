package service

import (
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

// 获取Visual列表
func AllVisual() []model.MVisual {
	m := []model.MVisual{}
	interdb.DB().Find(&m)
	return m

}
func GetVisualGroup(rid string) model.MGenericGroup {
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
func GetVisualWithUUID(uuid string) (*model.MVisual, error) {
	m := model.MVisual{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Visual
func DeleteVisual(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MVisual{}).Error
}

// 创建Visual
func InsertVisual(Visual model.MVisual) error {
	return interdb.DB().Create(&Visual).Error
}

// 更新Visual
func UpdateVisual(Visual model.MVisual) error {
	m := model.MVisual{}
	if err := interdb.DB().Where("uuid=?", Visual.UUID).First(&m).Error; err != nil {
		return err
	}
	return interdb.DB().Model(m).Updates(Visual).Error
}
