package rulexlib

/*
*
* 字节序处理器
*
 */
import (
	"encoding/hex"
	"fmt"
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
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
		data := l.ToString(2)
		start := l.ToInt(3)
		end := l.ToInt(4)
		subStr, err := SubStr(data, start, end)
		if err != nil {
			l.Push(lua.LString(""))
			l.Push(lua.LString(err.Error()))
			return 2
		}

		l.Push(lua.LString(ReverseString(subStr)))
		l.Push(lua.LNil)
		return 2
	}
}

// DCBA 将十六进制串里的N个字节"ABCD"逆序成"DCBA"
func DCBA(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		data := l.ToString(2)
		start := l.ToInt(3)
		end := l.ToInt(4)
		subStr, err := SubStr(data, start, end)
		if err != nil {
			l.Push(lua.LString(""))
			l.Push(lua.LString(err.Error()))
			return 2
		}

		l.Push(lua.LString(ReverseString(subStr)))
		l.Push(lua.LNil)
		return 2
	}
}

// BADC 将十六进制串里的N个字节"ABCD"逆序成"BADC"
func BADC(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		data := l.ToString(2)
		start := l.ToInt(3)
		end := l.ToInt(4)
		subStr, err := SubStr(data, start, end)
		if err != nil {
			l.Push(lua.LString(""))
			l.Push(lua.LString(err.Error()))
			return 2
		}
		bytes, _ := hex.DecodeString(ReverseString(subStr))
		l.Push(lua.LString(hex.EncodeToString(ReverseByteOrder(bytes))))
		l.Push(lua.LNil)
		return 2
	}
}

// CDAB 将十六进制串里的2个字节"ABCD"转化成"CDAB"
func CDAB(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		data := l.ToString(2)
		start := l.ToInt(3)
		end := l.ToInt(4)
		subStr, err := SubStr(data, start, end)
		if err != nil {
			l.Push(lua.LString(""))
			l.Push(lua.LString(err.Error()))
			return 2
		}
		bytes, _ := hex.DecodeString(subStr)
		l.Push(lua.LString(hex.EncodeToString(ReverseByteOrder(bytes))))
		l.Push(lua.LNil)
		return 2
	}
}

// SubStr 字符串截取子串，遵循前闭后开原则[start:end)
func SubStr(data string, start, end int) (string, error) {
	if start < 0 || end < 0 || start > end || end > len(data) {
		return "", fmt.Errorf("slice bounds out of range [%d:%d] with length %d", start, end, len(data))
	}
	return data[start:end], nil
}
