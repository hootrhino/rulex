package test

import (
	"time"

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
func Test_modbus_485_yk8(t *testing.T) {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.GlobalConfig.AppDebugMode = true
	glogger.StartGLogger(true, core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// YK8 Inend
	YK8Device := typex.NewDevice(typex.YK08_RELAY,
		"继电器", "继电器", "", map[string]interface{}{
			"timeout":   5,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM3",
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 9600,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "node1",
					"function": 3,
					"slaverId": 3,
					"address":  0,
					"quantity": 1,
				},
			},
		})
	YK8Device.UUID = "YK8Device1"
	if err := engine.LoadDevice(YK8Device); err != nil {
		t.Fatal("YK8Device load failed:", err)
	}

	mqttInend := typex.NewInEnd(typex.MQTT, "MQTT", "MQTT", map[string]interface{}{
		"host":     "broker.emqx.io",
		"port":     1883,
		"clientId": "yk8001",
		"username": "yk8001",
		"password": "yk8001",
		"pubTopic": "$thing/up/property/yk8/yk8001",
		"subTopic": "$thing/down/property/yk8/yk8001",
	})

	if err := engine.LoadInEnd(mqttInend); err != nil {
		t.Fatal("mqttOutEnd load failed:", err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"数据推送至IOTHUB",
		"数据推送至IOTHUB",
		[]string{mqttInend.UUID}, // 数据来自MQTT Server
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
    Actions = {function(data)
			print(data)
			n, err = rulexlib:WriteDevice('YK8Device1', rulexlib:T2J({{
				['function'] = 15,
				['slaverId'] = 3,
				['address'] = 0,
				['quantity'] = 1,
				['value'] = rulexlib:T2Str({0, 0, 0, 0, 1, 1, 1, 1})
			}}))
			if (err) then
				throw()
			end
	return true, data
end}
`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Second)
	engine.Stop()
}
