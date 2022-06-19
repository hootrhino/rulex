package test

import (
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexlib"

	"rulex/typex"
	"testing"
	"time"
)

/*
*
* Test 485 sensor gateway
*
 */
func Test_modbus_485_sensor_gateway(t *testing.T) {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.StartLogWatcher(core.GlobalConfig.LogPath)
	rulexlib.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "../rulex-test_"+time.Now().Format("2006-01-02-15_04_05")+".db", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// RTU485_THER Inend
	RTU485_THERInend := typex.NewDevice("RTU485_THER", "RTU485_THER", "RTU485_THER", "", map[string]interface{}{
		"slaverIds": []int{1},
		"timeout":   5,
		"frequency": 5,
		"config": map[string]interface{}{
			"uart":     "COM3",
			"baudRate": 115200,
			"dataBits": 8,
			"parity":   "N",
			"stopBits": 1,
		},
	})

	if err := engine.LoadDevice(RTU485_THERInend); err != nil {
		t.Error("grpcInend load failed:", err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{RTU485_THERInend.UUID}, // 数据来自设备
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				local table = rulexlib:J2T(data)
				local value = table['value']
				local t = rulexlib:HsubToN(value, 5, 8)
				local h = rulexlib:HsubToN(value, 0, 4)
				local t1 = rulexlib:HToN(string.sub(value, 5, 8))
				local h2 = rulexlib:HToN(string.sub(value, 0, 4))
				print('Data ========> ', rulexlib:T2J({
					Device = "TH00000001",
					Ts = rulexlib:TsUnix(),
					T = t,
					H = h,
					T1 = t1,
					H2 = h2
				}))
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	time.Sleep(3 * time.Second)
	engine.Stop()
}
