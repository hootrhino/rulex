package main

import (
	"encoding/json"
	"github.com/ngaut/log"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"rulex/core"
	"rulex/plugin/demo_plugin"
	httpserver "rulex/plugin/http_server"
	"strings"
	"syscall"
)

//
func main() {
	core.InitGlobalConfig()
	app := &cli.App{
		Name:  "RULEX, a lightweight iot data rule gateway",
		Usage: "http://rulex.ezlinker.cn",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run rulex immediately",
				Action: func(c *cli.Context) error {
					Run()
					log.Debug("Run rulex successfully.")
					return nil
				},
			},
			{
				Name:  "install",
				Usage: "Install rulex to your path",
				Action: func(c *cli.Context) error {
					log.Debug("Install to: /usr/bin/rule core.")
					os.Exit(0)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

//
func Run() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := core.NewRuleEngine()
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
	grpcInend := core.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", &map[string]interface{}{
		"port": "2581",
	})
	if err := engine.LoadInEnd(grpcInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	// CoAP Inend
	coapInend := core.NewInEnd("COAP", "Rulex COAP InEnd", "Rulex COAP InEnd", &map[string]interface{}{
		"port": "2582",
	})
	if err := engine.LoadInEnd(coapInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	// Http Inend
	httpInend := core.NewInEnd("HTTP", "Rulex HTTP InEnd", "Rulex HTTP InEnd", &map[string]interface{}{
		"port": "2583",
	})
	if err := engine.LoadInEnd(httpInend); err != nil {
		log.Error("Rule load failed:", err)
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range hh.AllMInEnd() {
		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(minEnd.Config), &config); err != nil {
			log.Error(err)
		}
		in1 := core.NewInEnd(minEnd.Type, minEnd.Name, minEnd.Description, &config)
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
		rule := core.NewRule(engine,
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
		newOutEnd := core.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, &config)
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
