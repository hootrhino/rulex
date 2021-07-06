package x

import (
	"encoding/json"
	"rulex/statistics"

	"github.com/ngaut/log"
	"github.com/yuin/gopher-lua"
)

// Loader
func LoadDbLib(e *RuleEngine, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"dataToMongo": func(l *lua.LState) int {
			id := l.ToString(1)
			data := l.ToString(2)
			toMongo(e, id, data)
			return 0
		},
	})
	vm.Push(mod)
	return 1
}

//
//
//
func toMongo(e *RuleEngine, id string, data interface{}) {
	bsonf := &map[string]interface{}{}
	err := json.Unmarshal([]byte(data.(string)), bsonf)
	if err != nil {
		statistics.IncOutFailed()
		log.Errorf("Mongo data must be JSON format:%#v", data, err)
	} else {
		statistics.IncOut()
		(*e.OutEnds)[id].Target.To(bsonf)
	}
}
