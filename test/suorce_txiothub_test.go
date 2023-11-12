package test

import (
	"time"

	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"

	"github.com/hootrhino/rulex/typex"
)

/*
*
* 腾讯云iothub测试
*
 */
func Test_tencent_iothub(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer(engine)); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}

	// YK8 Inend
	YK8Device := typex.NewDevice(typex.YK08_RELAY,
		"继电器", "继电器", map[string]interface{}{
			"timeout":   5,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM10",
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
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadDeviceWithCtx(YK8Device, ctx, cancelF); err != nil {
		t.Fatal("YK8Device load failed:", err)
	}

	tencentIothub := typex.NewInEnd(typex.GENERIC_IOT_HUB,
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
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(tencentIothub, ctx1, cancelF1); err != nil {
		t.Fatal("mqttOutEnd load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"数据推送至IOTHUB",
		"数据推送至IOTHUB",
		[]string{tencentIothub.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(args)
		    print('From IotHUB:', data)
			-- property
		    iothub:PropertySuccess('tencentIothub', 'test-is-lua-001')
		    iothub:PropertyFailed('tencentIothub', 'test-is-lua-001')
			-- action
		    iothub:ActionSuccess('tencentIothub', 'test-is-lua-001', rulexlib:T2J({sw1=0}))
		    iothub:ActionFailed('tencentIothub', 'test-is-lua-001', rulexlib:T2J({sw1=0}))
			return true, args
		end
	}
`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Second)
	engine.Stop()
}
