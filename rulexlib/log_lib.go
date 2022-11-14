package rulexlib

import (
	"time"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

func Log(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		glogger.LuaLog([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + content + "\n"))
		return 0
	}
}
