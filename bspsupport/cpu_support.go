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

package archsupport

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

/*
*
* Windows
*
 */
func GetWindowsCPUName() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "name")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get CPU name: %w", err)
	}
	// Name
	// Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
	cpuNames := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(cpuNames) == 2 {
		return strings.TrimLeft(cpuNames[1], " "), nil
	}
	return "", fmt.Errorf("no CPU names found")
}

/*
*
* Linux
*
 */
func GetLinuxCPUName() (string, error) {
	cmd := exec.Command("bash", "-c", "lscpu | grep 'Model name'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get CPU name: %w", err)
	}
	//Model name:                         Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
	cpuNames := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(cpuNames) > 0 && (len(cpuNames[0]) > 10) {
		return strings.TrimLeft(cpuNames[0][11:], " "), nil
	}
	return "", fmt.Errorf("no CPU names found")
}
