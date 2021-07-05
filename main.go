package main

import (
	"os"
	"os/signal"
	"rulenginex/plugin/http_server"
	"rulenginex/x"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

//
var engine *x.RuleEngine

//
func main() {
	gin.SetMode(gin.ReleaseMode)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	engine = x.NewRuleEngine()

	in1 := x.NewInEnd("MQTT", "MQTT Stream", "MQTT Input Stream", &map[string]interface{}{
		"server":   "127.0.0.1",
		"port":     1883,
		"username": "test",
		"password": "test",
		"clientId": "test",
	})
	in1.Id = "INEND1"
	if err0 := engine.LoadInEnd(in1); err0 != nil {
		log.Fatal("InEnd load failed:", err0)

	}
	in2 := x.NewInEnd("COAP", "COAP Stream", "COAP Input Stream", &map[string]interface{}{
		"server": "127.0.0.1",
		"port":   1883,
	})
	in2.Id = "INEND2"
	if err := engine.LoadInEnd(in2); err != nil {
		log.Fatal("InEnd load failed:", err)
	}
	out1 := x.NewOutEnd("mongo", "Data to mongodb", "Save data to mongodb",
		&map[string]interface{}{
			"mongourl": "mongodb+srv://rulenginex:rulenginex@cluster0.rsdmb.mongodb.net/test",
		})
	out1.Id = "MongoDB001"
	if err1 := engine.LoadOutEnds(out1); err1 != nil {
		log.Fatal("OutEnd load failed:", err1)
	}
	actions := `
local json = require("json")
Actions = {
	function(data)
	    dataToMongo("MongoDB001", data)
	    print("[LUA Actions Callback]:dataToMongo Mqtt payload:", data)
		return true, data
	end
}
`
	from := []string{in1.Id}
	failed := `
function Failed(error)
  -- print("[LUA Callback] call failed from lua:", error)
end
`
	success := `
function Success()
  -- print("[LUA Callback] call success from lua")
end
`
	rule1 := x.NewRule(engine, "just_a_test_rule", "just_a_test_rule", from, success, actions, failed)
	rule1.Id = "just_a_test_rule"

	//
	if e := engine.LoadRule(rule1); e != nil {
		log.Fatal("rule load failed:", e)
	}
	httpServer := plugin.HttpApiServer{}
	if e := engine.LoadPlugin(&httpServer); e != nil {
		log.Fatal("rule load failed:", e)
	}
	engine.Start()
	<-c
	os.Exit(0)
}
