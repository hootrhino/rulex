package engine

import (
	"encoding/json"
	"os"
	"os/signal"
	httpserver "rulex/plugin/http_server"
	"rulex/typex"
	"strings"
	"syscall"

	"github.com/ngaut/log"
)

//
//
//
func RunRulex(dbPath string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := NewRuleEngine()
	engine.Start()
	hh := httpserver.NewHttpApiServer(2580, "plugin/http_server/templates", dbPath, engine)
	engine.LoadPlugin(hh)

	//
	// Load inend from sqlite
	//
	for _, minEnd := range hh.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Error(err)
		}
		in1 := typex.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, &config)
		// Important !!!!!!!!
		in1.Id = minEnd.UUID
		if err := engine.LoadInEnd(in1); err != nil {
			log.Error("InEnd load failed:", err)
		}
	}

	//
	// Load rule from sqlite
	//
	for _, mRule := range hh.AllMRules() {
		rule := typex.NewRule(engine,
			mRule.Name,
			mRule.Description,
			strings.Split(mRule.From, ","),
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		if err := engine.LoadRule(rule); err != nil {
			log.Error(err)
		}
	}
	//
	// Load out end from sqlite
	//
	for _, mOutEnd := range hh.AllMOutEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
			log.Error(err)
		}
		newOutEnd := typex.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, &config)
		// Important !!!!!!!!
		newOutEnd.Id = mOutEnd.UUID
		if err := engine.LoadOutEnd(newOutEnd); err != nil {
			log.Error("OutEnd load failed:", err)
		}
	}
	signal := <-c
	log.Info("Received stop signal:", signal)
	engine.Stop()
	os.Exit(0)
}

//
//
//
func InitData() {
	engine := NewRuleEngine()
	hh := httpserver.NewHttpApiServer(3580, "plugin/http_server/templates", "rulex.db", engine)
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
		From:        rule.From[0],
		Actions:     rule.Actions,
		Success:     rule.Success,
		Failed:      rule.Failed,
	})

}
