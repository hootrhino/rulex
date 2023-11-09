package test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/rulexrpc"
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type _rpcCodecServer struct {
	rulexrpc.UnimplementedCodecServer
}

func (_rpcCodecServer) Decode(c context.Context, req *rulexrpc.CodecRequest) (resp *rulexrpc.CodecResponse, err error) {
	glogger.GLogger.Debug("[REQUEST]=====================> ", req.String())
	resp = new(rulexrpc.CodecResponse)
	resp.Data = []byte("DecodeOK")
	return resp, nil
}
func (_rpcCodecServer) Encode(c context.Context, req *rulexrpc.CodecRequest) (resp *rulexrpc.CodecResponse, err error) {
	glogger.GLogger.Debug("[REQUEST]=====================> ", req.String())
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
		glogger.GLogger.Fatal(err)
		return
	}
	rpcServer := grpc.NewServer()
	rulexrpc.RegisterCodecServer(rpcServer, new(_rpcCodecServer))
	go func(c context.Context) {
		defer listener.Close()
		glogger.GLogger.Info("rpcCodecServer started on", listener.Addr())
		rpcServer.Serve(listener)
	}(context.TODO())

}
func Test_Codec(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	_startServer()
	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC",
		"Rulex Grpc InEnd",
		"Rulex Grpc InEnd", map[string]interface{}{
			"port": 2581,
		})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	grpcCodec1 := typex.NewOutEnd("GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 1998,
			"type": "DECODE",
		})
	grpcCodec1.UUID = "grpcCodec001"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(grpcCodec1, ctx1, cancelF1); err != nil {
		glogger.GLogger.Fatal("grpcCodec load failed:", err)
	}
	grpcCodec2 := typex.NewOutEnd("GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET",
		"GRPC_CODEC_TARGET", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 1998,
			"type": "ENCODE",
		})
	grpcCodec2.UUID = "grpcCodec002"
	ctx2, cancelF2 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(grpcCodec2, ctx2, cancelF2); err != nil {
		glogger.GLogger.Fatal("grpcCodec load failed:", err)
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
			print('rulexlib:RPCDEC ==> ', rulexlib:RPCDEC('grpcCodec001', data))
			print('rulexlib:RPCENC ==> ', rulexlib:RPCENC('grpcCodec002', data))
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Fatal(err)
	}
	grpcConnection, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer grpcConnection.Close()
	client := rulexrpc.NewRulexRpcClient(grpcConnection)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: string([]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16}),
	})
	if err != nil {
		glogger.GLogger.Error(err)
	}
	glogger.GLogger.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
