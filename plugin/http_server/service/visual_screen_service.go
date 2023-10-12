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
func GetVisualWithUUID(uuid string) (model.MVisual, error) {
	m := model.MVisual{}
	err := interdb.DB().Where("uuid=?", uuid).First(&m).Error
	return m, err
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
	return interdb.DB().
		Model(Visual).
		Where("uuid=?", Visual.UUID).
		Updates(&Visual).Error
}
