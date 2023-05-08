package rulexlib

import (
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

// rulexlib:WriteInStream('INEND', rulexlib:T2J(t))
var sourceReadBuffer []byte = []byte{}

/*
*
* 向资源写入数据
*
 */
func WriteSource(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		data := l.ToString(3)

		if value, ok := rx.AllInEnd().Load(uuid); ok {
			n, err := value.(*typex.InEnd).Source.DownStream([]byte(data))
			if err != nil {
				glogger.GLogger.Error(err)
				l.Push(lua.LNil)
				l.Push(lua.LString(err.Error()))
				return 2
			}
			l.Push(lua.LNumber(n))
			l.Push(lua.LNil)
			return 2
		}
		l.Push(lua.LNil)
		l.Push(lua.LString("source not exists:" + uuid))
		return 2
	}
}

/*
*
* 从资源里面读数据出来
*
 */
func ReadSource(rx typex.RuleX) func(*lua.LState) int {

	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		InEnd := rx.GetInEnd(uuid)
		if InEnd != nil {
			n, err := InEnd.Source.UpStream(sourceReadBuffer)
			if err != nil {
				glogger.GLogger.Error(err)
				l.Push(lua.LNil)
				l.Push(lua.LString(err.Error()))
				return 2
			}
			l.Push(lua.LString(sourceReadBuffer[:n]))
			l.Push(lua.LNil)
			return 2
		}
		l.Push(lua.LNil)
		l.Push(lua.LString("source not exists:" + uuid))
		return 2
	}
}
