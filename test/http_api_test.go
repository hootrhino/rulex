package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rulex/plugin/http_server"
	"rulex/x"

	"syscall"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
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
	//
	// POST
	inendsJson, errs3 := json.Marshal(map[string]interface{}{
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
	})
	if errs3 != nil {
		log.Fatal(errs3)
	}
	//
	r2, errs4 := http.Post("http://127.0.0.1:2580/api/v1/inends",
		"application/json",
		bytes.NewBuffer(inendsJson))
	if errs4 != nil {
		log.Fatal(errs4)
	}
	body2, errs5 := ioutil.ReadAll(r2.Body)
	if errs5 != nil {
		log.Fatal(errs5)
	}
	///////
	in111, errs2 := hh.GetMInEnd(1)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	assert.Equal(t, len(hh.AllMInEnd()), int(1))
	assert.Equal(t, in111.ID, uint(1))
	assert.Equal(t, in111.Type, "MQTT")
	assert.Equal(t, in111.Name, "MQTT test")
	assert.Equal(t, in111.Description, "MQTT Test Resource")
	config := map[string]interface{}{}
	uerror := json.Unmarshal([]byte(in111.Config), &config)
	if uerror != nil {
		log.Fatal(uerror)
	}
	in1 := x.NewInEnd(in111.Type, in111.Name, in111.Description, &config)
	if err0 := engine.LoadInEnd(in1); err0 != nil {
		log.Fatal("InEnd load failed:", err0)

	}
	t.Log(string(body2))
	time.Sleep(2 * time.Second)

	// Get list
	r1, errs1 := http.Get("http://127.0.0.1:2580/api/v1/inends")
	if errs1 != nil {
		log.Fatal(errs1)
	}
	defer r1.Body.Close()
	body1, errs2 := ioutil.ReadAll(r1.Body)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	assert.Equal(t, len(string(body1)) > 100, true)
}
