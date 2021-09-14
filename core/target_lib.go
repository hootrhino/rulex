package core

import (
	"encoding/json"
	"rulex/statistics"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// Loader
func LoadTargetLib(e RuleX, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"DataToMongo": func(l *lua.LState) int {
			id := l.ToString(1)
			data := l.ToString(2)
			DataToMongo(e, id, data)
			return 0
		},
		"DataToHttpServer": func(l *lua.LState) int {
			id := l.ToString(1)
			data := l.ToString(2)
			DataToHttpServer(e, id, data)
			return 0
		},
	})
	vm.Push(mod)
	return 1
}
func handleDataFormat(e RuleX, id string, data string) {
	bson := &map[string]interface{}{}
	err := json.Unmarshal([]byte(data), bson)
	if err != nil {
		statistics.IncOutFailed()
		log.Error("Data must be JSON format:", data, ", But current is: ", err)
	} else {
		statistics.IncOut()
		(*e.AllOutEnd()[id]).Target.To(bson)
	}
}

//
//
//
func DataToMongo(e RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

//
//
//
func DataToHttpServer(e RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

//
//
//
func DataToMqttBroker(e RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}
