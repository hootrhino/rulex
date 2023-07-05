package engine

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	icmpsender "github.com/hootrhino/rulex/plugin/icmp_sender"
	mqttserver "github.com/hootrhino/rulex/plugin/mqtt_server"
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
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
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
	// Load Http api Server
	httpServer := httpserver.NewHttpApiServer()
	if err := engine.LoadPlugin("plugin.http_server", httpServer); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	mqttServer := mqttserver.NewMqttServer()
	if err := engine.LoadPlugin("plugin.mqtt_server", mqttServer); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	icmpSender := icmpsender.NewICMPSender()
	if err := engine.LoadPlugin("plugin.icmpsender", icmpSender); err != nil {
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
