package rulexlib

import (
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
)

/*
*
* Lua的日志打印到文件里面
*
 */
func Log(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		glogger.LuaLog([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + content + "\n"))
		return 0
	}
}

/*
*
* APP debug输出, applib:debug(".....")
*
 */
func DebugAPP(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "app/console/" + uuid,
		}).Info(content)
		return 0
	}
}

/*
*
* 辅助Debug使用, 用来向前端Dashboard打印日志的时候带上ID
*
 */
func Debug(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(1)
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "rule/log/" + uuid,
		}).Debug(content)
		return 0
	}
}
