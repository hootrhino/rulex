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
	"math"
	"os/exec"
	"strings"
)

// GetMemPercent 获取Windows系统内存使用百分比
// wmic OS get FreePhysicalMemory,TotalVisibleMemorySize /value
// FreePhysicalMemory=14309944
// TotalVisibleMemorySize=25087752
func GetMemPercent() (float64, error) {
	cmd := exec.Command("wmic", "OS", "get", "FreePhysicalMemory,TotalVisibleMemorySize", "/value")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("执行WMIC命令时出错: %v", string(output))
	}

	result := string(output)

	// 解析输出以获取可用物理内存和总可见内存的值
	lines := strings.Split(result, "\n")
	var freeMemory, totalMemory uint64
	for _, line := range lines {
		fields := strings.Split(line, "=")
		if len(fields) == 2 {
			fieldName := strings.TrimSpace(fields[0])
			fieldValue := strings.TrimSpace(fields[1])
			switch fieldName {
			case "FreePhysicalMemory":
				fmt.Sscanf(fieldValue, "%d", &freeMemory)
			case "TotalVisibleMemorySize":
				fmt.Sscanf(fieldValue, "%d", &totalMemory)
			}
		}
	}

	if totalMemory > 0 {
		// 计算内存使用百分比
		memUsagePercent := float64(totalMemory-freeMemory) / float64(totalMemory) * 100
		return math.Round(memUsagePercent*100) / 100, nil
	}

	return 0, fmt.Errorf("Can not readable")
}
