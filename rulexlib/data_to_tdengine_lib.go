package rulexlib

import (
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
func DataToTdEngine(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		//
		// SQL: INSERT INTO meter VALUES (NOW, %v, %v....);
		//
		data := l.ToString(3) // Data must arrays [1,2,3,4....]
		handleDataFormat(rx, id, data)
		return 0
	}
}
