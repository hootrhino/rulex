package engine

import (
	"encoding/json"
	"os"
	"os/signal"
	"rulex/core"
	httpserver "rulex/plugin/http_server"
	// mqttserver "rulex/plugin/mqtt_server"
	"rulex/typex"
	"syscall"

	"github.com/ngaut/log"
)

//
//
//
func RunRulex(dbPath string) {
	core.InitGlobalConfig()
	core.SetLogLevel()
	core.SetPerformance()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := NewRuleEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates", dbPath, engine)
	// Load Http api Server
	engine.LoadPlugin(hh)
	// Load Mqtt Server

	// if err := engine.LoadPlugin(mqttserver.NewMqttServer()); err != nil {
	// 	log.Error(err)
	// 	panic(err)
	// }
	//
	// Load inend from sqlite
	//
	for _, minEnd := range hh.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Error(err)
		}
		in1 := typex.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, config)
		// Important !!!!!!!!
		in1.Id = minEnd.UUID
		if err := engine.LoadInEnd(in1); err != nil {
			log.Error("InEnd load failed:", err)
		}
	}

	//
	// Load rule from sqlite
	//
	for _, mRule := range hh.AllMRules() {
		rule := typex.NewRule(engine,
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
	for _, mOutEnd := range hh.AllMOutEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
			log.Error(err)
		}
		newOutEnd := typex.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, config)
		// Important !!!!!!!!
		newOutEnd.Id = mOutEnd.UUID
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			log.Error("OutEnd load failed:", err)
		}
	}
	signal := <-c
	log.Warn("Received stop signal:", signal)
	engine.Stop()
	os.Exit(0)
}
