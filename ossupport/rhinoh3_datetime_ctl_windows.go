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
	"strings"
	"time"

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
		// 20231226170151.500000+480
		wmicTime := parseWinWmicTime(result[0].LastBootUpTime)
		Seconds := wmicTime.Abs().Seconds()
		hour := int(Seconds / 3600)
		minute := int(Seconds/60) % 60
		second := int(Seconds) % 60
		return fmt.Sprintf("%d Hours %02d Minutes %02d Seconds", hour, minute, second), nil
	}

	return "0 Year 0 Month 0 Days 0 Hours 0 Minutes 0 Seconds",
		fmt.Errorf("failed to retrieve system uptime")
}

/*
*
* 解析ISO时间戳
*
 */
func parseWinWmicTime(timestamp string) time.Duration {
	parts1 := strings.Split(timestamp, "+")
	if len(parts1) != 2 {
		return 0
	}
	//2023 12 26 17 01 51.500000+480
	mainTimestamp, err := time.Parse("20060102150405", parts1[0])
	if err != nil {
		return 0
	}
	return time.Until(mainTimestamp)
}
