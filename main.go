package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"rulex/plugin/http_server"
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
	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates/*")
	if e := engine.LoadPlugin(hh); e != nil {
		log.Fatal("rule load failed:", e)
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range hh.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Fatal(err)
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
		rule1 := x.NewRule(engine,
			mRule.Name,
			mRule.Description,
			strings.Split(mRule.From, ","),
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		if err := engine.LoadRule(rule1); err != nil {
			log.Error(err)
		}
	}
	<-c
	os.Exit(0)
}
