package test

import (
	"testing"
	"time"

	httpserver "github.com/i4de/rulex/plugin/http_server"
	"github.com/i4de/rulex/typex"
)

/*
*
* Test_UART_Device
*
 */
func Test_g7776_Device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}

	// YK8 Inend
	GUART := typex.NewDevice(typex.USER_G776,
		"UART", "UART", "UART", map[string]interface{}{
			"baudRate":  9600,
			"dataBits":  8,
			"frequency": 5,
			"parity":    "N",
			"stopBits":  1,
			"tag":       "tag1",
			"timeout":   5,
			"uart":      "COM2",
		})
	GUART.UUID = "GUART1"
	if err := engine.LoadDevice(GUART); err != nil {
		t.Fatal("GUART load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"test",
		"test",
		[]string{},
		[]string{GUART.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(data)
			print("received: ",data)
			local n1, err = rulexlib:WriteDevice("GUART1", 0, "GUART1 data")
			print("write size: ",n1, "error: ",err)
			return true, data
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
