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
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
*
* Stop RULEX
*
 */
func StopRulex() error {
	cmd := exec.Command("service", "rulex", "stop")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s,%s", err, string(out))
	}
	return nil
}

/*
*
* 重启
*
 */
func Restart() error {
	{
		cmd := exec.Command("systemctl", "daemon-reload")
		cmd.SysProcAttr = NewSysProcAttr()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s,%s", err, string(out))
		}
		log.Println("[Prepare Stage] systemctl daemon-reload:", string(out))

	}
	{
		cmd := exec.Command("service", "rulex", "restart")
		cmd.SysProcAttr = NewSysProcAttr()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s,%s", err, string(out))
		}
		log.Println("[Prepare Stage] service start:", string(out))

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
	file, err1 := os.Create("/usr/local/local-recover-log.txt")
	if err1 != nil {
		return
	}
	defer file.Close()
	cmd := exec.Command("/usr/local/rulex", "recover", "-recover=true")
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = []string{}
	cmd.Stdout = file
	cmd.Stderr = file
	err := cmd.Start()
	log.Printf("Start Recover Process Pid=%d, Cmd:%s", cmd.Process.Pid, cmd.String())
	if err != nil {
		log.Println("Start Recover Process Failed:", err)
		return
	}
	log.Println("Old Process Exited:", os.Getpid())
	os.Exit(0)
}

/*
*
* 启用升级进程
*
 */
func StartUpgradeProcess(path string, args []string) {
	file, err1 := os.Create("/usr/local/local-upgrade-log.txt")
	if err1 != nil {
		return
	}
	defer file.Close()

	cmd := exec.Command("bash", "-c", path+" "+strings.Join(args, " "))
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = []string{}
	cmd.Stdout = file
	cmd.Stderr = file
	log.Printf("Start Upgrade Process Pid=%d, Gid=%d", os.Getpid(), os.Getegid())
	if cmd.Process != nil {
		cmd.Process.Release() // 用来分离进程用,简直天坑参数！！！
	}
	err := cmd.Start()
	if cmd.Process != nil {
		cmd.Process.Release() // 用来分离进程用,简直天坑参数！！！
	}
	log.Printf("Start Upgrade Process Pid=%d, Cmd:%s", cmd.Process.Pid, cmd.String())
	if err != nil {
		log.Println("Start Upgrade Process Failed:", err)
		return
	}
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
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
