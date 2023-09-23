package rulexlib

import (
	"encoding/json"
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/component/interqueue"
	"github.com/hootrhino/rulex/typex"
)

func DataToMqtt(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		err := handleDataFormat(rx, id, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}
func DataToMqttTopic(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		topic := l.ToString(3)
		data := l.ToString(4)
		err := handleMqttFormat(rx, id, topic, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

type mqtt_data struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// 处理MQTT消息
// 支持自定义MQTT Topic, 需要在Target的to接口来实现这个
func handleMqttFormat(e typex.RuleX,
	uuid string,
	topic string,
	incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		bytes, _ := json.Marshal(mqtt_data{
			Topic: topic, Payload: incoming,
		})
		return interqueue.DefaultDataCacheQueue.PushOutQueue(outEnd, string(bytes))
	}
	msg := "target not found:" + uuid
	glogger.GLogger.Error(msg)
	return errors.New(msg)

}
