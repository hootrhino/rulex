package test

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"

	"github.com/i4de/rulex/typex"
)

/*
*
* Test 485 sensor gateway
*
 */
func Test_modbus_485_sensor_gateway(t *testing.T) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "./rulex.db", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// RTU485_THER Inend
	RTU485Device := typex.NewDevice("RTU485_THER",
		"温湿度采集器", "温湿度采集器", "", map[string]interface{}{
			"slaverIds": []uint8{1, 2},
			"timeout":   5,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM3",
				"baudRate": 4800,
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
			},
		})
	RTU485Device.UUID = "RTU485Device1"
	if err := engine.LoadDevice(RTU485Device); err != nil {
		t.Error("RTU485Device load failed:", err)
	}
	mqttOutEnd := typex.NewOutEnd(
		"MQTT",
		"MQTT桥接",
		"MQTT桥接", map[string]interface{}{
			"Host":      "106.15.225.172",
			"Port":      1883,
			"DataTopic": "iothub/upstream/IGW00000001",
			"ClientId":  "IGW00000001",
			"Username":  "IGW00000001",
			"Password":  "IGW00000001",
		},
	)
	mqttOutEnd.UUID = "mqttOutEnd"
	if err := engine.LoadOutEnd(mqttOutEnd); err != nil {
		t.Error("mqttOutEnd load failed:", err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"数据推送至IOTHUB",
		"数据推送至IOTHUB",
		[]string{},
		[]string{RTU485Device.UUID}, // 数据来自网关设备,所以这里需要配置设备ID
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
				local t = rulexlib:J2T(data)
				t['type'] = 'sub_device'
				t['sn'] = 'IGW00000001'
				local jsons = rulexlib:T2J(t)
				rulexlib:log(jsons)
				rulexlib:DataToMqtt('mqttOutEnd', jsons)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	s := <-c
	glogger.GLogger.Warn("Received stop signal:", s)
	engine.Stop()
	os.Exit(0)
}
