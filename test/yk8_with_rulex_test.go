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

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
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

	tencentIothub := typex.NewInEnd(typex.TENCENT_IOT_HUB,
		"MQTT", "MQTT", map[string]interface{}{
			"host":       "127.0.0.1",
			"port":       1883,
			"clientId":   "RULEX-001",
			"username":   "RULEX-001",
			"password":   "RULEX-001",
			"productId":  "RULEX-001",
			"deviceName": "RULEX-001",
		})
	tencentIothub.UUID = "tencentIothub"

	if err := engine.LoadInEnd(tencentIothub); err != nil {
		t.Fatal("mqttOutEnd load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"数据推送至IOTHUB",
		"数据推送至IOTHUB",
		[]string{tencentIothub.UUID}, // 数据来自MQTT Server
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(data)
			print(data)
			local dataT, err = rulexlib:J2T(data)
			if dataT['method'] == 'control' then
				local params = dataT['params']
				local cmd = {
					[1] = params['sw8'],
					[2] = params['sw7'],
					[3] = params['sw6'],
					[4] = params['sw5'],
					[5] = params['sw4'],
					[6] = params['sw3'],
					[7] = params['sw2'],
					[8] = params['sw1']
				}
				local n1, err1 = rulexlib:WriteDevice('YK8Device1', rulexlib:T2J({{
					['function'] = 15,
					['slaverId'] = 3,
					['address'] = 0,
					['quantity'] = 1,
					['value'] = rulexlib:T2Str(cmd)
				}}))

				local n2, err2 = rulexlib:WriteSource('tencentIothub', rulexlib:T2J({
					method = 'control_reply',
					clientToken = dataT['clientToken'],
					code = 0,
					status = 'OK'
				}))

			end
			return true, data
		end
	}
`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Second)
	engine.Stop()
}
