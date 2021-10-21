package test

import (
	"context"
	"io/ioutil"
	"net/http"
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

func TestFullyRun(t *testing.T) {
	Run()
}

//
func Run() {

	core.InitGlobalConfig()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine()
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
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    local V1 = stdlib:JqSelect(".[] | select(.temp > 50000000)", data)
                print("[LUA Actions Callback 1 ===> Data is:", data)
			    print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", stdlib:JqSelect(".[] | select(.temp > 50000000)", data))
				return true, data
			end,
			function(data)
			    local V2 = stdlib:JqSelect(".[] | select(.hum < 20)", data)
			    print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", stdlib:JqSelect(".[] | select(.hum < 20)", data))
				return true, data
			end,
			function(data)
			    local V3 = stdlib:JqSelect(".[] | select(.co2 > 50)", data)
			    print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", stdlib:JqSelect(".[] | select(.co2 > 50)", data))
				return true, data
			end,
			function(data)
			    local V4 = stdlib:JqSelect(".[] | select(.lex > 50)", data)
			    print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", stdlib:JqSelect(".[] | select(.lex > 50)", data))
				return true, data
			end,
			function(data)
				local json = require("json")
				print("[LUA Actions Callback 5, json.decode] ==>",json.decode(data))
				print("[LUA Actions Callback 5, json.encode] ==>",json.encode(json.decode(data)))
				return true, data
			end,
			function(data)
				local json = require("json")
				-- 0110_0001 0110_0001 0110_0010
				-- <a:5 b:3 c:1 => a:00001100 b:00000001 c:0
				local V6 = json.encode(stdlib:MatchBinary("<a:5 b:3 c:1", "aab", false))
				print("[LUA Actions Callback 5, stdlib.MatchBinary] ==>", V6)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
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
		Value: `
[
	{"co2":10,"hum":30,"lex":22,"temp":100},
	{"co2":100,"hum":300,"lex":220,"temp":1000},
	{"co2":1000,"hum":3000,"lex":2200,"temp":10000}
]
`,
	})
	if err != nil {
		log.Error("grpc.Dial err: %v", err)
	}
	log.Infof("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())
	time.Sleep(1 * time.Second)
	log.Info("Test Http Api===> " + HttpGet("http://127.0.0.1:2580/api/v1/system"))
	engine.Stop()
}

func HttpGet(api string) string {
	var err error
	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Error(err)
		return ""
	}

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(body)
}
