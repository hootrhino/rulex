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
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// GetMemPercent 获取Linux内存使用百分比
func GetMemPercent() (float64, error) {
	// 打开 /proc/meminfo 文件
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("Open /proc/meminfo error: %v", err)
	}
	defer file.Close()

	// 初始化变量用于存储内存信息
	var totalMem, freeMem int64

	// 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 3 {
			// 提取字段名、值和单位
			fieldName := fields[0]
			fieldValue := fields[1]

			// 将值转换为整数
			value, err := parseMemInfoValue(fieldValue)
			if err != nil {
				return 0, fmt.Errorf("parse MemInfo Value error: %v", err)
			}

			// 根据字段名更新内存信息
			switch fieldName {
			case "MemTotal:":
				totalMem = value
			case "MemFree:":
				freeMem = value
			}
		}
	}

	// 计算已使用内存
	usedMem := totalMem - freeMem

	// 计算内存使用百分比
	memUsagePercent := float64(usedMem) / float64(totalMem) * 100

	return math.Round(memUsagePercent*100) / 100, nil
}

// parseMemInfoValue 解析 /proc/meminfo 文件中的内存值（以 KB 为单位）
func parseMemInfoValue(valueStr string) (int64, error) {
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return value * 1024, nil
}
