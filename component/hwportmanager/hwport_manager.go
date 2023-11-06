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

/*
*
* RhinoPi 硬件接口相关管理
* 警告：此处会随着硬件不同而不兼容，要移植的时候记得统一一下目标硬件的端口
*
 */
package hwportmanager

import (
	"encoding/json"
	"fmt"

	"github.com/hootrhino/rulex/typex"
)

var __HwPortsManager *HwPortsManager

type HwPortsManager struct {
	Interfaces map[string]RhinoH3HwPort
	rulex      typex.RuleX
}

func InitHwPortsManager(rulex typex.RuleX) *HwPortsManager {
	__HwPortsManager = &HwPortsManager{
		Interfaces: map[string]RhinoH3HwPort{},
		rulex:      rulex,
	}
	return __HwPortsManager
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
type HwPortOccupy struct {
	UUID string // UUID
	Type string // DEVICE, Other......
}
type RhinoH3HwPort struct {
	UUID        string       // 接口名称
	Name        string       // 接口名称
	Alias       string       // 别名
	Busy        bool         // 运行时数据，是否被占
	OccupyBy    HwPortOccupy // 运行时数据，被谁占用了 UUID
	Type        string       // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Description string       // 额外备注
	Config      interface{}  // 配置, 串口配置、或者网卡、USB等
}

func (v RhinoH3HwPort) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

/*
*
* 加载配置到运行时, 需要刷新与配置相关的所有设备
*
 */
func SetHwPort(Port RhinoH3HwPort) {
	__HwPortsManager.Interfaces[Port.Name] = Port
	refreshHwPort(Port.Name)
}
func RefreshPort(Port RhinoH3HwPort) {
	__HwPortsManager.Interfaces[Port.Name] = Port
	refreshHwPort(Port.Name)
}

/*
*
* 刷新所有关联的设备, 也就是 OccupyBy=UUID 需要重载
*
 */
func refreshHwPort(name string) {
	Port, ok := __HwPortsManager.Interfaces[name]
	if !ok {
		return
	}
	if Port.Busy {
		if Port.OccupyBy.Type == "DEVICE" {
			UUID := Port.OccupyBy.UUID
			if Device := __HwPortsManager.rulex.GetDevice(UUID); Device != nil {
				// 拉闸 DEV_DOWN 以后就重启了, 然后就会拉取最新的配置
				Device.Device.SetState(typex.DEV_DOWN)
			}
		}
	}

}

/*
*
* 获取一个运行时状态
*
 */
func GetHwPort(name string) (RhinoH3HwPort, error) {
	if Port, ok := __HwPortsManager.Interfaces[name]; ok {
		return Port, nil
	}
	return RhinoH3HwPort{}, fmt.Errorf("interface not exists:%s", name)
}

/*
*
* 所有的接口
*
 */
func AllHwPort() []RhinoH3HwPort {
	result := []RhinoH3HwPort{}
	for _, v := range __HwPortsManager.Interfaces {
		result = append(result, v)
	}
	return result
}

/*
*
* 忙碌
*
 */
func SetInterfaceBusy(name string, OccupyBy HwPortOccupy) {
	if Port, ok := __HwPortsManager.Interfaces[name]; ok {
		Port.Busy = true
		Port.OccupyBy = OccupyBy
		__HwPortsManager.Interfaces[name] = Port
	}
}

/*
*
* 释放
*
 */
func FreeInterfaceBusy(name string) {
	if Port, ok := __HwPortsManager.Interfaces[name]; ok {
		Port.Busy = false
		Port.OccupyBy = HwPortOccupy{
			"", "",
		}
		__HwPortsManager.Interfaces[name] = Port
	}
}
