package test

import (
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

func Test_TSS200_ReadData(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}

	tss200 := typex.NewDevice(typex.TSS200V02,
		"TSS200V02", "TSS200V02", map[string]interface{}{
			"mode":      "UART",
			"timeout":   10,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM2", // 虚拟串口测试, COM2上连了个MODBUS-POOL测试器
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 9600,
				"ip":       "127.0.0.1",
				"port":     502,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "node1",
					"function": 3,
					"slaverId": 1,
					"address":  17,
					"quantity": 9,
				},
			},
		})
	tss200.UUID = "TSS200V02"
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(tss200, ctx, cancelF); err != nil {
		t.Log(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{"TSS200V02"},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
				print('data ==> ', data)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
	}
	time.Sleep(20 * time.Second)
	engine.Stop()
}
