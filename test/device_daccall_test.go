package test

import (
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/i4de/rulex/typex"
)

func Test_dac_call_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
		t.Fatal(err)
	}
	GMODBUS := typex.NewDevice(typex.GENERIC_MODBUS,
		"GENERIC_MODBUS", "GENERIC_MODBUS", "", map[string]interface{}{
			"mode": "TCP",
			// "mode":      "RTU",
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
	if err := engine.LoadDevice(GMODBUS); err != nil {
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
			function(data)
				print("data=", data)
				local r, error = device:DCACall("GMODBUS", "get_status", {1, 2, 3})
				print("r=", r, ", error=", error)
				return true, data
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
