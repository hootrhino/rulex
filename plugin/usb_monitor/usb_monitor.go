package usbmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strings"

	"github.com/i4de/rulex/glogger"
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

type _info struct {
	Type   string `json:"type"`
	Device string `json:"device"`
}

func (usbm *usbMonitor) Start(_ typex.RuleX) error {
	// 为了减小问题, 直接把Windows给限制了不支持, 实际上大部分情况下都是Arm-Linux场景
	if runtime.GOOS == "windows" {
		return errors.New("USB monitor plugin not support windows")
	}

	fd, err := unix.Socket(
		unix.AF_NETLINK,
		unix.SOCK_DGRAM,
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
		data := make([]byte, 1024)
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
			// add@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			// ACTION=add
			// DEVPATH=/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			// SUBSYSTEM=usb-serial
			// SEQNUM=10822

			n, _, _ := unix.Recvfrom(fd, data, 0)
			if n > 16 {
				info := []string{}
				if string(data[:4]) == "add@" {
					l := 0
					for i := 0; i < n; i++ {
						if data[i] == 0 {
							if strings.Contains(string(data[:l]), "tty") {
								msg := string(data[strings.Index(string(data[:l]), "tty"):l])
								jsonBytes, _ := json.Marshal(_info{
									Type:   "add",
									Device: msg,
								})
								info = append(info, string(jsonBytes))
							}
							break
						} else {
							l++
						}
					}
				}
				if string(data[:7]) == "remove@" {
					l := 0
					for i := 0; i < n; i++ {
						if data[i] == 0 {
							if strings.Contains(string(data[:l]), "tty") {
								msg := string(data[strings.Index(string(data[:l]), "tty"):l])
								jsonBytes, _ := json.Marshal(_info{
									Type:   "add",
									Device: msg,
								})
								info = append(info, string(jsonBytes))
							}
							break
						} else {
							l++
						}
					}
				}
				if len(info) > 0 {
					// remove@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
					// add@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
					for _, ii := range info {
						glogger.GLogger.Info(ii)
					}
				}
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
