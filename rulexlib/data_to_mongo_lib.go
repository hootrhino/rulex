package rulexlib

import (
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
func DataToMongo(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		handleDataFormat(rx, id, data)
		return 0
	}
}
