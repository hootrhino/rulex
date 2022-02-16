package rulexlib

import (
	"fmt"
	"rulex/typex"
	"time"

	lua "github.com/yuin/gopher-lua"
)

/*
*
* Unix 时间戳
*
 */
func TsUnix(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(fmt.Sprintf("%v", time.Now().Unix())))
		return 1
	}
}

/*
*
* Unix 纳秒时间戳
*
 */
func TsUnixNano(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(fmt.Sprintf("%v", time.Now().UnixNano())))
		return 1
	}
}

/*
*
* 时间字符串 2006-01-02 15:04:05
*
 */
func Time(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(time.Now().Format("2006-01-02 15:04:05")))
		return 1
	}
}
