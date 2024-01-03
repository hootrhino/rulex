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

/*
*
* 物模型
*
 */
type MIotSchema struct {
	RulexModel
	UUID        string `gorm:"not null"`
	Name        string `gorm:"not null"` // 名称
	Description string // 额外信息
}

/*
*
* 属性
*
 */
type MIotProperty struct {
	RulexModel
	SchemaId    string `gorm:"not null"`
	UUID        string `gorm:"not null"`
	Label       string `gorm:"not null"` // UI显示的那个文本
	Name        string `gorm:"not null"` // 变量关联名
	Type        string `gorm:"not null"` // 类型, 只能是上面几种
	Rw          string `gorm:"not null"` // R读 W写 RW读写
	Unit        string `gorm:"not null"` // 单位 例如：摄氏度、米、牛等等
	Rule        string `gorm:"not null"` // 规则,IoTPropertyRule
	Description string // 额外信息
}
