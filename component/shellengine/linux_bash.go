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

package shellengine

import (
	"context"
	"log"
	"os/exec"

	"github.com/hootrhino/rulex/typex"
)

type LinuxBashShell struct {
	rulex typex.RuleX
}

func InitLinuxBashShell(rulex typex.RuleX) *LinuxBashShell {
	return &LinuxBashShell{rulex: rulex}

}
func (lsh *LinuxBashShell) JustRun(ctx context.Context, cmd string) error {
	Cmd := exec.CommandContext(ctx, "sh", "-c", cmd)
	err := Cmd.Start()
	if err != nil {
		O, _ := Cmd.Output()
		log.Println("[Linux Bash Shell] error", err, ", Output:", string(O))
		return err
	}
	return nil
}
func (lsh *LinuxBashShell) JustRunDetach(ctx context.Context, cmd string) error {
	Cmd := exec.CommandContext(ctx, "sh", "-c", cmd)
	if Cmd != nil {
		Cmd.Process.Release()
	}
	O, err := Cmd.CombinedOutput()
	if err != nil {
		log.Println("[Linux Bash Shell] error", err, ", Output:", string(O))
		return err
	}

	return nil

}
func (lsh *LinuxBashShell) RunAndWaitResult(ctx context.Context, cmd string) ([]byte, error) {
	Cmd := exec.CommandContext(ctx, "sh", "-c", cmd)
	O, err := Cmd.CombinedOutput()
	if err != nil {
		log.Println("[Linux Bash Shell] error", err, ", Output:", string(O))
		return nil, err
	}
	return O, nil
}
func (lsh *LinuxBashShell) RunWithInteractive(ctx context.Context, cmd string,
	in, out chan []byte) error {
	return nil

}
func (lsh *LinuxBashShell) ReadFD(fd int, p []byte) (n int, err error) {
	return 0, nil

}
func (lsh *LinuxBashShell) IOctl(fd int, request, arg uintptr) error {
	return nil
}
