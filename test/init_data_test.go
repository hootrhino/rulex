package test

import (
	"encoding/json"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/typex"
	"testing"

	"github.com/ngaut/log"
)

//
// 初始化一些测试数据
//
func TestInitData(t *testing.T) {
	engine := engine.NewRuleEngine()
	hh := httpserver.NewHttpApiServer(3580, "plugin/http_server/templates", "rulex-defaule-data.db", engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin(hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	hh.Truncate()
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", &map[string]interface{}{
		"port": "2581",
	})
	b1, _ := json.Marshal(grpcInend.Config)
	hh.InsertMInEnd(&httpserver.MInEnd{
		UUID:        grpcInend.Id,
		Type:        grpcInend.Type.String(),
		Name:        grpcInend.Name,
		Config:      string(b1),
		Description: grpcInend.Description,
	})
	// CoAP Inend
	coapInend := typex.NewInEnd("COAP", "Rulex COAP InEnd", "Rulex COAP InEnd", &map[string]interface{}{
		"port": "2582",
	})
	b2, _ := json.Marshal(coapInend.Config)
	hh.InsertMInEnd(&httpserver.MInEnd{
		UUID:        coapInend.Id,
		Type:        coapInend.Type.String(),
		Name:        coapInend.Name,
		Config:      string(b2),
		Description: coapInend.Description,
	})
	// Http Inend
	httpInend := typex.NewInEnd("HTTP", "Rulex HTTP InEnd", "Rulex HTTP InEnd", &map[string]interface{}{
		"port": "2583",
	})
	b3, _ := json.Marshal(httpInend.Config)
	hh.InsertMInEnd(&httpserver.MInEnd{
		UUID:        httpInend.Id,
		Type:        httpInend.Type.String(),
		Name:        httpInend.Name,
		Config:      string(b3),
		Description: httpInend.Description,
	})

	// Udp Inend
	udpInend := typex.NewInEnd("UDP", "Rulex UDP InEnd", "Rulex UDP InEnd", &map[string]interface{}{
		"port": "2584",
	})
	b4, _ := json.Marshal(udpInend.Config)
	hh.InsertMInEnd(&httpserver.MInEnd{
		UUID:        udpInend.Id,
		Type:        udpInend.Type.String(),
		Name:        udpInend.Name,
		Config:      string(b4),
		Description: udpInend.Description,
	})

	rule := typex.NewRule(engine,
		"Just a test",
		"Just a test",
		[]string{grpcInend.Id},
		`function Success() print("[LUA Success]OK") end`,
		`
			Actions = {
				function(data)
					print("[LUA Actions Callback]", data)
					return true, data
				end
			}`,
		`function Failed(error) print("[LUA Failed]OK", error) end`)
	hh.InsertMRule(&httpserver.MRule{
		Name:        rule.Name,
		Description: rule.Description,
		From:        rule.From,
		Actions:     rule.Actions,
		Success:     rule.Success,
		Failed:      rule.Failed,
	})

}
