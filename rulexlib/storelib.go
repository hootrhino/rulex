package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
	"time"
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
			l.Push(lua.LNil)
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

func StoreSetWithDuration(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		k := l.ToString(2)
		v := l.ToString(3)
		d := l.ToInt64(4) // second
		duration := time.Duration(d) * time.Second
		core.GlobalStore.SetWithDuration(k, v, duration)
		return 0
	}
}
