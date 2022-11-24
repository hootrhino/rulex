package test

import (
	"time"

	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"

	"github.com/i4de/rulex/typex"
)

func Test_ithings(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}
	ithingsIothub := typex.NewInEnd(typex.TENCENT_IOT_HUB,
		"MQTT", "MQTT", map[string]interface{}{
			"host":       "106.15.225.172",
			"port":       1883,
			"clientId":   "25aZkpfuSfCRULEX-MQTT网关",
			"username":   "25aZkpfuSfCRULEX-MQTT网关;12010126;EN8LX;1680257154",
			"password":   "2ee0a8f7a36619ae28945948e5c60ce4da5468596057af58747274edecc9aa5e;hmacsha256",
			"productId":  "25aZkpfuSfC",
			"deviceName": "RULEX-MQTT网关",
		})
	ithingsIothub.UUID = "ithingsIothub"
	if err := engine.LoadInEnd(ithingsIothub); err != nil {
		t.Fatal("mqttOutEnd load failed:", err)
	}

	callback :=
		`Actions = {
			function(data)
				print("From ithingsIothub===>", data)
				return false, data
			end
		}`
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{ithingsIothub.UUID},
		[]string{},
		`function Success() print("[ithingsIothub Success Callback]=> OK") end`,
		callback,
		`function Failed(error) print("[ithingsIothub Failed Callback]", error) end`)

	if err := engine.LoadRule(rule1); err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Second)
	engine.Stop()
}
