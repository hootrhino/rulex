package test

import (
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/i4de/rulex/typex"
)

func Test_TSS200_ReadData(t *testing.T) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
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
		glogger.GLogger.Error(err)
	}
	tss200 := &typex.Device{
		UUID:         "TSS200V02",
		Name:         "TSS200V02",
		Type:         "TSS200V02",
		ActionScript: "{}",
		Description:  "TSS200V02",
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

	if err := engine.LoadDevice(tss200); err != nil {
		t.Log(err)
	}
	time.Sleep(20 * time.Second)
	engine.Stop()
}
