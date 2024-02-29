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
* Hnc8点位表管理
*
 */
// InsertHnc8PointPosition 插入Hnc8点位表
func InsertHnc8PointPositions(list []model.MHnc8DataPoint) error {
	m := model.MHnc8DataPoint{}
	return interdb.DB().Model(m).Create(list).Error
}

// InsertHnc8PointPosition 插入Hnc8点位表
func InsertHnc8PointPosition(P model.MHnc8DataPoint) error {
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

// DeleteHnc8PointByDevice 删除Hnc8点位与设备
func DeleteHnc8PointByDevice(uuids []string, deviceUuid string) error {
	return interdb.DB().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MHnc8DataPoint{}).Error
}

// DeleteAllHnc8PointByDevice 删除Hnc8点位与设备
func DeleteAllHnc8PointByDevice(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MHnc8DataPoint{}).Error
}

// 更新DataSchema
func UpdateHnc8Point(MHnc8DataPoint model.MHnc8DataPoint) error {
	return interdb.DB().Model(model.MHnc8DataPoint{}).
		Where("device_uuid=? AND uuid=?",
			MHnc8DataPoint.DeviceUuid, MHnc8DataPoint.UUID).
		Updates(MHnc8DataPoint).Error
}
