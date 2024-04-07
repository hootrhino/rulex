package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
)

// Topic
// app:         app/console/$uuid
// rule:        rule/$uuid
// Test device: device/rule/test/$uuid
// Test inend:  inend/rule/test/$uuid
// Test outend: outend/rule/test/$uuid
/*
*
* APP debug输出, Debug(".....")
*
 */
func DebugAPP(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(1)
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
func DebugRule(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(1)
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "rule/log/" + uuid,
		}).Info(content)
		return 0
	}
}
