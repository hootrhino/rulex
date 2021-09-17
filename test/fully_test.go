package test

import (
	"context"
	"github.com/ngaut/log"
	"google.golang.org/grpc"
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
)

func TestFullyRun(t *testing.T) {
	runTest()
}

//
func runTest() {
	core.InitGlobalConfig()
	Run()
}

//
func Run() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine()
	engine.Start()
	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin(hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	// Load a demo plugin
	if err := engine.LoadPlugin(demo_plugin.NewDemoPlugin()); err != nil {
		log.Error("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", &map[string]interface{}{
		"port": "2581",
	})
	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	// CoAP Inend
	coapInend := typex.NewInEnd("COAP", "Rulex COAP InEnd", "Rulex COAP InEnd", &map[string]interface{}{
		"port": "2582",
	})
	if err := engine.LoadInEnd(coapInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	// Http Inend
	httpInend := typex.NewInEnd("HTTP", "Rulex HTTP InEnd", "Rulex HTTP InEnd", &map[string]interface{}{
		"port": "2583",
	})
	if err := engine.LoadInEnd(httpInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	// Udp Inend
	udpInend := typex.NewInEnd("UDP", "Rulex UDP InEnd", "Rulex UDP InEnd", &map[string]interface{}{
		"port": "2584",
	})
	if err := engine.LoadInEnd(udpInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	//
	// Load Rule
	//
	rule := typex.NewRule(engine,
		"Just a test",
		"Just a test",
		[]string{grpcInend.Id},
		`function Success() print("[LUA Success]==========================> OK") end`,
		`
		Actions = {
			function(data)
			    print("[LUA Actions Callback]==========================> ", data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed]==========================> OK", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		log.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithInsecure())
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `{"temp":100,"hum":30, "co2":123.4, "lex":22.56}`,
	})
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	log.Debugf("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	time.Sleep(2 * time.Second)
	engine.Stop()
}
