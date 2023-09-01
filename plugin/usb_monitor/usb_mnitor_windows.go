package usbmonitor

import (
	"errors"

	"github.com/hootrhino/rulex/typex"

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
	return &usbMonitor{
		uuid: "USB-MONITOR",
	}
}
func (usbm *usbMonitor) Init(_ *ini.Section) error {
	return nil

}

func (usbm *usbMonitor) Start(_ typex.RuleX) error {
	return errors.New("USB monitor plugin not support windows")
}

func (usbm *usbMonitor) Stop() error {
	return nil
}

func (usbm *usbMonitor) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     usbm.uuid,
		Name:     "USB Monitor",
		Version:  "v0.0.1",
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
func (usbm *usbMonitor) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
