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

package hwportmanager

/*
*
* 硬件设备上的端口
*
 */
type HardwarePort struct {
	UUID        int64  `json:"uuid"`        // ID
	Name        string `json:"name"`        // 名称
	Alias       string `json:"alias"`       // 别名
	Type        string `json:"type"`        // 类型,主要有 USER、SYSTEM两种
	Path        string `json:"path"`        // 路径
	Description string `json:"description"` // 额外信息
}

/*
*
* 统一管理不同型号设备的硬件端口
*
 */
type HardwarePortClass struct {
	Audio HardwarePort `json:"audio"`
	Can   HardwarePort `json:"can"`
	Di    HardwarePort `json:"di"`
	Do    HardwarePort `json:"do"`
	Hdmi  HardwarePort `json:"hdmi"`
	Ir    HardwarePort `json:"ir"`
	Relay HardwarePort `json:"relay"`
	Rs232 HardwarePort `json:"rs232"`
	Rs485 HardwarePort `json:"rs485"`
	Uart  HardwarePort `json:"uart"`
	Video HardwarePort `json:"video"`
}
