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
	"os"

	"github.com/hootrhino/rulex/component/interdb"
)

/*
*
* 清理垃圾进程包
*
 */
func CleanGoodsUpload() error {
	sql1 := `SELECT local_path FROM m_goods;`
	m_goods_local_paths := []string{}
	err1 := interdb.DB().Raw(sql1).Find(&m_goods_local_paths).Error
	if err1 != nil {
		return err1
	}
	for _, m_goods_local_path := range m_goods_local_paths {
		if err := os.Remove(m_goods_local_path); err != nil {
			return err
		}
	}
	return nil
}

/*
*
* 清理缩略图
*
 */
func CleanThumbnailUpload() error {
	sql1 := `SELECT thumbnail FROM m_visuals;`
	thumbnail_local_paths := []string{}
	err1 := interdb.DB().Raw(sql1).Find(&thumbnail_local_paths).Error
	if err1 != nil {
		return err1
	}
	for _, thumbnail_local_path := range thumbnail_local_paths {
		if err := os.Remove(thumbnail_local_path); err != nil {
			return err
		}
	}
	return nil
}
