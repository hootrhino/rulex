package rulexlib

import (
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

/*
*
* 读: rulexlib:ReadDevice(ID) -> data, err
* 写: rulexlib:WriteDevice(ID, []byte{}) -> data, err
*
 */
var __readBuffer []byte = []byte{}

func ReadDevice(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			n, err := Device.Device.OnRead(__readBuffer)
			if err != nil {
				glogger.GLogger.Error(err)
				l.Push(lua.LNil)
				l.Push(lua.LString(err.Error()))
				return 2
			}
			l.Push(lua.LString(__readBuffer[:n]))
			l.Push(lua.LNil)
			return 2
		}
		l.Push(lua.LNil)
		l.Push(lua.LString("device not exists:" + devUUID))
		return 2
	}
}

/*
*
* 写数据
*
 */
func WriteDevice(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		data := l.ToString(3)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			n, err := Device.Device.OnWrite([]byte(data))
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
		l.Push(lua.LString("device not exists:" + devUUID))
		return 0
	}
}
