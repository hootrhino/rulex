package x

import (
	"github.com/ngaut/log"
	"github.com/yuin/gopher-lua"
)

// Loader
func LoadDecodeLib(e *RuleEngine,L *lua.LState) int {
	mod := L.SetFuncs(L.G.Global, map[string]lua.LGFunction{
		"LoadDecodeLibOk": func(L *lua.LState) int {
			log.Debug("LoadDecodeLibOk")
			return 0
		},
	})
	L.Push(mod)
	return 1
}
