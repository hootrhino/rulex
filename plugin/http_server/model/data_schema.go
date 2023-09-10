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
* Type: SIMPLE_LINE(简单一线),COMPLEX_LINE(复杂多线)
*
 */
type MDataSchema struct {
	RulexModel
	UUID   string `gorm:"not null"`
	Name   string `gorm:"not null"` // 名称
	Type   string `gorm:"not null"` // 类型, LINE,PINE,BAR,TXT
	Schema string `gorm:"not null"` // 数据规范
}
