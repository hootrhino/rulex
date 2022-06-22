package rulexlib

import (
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

func SelfRuleUUID(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(uuid))
		return 1
	}
}
