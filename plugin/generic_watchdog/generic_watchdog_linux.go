// Copyright (C) 2023 wwhai
//
// This program is free genericWatchDogare: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free genericWatchDogare Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package genericwatchdog

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* C语言驱动相关
*
 */
const (
	WATCHDOG  = "/dev/watchdog"
	WATCHDOG0 = "/dev/watchdog0"
)
const (
	WATCHDOG_IOCTL_BASE = 'W'
	WDIOC_GETSUPPORT    = 2150127360
	WDIOC_SETTIMEOUT    = 3221509894
	WDIOC_GETTIMEOUT    = 2147768071
	WDIOS_DISABLECARD   = 1
	WDIOS_ENABLECARD    = 2
	WDIOC_SETOPTIONS    = 2147768068
	WDIOC_KEEPALIVE     = 2147768069
)

type watchdogInfo struct {
	options          uint32
	firmware_version uint32
	identity         [32]byte
}

func (w watchdogInfo) ToString() string {
	return fmt.Sprintf("Options: %d\nFirmware Version: %d\nIdentity: %s\n",
		w.options, w.firmware_version, string(w.identity[:]))
}

/*
*
* 软件看门狗
*
 */
type genericWatchDog struct {
	uuid string
}

func NewGenericWatchDog() *genericWatchDog {
	return &genericWatchDog{
		uuid: "SOFT_WATCHDOG",
	}
}

func (dog *genericWatchDog) Init(config *ini.Section) error {
	info, err := getWdogInfo()
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	glogger.GLogger.Info(info.ToString())
	return nil
}

func (dog *genericWatchDog) Start(typex.RuleX) error {
	go func() {
		defer stopWatchdog()
		for {
			select {
			case <-context.Background().Done():
				return
			default:
				feedWatchdog()
				time.Sleep(9 * time.Second)
			}
		}
	}()
	return nil
}
func (dog *genericWatchDog) Stop() error {
	return stopWatchdog()
}

func (hh *genericWatchDog) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "Linux Soft WatchDog",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *genericWatchDog) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

/*
*
* 停止
*
 */
func stopWatchdog() error {
	watchdogFile, err := os.OpenFile(WATCHDOG, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer watchdogFile.Close()

	_, err = fmt.Fprint(watchdogFile, "V")
	if err != nil {
		return err
	}

	return nil
}

/*
*
* 喂狗
*
 */
func feedWatchdog() error {
	// 打开 watchdog 设备文件以进行写入
	watchdogFile, err := os.OpenFile(WATCHDOG, os.O_WRONLY, 0)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	defer watchdogFile.Close()

	_, err = fmt.Fprint(watchdogFile, "W")
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

/*
*
* 获取看门狗的参数
*
 */
func getWdogInfo() (watchdogInfo, error) {
	wi := watchdogInfo{}
	fd, err := syscall.Open(WATCHDOG, syscall.O_RDWR, 0)
	if err != nil {
		glogger.GLogger.Error(err)
		return wi, err
	}
	defer syscall.Close(fd)
	if err := ioctl(fd, WDIOC_GETSUPPORT, uintptr(unsafe.Pointer(&wi))); err != nil {
		glogger.GLogger.Error(err)
		return wi, err

	}
	return wi, err

}
func ioctl(fd int, request, arg uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), request, arg)
	if errno != 0 {
		return os.NewSyscallError(fmt.Sprintf("ioctl error:%v,%v,%v", fd, request, arg), errno)
	}
	return nil
}
