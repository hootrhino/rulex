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

package ossupport

import (
	"fmt"
	"github.com/hootrhino/wmi"
)

/*
*
* 获取开机时间
*
 */
type Win32OperatingSystem struct {
	LastBootUpTime string
}

func GetUptime() (string, error) {
	var result []Win32OperatingSystem
	query := "SELECT LastBootUpTime FROM Win32_OperatingSystem"

	err := wmi.Query(query, &result)
	if err != nil {
		return "", err
	}

	if len(result) > 0 {
		return result[0].LastBootUpTime, nil
	}

	return "0:0:0", fmt.Errorf("Failed to retrieve system uptime")
}
