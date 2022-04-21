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

func Test_JQ_Parse(t *testing.T) {
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

	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				local V2 = rulexlib:JQ(".[] | select(.hum < 20)", data)
				print("rulexlib:JQ(".[] | select(.hum < 20)", data))
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
		Value: `[{"co2":10,"hum":30,"lex":22,"temp":100},{"co2":100,"hum":300,"lex":220,"temp":1000},{"co2":1000,"hum":3000,"lex":2200,"temp":10000}]`,
	})
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
