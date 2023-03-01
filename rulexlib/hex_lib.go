package rulexlib

import (
	"encoding/hex"

	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* 十六进制字符串转byte数组
*
 */
func Hexs2Bytes(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		hexs := l.ToString(2)
		s, e := hex.DecodeString(hexs)
		if e != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(e.Error()))
		} else {
			table := lua.LTable{}
			for i, v := range s {
				table.RawSetInt(i, lua.LNumber(v))
			}
			l.Push(&table)
			l.Push(lua.LNil)
		}
		return 2
	}
}

/*
*
* byte数组转十六进制字符串
*
 */
func Bytes2Hexs(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		bytes := l.ToString(2)
		l.Push(lua.LString(hex.EncodeToString([]byte(bytes))))
		l.Push(lua.LNil)
		return 2
	}
}
