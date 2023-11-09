package test

import (
	"context"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/rulexrpc"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/engine"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/demo_plugin"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"

	"google.golang.org/grpc"
)

func Test_Binary_LUA_Parse(t *testing.T) {
	engine := engine.InitRuleEngine(core.InitGlobalConfig("conf/rulex.ini"))
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Load a demo plugin
	if err := engine.LoadPlugin("plugin.demo", demo_plugin.NewDemoPlugin()); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": "2581",
	})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}
	//
	// Load Rule
	//
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			--        ┌───────────────────────────────────────────────┐
			-- data = |00 00 00 01|00 00 00 02|00 00 00 03|00 00 00 04|
			--        └───────────────────────────────────────────────┘
			function(args)
				local json = require("json")
				local V6 = json.encode(rulexlib:MB("<a:8 b:8 c:8 d:8", data, false))
				print("[LUA Actions Callback 5, rulex.MatchBinary] ==>", V6)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581")
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	for i := 0; i < 2; i++ {
		glogger.GLogger.Infof("Rulex Rpc Call ==========================>>: %v", i)
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: string([]byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9,
				10, 11, 12, 13, 14, 15, 16}),
		})
		if err != nil {
			glogger.GLogger.Error(err)

		}
		glogger.GLogger.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(1 * time.Second)
	engine.Stop()
}
