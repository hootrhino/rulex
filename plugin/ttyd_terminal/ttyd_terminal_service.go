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
		// 允许忙碌中停止
		if arg.Name == "stop" {
			tty.stop()
			return typex.ServiceResult{Out: "Terminal Stop Success"}
		}
		// 禁止多开
		return typex.ServiceResult{Out: "Terminal already running now"}
	}

	if arg.Name == "stop" {
		tty.stop()
		return typex.ServiceResult{Out: "Terminal Stop Success"}

	}
	if arg.Name == "start" {
		tty.busying = true
		ctx, cancel := context.WithCancel(typex.GCTX)
		tty.ctx = ctx
		tty.cancel = cancel

		tty.ttydCmd = exec.CommandContext(typex.GCTX,
			"ttyd", "-W",
			"-p", fmt.Sprintf("%d", tty.mainConfig.ListenPort),
			"-c", "rhino:hoot", // username"password
			"-O",      // Not allow Cros
			"-m", "1", // Max allow 1 client
			"-o", "-6", "bash")
		inOut := wsInOut{}
		tty.ttydCmd.Stdout = inOut
		tty.ttydCmd.Stderr = inOut
		if err1 := tty.ttydCmd.Start(); err1 != nil {
			glogger.GLogger.Infof("cmd.Start error: %v", err1)
			return typex.ServiceResult{Out: err1.Error()}
		}
		// 如果5分钟没人操作就结束
		go func(tty *WebTTYPlugin) {
			defer func() {
				tty.stop()
			}()
			glogger.GLogger.Info("ttyd started successfully on port:", tty.mainConfig.ListenPort)
			tty.ttydCmd.Process.Wait()
			glogger.GLogger.Info("ttyd stopped with state:", tty.ttydCmd.ProcessState.String())
		}(tty)
		return typex.ServiceResult{Out: "Terminal Start Success"}
	}
	return typex.ServiceResult{Out: "Unknown service name:" + arg.Name}
}

func (tty *WebTTYPlugin) stop() error {
	if tty.cancel != nil {
		tty.cancel()
	}
	if tty.ttydCmd.ProcessState != nil {
		tty.ttydCmd.Process.Kill()
		tty.ttydCmd.Process.Signal(os.Kill)
	}
	tty.busying = false
	return nil
}

type wsInOut struct {
}

func NewWSStdInOut() wsInOut {
	return wsInOut{}
}

func (hk wsInOut) Write(p []byte) (n int, err error) {
	glogger.Logrus.Info(string(p))
	return len(p), nil
}
func (hk wsInOut) Read(p []byte) (n int, err error) {
	return len(p), nil
}
