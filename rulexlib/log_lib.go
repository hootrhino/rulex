package rulexlib

import (
	"rulex/typex"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func NewLuaLogger(filepath string, maxSlotCount int) *typex.LogWriter {
	return typex.NewLogWriter(filepath, maxSlotCount)
}

func Log(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		typex.LUA_LOGGER.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + content + "\n"))
		return 0
	}
}

/*
*
* StartLuaLogger
*
 */
func StartLuaLogger(path string) {
	typex.LUA_LOGGER = NewLuaLogger("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
}
