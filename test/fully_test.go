package test

import (
	"context"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"
	"github.com/i4de/rulex/rulexrpc"
	"github.com/i4de/rulex/typex"

	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestFullyRun(t *testing.T) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)

	engine.Start()
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer(
		2580,
		"../"+GenDate()+".db",
		engine,
	)); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": 2581,
	})
	if err := engine.LoadInEnd(grpcInend); err != nil {
		glogger.GLogger.Error("grpcInend load failed:", err)
	}
	//
	// Load Rule
	//
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    local V1 = rulexlib:JQ(".[] | select(.temp > 50000000)", data)
                print("[LUA Actions Callback 1 ===> Data is:", data)
			    print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", rulexlib:JQ(".[] | select(.temp > 50000000)", data))
				return true, data
			end,
			function(data)
			    local V2 = rulexlib:JQ(".[] | select(.hum < 20)", data)
			    print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", rulexlib:JQ(".[] | select(.hum < 20)", data))
				return true, data
			end,
			function(data)
			    local V3 = rulexlib:JQ(".[] | select(.co2 > 50)", data)
			    print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", rulexlib:JQ(".[] | select(.co2 > 50)", data))
				return true, data
			end,
			function(data)
			    local V4 = rulexlib:JQ(".[] | select(.lex > 50)", data)
			    print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", rulexlib:JQ(".[] | select(.lex > 50)", data))
				return true, data
			end,
			function(data)
				--
				print("[LUA Actions Callback 5, rulexlib:J2T] ==>",rulexlib:J2T(data))
				print("[LUA Actions Callback 5, rulexlib:T2J] ==>",rulexlib:T2J(rulexlib:J2T(data)))
				return true, data
			end,
			function(data)
			    --
				-- 0110_0001 0110_0001 0110_0010
				-- <a:5 b:3 c:1 => a:00001100 b:00000001 c:0
				local V6 = rulexlib:T2J(rulexlib:MB("<a:5 b:3 c:1", "aab", false))
				print("[LUA Actions Callback 6, rulex.MatchBinary] ==>", V6)
				print("[LUA Actions Callback 6, rulex.MatchBinary] ==>", V6)
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
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				print("[LUA Actions Callback RULE ==================> 1] ==>", data)
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
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    print("[LUA Actions Callback RULE ==================> 2] ==>", data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	rule4 := typex.NewRule(engine,
		"uuid4",
		"rule4",
		"rule4",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[rulexlib:J2T(data) Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    local t1 = rulexlib:J2T(data)
			    print("[rulexlib:J2T(data)] ==>", rulexlib:T2J(t1))
				return true, data
			end,
			function(data)
			    print("[rulexlib:Time()] ==>", rulexlib:Time())
			    print("[rulexlib:TsUnix()] ==>", rulexlib:TsUnix())
			    print("[rulexlib:TsUnixNano()] ==>", rulexlib:TsUnixNano())
			    rulexlib:VSet('k', 'v')
			    print("[rulexlib:VGet(k)] ==>", rulexlib:VGet('k'))

				return true, data
			end
		}`,
		`function Failed(error) print("[rulexlib:J2T(data) Failed Callback]", error) end`)
	if err := engine.LoadRule(rule1); err != nil {
		glogger.GLogger.Error(err)
	}
	if err := engine.LoadRule(rule2); err != nil {
		glogger.GLogger.Error(err)
	}
	if err := engine.LoadRule(rule3); err != nil {
		glogger.GLogger.Error(err)
	}
	if err := engine.LoadRule(rule4); err != nil {
		glogger.GLogger.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	for i := 0; i < 30; i++ {
		glogger.GLogger.Infof("Test count ==========================>>: %v", i)
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: `
					[
						{"co2":10,"hum":30,"lex":22,"temp":100},
						{"co2":100,"hum":300,"lex":220,"temp":1000},
						{"co2":1000,"hum":3000,"lex":2200,"temp":10000}
					]
	`})

		if err != nil {
			glogger.GLogger.Error(err)
		}
		glogger.GLogger.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(1 * time.Second)
	glogger.GLogger.Info("Test Http system Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/system"))
	glogger.GLogger.Info("Test Http inends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/inends"))
	glogger.GLogger.Info("Test Http outends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/outends"))
	glogger.GLogger.Info("Test Http rules Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/rules"))

	engine.Stop()
}
