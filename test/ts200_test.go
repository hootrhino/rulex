package test

import (
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexlib"

	"rulex/typex"
	"testing"
	"time"

	"github.com/ngaut/log"
)

func Test_TS200_ReadData(t *testing.T) {
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
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{"TS200"},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				print('data ==> ', data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		log.Error(err)
	}
	ts200 := &typex.Device{
		UUID:         "TS200V02",
		Name:         "TS200V02",
		Type:         "TS200V02",
		ActionScript: "{}",
		Description:  "TS200V02",
		Config: map[string]interface{}{
			"timeout":   5,
			"slaverIds": []int{1},
			"config": map[string]interface{}{
				"uart":     "com10",
				"baudRate": 9600,
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
			},
		},
	}

	if err := engine.LoadDevice(ts200); err != nil {
		t.Log(err)
	}
	time.Sleep(20 * time.Second)
	engine.Stop()
}
