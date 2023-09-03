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

package ttyd_terminal

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 服务调用接口
*
 */

func (tty *WebTTYPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	if tty.busying {
		if arg.Name == "stop" {
			if tty.cancel != nil {
				tty.cancel()
				tty.busying = false
				return typex.ServiceResult{Out: "Stop Success"}
			}
		}
		return typex.ServiceResult{Out: "Modbus Scanner Busing now"}
	}

	if arg.Name == "stop" {
		if tty.cancel != nil {
			tty.cancel()
		}
		if tty.ttydCmd != nil {
			if tty.ttydCmd.ProcessState != nil {
				tty.ttydCmd.Process.Kill()
				tty.ttydCmd.Process.Signal(os.Kill)
			}
		}
		tty.busying = false
	}
	if arg.Name == "start" {
		tty.busying = true
		ctx, cancel := context.WithCancel(typex.GCTX)
		tty.ctx = ctx
		tty.cancel = cancel
		tty.ttydCmd = exec.CommandContext(typex.GCTX,
			"ttyd", "-W", "-p", fmt.Sprintf("%d", tty.mainConfig.ListenPort),
			"-o", "-6", "bash")
		if err1 := tty.ttydCmd.Start(); err1 != nil {
			glogger.GLogger.Infof("cmd.Start error: %v", err1)
			return typex.ServiceResult{}
		}
		go func(tty *WebTTYPlugin) {
			defer func() {
				tty.busying = false
			}()
			glogger.GLogger.Info("ttyd started successfully on port:", tty.mainConfig.ListenPort)
			tty.ttydCmd.Process.Wait()
			glogger.GLogger.Info("ttyd stopped")
		}(tty)
	}
	return typex.ServiceResult{}
}
