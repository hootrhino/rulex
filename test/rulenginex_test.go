package test

import (
	"rulenginex/x"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"

	"testing"
)

func publish(t *testing.T) {

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		t.Log("Mqtt InEnd Connected Success")
		client.Publish("$x_IN_END", 1, false, "test is ok1")
		client.Publish("$x_IN_END", 1, false, "test is ok2")
		client.Publish("$x_IN_END", 1, false, "test is ok3")
		t.Log("Publish to x_IN_END ok")

	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		t.Log("Connect lost:", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	//
	opts.SetClientID("x_IN_END_TEST1")
	opts.SetUsername("x_IN_END_TEST1")
	opts.SetPassword("x_IN_END_TEST1")
	//
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Log("error:", token.Error())
	}
}

func TestRunCallback(t *testing.T) {
	id := x.MakeUUID("RULE")
	actions := `
	    local json = require("json")
		Actions = {
			function(data)
				LoadEncodeLibOk()
				LoadDecodeLibOk()
				LoadSqlLibOk()
				print("[LUA Actions Callback] data:", data)
				print("json.encode ====>", json.encode({"a",1,"b",2,"c",3}))
				print("json.decode ====>", json.decode(data))
				return true, data .. " A1"
			end
		}`
	from := []string{"mqtt1", "mqtt2"}
	failed := `function Failed(error) print("[LUA Callback] call failed from lua:", error) end`
	success := `function Success() print("[LUA Callback] call success from lua") end`
	rule1 := x.NewRule(id, "just_a_test_rule", "just_a_test_rule", from, success, actions, failed)

	x.SaveRule(rule1)
	rule2 := x.GetRule(id)
	assert.Equal(t, rule2.Id, id)
	assert.Equal(t, rule2.Name, "just_a_test_rule")
	assert.Equal(t, rule2.Description, "just_a_test_rule")
	assert.Equal(t, rule2.From[0], from[0])
	assert.Equal(t, rule2.From[1], from[1])
	assert.Equal(t, rule2.Actions[0], actions[0])
	assert.Equal(t, rule2.Actions[1], actions[1])
	//
	err0 := x.VerifyCallback(rule2)
	if err0 != nil {
		t.Error("VerifySyntax Failed:", err0)
	} else {
		t.Log("VerifyCallback Success")
	}
	t1 := time.Now().UnixNano()
	for i := 0; i < 1; i++ {
		rule1.ExecuteActions(lua.LString(`{"t":100,"h":30}`))
	}
	t2 := time.Now().UnixNano()
	t.Log("ExecuteActions time:", (t2-t1)/int64(time.Millisecond), "ms")
}
