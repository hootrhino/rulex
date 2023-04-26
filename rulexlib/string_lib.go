package rulexlib

import (
	"strings"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* Table 转成 String
*
 */
func T2Str(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		table := l.ToTable(2)
		args := []string{}
		table.ForEach(func(l1, value lua.LValue) {
			args = append(args, value.String())
		})
		r := strings.Join(args, "")
		l.Push(lua.LString(r))
		return 1
	}
}

/*
*
* 字节数组转字符串: {1, 2, 3, 4, 5} => "****", 列表必须是合法字节！
*
 */
var result1 = [1024 * 100]byte{}

func Bin2Str(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		table := l.ToTable(2)
		acc := 0
		ok := true
		errMsg := ""
		table.ForEach(func(l1, value lua.LValue) {
			switch t := value.(type) {
			case lua.LNumber:
				result1[acc] = byte(t)
				acc++
			default:
				errMsg = "Bin2Str error:" + t.String()
				glogger.GLogger.Error(errMsg)
				ok = false
				return
			}
		})
		if !ok {
			l.Push(lua.LNil)
			l.Push(lua.LString(errMsg))
		} else {
			l.Push(lua.LString(result1[:acc]))
			l.Push(lua.LNil)
		}
		return 2
	}
}
