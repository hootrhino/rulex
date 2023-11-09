package test

import (
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

func Test_Generic_snmp_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
		t.Fatal(err)
	}
	GENERIC_SNMP := typex.NewDevice(typex.GENERIC_SNMP,
		"GENERIC_SNMP", "GENERIC_SNMP", map[string]interface{}{
			"timeout":   10,
			"frequency": 5,
			"target":    "127.0.0.1",
			"port":      161,
			"community": "public",
			"transport": "udp",
			"version":   3,
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(GENERIC_SNMP, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{GENERIC_SNMP.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
			    print(data)
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
