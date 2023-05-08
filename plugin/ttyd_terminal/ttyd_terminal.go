package ttyd_terminal

import (
	"fmt"
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
}

func NewWebTTYPlugin() *WebTTYPlugin {
	return &WebTTYPlugin{
		mainConfig: _ttydConfig{},
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
	if tty.mainConfig.ListenPort == 0 {
		tty.mainConfig.ListenPort = 7681
	}

	return nil
}

/*
*
* 这里拉起来一个 ttyd 进程，默认运行 bash
*
 */
func (tty *WebTTYPlugin) Start(typex.RuleX) error {
	tty.ttydCmd = exec.CommandContext(typex.GCTX,
		"ttyd", "-p", fmt.Sprintf("%d", tty.mainConfig.ListenPort),
		"-o", "-6", "bash")
	// tty.ttydCmd.Stdout = glogger.GLogger.Out
	// tty.ttydCmd.Stderr = glogger.GLogger.Out
	if err1 := tty.ttydCmd.Start(); err1 != nil {
		glogger.GLogger.Info("cmd.Start error: %v", err1)
		return err1
	}
	go func(cmd *exec.Cmd) {
		glogger.GLogger.Info("ttyd started successfully on port:", tty.mainConfig.ListenPort)
		cmd.Process.Wait() // blocked until exited
		glogger.GLogger.Info("ttyd stopped")
	}(tty.ttydCmd)
	return nil
}
func (tty *WebTTYPlugin) Stop() error {
	if tty.ttydCmd == nil {
		return nil
	}
	if tty.ttydCmd.ProcessState != nil {
		tty.ttydCmd.Process.Kill()
	}
	return nil
}

func (hh *WebTTYPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "WebTTYPlugin",
		Version:  "0.0.1",
		Homepage: "https://github.com/tsl0922/ttyd",
		HelpLink: "https://github.com/tsl0922/ttyd",
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
func (cs *WebTTYPlugin) Service(arg typex.ServiceArg) error {
	return nil
}
