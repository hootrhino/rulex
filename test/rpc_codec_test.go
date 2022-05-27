package test

import (
	"context"
	"net"
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

type _rpcCodecServer struct {
	rulexrpc.UnimplementedCodecServer
}

func (_rpcCodecServer) Decode(c context.Context, req *rulexrpc.CodecRequest) (resp *rulexrpc.CodecResponse, err error) {
	log.Debug("[REQUEST]=====================> ", req.String())
	resp = new(rulexrpc.CodecResponse)
	resp.Data = []byte("DecodeOK")
	return resp, nil
}
func (_rpcCodecServer) Encode(c context.Context, req *rulexrpc.CodecRequest) (resp *rulexrpc.CodecResponse, err error) {
	log.Debug("[REQUEST]=====================> ", req.String())
	resp = new(rulexrpc.CodecResponse)
	resp.Data = []byte("EncodeOK")
	return resp, nil
}

/*
*
*
*
 */
func _startServer() {
	listener, err := net.Listen("tcp", ":1998")
	if err != nil {
		log.Fatal(err)
		return
	}
	rpcServer := grpc.NewServer()
	rulexrpc.RegisterCodecServer(rpcServer, new(_rpcCodecServer))
	go func(c context.Context) {
		log.Info("rpcCodecServer started on", listener.Addr())
		rpcServer.Serve(listener)
	}(context.TODO())

}
func Test_Codec(t *testing.T) {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.StartLogWatcher(core.GlobalConfig.LogPath)
	rulexlib.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()
	_startServer()
	hh := httpserver.NewHttpApiServer(2580, "./rulex.db", engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC",
		"Rulex Grpc InEnd",
		"Rulex Grpc InEnd", map[string]interface{}{
			"port": 2581,
		})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	grpcCodec1 := typex.NewOutEnd("GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 1998,
			"type": "DECODE",
		})
	grpcCodec1.UUID = "grpcCodec001"
	if err := engine.LoadOutEnd(grpcCodec1); err != nil {
		log.Fatal("grpcCodec load failed:", err)
	}
	grpcCodec2 := typex.NewOutEnd("GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 1998,
			"type": "ENCODE",
		})
	grpcCodec2.UUID = "grpcCodec002"
	if err := engine.LoadOutEnd(grpcCodec2); err != nil {
		log.Fatal("grpcCodec load failed:", err)
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
			print('rulexlib:RPCDEC ==> ', rulexlib:RPCDEC('grpcCodec001', data))
			print('rulexlib:RPCENC ==> ', rulexlib:RPCENC('grpcCodec002', data))
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		log.Fatal(err)
	}
	grpcConnection, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc.Dial err: %v", err)
	}
	defer grpcConnection.Close()
	client := rulexrpc.NewRulexRpcClient(grpcConnection)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: string([]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16}),
	})
	if err != nil {
		log.Fatal("grpc.Dial err: %v", err)
	}
	log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
