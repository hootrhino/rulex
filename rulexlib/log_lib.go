package rulexlib

import (
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

func Log(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		glogger.LuaLog([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + content + "\n"))
		return 0
	}
}
