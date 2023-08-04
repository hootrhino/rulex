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
func GetGenericGroupWithUUID(uuid string) (*model.MGenericGroup, error) {
	m := model.MGenericGroup{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除GenericGroup
func DeleteGenericGroup(uuid string) error {
	return sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MGenericGroup{}).Error
}

// 创建GenericGroup
func InsertGenericGroup(GenericGroup *model.MGenericGroup) error {
	return sqlitedao.Sqlite.DB().Create(GenericGroup).Error
}

// 更新GenericGroup
func UpdateGenericGroup(GenericGroup *model.MGenericGroup) error {
	m := model.MGenericGroup{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", GenericGroup.UUID).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*GenericGroup)
		return nil
	}
}
