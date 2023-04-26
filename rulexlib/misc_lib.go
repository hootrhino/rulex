package rulexlib

import (
	"encoding/hex"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* XOR 校验
*
 */
func XOR(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		hexs := l.ToString(2)
		vb, err := hex.DecodeString(hexs)
		vv := l.ToNumber(3)
		if err != nil {
			l.Push(lua.LFalse)
			return 1
		}
		if utils.XOR(vb) == int(vv) {
			l.Push(lua.LTrue)
		} else {
			l.Push(lua.LFalse)
		}
		return 1
	}
}

/*
*
* CRC16校验
*
 */
func CRC16(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		hexs := l.ToString(2)
		vb, err := hex.DecodeString(hexs)
		if err != nil {
			l.Push(lua.LFalse)
			return 1
		}
		vv := l.ToNumber(3)
		if utils.CRC16(vb) == uint16(vv) {
			l.Push(lua.LTrue)
		} else {
			l.Push(lua.LFalse)
		}
		return 1
	}
}
