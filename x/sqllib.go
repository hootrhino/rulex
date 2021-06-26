package x

import (

	// "github.com/marianogappa/sqlparser"
	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// LoadSqlLib
func LoadSqlLib(e *RuleEngine, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"LoadSqlLibOk": func(vm *lua.LState) int {
			log.Debug("LoadSqlLibOk")
			return 0
		},
	})
	vm.Push(mod)
	return 1
}
