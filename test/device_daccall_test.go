package test

import (
	"testing"
	"time"

	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	mqttserver "github.com/hootrhino/rulex/plugin/mqtt_server"
	"github.com/hootrhino/rulex/typex"
)

func Test_dac_call_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("NewHttpApiServer load failed:", err)
		t.Fatal(err)
	}
	qq := mqttserver.NewMqttServer()
	if err := engine.LoadPlugin("plugin.mqtt_server", qq); err != nil {
		glogger.GLogger.Fatal("NewMqttServer load failed:", err)
		t.Fatal(err)
	}
	GMODBUS := typex.NewDevice(typex.GENERIC_MODBUS,
		"GENERIC_MODBUS", "GENERIC_MODBUS", map[string]interface{}{
			"mode": "TCP",
			// "mode":      "UART",
			"timeout":   10,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM4", // 虚拟串口测试, COM2上连了个MODBUS-POOL测试器
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 4800,
				"host":     "127.0.0.1",
				"port":     502,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "node1",
					"function": 3,
					"slaverId": 1,
					"address":  0,
					"quantity": 1,
				},
			},
		})
	GMODBUS.UUID = "GMODBUS"
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(GMODBUS, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{GMODBUS.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
				print("data=", data)
				local r, error = device:DCACall("GMODBUS", "get_status", {1, 2, 3})
				print("r=", r, ", error=", error)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
		t.Fatal(err)
	}

	time.Sleep(25 * time.Second)
	engine.Stop()
}
