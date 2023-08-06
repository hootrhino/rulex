package rulexlib

import (
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* 读: rulexlib:ReadDevice(ID, cmd, buffer) -> data, err
* 写: rulexlib:WriteDevice(ID, cmd, []byte{}) -> data, err
*
 */

var deviceReadBuffer []byte = make([]byte, common.T_4KB)

func ReadDevice(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		cmd := l.ToString(3)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Device.Status() == typex.DEV_UP {
				n, err := Device.Device.OnRead([]byte(cmd), deviceReadBuffer)
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LNil)
					l.Push(lua.LString(err.Error()))
					return 2
				} else {
					l.Push(lua.LString(deviceReadBuffer[:n]))
					l.Push(lua.LNil)
					return 2
				}
			} else {
				l.Push(lua.LNil)
				l.Push(lua.LString("device down:" + devUUID))
				return 2
			}
		} else {
			l.Push(lua.LNil)
			l.Push(lua.LString("device not exists:" + devUUID))
			return 2
		}

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
		cmd := l.ToString(3)
		data := l.ToString(4)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Device.Status() == typex.DEV_UP {
				n, err := Device.Device.OnWrite([]byte(cmd), []byte(data))
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LNil)
					l.Push(lua.LString(err.Error()))
					return 2
				} else {
					l.Push(lua.LNumber(n))
					l.Push(lua.LNil)
					return 2
				}
			} else {
				l.Push(lua.LNil)
				l.Push(lua.LString("device down:" + devUUID))
				return 2
			}
		}
		l.Push(lua.LNil)
		l.Push(lua.LString("device not exists:" + devUUID))
		return 2
	}
}

/*
*
* 控制操作[2023年4月16日新增, 需求来自总线控制多个不对等设备]
*
 */
func CtrlDevice(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// write(uuid,cmd,data)
		devUUID := l.ToString(2)
		cmd := l.ToString(3)
		data := l.ToString(4)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Device.Status() == typex.DEV_UP {
				result, err := Device.Device.OnCtrl([]byte(cmd), []byte(data))
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LNil)
					l.Push(lua.LString(err.Error()))
					return 2
				} else {
					l.Push(lua.LString(string(result)))
					l.Push(lua.LNil)
					return 2
				}
			} else {
				l.Push(lua.LNil)
				l.Push(lua.LString("device down:" + devUUID))
				return 2
			}

		}
		l.Push(lua.LNil)
		l.Push(lua.LString("device not exists:" + devUUID))
		return 2
	}
}
