package stdlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
// Loader
//
func LoadStreamLib(e typex.RuleX, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"WriteInStream": func(l *lua.LState) int {
			id := l.ToString(2)
			data := l.ToString(3)
			WriteInStream(e, id, data)
			return 0
		},
		"WriteOutStream": func(l *lua.LState) int {
			id := l.ToString(2)
			data := l.ToString(3)
			WriteOutStream(e, id, data)
			return 0
		},
	})
	vm.Push(mod)
	return 1
}
func WriteInStream(e typex.RuleX, id string, data string) {
	e.GetInEnd(id).Resource.OnStreamApproached(data)
}
func WriteOutStream(e typex.RuleX, id string, data string) {
	e.GetOutEnd(id).Target.OnStreamApproached(data)
}
