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

// modbus数据点位表
type MModbusDataPoint struct {
	RulexModel
	UUID       string  `gorm:"not null" json:"uuid"`
	DeviceUuid string  `gorm:"not null" json:"device_uuid"`
	Tag        string  `gorm:"not null" json:"tag"`
	Alias      string  `gorm:"not null" json:"alias"`
	Function   *int    `gorm:"not null" json:"function"`
	SlaverId   *byte   `gorm:"not null" json:"slaver_id"`
	Address    *uint16 `gorm:"not null" json:"address"`
	Frequency  *int64  `gorm:"not null" json:"frequency"`
	Quantity   *uint16 `gorm:"not null" json:"quantity"`
}
