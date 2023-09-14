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
	"fmt"
	"os/exec"
	"strings"
)

func GetBSPNetIfaces() ([]string, error) {
	// 执行命令
	cmd := exec.Command("sh", "-c", "ip -o link show | awk -F': ' '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error:%s", string(output))
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	netIfaces := make([]string, 0, len(lines))

	for _, line := range lines {
		if line != "" {
			netIfaces = append(netIfaces, line)
		}
	}
	return netIfaces, nil
}
