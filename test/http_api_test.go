package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"rulex/plugin/http_server"
	"rulex/x"
	"strings"

	"syscall"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/ngaut/log"
)

func TestHttpAPi(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	engine := x.NewRuleEngine()
	engine.Start()
	////////
	hh := httpserver.NewHttpApiServer(2580, "../plugin/http_server/templates/*")
	if e := engine.LoadPlugin(hh); e != nil {
		log.Fatal("rule load failed:", e)
	}
	hh.Truncate()
	//
	rrr := post(map[string]interface{}{
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
	}, "inends")
	t.Log("Create inend =======> ", rrr)
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
	config := map[string]interface{}{}
	uerror := json.Unmarshal([]byte(mIn_id_1.Config), &config)
	if uerror != nil {
		log.Fatal(uerror)
	}
	//
	newInEnd1 := x.NewInEnd(mIn_id_1.Type, mIn_id_1.Name, mIn_id_1.Description, &config)
	newInEnd1.Id = mIn_id_1.UUID
	if err0 := engine.LoadInEnd(newInEnd1); err0 != nil {
		log.Fatal("InEnd load failed:", err0)
	}
	//
	// Create rule
	//
	rrr22 := post(map[string]interface{}{
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
		"from": newInEnd1.Id,
		"failed": `
		function Failed(error)
		  -- print("[LUA Callback] call failed from lua:", error)
		end`,
		"success": `
		function Success()
		  -- print("[LUA Callback] call success from lua")
		end`,
	}, "rules")
	log.Debug("Create rule ====> ", rrr22)
	//
	//
	//
	for _, mRule := range hh.AllMRules() {
		rule1 := x.NewRule(engine,
			mRule.Name,
			mRule.Description,
			strings.Split(mRule.From, ","),
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		if err := engine.LoadRule(rule1); err != nil {
			log.Fatal("rule load failed:", err)
		}
	}
	//
	//
	//
	assert.Equal(t, len((get("inends"))) > 100, true)
	time.Sleep(2 * time.Second)

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
