package test

import (
	"net/http"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/utils"

	"github.com/hootrhino/rulex/typex"
)

/*
*
* Test_data_to_http
*
 */
func Test_http_source(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// http Inend
	httpInend := typex.NewInEnd(
		"HTTP",
		"Test",
		"Test", map[string]interface{}{
			"port": 8088,
			"host": "127.0.0.1",
		},
	)
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(httpInend, ctx, cancelF); err != nil {
		t.Fatal("httpInend load failed:", err)
	}
	//
	// Load Rule [{"co2":10,"hum":30,"lex":22,"temp":100}]
	//
	callback :=
		`Actions = {
			function(args)
				print("From http===>", data)
				return false, data
			end
		}`
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{httpInend.UUID},
		[]string{},
		`function Success() print("[Test_data_to_http Success Callback]=> OK") end`,
		callback,
		`function Failed(error) print("[Test_data_to_http Failed Callback]", error) end`)

	if err := engine.LoadRule(rule1); err != nil {
		t.Fatal(err)
	}
	res, err := utils.Post(*http.DefaultClient,
		map[string]interface{}{"data": "ok, let's go!"},
		"http://127.0.0.1:8088/in",
		map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
	//
	//
	//
	time.Sleep(3 * time.Second)
	engine.Stop()
}
