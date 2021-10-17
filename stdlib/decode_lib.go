package stdlib

import (
	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
	"rulex/typex"
)

// Loader
func LoadDecodeLib(e typex.RuleX, L *lua.LState) {
	L.SetFuncs(L.G.Global, map[string]lua.LGFunction{
		"LoadDecodeLibOk": func(L *lua.LState) int {
			log.Debug("LoadDecodeLibOk")
			return 0
		},
	})
}
