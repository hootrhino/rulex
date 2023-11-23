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
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

/*
*
* Stop RULEX
*
 */
func StopRulex() error {
	pid, err1 := GetEarliestProcessPID("rulex")
	if err1 != nil {
		return err1
	}
	err2 := KillProcess(pid)
	if err2 != nil {
		return err2
	}
	return nil
}

/*
*
* 重启
*
 */
func RestartRulex() error {
	cmd := exec.Command("/etc/init.d/rulex.service", "restart")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s,%s", err, string(out))
	}
	return nil
}

/*
*
* 恢复上传的DB
1 停止RULEX
2 删除DB
3 复制DB过去
4 重启
- path: /usr/local/rulex, args: recover=true
*
*/
func StartRecoverProcess() {
	cmd := exec.Command("./rulex", "recover", "-recover=true")
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Println("Start Recover Process Failed:", err)
		return
	}
	log.Printf("Start Recover Process Pid=%d, Cmd:%s\n", cmd.Process.Pid, cmd.String())
	log.Println("Old Process Exited:", os.Getpid())
	os.Exit(0)
}

/*
*
* 启用升级进程
*
 */
func StartUpgradeProcess() {
	cmd := exec.Command("./rulex", "upgrade", "-oldpid=", fmt.Sprintf("%d", os.Getpid()))
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Println("Start Upgrade Process Failed:", err)
		return
	}
	log.Printf("Start Upgrade Process Pid=%d, Cmd:%s", cmd.Process.Pid, cmd.String())
	log.Println("Old Process Exited:", os.Getpid())
	os.Exit(0)
}

/*
*
* 直接重启Linux
*
 */
func Reboot() error {
	cmd := exec.Command("reboot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

/*
*
* 解压安装包
*
 */
func UnzipFirmware(zipFile, destDir string) error {
	cmd := exec.Command("unzip", "-o", zipFile, "-d", destDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unzip file: %s, %s", err.Error(), string(out))
	}
	return nil
}

/*
*
* 移动文件
*
 */
func MoveFile(sourcePath, destPath string) error {

	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}
