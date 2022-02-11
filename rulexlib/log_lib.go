package rulexlib

import (
	"rulex/core"
	"rulex/typex"
	"time"

	lua "github.com/yuin/gopher-lua"
)

var LUA_LOGGER *typex.LogWriter

func NewLuaLogger(filepath string, maxSlotCount int) *typex.LogWriter {
	return typex.NewLogWriter(filepath, maxSlotCount)
}

type LogLib struct {
}

func NewLogLib() typex.XLib {

	return &LogLib{}
}
func (l *LogLib) Name() string {
	return "log"
}
func (l *LogLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		LUA_LOGGER.Write([]byte("[CALLBACK]" + content + "\n"))
		return 0
	}
}

/*
*
* StartLuaLogger
*
 */
func StartLuaLogger() {
	LUA_LOGGER = NewLuaLogger("./"+time.Now().Format("2006-01-02_15-04-05-")+core.GlobalConfig.LuaLogPath, 1000)
}
