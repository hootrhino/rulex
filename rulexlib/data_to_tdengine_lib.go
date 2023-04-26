package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

// 数据推送到Tdengine
func DataToTdEngine(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		//
		// SQL: INSERT INTO meter VALUES (NOW, %v, %v....);
		//
		data := l.ToString(3) // Data must arrays [1,2,3,4....]
		err := handleDataFormat(rx, id, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}
