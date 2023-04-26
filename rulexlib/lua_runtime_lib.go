package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* 获取当前的规则UUID
*
 */
func SelfRuleUUID(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(uuid))
		return 1
	}
}
