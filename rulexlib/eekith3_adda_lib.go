package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	archsupport "github.com/hootrhino/rulex/bspsupport"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* DI2(0/1)
*
 */
func H3DO1Set(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		value := l.ToNumber(2)
		if value == 0 || value == 1 {
			e := archsupport.EEKIT_GPIOSetDO1((int(value)))
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
func H3DO1Get(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		v, e := archsupport.EEKIT_GPIOGetDO1()
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

func H3DO2Set(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		value := l.ToNumber(2)
		if value == 0 || value == 1 {
			e := archsupport.EEKIT_GPIOSetDO2(int(value))
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
func H3DO2Get(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		v, e := archsupport.EEKIT_GPIOGetDO2()
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
* DI 1,2,3 -> gpio 8-9-10
*
 */
func H3DI1Get(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		Value, e := archsupport.EEKIT_GPIOGetDI1()
		if e != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNumber(Value))
			l.Push(lua.LNil)
		}
		return 2
	}
}
func H3DI2Get(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		Value, e := archsupport.EEKIT_GPIOGetDI2()
		if e != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNumber(Value))
			l.Push(lua.LNil)
		}
		return 2
	}
}
func H3DI3Get(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		v, e := archsupport.EEKIT_GPIOGetDI3()
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
