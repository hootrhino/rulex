package extralib

import (

	// "github.com/marianogappa/sqlparser"
	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// LoadSqlLib
func LoadSqlLib(L *lua.LState) int {
	mod := L.SetFuncs(L.G.Global, map[string]lua.LGFunction{
		"LoadSqlLibOk": func(L *lua.LState) int {
			log.Debug("LoadSqlLibOk")
			return 0
		},
	})
	L.Push(mod)
	return 1
}
