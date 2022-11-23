package rulexlib

import (
	"github.com/i4de/rulex/common"
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

var deviceReadBuffer []byte = make([]byte, common.T_4KB)

func ReadDevice(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// read(uuid,cmd)
		devUUID := l.ToString(2)
		cmd := l.ToInt(3)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			n, err := Device.Device.OnRead(cmd, deviceReadBuffer)
			if err != nil {
				glogger.GLogger.Error(err)
				l.Push(lua.LNil)
				l.Push(lua.LString(err.Error()))
				return 2
			}
			l.Push(lua.LString(deviceReadBuffer[:n]))
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
		// write(uuid,cmd,data)
		devUUID := l.ToString(2)
		cmd := l.ToInt(3)
		data := l.ToString(4)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			n, err := Device.Device.OnWrite(cmd, []byte(data))
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
		return 2
	}
}
