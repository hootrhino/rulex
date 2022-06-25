package test

import (
	"context"
	"testing"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"
	"github.com/i4de/rulex/rulexrpc"
	"github.com/i4de/rulex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_JQ_Parse(t *testing.T) {
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

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": 2581,
	})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				print("[LUA rulexlib:J2T] ==>",rulexlib:J2T(data))
				print("[LUA rulexlib:T2J] ==>",rulexlib:T2J(rulexlib:J2T(data)))
				print("[LUA rulexlib:T2J] ==>",rulexlib:T2J(rulexlib:J2T(data)) == data)
				return true, data
			end,
			function(data)
			print(data)
				return true, data
			end,
			function(data)
			print(data)
				return true, data
			end,
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `[{"co2":10,"hum":30,"lex":22,"temp":100}]`,
	})
	if err != nil {
		glogger.GLogger.Error(err)
	}
	glogger.GLogger.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
