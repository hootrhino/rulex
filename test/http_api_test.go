package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"rulex/plugin/http_server"
	"rulex/x"

	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/ngaut/log"
)

//
// TestHttpAPi
//
func TestHttpAPi(t *testing.T) {

	engine := x.NewRuleEngine()
	engine.Start()
	////////
	hh := httpserver.NewHttpApiServer(2580, "../plugin/http_server/templates/*", engine)
	if e := engine.LoadPlugin(hh); e != nil {
		log.Fatal("rule load failed:", e)
	}
	hh.Truncate()
	//---------------------------------------------------------------------------------------
	//
	//---------------------------------------------------------------------------------------
	log.Debug(
		post(map[string]interface{}{
			"type":        "MQTT",
			"name":        "MQTT test",
			"description": "MQTT Test Resource",
			"config": map[string]interface{}{
				"server":   "127.0.0.1",
				"port":     1883,
				"username": "test",
				"password": "test",
				"clientId": "test",
			},
		}, "inends"),
	)
	///////
	mIn_id_1, errs2 := hh.GetMInEnd(1)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	assert.Equal(t, len(hh.AllMInEnd()), int(1))
	assert.Equal(t, mIn_id_1.ID, uint(1))
	assert.Equal(t, mIn_id_1.Type, "MQTT")
	assert.Equal(t, mIn_id_1.Name, "MQTT test")
	assert.Equal(t, mIn_id_1.Description, "MQTT Test Resource")

	//---------------------------------------------------------------------------------------
	// Create outend
	//---------------------------------------------------------------------------------------
	log.Debug(
		post(map[string]interface{}{
			"type":        "mongo",
			"name":        "data to mongo",
			"description": "data to mongo",
			"config": map[string]interface{}{
				"mongourl": "mongodb+srv://rulenginex:rulenginex@cluster0.rsdmb.mongodb.net/rulex_test_db?retryWrites=true&w=majority",
			},
		}, "outends"),
	)

	m_Out_id_1, errs2 := hh.GetMOutEnd(1)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	assert.Equal(t, len(hh.AllMInEnd()), int(1))
	assert.Equal(t, m_Out_id_1.ID, uint(1))
	assert.Equal(t, m_Out_id_1.ID, uint(1))
	assert.Equal(t, m_Out_id_1.Type, "mongo")
	assert.Equal(t, m_Out_id_1.Name, "data to mongo")
	assert.Equal(t, m_Out_id_1.Description, "data to mongo")

	//
	// Create rule
	//
	log.Debug(
		post(map[string]interface{}{
			"name":        "just_a_test",
			"description": "just_a_test",
			"actions": `
		local json = require("json")
		Actions = {
			function(data)
				dataToMongo("MongoDB001", data)
				print("[LUA Actions Callback]:dataToMongo Mqtt payload:", data)
				return true, data
			end
		}`,
			"from": m_Out_id_1.UUID,
			"failed": `
		function Failed(error)
		  -- print("[LUA Callback] call failed from lua:", error)
		end`,
			"success": `
		function Success()
		  -- print("[LUA Callback] call success from lua")
		end`,
		}, "rules"),
	)
	//
	assert.Equal(t, len((get("inends"))) > 100, true)
	time.Sleep(3 * time.Second)

}

func post(data map[string]interface{}, api string) string {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Fatal(errs1)
	}
	r, errs2 := http.Post("http://127.0.0.1:2580/api/v1/"+api,
		"application/json",
		bytes.NewBuffer(p))
	if errs2 != nil {
		log.Fatal(errs2)
	}
	defer r.Body.Close()

	body, errs5 := ioutil.ReadAll(r.Body)
	if errs5 != nil {
		log.Fatal(errs5)
	}
	return string(body)
}
func get(api string) string {
	// Get list
	r, errs := http.Get("http://127.0.0.1:2580/api/v1/" + api)
	if errs != nil {
		log.Fatal(errs)
	}
	defer r.Body.Close()
	body, errs2 := ioutil.ReadAll(r.Body)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	return string(body)
}
