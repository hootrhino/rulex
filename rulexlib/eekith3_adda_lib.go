package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/vendor3rd"
)

/*
*
* 读GPIO， lua的函数调用应该是这样: eekith3:GPIOGet(pin) -> v,error
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
* 写GPIO， lua的函数调用应该是这样: eekith3:GPIOSet(pin, v) -> error
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

/*
*
*  DI1(0/1)
*
 */
func H3DO1Set(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		value := l.ToNumber(2)
		if value == 0 || value == 1 {
			_, e := vendor3rd.EEKIT_GPIOSet(6, int(value))
			if e != nil {
				l.Push(lua.LString(e.Error()))
			} else {
				l.Push(lua.LNil)
			}
		} else {
			l.Push(lua.LString("DO1 Only can set '0' or '1'."))
		}
		return 1
	}
}

/*
*
* DI2(0/1)
*
 */
func H3DO2Set(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		value := l.ToNumber(2)
		if value == 0 || value == 1 {
			_, e := vendor3rd.EEKIT_GPIOSet(7, int(value))
			if e != nil {
				l.Push(lua.LString(e.Error()))
			} else {
				l.Push(lua.LNil)
			}
		} else {
			l.Push(lua.LString("DO2 Only can set '0' or '1'."))
		}
		return 1
	}
}
