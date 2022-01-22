package test

import (
	"context"
	"os"
	"os/signal"
	"rulex/core"
	"rulex/engine"
	"rulex/plugin/demo_plugin"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexrpc"
	"rulex/typex"
	"syscall"
	"testing"
	"time"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
)

func Test_Binary_LUA_Parse(t *testing.T) {
	core.InitGlobalConfig()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine(core.InitGlobalConfig())
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates", "./rulex.db", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin(hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	// Load a demo plugin
	if err := engine.LoadPlugin(demo_plugin.NewDemoPlugin()); err != nil {
		log.Error("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": "2581",
	})
	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	//
	// Load Rule
	//
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			--        ┌───────────────────────────────────────────────┐
			-- data = |00 00 00 01|00 00 00 02|00 00 00 03|00 00 00 04|
			--        └───────────────────────────────────────────────┘
			function(data)
				local json = require("json")
				local V6 = json.encode(rulexlib:MatchBinary("<a:8 b:8 c:8 d:8", data, false))
				print("[LUA Actions Callback 5, rulex.MatchBinary] ==>", V6)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		log.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581")
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	for i := 0; i < 2; i++ {
		log.Infof("Rulex Rpc Call ==========================>>: %v", i)
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: string([]byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9,
				10, 11, 12, 13, 14, 15, 16}),
		})
		if err != nil {
			log.Error("grpc.Dial err: %v", err)
		}
		log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(1 * time.Second)
	engine.Stop()
}
