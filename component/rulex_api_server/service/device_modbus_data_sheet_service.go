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
	"fmt"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
)

/*
*
* Siemens点位表管理
*
 */
// InsertSiemensPointPosition 插入Siemens点位表
func InsertSiemensPointPositions(list []model.MSiemensDataPoint) error {
	m := model.MSiemensDataPoint{}
	return interdb.DB().Model(m).Create(list).Error
}

// InsertSiemensPointPosition 插入Siemens点位表
func InsertSiemensPointPosition(P model.MSiemensDataPoint) error {
	IgnoreUUID := P.UUID
	Count := int64(0)
	P.UUID = ""
	interdb.DB().Model(P).Where(P).Count(&Count)
	if Count > 0 {
		return fmt.Errorf("already exists same record:%s", IgnoreUUID)
	}
	P.UUID = IgnoreUUID
	return interdb.DB().Model(P).Create(&P).Error
}

// DeleteSiemensPointByDevice 删除Siemens点位与设备
func DeleteSiemensPointByDevice(uuids []string, deviceUuid string) error {
	return interdb.DB().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MSiemensDataPoint{}).Error
}

// DeleteAllSiemensPointByDevice 删除Siemens点位与设备
func DeleteAllSiemensPointByDevice(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MSiemensDataPoint{}).Error
}

// 更新DataSchema
func UpdateSiemensPoint(MSiemensDataPoint model.MSiemensDataPoint) error {
	return interdb.DB().Model(model.MSiemensDataPoint{}).
		Where("device_uuid=? AND uuid=?",
			MSiemensDataPoint.DeviceUuid, MSiemensDataPoint.UUID).
		Updates(MSiemensDataPoint).Error
}
