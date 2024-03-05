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

package ossupport

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
*
* Windows 获取视频设备
*
 */
type CameraDevice struct {
	Node     string
	DeviceID string
	Name     string
}

/*
*
PS C:\Users>  WMIC Path Win32_PnPEntity WHERE "Caption LIKE '%CAMERA%'" GET DeviceID,Name /FORMAT:CSV
------
Node,DeviceID,Name
DESKTOP-EMD3M3C,USB\VID_1908&amp;PID_2311&amp;MI_00\6&amp;1666943&amp;0&amp;0000,USB2.0 PC CAMERA
*
*/
func GetWindowsVideos() ([]CameraDevice, error) {
	cmd := exec.Command("powershell.exe", "-Command",
		`WMIC Path Win32_PnPEntity WHERE "Caption LIKE '%CAMERA%'" GET DeviceID,Name /FORMAT:CSV`)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := []string{}
	for _, v := range strings.Split(string(output), "\n") {
		if v != "" && v != "\r\r" {
			lines = append(lines, v)
		}
	}
	videos := []CameraDevice{}
	if len(lines) < 2 {
		return videos, fmt.Errorf("camera device not found")
	}
	for _, line := range lines[1:] {
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) < 3 {
			return videos, fmt.Errorf("parse output line error:%s", line)
		}
		videos = append(videos, CameraDevice{
			Node:     fields[0],
			DeviceID: fields[1],
			Name:     fields[2],
		})
	}
	return videos, nil
}
