package test

import (
	"time"

	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"

	"github.com/i4de/rulex/typex"
)

/*
*
* Test_IcmpSender_Device
*
 */
func Test_IcmpSender_Device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}

	ICMP_SENDER := typex.NewDevice(typex.ICMP_SENDER,
		"ICMP_SENDER", "ICMP_SENDER", "ICMP_SENDER", map[string]interface{}{
			"autoRequest": true,
			"timeout":     5,
			"frequency":   5,
			"hosts":       []string{"127.0.0.1", "8.8.8.8"},
		})
	ICMP_SENDER.UUID = "ICMP_SENDER1"
	if err := engine.LoadDevice(ICMP_SENDER); err != nil {
		t.Fatal("ICMP_SENDER load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"test",
		"test",
		[]string{},
		[]string{ICMP_SENDER.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(data)
			print(data)
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
