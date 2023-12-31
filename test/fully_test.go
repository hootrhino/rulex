package test

import (
	"context"

	httpserver "github.com/hootrhino/rulex/component/rulex_api_server"
	"github.com/hootrhino/rulex/component/rulexrpc"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestFullyRun(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer(engine)); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"host": "127.0.0.1",
		"port": 2581,
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Error("grpcInend load failed:", err)
	}
	//
	// Load Rule
	//
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		grpcInend.UUID,
		"",
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
			    return true, args
			end,
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	//--------------------------------------------------
	rule2 := typex.NewRule(engine,
		"uuid2",
		"rule2",
		"rule2",
		grpcInend.UUID,
		"",
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
				print("[LUA Actions Callback RULE ==================> 1] ==>", data)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	//--------------------------------------------------
	rule3 := typex.NewRule(engine,
		"uuid3",
		"rule3",
		"rule3",
		grpcInend.UUID,
		"",
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
			    print("[LUA Actions Callback RULE ==================> 2] ==>", data)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	rule4 := typex.NewRule(engine,
		"uuid4",
		"rule4",
		"rule4",
		grpcInend.UUID,
		"",
		`function Success() print("[rulexlib:J2T(data) Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
				return true, args
			end,
			function(args)
			    print("[HelloLib] ==>", Hello())
			    return true, args
			end,
			function(args)
			    print(rulexlib:Time())
			    return true, args
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

	glogger.GLogger.Info("Test Http system Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/system"))
	glogger.GLogger.Info("Test Http inends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/inends"))
	glogger.GLogger.Info("Test Http outends Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/outends"))
	glogger.GLogger.Info("Test Http rules Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/rules"))

	time.Sleep(5 * time.Second)
	engine.Stop()
}
