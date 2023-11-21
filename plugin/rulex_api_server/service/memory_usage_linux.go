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
	"os"
	"strings"
	"math"
)

func GetMemPercent() (float64, error) {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0.0, err
	}

	meminfo := string(content)

	var memTotal, memAvailable float64
	lines := strings.Split(meminfo, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			fmt.Sscan(fields[1], &memTotal)
		case "MemAvailable:":
			fmt.Sscan(fields[1], &memAvailable)
		}
	}
	memPercent := 100.0 * (1.0 - memAvailable/memTotal)
	return math.Round(memPercent*100) / 100, nil
}
