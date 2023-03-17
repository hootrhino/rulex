package rulexlib

import (
	lua "github.com/i4de/gopher-lua"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/vendor3rd"
)

/*
*
* 读GPIO， lua的函数调用应该是这样: eekit:GPIOGet(pin) -> v,error
*
 */
func EEKIT_GPIOGet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToNumber(2)
		v, e := vendor3rd.EEKIT_GPIOGet(int(pin))
		if e != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNumber(v))
			l.Push(lua.LNil)
		}
		return 2
	}
}

/*
*
* 写GPIO， lua的函数调用应该是这样: eekit:GPIOSet(pin, v) -> error
*
 */
func EEKIT_GPIOSet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToNumber(2)
		value := l.ToNumber(3)
		_, e := vendor3rd.EEKIT_GPIOSet(int(pin), int(value))
		if e != nil {
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
