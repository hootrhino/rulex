package test

import (
	"rulenginex/x"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"

	"testing"
)

//
func TestRuleEngingAPI(t *testing.T) {
	testId := x.MakeUUID("INENDS")
	ruleEngine := x.RuleEngine{}
	// Input
	in1 := x.InEnd{
		Id:          testId,
		Type:        "UDP",
		Name:        "UDP Stream",
		Description: "UDP Input Stream",
		Config: &map[string]interface{}{
			"packet_length": 1024,
		},
	}
	ruleEngine.LoadInEnds(&in1)
	// 加载输出端
	out1 := x.OutEnd{
		Id:          testId,
		Type:        "Mysql",
		Name:        "Mysql",
		Description: "Insert to Mysql",
		Config: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     3306,
			"username": "root",
			"password": "root",
		},
	}
	ruleEngine.LoadOutEnds(&out1)
	ruleEngine.Start(func() {

	})
	ruleEngine.Work(&in1, `{"t":100,"h":30}`)
	ruleEngine.Stop()
}

/**
测试计划：
1 规则引擎启动
2 创建一个MQTT输入资源
3 尝试连接MQTT Server
4 尝试订阅Topic
5 进入工作状态
6 MQTT 发送一个数据
7 检查这个MQTT InEnd 有没有关联规则
8 如果有关联，查出这个关联的相关参数，执行回调：actions
9 成功后回调 success；失败后回调 failed
10 统计性能指标
*/
func TestInitMqttResource(t *testing.T) {
	ruleEngine := x.RuleEngine{}
	ruleEngine.Start(func() {

	})
	in1 := x.InEnd{
		Id:          x.MakeUUID("INENDS"),
		Type:        "MQTT",
		Name:        "MQTT Stream",
		Description: "MQTT Input Stream",
		Config: &map[string]interface{}{
			"server":   "127.0.0.1",
			"port":     1883,
			"username": "test",
			"password": "test",
			"clientId": "test",
		},
	}
	ruleEngine.LoadInEnds(&in1)

	// ruleEngine.LoadRules()
	ruleEngine.Work(&in1, `{"t":100,"h":30}`)

	//
	time.Sleep(time.Duration(3) * time.Second)
	publish(t)
	time.Sleep(time.Duration(3) * time.Second)

	//
	ruleEngine.Stop()
}

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
