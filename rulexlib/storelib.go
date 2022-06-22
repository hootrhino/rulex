package rulexlib

import (
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

func StoreSet(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		k := l.ToString(2)
		v := l.ToString(3)
		core.GlobalStore.Set(k, v)
		return 0
	}
}
func StoreGet(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		k := l.ToString(2)
		v := core.GlobalStore.Get(k)
		if v == "" {
			l.Push(nil)
		} else {
			l.Push(lua.LString(v))
		}
		return 1
	}

}
func StoreDelete(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		k := l.ToString(2)
		core.GlobalStore.Delete(k)
		return 0
	}
}
