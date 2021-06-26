package extralib

import (
	"github.com/ngaut/log"
	"github.com/yuin/gopher-lua"
)

// Loader
func LoadDbLib(vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"dataToMongo": func(l *lua.LState) int {
			id := l.ToString(1)
			data := l.ToString(2)
			toMongo(id, data)
			return 0
		},
	})
	vm.Push(mod)
	return 1
}

//
//
//
func toMongo(id string, data interface{}) {
	log.Debug("dataToMongo:", id, data)
}
