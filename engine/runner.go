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
func RunRulex(dbPath string) {
	core.InitGlobalConfig()
	engine := NewRuleEngine()
	core.StartLogWatcher()
	rulexlib.StartLuaLogger()
	core.SetLogLevel()
	core.SetPerformance()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine.Start()

	httpServer := httpserver.NewHttpApiServer(2580, "/plugin/http_server/www/", dbPath, engine)
	// Load Http api Server
	engine.LoadPlugin(httpServer)
	// Load Mqtt Server

	if err := engine.LoadPlugin(mqttserver.NewMqttServer()); err != nil {
		log.Error(err)
		panic(err)
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
		newOutEnd := typex.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, config)
		// Important !!!!!!!!
		newOutEnd.UUID = mOutEnd.UUID
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			log.Error("OutEnd load failed:", err)
		}
	}
	signal := <-c
	log.Warn("Received stop signal:", signal)
	engine.Stop()
	//
	// 关闭日志器
	//
	core.GLOBAL_LOGGER.Close()
	rulexlib.LUA_LOGGER.Close()
	os.Exit(0)
}
