package x

import (
	"github.com/ngaut/log"
	"github.com/yuin/gopher-lua"
)

// Loader
func LoadEncodeLib(e *RuleEngine, L *lua.LState) int {
	mod := L.SetFuncs(L.G.Global, map[string]lua.LGFunction{
		"LoadEncodeLibOk": func(L *lua.LState) int {
			log.Debug("LoadEncodeLibOk")
			return 0
		},
	})
	L.Push(mod)
	return 1
}
