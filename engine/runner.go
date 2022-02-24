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
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
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
		in1 := typex.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, config)
		// Important !!!!!!!!
		in1.UUID = minEnd.UUID
		if err := engine.LoadInEnd(in1); err != nil {
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
			mRule.From,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		if err := engine.LoadRule(rule); err != nil {
			log.Error(err)
		}
	}
	//
	// Load out end from sqlite
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
	s := <-c
	log.Warn("Received stop signal:", s)
	engine.Stop()
	//
	// 关闭日志器
	//
	if err := core.GLOBAL_LOGGER.Close(); err != nil {
		return
	}
	if err := rulexlib.LUA_LOGGER.Close(); err != nil {
		return
	}
	os.Exit(0)
}
