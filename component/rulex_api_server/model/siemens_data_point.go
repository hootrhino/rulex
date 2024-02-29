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

package model

// 西门子数据点位表
type MSiemensDataPoint struct {
	RulexModel
	UUID           string   `gorm:"not null"`
	DeviceUuid     string   `gorm:"not null"` // 所属设备
	SiemensAddress string   `gorm:"not null"` // 西门子的地址字符串
	Tag            string   `gorm:"not null"` // Tag
	Alias          string   `gorm:"not null"` // 别名
	DataBlockType  string   `gorm:"not null"` // 类型 INT UINT....
	DataBlockOrder string   `gorm:"not null"` // 字节序
	Weight         *float64 `gorm:"not null"` // 权重
	Frequency      *int64   `gorm:"not null"` // 采集频率
}
