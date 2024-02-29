// Copyright (C) 2024 wwhai
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

// 华中数控
type MHnc8DataPoint struct {
	RulexModel  `json:"-"`
	UUID        string `gorm:"not null"`
	DeviceUuid  string `gorm:"not null"`
	Name        string `gorm:"not null"` // 点位名称
	Alias       string `gorm:"not null"` // 别名
	ApiFunction string `gorm:"not null"` // API路径
	Group       int    `gorm:"not null"` // 分组采集
	Address     string `gorm:"not null"` // 地址
}
