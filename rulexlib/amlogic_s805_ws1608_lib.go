package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/vendor3rd"
	"github.com/plgd-dev/kit/v2/strings"
)

/*
*
* 读GPIO， lua的函数调用应该是这样: ws1608:GPIOGet(pin) -> v,error
*
 */
func WKYWS1608_GPIOGet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToString(2)
		if !strings.SliceContains([]string{"red", "green", "blue"}, pin) {
			l.Push(lua.LNumber(0))
			l.Push(lua.LNil)
			return 1
		}
		v, e := vendor3rd.AmlogicWKYS805_RGBGet(pin)
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
* 写GPIO， lua的函数调用应该是这样: ws1608:GPIOSet(pin, v) -> error
*
 */
func WKYWS1608_GPIOSet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		pin := l.ToString(2)
		value := l.ToNumber(3)
		if !strings.SliceContains([]string{"red", "green", "blue"}, pin) {
			l.Push(lua.LNil)
			return 1
		}
		_, e := vendor3rd.AmlogicWKYS805_RGBSet((pin), int(value))
		if e != nil {
			l.Push(lua.LString(e.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
