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
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

/*
*
* kill -9
*
 */
func KillProcess(processID int) error {
	cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", processID))
	err := cmd.Run()
	return err
}

/*
*
* pgrep rulex -> 38506\n
*
 */
func GetProcessPID(processName string) (int, error) {
	cmd := exec.Command("pgrep", processName)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	pidString := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

/*
*
* 取最老的那个进程Id
*
 */
func GetEarliestProcessPID(processName string) (int, error) {
	cmd := exec.Command("pgrep", processName)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	pidStrings := strings.Fields(string(output))
	var pids []int
	for _, pidString := range pidStrings {
		pid, err := strconv.Atoi(pidString)
		if err == nil {
			pids = append(pids, pid)
		}
	}
	if len(pids) == 0 {
		return 0, fmt.Errorf("No process found with name %s", processName)
	}
	sort.Ints(pids)
	earliestPID := pids[0]
	return earliestPID, nil
}
