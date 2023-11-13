package rulexlib

import (
	"encoding/hex"
	"encoding/json"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

//  --------------------------------------------
// |Function | Register Type
//  --------------------------------------------
// |	1	 | Read Coil
// |	2	 | Read Discrete Input
// |	3	 | Read Holding Registers
// |	4	 | Read Input Registers
// |	5	 | Write Single Coil
// |	6	 | Write Single Holding Register
// |	15	 | Write Multiple Coils
// |	16	 | Write Multiple Holding Registers
//  --------------------------------------------
/*
*
* Modbus Function1
*
 */
func F1(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function2
*
 */
func F2(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function3
*
 */
func F3(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function4
*
 */
func F4(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
  - Modbus Function5
    local error = modbus:F5("uuid1", 0, 1, "0001020304")

*
*/

func F5(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Values := l.ToString(5)
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		for _, v := range HexValues {
			if v > 1 {
				l.Push(lua.LString("Value Only Support '00' or '01'"))
				return 1
			}
		}
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Type != typex.GENERIC_MODBUS {
				l.Push(lua.LString("Only support GENERIC_MODBUS device"))
				return 1
			}
			if Device.Device.Status() == typex.DEV_UP {
				args, _ := json.Marshal([]common.RegisterW{
					{
						Function: 5,
						SlaverId: byte(slaverId),
						Address:  uint16(Address),
						Values:   HexValues,
					},
				})
				_, err := Device.Device.OnWrite([]byte("F5"), args)
				if err != nil {
					l.Push(lua.LString(err.Error()))
					return 1
				}
			}
			l.Push(lua.LString("Device down:" + devUUID))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
*     local error = modbus:F6("uuid1", 0, 1, "0001020304")

*
 */
func F6(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Values := l.ToString(5) // 必须是单个字节: 000100010001
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Type != typex.GENERIC_MODBUS {
				l.Push(lua.LString("Only support GENERIC_MODBUS device"))
				return 1
			}
			if Device.Device.Status() == typex.DEV_UP {
				args, _ := json.Marshal(common.RegisterW{
					Function: 6,
					SlaverId: byte(slaverId),
					Address:  uint16(Address),
					Quantity: uint16(1), //2字节
					Values:   HexValues,
				})
				_, err := Device.Device.OnWrite([]byte("F6"), args)
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LString(err.Error()))
					return 1
				}
			}
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
  - Modbus Function15
    local error = modbus:F15("uuid1", 0, 1, "0001020304")

*
*/
func F15(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Quantity := l.ToNumber(5) // 必须是单个字节: 000100010001
		Values := l.ToString(6)   // 必须是单个字节: 000100010001
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Type != typex.GENERIC_MODBUS {
				l.Push(lua.LString("Only support GENERIC_MODBUS device"))
				return 1
			}
			if Device.Device.Status() == typex.DEV_UP {
				args, _ := json.Marshal(common.RegisterW{
					Function: 15,
					SlaverId: byte(slaverId),
					Address:  uint16(Address),
					Quantity: uint16(Quantity),
					Values:   HexValues,
				})
				_, err := Device.Device.OnWrite([]byte("F15"), args)
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LString(err.Error()))
					return 1
				}
			}
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* Modbus Function16
*    local error = modbus:F16("uuid1", 0, 1, "0001020304")
 */
func F16(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Quantity := l.ToNumber(5) //
		Values := l.ToString(6)   //
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Type != typex.GENERIC_MODBUS {
				l.Push(lua.LString("Only support GENERIC_MODBUS device"))
				return 1
			}
			if Device.Device.Status() == typex.DEV_UP {
				args, _ := json.Marshal(common.RegisterW{
					Function: 16,
					SlaverId: byte(slaverId),
					Address:  uint16(Address),
					Quantity: uint16(Quantity),
					Values:   HexValues,
				})
				_, err := Device.Device.OnWrite([]byte("F16"), args)
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		l.Push(lua.LNil)
		return 1
	}
}
