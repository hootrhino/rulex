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

package engine

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/rulexlib"
	"github.com/hootrhino/rulex/typex"
)

/*
*
  - 分组加入函数
*/
func AddRuleLibToGroup(r *typex.Rule, rx typex.RuleX,
	ModuleName string, funcs map[string]func(l *lua.LState) int) {
	table := r.LuaVM.NewTable()
	r.LuaVM.SetGlobal(ModuleName, table)
	for funcName, f := range funcs {
		table.RawSetString(funcName, r.LuaVM.NewClosure(f))
	}
	r.LuaVM.Push(table)
}

func LoadRuleLibGroup(r *typex.Rule, e typex.RuleX) {
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
		AddRuleLibToGroup(r, e, "data", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Debug":   rulexlib.DebugAPP(e, r.UUID),
			"Throw":   rulexlib.Throw(e),
			"Println": rulexlib.Println(e),
		}
		AddRuleLibToGroup(r, e, "_G", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"VSet": rulexlib.StoreSet(e),
			"VGet": rulexlib.StoreGet(e),
			"VDel": rulexlib.StoreDelete(e),
		}
		AddRuleLibToGroup(r, e, "kv", Funcs)
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
		AddRuleLibToGroup(r, e, "time", Funcs)
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
		AddRuleLibToGroup(r, e, "hex", Funcs)
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
		AddRuleLibToGroup(r, e, "binary", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2J": rulexlib.JSONE(e),
			"J2T": rulexlib.JSOND(e),
		}
		AddRuleLibToGroup(r, e, "json", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ReadDevice":  rulexlib.ReadDevice(e),
			"WriteDevice": rulexlib.WriteDevice(e),
			"CtrlDevice":  rulexlib.CtrlDevice(e),
			"ReadSource":  rulexlib.ReadSource(e),
			"WriteSource": rulexlib.WriteSource(e),
		}
		AddRuleLibToGroup(r, e, "device", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2Str":   rulexlib.T2Str(e),
			"Bin2Str": rulexlib.Bin2Str(e),
		}
		AddRuleLibToGroup(r, e, "string", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5":  rulexlib.F5(e),
			"F6":  rulexlib.F6(e),
			"F15": rulexlib.F15(e),
			"F16": rulexlib.F16(e),
		}
		AddRuleLibToGroup(r, e, "modbus", Funcs)
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
		AddRuleLibToGroup(r, e, "rhinopi", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"XOR":   rulexlib.XOR(e),
			"CRC16": rulexlib.CRC16(e),
		}
		AddRuleLibToGroup(r, e, "misc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rulexlib.RASPI4_GPIOGet(e),
			"GPIOSet": rulexlib.RASPI4_GPIOSet(e),
		}
		AddRuleLibToGroup(r, e, "raspi4b", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rulexlib.WKYWS1608_GPIOGet(e),
			"GPIOSet": rulexlib.WKYWS1608_GPIOSet(e),
		}
		AddRuleLibToGroup(r, e, "ws1608", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"TFloat": rulexlib.TruncateFloat(e),
		}
		AddRuleLibToGroup(r, e, "math", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"PlayMusic": rulexlib.PlayMusic(e),
		}
		AddRuleLibToGroup(r, e, "audio", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Request": rulexlib.Request(e),
		}
		AddRuleLibToGroup(r, e, "rpc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Execute": rulexlib.JqSelect(e),
		}
		AddRuleLibToGroup(r, e, "jq", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Ping": rulexlib.PingIp(e),
		}
		AddRuleLibToGroup(r, e, "network", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Get":  rulexlib.HttpGet(e),
			"Post": rulexlib.HttpPost(e),
		}
		AddRuleLibToGroup(r, e, "http", Funcs)
	}
	{
		// Just For test
		Func1 := map[string]func(l *lua.LState) int{
			"Time": rulexlib.Time(e),
		}
		AddRuleLibToGroup(r, e, "time1", Func1)
		Func2 := map[string]func(l *lua.LState) int{
			"Time": rulexlib.TsUnixNano(e),
		}
		AddRuleLibToGroup(r, e, "time2", Func2)
	}
}

/*
*
* 加载外部扩展库
*
 */
func LoadExtLuaLib(e typex.RuleX, r *typex.Rule) error {
	for _, s := range core.GlobalConfig.Extlibs.Value {
		err := r.LoadExternLuaLib(s)
		if err != nil {
			return err
		}
	}
	return nil
}
