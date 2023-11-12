package test

import (
	"testing"
	"time"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"
)

func Test_Generic_opcua_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
		t.Fatal(err)
	}
	GENERIC_OPCUA := typex.NewDevice(typex.GENERIC_OPCUA,
		"GENERIC_OPCUA", "GENERIC_OPCUA", map[string]interface{}{
			"commonConfig": map[string]interface{}{
				"endpoint":  "opc.tcp://NOAH:53530/OPCUA/SimulationServer",
				"policy":    device.POLICY_BASIC128RSA15,
				"mode":      device.MODE_SIGN,
				"auth":      device.AUTH_ANONYMOUS,
				"username":  "1",
				"password":  "1",
				"timeout":   10,
				"frequency": 500,
				"retryTime": 10,
			},
			"opcuaNodes": []map[string]interface{}{
				{
					"tag":         "node1",
					"description": "node 1",
					"nodeId":      "ns=3;i=1013",
					"dataType":    "String",
					"value":       "",
				},
				{
					"tag":         "node2",
					"description": "node 2",
					"nodeId":      "ns=3;i=1001",
					"dataType":    "String",
					"value":       "",
				},
			},
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(GENERIC_OPCUA, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{GENERIC_OPCUA.UUID},
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

	time.Sleep(1 * time.Second)
	engine.Stop()
}
