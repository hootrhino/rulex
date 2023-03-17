package rulexlib

import (
	"github.com/i4de/rulex/typex"

	lua "github.com/i4de/gopher-lua"
)

/*
*
* 数据转发到 UDP：local err: = rulexlib:DataToUdp(uuid, data)
*
 */
func DataToUdp(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		err := handleDataFormat(rx, id, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}
