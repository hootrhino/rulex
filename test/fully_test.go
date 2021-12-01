package test

import (
	"context"
	"os"
	"os/signal"
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexrpc"
	"rulex/typex"
	"syscall"
	"testing"
	"time"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
)

func TestFullyRun(t *testing.T) {
	core.InitGlobalConfig()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "/../plugin/http_server/www/", "./rulex.db", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin(hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": 2581,
	})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Error("grpcInend load failed:", err)
	}
	//
	// Load Rule
	//
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    local V1 = rulex:JqSelect(".[] | select(.temp > 50000000)", data)
                print("[LUA Actions Callback 1 ===> Data is:", data)
			    print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", rulex:JqSelect(".[] | select(.temp > 50000000)", data))
				return true, data
			end,
			function(data)
			    local V2 = rulex:JqSelect(".[] | select(.hum < 20)", data)
			    print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", rulex:JqSelect(".[] | select(.hum < 20)", data))
				return true, data
			end,
			function(data)
			    local V3 = rulex:JqSelect(".[] | select(.co2 > 50)", data)
			    print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", rulex:JqSelect(".[] | select(.co2 > 50)", data))
				return true, data
			end,
			function(data)
			    local V4 = rulex:JqSelect(".[] | select(.lex > 50)", data)
			    print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", rulex:JqSelect(".[] | select(.lex > 50)", data))
				return true, data
			end,
			function(data)
				local json = require("json")
				print("[LUA Actions Callback 5, json.decode] ==>",json.decode(data))
				print("[LUA Actions Callback 5, json.encode] ==>",json.encode(json.decode(data)))
				return true, data
			end,
			function(data)
				local json = require("json")
				-- 0110_0001 0110_0001 0110_0010
				-- <a:5 b:3 c:1 => a:00001100 b:00000001 c:0
				local V6 = json.encode(rulex:MatchBinary("<a:5 b:3 c:1", "aab", false))
				print("[LUA Actions Callback 5, rulex.MatchBinary] ==>", V6)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	//--------------------------------------------------
	rule2 := typex.NewRule(engine,
		"uuid2",
		"rule2",
		"rule2",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				local json = require("json")
				-- 0110_0001 0110_0001 0110_0010
				-- <a:5 b:3 c:1 => a:00001100 b:00000001 c:0
				local V6 = json.encode(rulex:MatchBinary("<a:5 b:3 c:1", "aab", false))
				print("[LUA Actions Callback RULE ==================> uuid2] ==>", V6)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	//--------------------------------------------------
	rule3 := typex.NewRule(engine,
		"uuid3",
		"rule3",
		"rule3",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				local json = require("json")
				-- 0110_0001 0110_0001 0110_0010
				-- <a:5 b:3 c:1 => a:00001100 b:00000001 c:0
				local V6 = json.encode(rulex:MatchBinary("<a:5 b:3 c:1", "aab", false))
				print("[LUA Actions Callback RULE ==================> uuid3] ==>", V6)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule1); err != nil {
		log.Error(err)
	}
	if err := engine.LoadRule(rule2); err != nil {
		log.Error(err)
	}
	if err := engine.LoadRule(rule3); err != nil {
		log.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithInsecure())
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	for i := 0; i < 30; i++ {
		log.Infof("Rulex Rpc Call ==========================>>: %v", i)
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: `
					[
						{"co2":10,"hum":30,"lex":22,"temp":100},
						{"co2":100,"hum":300,"lex":220,"temp":1000},
						{"co2":1000,"hum":3000,"lex":2200,"temp":10000}
					]
	`})

		if err != nil {
			log.Error("grpc.Dial err: %v", err)
		}
		log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(1 * time.Second)
	log.Info("Test Http system Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/system"))
	log.Info("Test Http inends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/inends"))
	log.Info("Test Http outends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/outends"))
	log.Info("Test Http rules Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/rules"))

	engine.Stop()
}
