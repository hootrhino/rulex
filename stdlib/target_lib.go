package stdlib

import (
	"encoding/json"
	"rulex/core"
	"rulex/statistics"
	"rulex/typex"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// Loader
func LoadTargetLib(e typex.RuleX, vm *lua.LState) int {
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
func handleDataFormat(e typex.RuleX, id string, incoming string) {
	data := &map[string]interface{}{}
	err := json.Unmarshal([]byte(incoming), data)
	if err != nil {
		statistics.IncOutFailed()
		log.Error("Data must be JSON format:", incoming, ", But current is: ", err)
	} else {
		statistics.IncOut()
		core.DefaultDataCacheQueue.Push(typex.QueueData{
			In:   nil,
			Out:  &(*e.AllOutEnd()[id]),
			E:    e,
			Data: incoming,
		})

	}
}

//
//
//
func DataToMongo(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

//
//
//
func DataToHttpServer(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

//
//
//
func DataToMqttBroker(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}
