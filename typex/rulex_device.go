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

package typex

import "github.com/hootrhino/rulex/utils"

// 设备元数据, 本质是保存在配置里面的数据的一个内存映射实例
type Device struct {
	UUID        string                 `json:"uuid"`        // UUID
	Name        string                 `json:"name"`        // 设备名称，例如：灯光开关
	Type        DeviceType             `json:"type"`        // 类型,一般是设备-型号，比如 ARDUINO-R3
	AutoRestart bool                   `json:"autoRestart"` // 是否允许挂了的时候重启
	Description string                 `json:"description"` // 设备描述信息
	BindRules   map[string]Rule        `json:"-"`           // 与之关联的规则
	State       DeviceState            `json:"state"`       // 状态
	Config      map[string]interface{} `json:"config"`      // 配置
	Device      XDevice                `json:"-"`           // 实体设备
}

func NewDevice(t DeviceType,
	name string,
	description string,
	config map[string]interface{}) *Device {
	return &Device{
		UUID:        utils.DeviceUuid(),
		Name:        name,
		Type:        t,
		State:       3,
		Description: description,
		AutoRestart: true, // 0.5以前默认自动重启
		Config:      config,
		BindRules:   map[string]Rule{},
	}
}
