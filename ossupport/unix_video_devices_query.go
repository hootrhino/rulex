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
	"os"
	"path/filepath"
	"strings"
)

/*
*
* Unix摄像头
*
 */
type UnixCameraDevice struct {
	Node     string `json:"node"`
	DeviceID string `json:"deviceId"`
	Name     string `json:"name"`
}

/*
*
* Unix 获取视频设备
*
 */
func GetUnixVideos() ([]UnixCameraDevice, error) {
	Devices := []UnixCameraDevice{}
	files, err := os.ReadDir("/dev")
	if err != nil {
		return Devices, err
	}
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "video") {
			Name := filepath.Join("/dev", name)
			Devices = append(Devices, UnixCameraDevice{
				Node:     "Local",
				DeviceID: Name,
				Name:     Name,
			})
		}
	}
	return Devices, nil
}
