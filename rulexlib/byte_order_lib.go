package rulexlib

/*
*
* 字节序处理器
*
 */
import (
	"github.com/i4de/rulex/typex"

	lua "github.com/i4de/gopher-lua"
)

//--------------------------------------------------------------------------------------------------
// 字节序转换 TODO: 目前还没时间实现，等下个任务周期
//--------------------------------------------------------------------------------------------------

/*
*
* 处理ABCD序
*
 */
func ABCD(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// hexs := l.ToString(2)
		return 0
	}
}

/*
*
* 处理DCBA序
*
 */
func DCBA(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// hexs := l.ToString(2)
		return 0
	}
}

/*
*
* 处理DCBA序
*
 */
func BADC(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// hexs := l.ToString(2)
		return 0
	}
}

/*
*
* 处理CDAB序
*
 */
func CDAB(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// hexs := l.ToString(2)
		return 0
	}
}
