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

package archsupport

import (
	"encoding/json"
)

/*
*
* RhinoPi 硬件接口相关管理
* 警告：此处会随着硬件不同而不兼容，要移植的时候记得统一一下目标硬件的端口
*
 */
var __RhinoH3HwInterfaces map[string]RhinoH3HwInterface

func init() {
	__RhinoH3HwInterfaces = map[string]RhinoH3HwInterface{}
}

/*
*
* 这里记录一些H3网关的硬件接口信息,同时记录串口是否被占用等
*
 */
type UartConfig struct {
	Timeout  int
	Uart     string
	BaudRate int
	DataBits int
	Parity   string
	StopBits int
}
type RhinoH3HwInterface struct {
	Name     string      // 接口名称
	Type     string      // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias    string      // 别名
	Busy     bool        // 是否被占
	OccupyBy string      // 被谁占用了
	Config   interface{} // 配置, 串口配置、或者网卡、USB等
}

func (v RhinoH3HwInterface) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}
