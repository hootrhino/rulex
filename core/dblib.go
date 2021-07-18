package core

import (
	"encoding/json"
	"rulex/statistics"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// Loader
func LoadDbLib(e *RuleEngine, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"DataToMongo": func(l *lua.LState) int {
			id := l.ToString(1)
			data := l.ToString(2)
			DataToMongo(e, id, data)
			return 0
		},
	})
	vm.Push(mod)
	return 1
}

//
//
//
func DataToMongo(e *RuleEngine, id string, data string) {
	bson := &map[string]interface{}{}
	err := json.Unmarshal([]byte(data), bson)
	if err != nil {
		statistics.IncOutFailed()
		log.Error("Mongo data must be JSON format:", data, " ==> ", err)
	} else {
		statistics.IncOut()
		(*e.OutEnds)[id].Target.To(bson)
	}
}
