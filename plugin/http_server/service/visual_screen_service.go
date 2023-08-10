package service

import (
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

// 获取Visual列表
func AllVisual() []model.MVisual {
	m := []model.MVisual{}
	sqlitedao.Sqlite.DB().Find(&m)
	return m

}
func GetVisualWithUUID(uuid string) (*model.MVisual, error) {
	m := model.MVisual{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Visual
func DeleteVisual(uuid string) error {
	return sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MVisual{}).Error
}

// 创建Visual
func InsertVisual(Visual model.MVisual) error {
	return sqlitedao.Sqlite.DB().Create(&Visual).Error
}

// 更新Visual
func UpdateVisual(Visual model.MVisual) error {
	m := model.MVisual{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", Visual.UUID).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(Visual)
		return nil
	}
}
