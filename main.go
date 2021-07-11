package main

import (
	"encoding/json"
	"os"
	"os/signal"
	httpserver "rulex/plugin/http_server"
	"rulex/x"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

//
func main() {
	gin.SetMode(gin.ReleaseMode)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	engine := x.NewRuleEngine()
	engine.Start()
	////////
	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates/*", engine)
	if e := engine.LoadPlugin(hh); e != nil {
		log.Fatal("rule load failed:", e)
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range hh.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Error(err)
		}
		in1 := x.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, &config)
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
		rule := x.NewRule(engine,
			mRule.Name,
			mRule.Description,
			strings.Split(mRule.From, ","),
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
		newOutEnd := x.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, &config)
		// Important !!!!!!!!
		newOutEnd.Id = mOutEnd.UUID
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			log.Error("OutEnd load failed:", err)
		}
	}
	<-c
	os.Exit(0)
}
