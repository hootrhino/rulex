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

func Test_S7_PLC_Parse(t *testing.T) {
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
				local V0 = rulexlib:MB(">a:16 b:16 c:16 d:16 e:16", data, false)
				local a = rulexlib:T2J(V0['a'])
				local b = rulexlib:T2J(V0['b'])
				local c = rulexlib:T2J(V0['c'])
				local d = rulexlib:T2J(V0['d'])
				local e = rulexlib:T2J(V0['e'])
				print('a ==> ', a, ' ->', rulexlib:BS2B(a), '==> ', rulexlib:B2I64('>', rulexlib:BS2B(a)))
				print('b ==> ', b, ' ->', rulexlib:BS2B(a), '==> ', rulexlib:B2I64('>', rulexlib:BS2B(b)))
				print('c ==> ', c, ' ->', rulexlib:BS2B(a), '==> ', rulexlib:B2I64('>', rulexlib:BS2B(c)))
				print('d ==> ', d, ' ->', rulexlib:BS2B(a), '==> ', rulexlib:B2I64('>', rulexlib:BS2B(d)))
				print('e ==> ', e, ' ->', rulexlib:BS2B(a), '==> ', rulexlib:B2I64('>', rulexlib:BS2B(e)))
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
		Value: string([]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16}),
	})
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
