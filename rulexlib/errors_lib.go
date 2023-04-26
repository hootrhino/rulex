package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

func Throw(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.RaiseError(l.ToString(2))
		return 0
	}
}
