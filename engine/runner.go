package engine

import (
	"encoding/json"
	"os"
	"os/signal"
	"rulex/core"
	httpserver "rulex/plugin/http_server"
	mqttserver "rulex/plugin/mqtt_server"
	"rulex/rulexlib"
	"rulex/typex"
	"syscall"

	"github.com/ngaut/log"
)

//
// 启动 Rulex
//
func RunRulex(dbPath string, iniPath string) {
	mainConfig := core.InitGlobalConfig(iniPath)
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.StartLogWatcher(core.GlobalConfig.LogPath)
	rulexlib.StartLuaLogger(core.GlobalConfig.LuaLogPath)
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
		log.Error(err)
		return
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range httpServer.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Error(err)
		}
		// :mInEnd: {k1 :{k1:v1}, k2 :{k2:v2}} --> InEnd: [{k1:v1}, {k2:v2}]
		var dataModelsMap map[string]typex.XDataModel
		if err := json.Unmarshal([]byte(minEnd.XDataModels), &dataModelsMap); err != nil {
			log.Error(err)
		}
		in := typex.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, config)
		// Important !!!!!!!!
		in.UUID = minEnd.UUID
		in.DataModelsMap = dataModelsMap
		if err := engine.LoadInEnd(in); err != nil {
			log.Error("InEnd load failed:", err)
		}
	}

	//
	// Load rule from sqlite
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
			log.Error(err)
		}
	}
	//
	// Load out from sqlite
	//
	for _, mOutEnd := range httpServer.AllMOutEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
			log.Error(err)
		}
		newOutEnd := typex.NewOutEnd(typex.TargetType(mOutEnd.Type), mOutEnd.Name, mOutEnd.Description, config)
		// Important !!!!!!!!
		newOutEnd.UUID = mOutEnd.UUID
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			log.Error("OutEnd load failed:", err)
		}
	}
	// 加载设备
	for _, mDevice := range httpServer.AllDevices() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
			log.Error(err)
		}
		newDevice := &typex.Device{
			UUID:         mDevice.UUID,
			Name:         mDevice.Name,
			Type:         typex.DeviceType(mDevice.Type),
			ActionScript: mDevice.ActionScript,
			Description:  mDevice.Description,
			Config:       config,
		}
		if err := engine.LoadDevice(newDevice); err != nil {
			log.Error("Device load failed:", err)
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
			log.Error("Goods load failed:", err)
		}
	}
	s := <-c
	log.Warn("Received stop signal:", s)
	engine.Stop()
	//
	// 关闭日志器
	//
	if err := typex.GLOBAL_LOGGER.Close(); err != nil {
		return
	}
	if err := typex.LUA_LOGGER.Close(); err != nil {
		return
	}
	os.Exit(0)
}
