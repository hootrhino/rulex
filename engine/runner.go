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

package engine

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	modbusscanner "github.com/hootrhino/rulex/plugin/modbus_scanner"
	modbusscrc "github.com/hootrhino/rulex/plugin/modbuscrc_tools"
	mqttserver "github.com/hootrhino/rulex/plugin/mqtt_server"
	netdiscover "github.com/hootrhino/rulex/plugin/net_discover"
	ttyterminal "github.com/hootrhino/rulex/plugin/ttyd_terminal"
	usbmonitor "github.com/hootrhino/rulex/plugin/usb_monitor"
	"gopkg.in/ini.v1"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	icmpsender "github.com/hootrhino/rulex/plugin/icmp_sender"
	"github.com/hootrhino/rulex/typex"
)

// 启动 Rulex
func RunRulex(iniPath string) {
	mainConfig := core.InitGlobalConfig(iniPath)
	//----------------------------------------------------------------------------------------------
	// Init logger
	//----------------------------------------------------------------------------------------------
	glogger.StartGLogger(
		core.GlobalConfig.LogLevel,
		mainConfig.EnableConsole,
		mainConfig.AppDebugMode,
		core.GlobalConfig.LogPath,
		mainConfig.AppId, mainConfig.AppName,
	)
	glogger.StartNewRealTimeLogger(core.GlobalConfig.LogLevel)
	//----------------------------------------------------------------------------------------------
	// Init Component
	//----------------------------------------------------------------------------------------------
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.SetDebugMode(mainConfig.EnablePProf)
	core.SetGomaxProcs(mainConfig.GomaxProcs)
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	engine := NewRuleEngine(mainConfig)
	engine.Start()

	// Load Plugin
	loadPlugin(engine)
	// Load Http api Server
	httpServer := httpserver.NewHttpApiServer()
	if err := engine.LoadPlugin("plugin.http_server", httpServer); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range httpServer.AllMInEnd() {
		if err := httpServer.LoadNewestInEnd(minEnd.UUID); err != nil {
			glogger.GLogger.Error("InEnd load failed:", err)
		}
	}

	//
	// Load out from sqlite
	//
	for _, mOutEnd := range httpServer.AllMOutEnd() {
		if err := httpServer.LoadNewestOutEnd(mOutEnd.UUID); err != nil {
			glogger.GLogger.Error("OutEnd load failed:", err)
		}
	}
	// 加载设备
	for _, mDevice := range httpServer.AllDevices() {
		glogger.GLogger.Debug("LoadNewestDevice mDevice.BindRules: ", mDevice.BindRules.String())
		if err := httpServer.LoadNewestDevice(mDevice.UUID); err != nil {
			glogger.GLogger.Error("Device load failed:", err)
		}

	}
	// 加载外挂
	for _, mGoods := range httpServer.AllGoods() {
		newGoods := typex.Goods{
			UUID:        mGoods.UUID,
			Addr:        mGoods.Addr,
			Description: mGoods.Description,
			Args:        mGoods.Args,
		}
		if err := engine.LoadGoods(newGoods); err != nil {
			glogger.GLogger.Error("Goods load failed:", err)
		}
	}
	//
	// APP stack
	//
	for _, mApp := range httpServer.AllApp() {
		app := typex.NewApplication(
			mApp.UUID,
			mApp.Name,
			mApp.Version,
			mApp.Filepath,
		)
		if err := engine.LoadApp(app); err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		if *mApp.AutoStart {
			glogger.GLogger.Debug("App autoStart allowed:", app.UUID, app.Version, app.Name)
			if err1 := engine.StartApp(app.UUID); err1 != nil {
				glogger.GLogger.Error("App autoStart failed:", err1)
			}
		}
	}
	s := <-c
	glogger.GLogger.Warn("Received stop signal:", s)
	engine.Stop()
	os.Exit(0)
}

// loadPlugin 根据Ini配置信息，加载插件
func loadPlugin(engine typex.RuleX) {
	cfg, _ := ini.ShadowLoad(core.INIPath)
	sections := cfg.ChildSections("plugin")
	for _, section := range sections {
		name := strings.TrimPrefix(section.Name(), "plugin.")
		enable, err := section.GetKey("enable")
		if err != nil {
			glogger.GLogger.Fatal(err)
		}
		if !enable.MustBool(false) {
			glogger.GLogger.Warnf("Plugin is disable:%s", name)
			continue
		}
		var plugin typex.XPlugin
		if name == "mqtt_server" {
			plugin = mqttserver.NewMqttServer()
		}
		if name == "usbmonitor" {
			plugin = usbmonitor.NewUsbMonitor()
		}
		if name == "icmpsender" {
			plugin = icmpsender.NewICMPSender()
		}
		if name == "netdiscover" {
			plugin = netdiscover.NewNetDiscover()
		}
		if name == "modbus_scanner" {
			plugin = modbusscanner.NewModbusScanner()
		}
		if name == "ttyd" {
			plugin = ttyterminal.NewWebTTYPlugin()
		}
		if name == "modbuscrc_tools" {
			plugin = modbusscrc.NewModbusCrcCalculator()
		}
		if plugin != nil {
			if err := engine.LoadPlugin(section.Name(), plugin); err != nil {
				glogger.GLogger.Error(err)
			}
		}
	}
}
