package rulexlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
func DataToMqttServer(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

type MqttLib struct {
}

func NewMqttLib() typex.XLib {
	return &MqttLib{}
}
func (l *MqttLib) Name() string {
	return "DataToMqttServer"
}
func (l *MqttLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		DataToMqttServer(rx, id, data)
		return 0
	}
}
