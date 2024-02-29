// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package appstack

import (
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/rulexlib"
	"github.com/hootrhino/rulex/typex"
)

// 临时校验语法
func ValidateLuaSyntax(bytes []byte) error {
	// 把虚拟机参数全部设置为0是为了防止误操作产生垃圾数据
	tempVm := lua.NewState(lua.Options{
		SkipOpenLibs:     true,
		RegistrySize:     0,
		RegistryMaxSize:  0,
		RegistryGrowStep: 0,
	})
	if err := tempVm.DoString(string(bytes)); err != nil {
		return err
	}
	// 检查函数入口
	AppMain := tempVm.GetGlobal("Main")
	if AppMain == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMain.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	tempVm.Close()
	tempVm = nil
	return nil
}

/*
*
  - 分组加入函数
*/
func AddAppLibToGroup(app *Application, rx typex.RuleX,
	ModuleName string, funcs map[string]func(l *lua.LState) int) {
	var table *lua.LTable
	if ModuleName == "_G" {
		table = app.vm.G.Global
	} else {
		table = app.vm.NewTable()
	}
	app.vm.SetGlobal(ModuleName, table)
	for funcName, f := range funcs {
		table.RawSetString(funcName, app.vm.NewClosure(f))
	}
	app.vm.Push(table)
}

func LoadAppLibGroup(app *Application, e typex.RuleX) {
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ToHttp":     rulexlib.DataToHttp(e),
			"ToMqtt":     rulexlib.DataToMqtt(e),
			"ToUdp":      rulexlib.DataToUdp(e),
			"ToTcp":      rulexlib.DataToTcp(e),
			"ToTdEngine": rulexlib.DataToTdEngine(e),
			"ToMongo":    rulexlib.DataToMongo(e),
			"ToScreen":   rulexlib.DataToUiComponent(e),
		}
		AddAppLibToGroup(app, e, "data", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Debug":   rulexlib.DebugAPP(e, app.UUID),
			"Throw":   rulexlib.Throw(e),
			"Println": rulexlib.Println(e),
		}
		AddAppLibToGroup(app, e, "_G", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"VSet": rulexlib.StoreSet(e),
			"VGet": rulexlib.StoreGet(e),
			"VDel": rulexlib.StoreDelete(e),
		}
		AddAppLibToGroup(app, e, "kv", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Time":       rulexlib.Time(e),
			"TimeMs":     rulexlib.TimeMs(e),
			"TsUnix":     rulexlib.TsUnix(e),
			"TsUnixNano": rulexlib.TsUnixNano(e),
			"NtpTime":    rulexlib.NtpTime(e),
			"Sleep":      rulexlib.Sleep(e),
		}
		AddAppLibToGroup(app, e, "time", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"HToN":       rulexlib.HToN(e),
			"HsubToN":    rulexlib.HsubToN(e),
			"MatchHex":   rulexlib.MatchHex(e),
			"MatchUInt":  rulexlib.MatchUInt(e),
			"Bytes2Hexs": rulexlib.Bytes2Hexs(e),
			"Hexs2Bytes": rulexlib.Hexs2Bytes(e),
			"ABCD":       rulexlib.ABCD(e),
			"DCBA":       rulexlib.DCBA(e),
			"BADC":       rulexlib.BADC(e),
			"CDAB":       rulexlib.CDAB(e),
		}
		AddAppLibToGroup(app, e, "hex", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"MB":            rulexlib.MatchBinary(e),
			"MBHex":         rulexlib.MatchBinaryHex(e),
			"B2BS":          rulexlib.ByteToBitString(e),
			"Bit":           rulexlib.GetABitOnByte(e),
			"B2I64":         rulexlib.ByteToInt64(e),
			"B64S2B":        rulexlib.B64S2B(e),
			"BS2B":          rulexlib.BitStringToBytes(e),
			"Bin2F32":       rulexlib.BinToFloat32(e),
			"Bin2F64":       rulexlib.BinToFloat64(e),
			"Bin2F32Big":    rulexlib.BinToFloat32(e),
			"Bin2F64Big":    rulexlib.BinToFloat64(e),
			"Bin2F32Little": rulexlib.BinToFloat32Little(e),
			"Bin2F64Little": rulexlib.BinToFloat64Little(e),
		}
		AddAppLibToGroup(app, e, "binary", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2J": rulexlib.JSONE(e),
			"J2T": rulexlib.JSOND(e),
		}
		AddAppLibToGroup(app, e, "json", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ReadDevice":  rulexlib.ReadDevice(e),
			"WriteDevice": rulexlib.WriteDevice(e),
			"CtrlDevice":  rulexlib.CtrlDevice(e),
			"ReadSource":  rulexlib.ReadSource(e),
			"WriteSource": rulexlib.WriteSource(e),
		}
		AddAppLibToGroup(app, e, "device", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2Str":   rulexlib.T2Str(e),
			"Bin2Str": rulexlib.Bin2Str(e),
		}
		AddAppLibToGroup(app, e, "string", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5":  rulexlib.F5(e),
			"F6":  rulexlib.F6(e),
			"F15": rulexlib.F15(e),
			"F16": rulexlib.F16(e),
		}
		AddAppLibToGroup(app, e, "modbus", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"DO1Set":  rulexlib.H3DO1Set(e),
			"DO1Get":  rulexlib.H3DO1Get(e),
			"DO2Set":  rulexlib.H3DO2Set(e),
			"DO2Get":  rulexlib.H3DO2Get(e),
			"DI1Get":  rulexlib.H3DI1Get(e),
			"DI2Get":  rulexlib.H3DI2Get(e),
			"DI3Get":  rulexlib.H3DI3Get(e),
			"Led1On":  rulexlib.Led1On(e),
			"Led1Off": rulexlib.Led1Off(e),
		}
		AddAppLibToGroup(app, e, "rhinopi", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"XOR":   rulexlib.XOR(e),
			"CRC16": rulexlib.CRC16(e),
		}
		AddAppLibToGroup(app, e, "misc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rulexlib.RASPI4_GPIOGet(e),
			"GPIOSet": rulexlib.RASPI4_GPIOSet(e),
		}
		AddAppLibToGroup(app, e, "raspi4b", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rulexlib.WKYWS1608_GPIOGet(e),
			"GPIOSet": rulexlib.WKYWS1608_GPIOSet(e),
		}
		AddAppLibToGroup(app, e, "ws1608", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"TFloat": rulexlib.TruncateFloat(e),
		}
		AddAppLibToGroup(app, e, "math", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"PlayMusic": rulexlib.PlayMusic(e),
		}
		AddAppLibToGroup(app, e, "audio", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Request": rulexlib.Request(e),
		}
		AddAppLibToGroup(app, e, "rpc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Execute": rulexlib.JqSelect(e),
		}
		AddAppLibToGroup(app, e, "jq", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Ping": rulexlib.PingIp(e),
		}
		AddAppLibToGroup(app, e, "network", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Get":  rulexlib.HttpGet(e),
			"Post": rulexlib.HttpPost(e),
		}
		AddAppLibToGroup(app, e, "http", Funcs)
	}
	{
		// Just For test
		Func1 := map[string]func(l *lua.LState) int{
			"Time": rulexlib.Time(e),
		}
		AddAppLibToGroup(app, e, "time1", Func1)
		Func2 := map[string]func(l *lua.LState) int{
			"Time": rulexlib.TsUnixNano(e),
		}
		AddAppLibToGroup(app, e, "time2", Func2)
	}
}
