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

package test

import (
	"testing"

	archsupport "github.com/hootrhino/rulex/bspsupport"
)

// go test -timeout 30s -run ^Test_Nvidia_SMI1 github.com/hootrhino/rulex/test -v -count=1

func Test_Nvidia_SMI1(t *testing.T) {
	t.Log(archsupport.GetGpuInfoWithNvidiaSmi())
}

// go test -timeout 30s -run ^Test_GetCpu_win11 github.com/hootrhino/rulex/test -v -count=1
func Test_GetCpu_win11(t *testing.T) {
	t.Log(archsupport.GetWindowsCPUName())
}

// go test -timeout 30s -run ^Test_GetCpu_linux github.com/hootrhino/rulex/test -v -count=1
func Test_GetCpu_linux(t *testing.T) {
	t.Log(archsupport.GetLinuxCPUName())
}
