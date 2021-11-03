package test

import (
	"os"
	"os/signal"
	"rulex/core"
	"rulex/engine"
	"rulex/typex"
	"syscall"
	"testing"
)

func TestLuaSyntax1(t *testing.T) {
	core.InitGlobalConfig()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine()
	engine.Start()
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": "2581",
	})
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
			end,
			function(data)
			    print("[LUA Actions Callback]==========================> ", data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed]==========================> OK", error) end`,
	)
	err := core.VerifyCallback(rule)
	t.Log(err)
}
func TestLuaSyntax2(t *testing.T) {
	core.InitGlobalConfig()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine()
	engine.Start()
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
		"port": "2581",
	})
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
		    end,
			function(data)
			    print("[LUA Actions Callback]==========================> ", data)
				return true, data
			end,,,,,++1122++33++44
			function(data)
			    print("[LUA Actions Callback]==========================> ", data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed]==========================> OK", error) end`,
	)
	err := core.VerifyCallback(rule)
	t.Log(err)
}
