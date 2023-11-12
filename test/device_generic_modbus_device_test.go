package test

import (
	"context"

	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	mbserver "github.com/tbrandon/mbserver"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

func Test_Generic_modbus_device_tcp_mode(t *testing.T) {
	//
	start_modbus_slaver_emu(t)
	//
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
		t.Fatal(err)
	}
	GMODBUS := typex.NewDevice(typex.GENERIC_MODBUS,
		"GENERIC_MODBUS", "GENERIC_MODBUS", map[string]interface{}{
			"mode": "TCP",
			// "mode":        "UART",
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

/*
*
* 启动一个用来单元测试的Modbus TCP模拟器
*
 */
func start_modbus_slaver_emu(t *testing.T) {
	server := mbserver.NewServer()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			{

			}

		}
		server.ListenTCP("0.0.0.0:1502")
		server.Debug = true
		// 模拟两个数： 37.5, 180.25
		server.HoldingRegisters = []uint16{0x2505, 0xB419} // 模拟数据
		//
		t.Log("Modbus Server started: 0.0.0.0:1502")
	}(context.Background())
}
