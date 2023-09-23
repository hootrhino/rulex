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
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

// 获取DataSchema列表
func AllDataSchema() []model.MDataSchema {
	m := []model.MDataSchema{}
	interdb.DB().Find(&m)
	return m

}
func GetDataSchemaWithUUID(uuid string) (*model.MDataSchema, error) {
	m := model.MDataSchema{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除DataSchema
func DeleteDataSchema(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MDataSchema{}).Error
}

// 创建DataSchema
func InsertDataSchema(DataSchema model.MDataSchema) error {
	return interdb.DB().Create(&DataSchema).Error
}

// 更新DataSchema
func UpdateDataSchema(DataSchema model.MDataSchema) error {
	m := model.MDataSchema{}
	if err := interdb.DB().Where("uuid=?", DataSchema.UUID).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(DataSchema)
		return nil
	}
}
