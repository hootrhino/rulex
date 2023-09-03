package usbmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strings"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	"golang.org/x/sys/unix"
	"gopkg.in/ini.v1"
)

/*
*
* USB 热插拔监控器, 方便观察USB状态, 本插件只支持Linux！！！
*
 */
type usbMonitor struct {
	uuid string
}

func NewUsbMonitor() typex.XPlugin {
	return &usbMonitor{uuid: "USB_EVENT_MONITOR"}
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
			//----
			// add@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			// ACTION=add
			// DEVPATH=/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			// SUBSYSTEM=usb-serial
			// SEQNUM=10822
			//----
			// remove@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			// add@/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/ttyUSB0
			n, _, _ := unix.Recvfrom(fd, data, 0)
			if n > 16 {

				Msg := parseType(data, n)
				if len(Msg) > 0 {
					glogger.GLogger.Info(Msg)
				}

			}
		}

	}(typex.GCTX)
	return nil
}
func parseType(data []byte, len int) string {
	offset := 0
	if string(data[:4]) == "add@" {
		for i := 0; i < len; i++ {
			if data[i] == 0 {
				return parseMsg("add", data, offset)
			} else {
				offset++
			}
		}
	}
	if string(data[:7]) == "remove@" {
		for i := 0; i < len; i++ {
			if data[i] == 0 {
				return parseMsg("remove", data, offset)
			} else {
				offset++
			}
		}
	}
	return ""
}

/*
*
* 只监控串口"/dev/tty*"设备, U盘不管
*
 */
func parseMsg(Type string, data []byte, offset int) string {
	if strings.Contains(string(data[:offset]), "tty") {
		msg := string(data[strings.Index(string(data[:offset]), "tty"):offset])
		nameTokens := strings.Split(msg, "/")
		info := _info{}
		info.Type = Type
		// 1 [ttyUSB0]
		if len(nameTokens) == 1 {
			info.Device = nameTokens[0]
			jsonBytes, _ := json.Marshal(info)
			return string(jsonBytes)
		}
		// 2 [tty ttyUSB0]
		if len(nameTokens) == 2 {
			info.Device = nameTokens[1]
			jsonBytes, _ := json.Marshal(info)
			return string(jsonBytes)
		}
		// 3 [ttyUSB0 tty ttyUSB0]
		if len(nameTokens) == 3 {
			if nameTokens[0] != nameTokens[2] {
				info.Device = nameTokens[2]
				jsonBytes, _ := json.Marshal(info)
				return string(jsonBytes)
			}
		}
	}
	return ""
}
func (usbm *usbMonitor) Stop() error {
	return nil
}

func (usbm *usbMonitor) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     usbm.uuid,
		Name:     "USB Hot Plugin Monitor",
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
func (cs *usbMonitor) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
