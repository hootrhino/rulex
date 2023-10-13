package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
* 注意：该接口是通用的，如果觉得这个不清晰，可以尝试具体的'DataToXXX'
* 数据转发到具体的目的地：local err: = data:ToTarget(uuid, data)
*
 */
func DataToTarget(rx typex.RuleX) func(*lua.LState) int {
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
