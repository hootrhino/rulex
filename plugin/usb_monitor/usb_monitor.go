package usbmonitor

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/i4de/rulex/typex"
	"golang.org/x/sys/unix"
	"gopkg.in/ini.v1"
)

/*
*
* USB 热插拔监控器, 方便观察USB状态, 本插件只支持Linux！！！
*
 */
type usbMonitor struct {
}

func NewUsbMonitor() typex.XPlugin {
	return &usbMonitor{}
}
func (usbm *usbMonitor) Init(_ *ini.Section) error {
	return nil

}

func (usbm *usbMonitor) Start(_ typex.RuleX) error {
	// 为了减小问题, 直接把Windows给限制了不支持, 实际上大部分情况下都是Arm-Linux场景
	if runtime.GOOS == "windows" {
		return errors.New("USB monitor plugin not support windows")
	}

	fd, err := unix.Socket(
		unix.AF_NETLINK,
		unix.SOCK_RAW,
		unix.NETLINK_KOBJECT_UEVENT,
	)

	if err != nil {
		return err
	}

	err = unix.Bind(fd, &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: 1,
		Pid:    0,
	})
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
			}
		}
		for {
			data := make([]byte, 1024)
			n, _, _ := unix.Recvfrom(fd, data, 0)
			if n != 0 {
				fmt.Println(string(data[:n]))
			}
		}

	}(typex.GCTX)
	return nil
}

func (usbm *usbMonitor) Stop() error {
	return nil

}

func (usbm *usbMonitor) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "USB Monitor",
		Version:  "0.0.1",
		Homepage: "www.github.com/i4de/rulex",
		HelpLink: "www.github.com/i4de/rulex",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
