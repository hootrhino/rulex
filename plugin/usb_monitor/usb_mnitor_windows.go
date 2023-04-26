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
	return errors.New("USB monitor plugin not support windows")
}

func (usbm *usbMonitor) Stop() error {
	return nil
}

func (usbm *usbMonitor) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "USB Monitor",
		Version:  "0.0.1",
		Homepage: "https://github.com/hootrhino/rulex.git",
		HelpLink: "https://github.com/hootrhino/rulex.git",
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
func (usbm *usbMonitor) Service(arg typex.ServiceArg) error {
	return nil
}
