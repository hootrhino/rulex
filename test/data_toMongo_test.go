package test

import (
	"context"
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexlib"
	"rulex/rulexrpc"
	"rulex/typex"
	"testing"
	"time"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_DataToMongo(t *testing.T) {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.StartLogWatcher(core.GlobalConfig.LogPath)
	rulexlib.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "./rulex.db", engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": 2581,
	})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	//
	mongoOut := typex.NewOutEnd("MONGO_SINGLE",
		"MONGO_SINGLE",
		"MONGO_SINGLE", map[string]interface{}{
			"mongoUrl":   "mongodb://admin:123456@127.0.0.1:27017/?connect=direct",
			"database":   "gateway_test",
			"collection": "env",
		})
	mongoOut.UUID = "mongoOut"
	if err := engine.LoadOutEnd(mongoOut); err != nil {
		log.Fatal(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    local r = rulexlib:DataToMongo('mongoOut', data)
				print("[LUA DataToMongo] ==>", r)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		log.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `[{"co2":10,"hum":30,"lex":22,"temp":100}]`,
	})
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
