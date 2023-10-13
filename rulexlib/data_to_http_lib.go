package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* 数据转发到HTTP：local err: = data:ToHttp(uuid, data)
*
 */
func DataToHttp(rx typex.RuleX) func(*lua.LState) int {
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
