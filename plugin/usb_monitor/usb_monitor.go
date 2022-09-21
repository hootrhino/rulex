package usbmonitor

import (
	"context"
	"errors"
	"runtime"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/rubiojr/go-usbmon"
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
	filter := &usbmon.ActionFilter{Action: usbmon.ActionAll}
	devs, err := usbmon.ListenFiltered(typex.GCTX, filter)
	if err != nil {
		return err
	}
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			{
				return
			}
		case dev := <-devs:
			{
				glogger.GLogger.Infof("USB device %s:%s has %s", dev.Major(), dev.Serial(), dev.Action())
			}
		default:
			{
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
