package rulexlib

/*
*
* 字节序处理器
*
 */
import (
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

/*
*
* 处理ABCD序
*
 */
func ABCD(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
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
		return 0
	}
}
