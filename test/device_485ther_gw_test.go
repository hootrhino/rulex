package test

import (
	"time"

	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"

	"github.com/hootrhino/rulex/typex"
)

/*
*
* Test 485 sensor gateway
*
 */
func Test_modbus_485_sensor_gateway(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// RTU485_THER Inend
	RTU485Device := typex.NewDevice("RTU485_THER",
		"温湿度采集器", "温湿度采集器", map[string]interface{}{
			"slaverIds": []uint8{1, 2},
			"timeout":   5,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM4",
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 4800,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "node1",
					"function": 3,
					"slaverId": 1,
					"address":  0,
					"quantity": 2,
				},
				{
					"tag":      "node2",
					"function": 3,
					"slaverId": 2,
					"address":  0,
					"quantity": 2,
				},
			},
		})
	RTU485Device.UUID = "RTU485Device1"
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(RTU485Device, ctx, cancelF); err != nil {
		t.Error("RTU485Device load failed:", err)
	}
	mqttOutEnd := typex.NewOutEnd(
		"MQTT",
		"MQTT桥接",
		"MQTT桥接", map[string]interface{}{
			"Host":     "127.0.0.1",
			"Port":     1883,
			"ClientId": "IGW00000001",
			"Username": "IGW00000001",
			"Password": "IGW00000001",
			"PubTopic": "iothub/up/IGW00000001",
			"SubTopic": "iothub/down/IGW00000001",
		},
	)
	mqttOutEnd.UUID = "mqttOutEnd-iothub"
	ctx1, cancelF1 := typex.NewCCTX()
	if err := engine.LoadOutEndWithCtx(mqttOutEnd, ctx1, cancelF1); err != nil {
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
Actions = {function(args)
	for tag, v in pairs(rulexlib:J2T(data)) do
		local ts = rulexlib:TsUnixNano()
		local value = rulexlib:J2T(v['value'])
		value['tag']= tag;
		local jsont = {
			method = 'report',
			requestId = ts,
			timestamp = ts,
			params = value
		}
		print('mqttOutEnd-iothub', rulexlib:T2J(jsont))
		data:ToMqtt('mqttOutEnd-iothub', rulexlib:T2J(jsont))
	end
	return true, args
end}
`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	time.Sleep(25 * time.Second)
	engine.Stop()
}
