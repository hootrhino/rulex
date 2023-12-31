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

// 获取DataSchema列表
func AllDataSchema() []model.MIotSchema {
	m := []model.MIotSchema{}
	interdb.DB().Find(&m)
	return m

}
func GetDataSchemaWithUUID(uuid string) (model.MIotSchema, error) {
	m := model.MIotSchema{}
	return m, interdb.DB().Where("uuid=?", uuid).First(&m).Error
}

// 删除DataSchema
func DeleteDataSchema(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MIotSchema{}).Error
}

// 创建DataSchema
func InsertDataSchema(DataSchema model.MIotSchema) error {
	return interdb.DB().Create(&DataSchema).Error
}

// 更新DataSchema
func UpdateDataSchema(DataSchema model.MIotSchema) error {
	return interdb.DB().
		Model(DataSchema).
		Where("uuid=?", DataSchema.UUID).
		Updates(&DataSchema).Error
}
