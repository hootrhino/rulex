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
* InsertInternalNotifies
*
 */
func InsertInternalNotify(m model.MInternalNotify) error {
	var count int64
	interdb.DB().Model(&m).Count(&count)
	// 超过100条记录就清空
	if count > 100 {
		if err := ClearInternalNotifies(); err != nil {
			return err
		}
	}
	return interdb.DB().Model(&m).Save(&m).Error
}

/*
*
* 右上角
*
 */
func AllInternalNotifiesHeader() []model.MInternalNotify {
	m := []model.MInternalNotify{}
	interdb.DB().Table("m_internal_notifies").Where("status=1").Limit(6).Find(&m)
	return m
}

/*
*
* 所有列表
*
 */
func AllInternalNotifies() []model.MInternalNotify {
	m := []model.MInternalNotify{}
	interdb.DB().Table("m_internal_notifies").Where("status=1").Limit(100).Find(&m)
	return m
}

/*
*
* 清空表
*
 */
func ClearInternalNotifies() error {
	return interdb.DB().Exec("DELETE FROM m_internal_notifies;VACUUM;").Error
}

/*
*
* 点击已读
*
 */
func ReadInternalNotifies(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MInternalNotify{}).Error
}
