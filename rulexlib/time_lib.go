package rulexlib

import (
	"fmt"
	"time"

	"github.com/hootrhino/rulex/typex"

	"github.com/wwhai/ntp"

	lua "github.com/hootrhino/gopher-lua"
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
* 等待毫秒
*
 */
func Sleep(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		ts := l.ToNumber(2)
		time.Sleep(time.Millisecond * time.Duration(ts))
		return 0
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
func TimeMs(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		l.Push(lua.LString(time.Now().Format("2006-01-02 15:04:05.000")))
		return 1
	}
}

/*
*
* NTP Server Time
* return: ntp time string, error
*
 */
func NtpTime(rx typex.RuleX) func(l *lua.LState) int {

	return func(l *lua.LState) int {
		// Ntp server:
		//   0.cn.pool.ntp.org
		//   1.cn.pool.ntp.org
		//   2.cn.pool.ntp.org
		//   3.cn.pool.ntp.org
		time, err := ntp.Time("0.cn.pool.ntp.org")
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LString(time.String()))
			l.Push(lua.LNil)
		}
		return 2

	}
}
