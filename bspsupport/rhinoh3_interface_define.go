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
	"fmt"
)

var __RhinoH3HwInterfaces map[string]*RhinoH3HwInterface

func init() {
	_InterfaceInit()

}

/*
*
* 这里记录一些H3网关的硬件接口信息,同时记录串口是否被占用等
*
 */
type RhinoH3HwInterface struct {
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Busy     bool   `json:"busy"`     // 是否被占
	Relation string `json:"relation"` // 被谁占用了
}

func (v RhinoH3HwInterface) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func UARTBusy(name string) {
	Port, ok := __RhinoH3HwInterfaces[name]
	if ok {
		Port.Busy = true
	}
}
func UARTFree(name string) {
	Port, ok := __RhinoH3HwInterfaces[name]
	if ok {
		Port.Busy = false
	}
}

/*
*
* 展示自己的信息
*
 */
func (v RhinoH3HwInterface) BusyingInfo() string {
	return fmt.Sprintf("Port [%s(%s)] is busying now, may be occupy by %s",
		v.Name, v.Alias, v.Relation)
}
func _InterfaceInit() {
	__RhinoH3HwInterfaces = map[string]*RhinoH3HwInterface{
		"/dev/ttyS1": {
			Name:     "/dev/ttyS1",
			Alias:    "RS4851(A1B1)",
			Busy:     false,
			Relation: "",
		},
		"/dev/ttyS2": {
			Name:     "/dev/ttyS2",
			Alias:    "RS4851(A2B2)",
			Busy:     false,
			Relation: "",
		},
	}

}
func AllUartInterfaces() map[string]RhinoH3HwInterface {
	r := map[string]RhinoH3HwInterface{}
	for k, v := range __RhinoH3HwInterfaces {
		r[k] = *v
	}
	return r
}
func GetUart(name string) RhinoH3HwInterface {
	Port, ok := __RhinoH3HwInterfaces[name]
	if ok {
		return *Port
	}
	return *Port
}
