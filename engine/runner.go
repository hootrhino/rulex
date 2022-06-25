package engine

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"
	mqttserver "github.com/i4de/rulex/plugin/mqtt_server"
	"github.com/i4de/rulex/sidecar"
	"github.com/i4de/rulex/typex"
)

//
// 启动 Rulex
//
func RunRulex(dbPath string, iniPath string) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig(iniPath)
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.SetLogLevel()
	core.SetPerformance()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	engine := NewRuleEngine(mainConfig)
	engine.Start()

	httpServer := httpserver.NewHttpApiServer(2580, dbPath, engine)
	// Load Http api Server
	if err := engine.LoadPlugin("plugin.http_server", httpServer); err != nil {
		return
	}
	// Load Mqtt Server

	if err := engine.LoadPlugin("plugin.mqtt_server", mqttserver.NewMqttServer()); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range httpServer.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			glogger.GLogger.Error(err)
		}
		// :mInEnd: {k1 :{k1:v1}, k2 :{k2:v2}} --> InEnd: [{k1:v1}, {k2:v2}]
		var dataModelsMap map[string]typex.XDataModel
		if err := json.Unmarshal([]byte(minEnd.XDataModels), &dataModelsMap); err != nil {
			glogger.GLogger.Error(err)
		}
		in := typex.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, config)
		in.UUID = minEnd.UUID // Important !!!!!!!!
		in.DataModelsMap = dataModelsMap
		if err := engine.LoadInEnd(in); err != nil {
			glogger.GLogger.Error("InEnd load failed:", err)
		}
	}

	//
	// Load out from sqlite
	//
	for _, mOutEnd := range httpServer.AllMOutEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
			glogger.GLogger.Error(err)
		}
		newOutEnd := typex.NewOutEnd(typex.TargetType(mOutEnd.Type), mOutEnd.Name, mOutEnd.Description, config)
		newOutEnd.UUID = mOutEnd.UUID // Important !!!!!!!!
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			glogger.GLogger.Error("OutEnd load failed:", err)
		}
	}
	// 加载设备
	for _, mDevice := range httpServer.AllDevices() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
			glogger.GLogger.Error(err)
		}
		newDevice := typex.NewDevice(typex.DeviceType(mDevice.Type), mDevice.Name, mDevice.Description, mDevice.ActionScript, config)
		newDevice.UUID = mDevice.UUID // Important !!!!!!!!
		if err := engine.LoadDevice(newDevice); err != nil {
			glogger.GLogger.Error("Device load failed:", err)
		}
	}
	// 加载外挂
	for _, mGoods := range httpServer.AllGoods() {
		newGoods := sidecar.Goods{
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
	// 规则最后加载
	//
	for _, mRule := range httpServer.AllMRules() {
		rule := typex.NewRule(engine,
			mRule.UUID,
			mRule.Name,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		if err := engine.LoadRule(rule); err != nil {
			glogger.GLogger.Error(err)
		}
	}
	s := <-c
	glogger.GLogger.Warn("Received stop signal:", s)
	engine.Stop()
	os.Exit(0)
}
