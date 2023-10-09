package rulexlib

import (
	"reflect"

	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* GRPC 解码
*
 */
func Request(rx typex.RuleX) func(*lua.LState) int {
	return request(rx)
}
func RPCDecode(rx typex.RuleX) func(*lua.LState) int {
	return request(rx)
}

/*
*
* GRPC 编码
*
 */
func RPCEncode(rx typex.RuleX) func(*lua.LState) int {
	return request(rx)
}
func request(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)               // UUID
		data := l.ToString(3)             // Data
		target := rx.GetOutEnd(id).Target // Codec Target
		// 两个返回值
		// () -> data ,err
		if target.Details().Type != typex.GRPC_CODEC_TARGET {
			l.Push(lua.LNil)                                             // Data
			l.Push(lua.LString("Only support 'GRPC_CODEC_TARGET' type")) // Error
			return 2
		}
		r, err := target.To(data)
		if err != nil {
			l.Push(lua.LNil)                 // Data
			l.Push(lua.LString(err.Error())) // Error
			return 2
		}
		switch t := r.(type) {
		case string:
			l.Push(lua.LString(t)) // Data
			l.Push(lua.LNil)       // Error
			return 2
		case []uint8:
			l.Push(lua.LString(t)) // Data
			l.Push(lua.LNil)       // Error
			return 2
		}
		l.Push(lua.LNil)                                                                        // Data
		l.Push(lua.LString("result must string, but current is:" + reflect.TypeOf(r).String())) // Error
		return 2

	}
}
