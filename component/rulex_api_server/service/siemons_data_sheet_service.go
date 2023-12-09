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

/*
*
* 西门子点位表管理
*
 */
// InsertSiemensPoint 插入modbus点位表
func InsertSiemensPoint(list []model.MModbusDataPoint) error {
	m := model.MModbusDataPoint{}
	return interdb.DB().Model(m).Create(list).Error
}

// DeleteSiemensPointAndDevice 删除modbus点位与设备
func DeleteSiemensPointByDevice(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MModbusDataPoint{}).Error
}

// AllSiemensPointByDeviceUuid 根据设备UUID查询设备点位
func AllSiemensPointByDeviceUuid(deviceUuid string,
	page, pageSize int) (list []model.MModbusDataPoint, err error) {
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	err = interdb.DB().Where("device_uuid=?", deviceUuid).
		Find(&list).Offset(int(page)).Limit(pageSize).Error
	return
}

// 更新DataSchema
func UpdateSiemensPoint(MModbusDataPoint model.MModbusDataPoint) error {
	return interdb.DB().Model(MModbusDataPoint).Updates(&MModbusDataPoint).Error
}
