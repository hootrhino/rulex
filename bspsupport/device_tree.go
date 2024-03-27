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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package archsupport

import "os"

type DeviceNode struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}
type DeviceTree struct {
	Network    []DeviceNode `json:"network"`
	Wlan       []DeviceNode `json:"wlan"`
	MNet4g     []DeviceNode `json:"net4g"`
	MNet5g     []DeviceNode `json:"net5g"`
	SoftRouter []DeviceNode `json:"soft_router"`
}

/*
*
* 获取设备树
*
 */
func GetDeviceCtrlTree() DeviceTree {
	env := os.Getenv("ARCHSUPPORT")
	if env == "EEKITH3" {
		return DeviceTree{
			Network: []DeviceNode{
				{"eth0", "ethernet", 1},
				{"eth1", "ethernet", 1},
			},
			Wlan: []DeviceNode{
				{"wlan0", "wlan", 1},
			},
			MNet4g: []DeviceNode{
				{"usb0", "4g", 1},
			},
			MNet5g: []DeviceNode{},
			SoftRouter: []DeviceNode{
				{"eth0", "ethernet", 1},
				{"eth1", "ethernet", 1},
			},
		}
	}
	if env == "RPI4B" {
		return DeviceTree{
			Network: []DeviceNode{
				{"eth0", "ethernet", 1},
			},
			Wlan: []DeviceNode{
				{"wlan0", "wlan", 1},
			},
			MNet4g:     []DeviceNode{},
			MNet5g:     []DeviceNode{},
			SoftRouter: []DeviceNode{},
		}
	}
	if env == "EN6400" {
		return DeviceTree{
			Network: []DeviceNode{
				{"eth0", "ethernet", 1},
			},
			Wlan:       []DeviceNode{},
			MNet4g:     []DeviceNode{},
			MNet5g:     []DeviceNode{},
			SoftRouter: []DeviceNode{},
		}
	}
	return DeviceTree{
		Network: []DeviceNode{
			{"eth0", "ethernet", 1},
		},
		Wlan:       []DeviceNode{},
		MNet4g:     []DeviceNode{},
		MNet5g:     []DeviceNode{},
		SoftRouter: []DeviceNode{},
	}
}
