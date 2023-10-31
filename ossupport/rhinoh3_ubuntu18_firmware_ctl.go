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
		cmd := exec.Command("sudo", "systemctl", "daemon-reload")
		cmd.SysProcAttr = NewSysProcAttr()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s,%s", err, string(out))
		}
	}
	{
		cmd := exec.Command("sudo", "systemctl", "restart", "rulex")
		cmd.SysProcAttr = NewSysProcAttr()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s,%s", err, string(out))
		}

	}
	return nil
}

/*
*
* 启用升级进程
*
 */
func StartUpgradeProcess(path string, args []string) {
	log.Printf("Start Upgrade Process Pid=%d, Gid=%d", os.Getpid(), os.Getegid())
	defer log.Println("Start Upgrade Process Exited")
	outputFile, err1 := os.Create("StartUpgradeProcess.txt")
	if err1 != nil {
		log.Println("Create Upgrade log error:", err1)
		return
	}
	defer outputFile.Close()
	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile
	cmd.Process.Pid = -1 // 用来分离进程用,简直天坑参数！！！
	err := cmd.Start()
	if err != nil {
		log.Println("Start Upgrade Process Failed:", err)
		return
	}
	log.Println("Start Upgrade Process:", cmd.String())
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
