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

import "gopkg.in/square/go-jose.v2/json"

type MHwPort struct {
	RulexModel
	UUID        string `gorm:"not null"`
	Name        string `gorm:"not null"` // 接口名称
	Type        string `gorm:"not null"` // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias       string `gorm:"not null"` // 别名
	Description string `gorm:"not null"` // 额外备注
	Config      string `gorm:"not null"` // 配置, 串口配置、或者网卡、USB等
}

func (md MHwPort) GetConfig() map[string]interface{} {
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}
