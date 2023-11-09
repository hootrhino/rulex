package test

import (
	"context"
	"os"
	"os/signal"
	"syscall"
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
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* 运行时镜像给dump出来
*
 */
func Test_snapshot_dump(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
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
		"port": 2581,
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
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
			function(args)
			    print(data)
				local json = require("json")
				local V6 = json.decode(data)
				local V7 = json.encode(rulexlib:MB(">a:16 b:8 c:8", data, false))
				-- {"a":"0000000000000001","b":"00000000","c":"00000001"}
				print("[LUA Actions Callback, rulex.MatchBinary] ==>", V7)
				return true, args
			end
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
		// lua 输出 {"a":"0000000000000001","b":"00000000","c":"00000001"}
		Value: string([]byte{0, 1, 0, 1}),
	})
	if err != nil {
		glogger.GLogger.Error(err)
	}
	glogger.GLogger.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	t.Log(engine.SnapshotDump())
	time.Sleep(1 * time.Second)
	engine.Stop()
}
