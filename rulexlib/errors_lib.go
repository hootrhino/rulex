package rulexlib

import (
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

func Throw(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.RaiseError(l.ToString(2))
		return 0
	}
}
