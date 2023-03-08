package rulexlib

import (
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/vendor3rd"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* 读GPIO， lua的函数调用应该是这样: eekit:GPIOGet(pin) -> v,error
*
 */
func RASPI4_GPIOGet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToNumber(2)
		v, e := vendor3rd.RASPI4_GPIOGet(int(pin))
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
func RASPI4_GPIOSet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToNumber(2)
		value := l.ToNumber(3)
		_, e := vendor3rd.RASPI4_GPIOSet(int(pin), int(value))
		if e != nil {
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
