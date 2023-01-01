package test

import (
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/i4de/rulex/typex"
)

func Test_Generic_modbus_device(t *testing.T) {
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
			// "mode": "TCP",
			"mode":      "RTU",
			"autoRequest": true,
			"timeout":     10,
			"frequency":   5,
			"config": map[string]interface{}{
				"uart":     "COM4", // 虚拟串口测试, COM2上连了个MODBUS-POOL测试器
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 4800,
				"host":     "127.0.0.1",
				"port":     1502,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "node1",
					"function": 3,
					"slaverId": 1,
					"address":  0,
					"quantity": 2,
				},
			},
		})

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
			    print("RawData --->",data)
			    local nodeT = rulexlib:J2T(data)
				local dataT = nodeT['node1']
				local rawbin, error = rulexlib:B64S2B(dataT['value'])
				local matchedData = rulexlib:MB(">a1:8 b2:8 c3:8 d4:8", rawbin, false)
				local a1 = rulexlib:B2I64('>', rulexlib:BS2B(matchedData["a1"]))
				local b2 = rulexlib:B2I64('>', rulexlib:BS2B(matchedData["b2"]))
				local c3 = rulexlib:B2I64('>', rulexlib:BS2B(matchedData["c3"]))
				local d4 = rulexlib:B2I64('>', rulexlib:BS2B(matchedData["d4"]))
				print('a1 --> ', matchedData["a1"], ' --> ', a1)
				print('b2 --> ', matchedData["b2"], ' --> ', b2)
				print('c3 --> ', matchedData["c3"], ' --> ', c3)
				print('d4 --> ', matchedData["d4"], ' --> ', d4)
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
