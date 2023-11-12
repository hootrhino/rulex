package rulexlib

import (
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
)

/*
*
* APP debug输出, stdlib:Debug(".....")
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
func DebugRule(rx typex.RuleX, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "rule/log/" + uuid,
		}).Info(content)
		return 0
	}
}

/*
*
* Println
*
 */
func Println(rx typex.RuleX) func(*lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		for i := 1; i <= top; i++ {
			fmt.Print(L.ToStringMeta(L.Get(i)).String())
			if i != top {
				fmt.Print("\t")
			}
		}
		fmt.Println("")
		return 0
	}
}
