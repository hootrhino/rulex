package ttyd_terminal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/ini.v1"
)

type _ttydConfig struct {
	Enable     bool   `ini:"enable"`
	Host       string `ini:"host"`
	ListenPort int    `ini:"listen_port"`
}

/*
*
* WEB终端: 本身不实现SSH协议, 而是通过控制一个外部进程(ttyd)来实现
* 相关资料: https://github.com/tsl0922/ttyd
*
 */
type WebTTYPlugin struct {
	ttydCmd    *exec.Cmd
	mainConfig _ttydConfig
	uuid       string
	busying    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWebTTYPlugin() *WebTTYPlugin {
	return &WebTTYPlugin{
		uuid:       "WEB_TTYD_TERMINAL",
		mainConfig: _ttydConfig{ListenPort: 7681},
		busying:    false,
	}
}

func (tty *WebTTYPlugin) Init(config *ini.Section) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not support current os:%s, only support linux at now", runtime.GOOS)
	}
	_, err := exec.LookPath("ttyd")
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if err := utils.InIMapToStruct(config, &tty.mainConfig); err != nil {
		return err
	}
	return nil
}

/*
*
* 这里拉起来一个 ttyd 进程，默认运行 bash
*
 */
func (tty *WebTTYPlugin) Start(typex.RuleX) error {
	return nil
}
func (tty *WebTTYPlugin) Stop() error {
	if tty.cancel != nil {
		tty.cancel()
	}
	if tty.ttydCmd == nil {
		return nil
	}
	if tty.ttydCmd.ProcessState != nil {
		tty.ttydCmd.Process.Kill()
		tty.ttydCmd.Process.Signal(os.Kill)
	}
	return nil
}

func (hh *WebTTYPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "Web Terminal",
		Version:  "v0.0.1",
		Homepage: "https://github.com/tsl0922/ttyd",
		HelpLink: "https://github.com/tsl0922/ttyd",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
