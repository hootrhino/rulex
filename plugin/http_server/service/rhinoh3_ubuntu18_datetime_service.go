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

package service

import (
	"fmt"
	"os/exec"
	"time"
)

/*
*
* NTP 用于启用NTP时间同步
*
 */

func UpdateTimeByNtp() error {
	err2 := setNtp(false)
	if err2 != nil {
		return err2
	}
	err1 := setNtp(true)
	if err1 != nil {
		return err1
	}
	return nil
}

/*
*
* 验证时间格式 YYYY-MM-DD HH:MM:SS
*
 */
func isValidTimeFormat(input string) bool {
	expectedFormat := "2006-01-02 15:04:05"
	_, err := time.Parse(expectedFormat, input)
	return err == nil
}

/*
*
* 获取当前系统时间
*
 */
func GetSystemTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*
*
*
设置时间，格式为 "YYYY-MM-DD HH:MM:SS"
*
*/
func SetSystemTime(newTime string) error {
	if !isValidTimeFormat(newTime) {
		return fmt.Errorf("Invalid time format:%s, must be 'YYYY-MM-DD HH:MM:SS'", newTime)
	}
	// newTime := "2023-08-10 15:30:00"
	cmd := exec.Command("date", "-s", newTime)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

/*
*
* v: true|false
*
 */
func setNtp(v bool) error {
	cmd := exec.Command("timedatectl", "set-ntp", func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}(v))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

/*
*
* 时区
*
 */
func GetTimeZone() string {
	z, _ := time.Now().Zone()
	return z

}

// SetTimeZone 设置系统时区
// timezone := "Asia/Shanghai"
func SetTimeZone(timezone string) error {
	cmd := exec.Command("timedatectl", "set-timezone", timezone)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error: %v\nOutput: %s", err, string(output))
	}
	return nil
}
