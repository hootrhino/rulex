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
* AddLib: 根据 KV形式加载库(推荐)
*  - Global: 命名空间
*   - funcName: 函数名称
 */
func addAppLib(app *Application,
	rx typex.RuleX, Global string, funcName string,
	f func(l *lua.LState) int) {
	rulexTb := app.VM().G.Global
	app.VM().SetGlobal(Global, rulexTb)
	mod := app.VM().SetFuncs(rulexTb, map[string]lua.LGFunction{
		funcName: f,
	})
	app.VM().Push(mod)
}

/*
*
* 加载app库函数, 注意这里的库函数和规则引擎的并不是完全一样的，有一些差别
*
 */
func LoadAppLib(app *Application, e typex.RuleX) {
	// 消息转发
	// 数据持久化
	{
		addAppLib(app, e, "data", "ToHttp", rulexlib.DataToHttp(e))
		addAppLib(app, e, "data", "ToMqtt", rulexlib.DataToMqtt(e))
		addAppLib(app, e, "data", "ToUdp", rulexlib.DataToUdp(e))
		addAppLib(app, e, "data", "ToTcp", rulexlib.DataToTcp(e))
		addAppLib(app, e, "data", "ToTdEngine", rulexlib.DataToTdEngine(e))
		addAppLib(app, e, "data", "ToMongo", rulexlib.DataToMongo(e))
		addAppLib(app, e, "data", "ToScreen", rulexlib.DataToUiComponent(e))

	}

	{
		addAppLib(app, e, "stdlib", "Debug", rulexlib.DebugAPP(e, app.UUID))
		addAppLib(app, e, "stdlib", "Throw", rulexlib.Throw(e))
		addAppLib(app, e, "stdlib", "Println", rulexlib.Println(e))
		addAppLib(app, e, "_G", "Debug", rulexlib.DebugAPP(e, app.UUID))
		addAppLib(app, e, "_G", "Throw", rulexlib.Throw(e))
	}

	// 二进制操作
	{
		addAppLib(app, e, "binary", "MB", rulexlib.MatchBinary(e))
		addAppLib(app, e, "binary", "MBHex", rulexlib.MatchBinaryHex(e))
		addAppLib(app, e, "binary", "B2BS", rulexlib.ByteToBitString(e))
		addAppLib(app, e, "binary", "Bit", rulexlib.GetABitOnByte(e))
		addAppLib(app, e, "binary", "B2I64", rulexlib.ByteToInt64(e))
		addAppLib(app, e, "binary", "B64S2B", rulexlib.B64S2B(e))
		addAppLib(app, e, "binary", "BS2B", rulexlib.BitStringToBytes(e))
		// 浮点数处理
		addAppLib(app, e, "binary", "Bin2F32", rulexlib.BinToFloat32(e))
		addAppLib(app, e, "binary", "Bin2F64", rulexlib.BinToFloat64(e))
	}
	{
		addAppLib(app, e, "hex", "HToN", rulexlib.HToN(e))
		addAppLib(app, e, "hex", "HsubToN", rulexlib.HsubToN(e))
		addAppLib(app, e, "hex", "MatchHex", rulexlib.MatchHex(e))
		addAppLib(app, e, "hex", "MatchUInt", rulexlib.MatchUInt(e))
		addAppLib(app, e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
		addAppLib(app, e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
		addAppLib(app, e, "hex", "ABCD", rulexlib.ABCD(e))
		addAppLib(app, e, "hex", "DCBA", rulexlib.DCBA(e))
		addAppLib(app, e, "hex", "BADC", rulexlib.BADC(e))
		addAppLib(app, e, "hex", "CDAB", rulexlib.CDAB(e))
	}
	{
		// URL处理
		addAppLib(app, e, "url", "UrlBuild", rulexlib.UrlBuild(e))
		addAppLib(app, e, "url", "UrlBuildQS", rulexlib.UrlBuildQS(e))
		addAppLib(app, e, "url", "UrlParse", rulexlib.UrlParse(e))
		addAppLib(app, e, "url", "UrlResolve", rulexlib.UrlResolve(e))
	}

	{
		// 时间库
		addAppLib(app, e, "time", "Time", rulexlib.Time(e))
		addAppLib(app, e, "time", "TimeMs", rulexlib.TimeMs(e))
		addAppLib(app, e, "time", "TsUnix", rulexlib.TsUnix(e))
		addAppLib(app, e, "time", "TsUnixNano", rulexlib.TsUnixNano(e))
		addAppLib(app, e, "time", "NtpTime", rulexlib.NtpTime(e))
		addAppLib(app, e, "time", "Sleep", rulexlib.Sleep(e))
	}
	{
		// 缓存器库
		addAppLib(app, e, "kv", "VSet", rulexlib.StoreSet(e))
		addAppLib(app, e, "kv", "VGet", rulexlib.StoreGet(e))
		addAppLib(app, e, "kv", "VDel", rulexlib.StoreDelete(e))
	}
	{
		// JSON
		addAppLib(app, e, "json", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
		addAppLib(app, e, "json", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	}
	{
		// LocalDBQuery -> data ,error
		addAppLib(app, e, "datacenter", "DBQuery", rulexlib.LocalDBQuery(e))
	}
	{
		// Device R/W
		addAppLib(app, e, "device", "ReadDevice", rulexlib.ReadDevice(e))
		addAppLib(app, e, "device", "WriteDevice", rulexlib.WriteDevice(e))
		// Ctrl Device: request --> response
		addAppLib(app, e, "device", "CtrlDevice", rulexlib.CtrlDevice(e))
		// Source R/W
		addAppLib(app, e, "device", "ReadSource", rulexlib.ReadSource(e))
		addAppLib(app, e, "device", "WriteSource", rulexlib.WriteSource(e))
	}
	{
		// String
		addAppLib(app, e, "string", "T2Str", rulexlib.T2Str(e))
		addAppLib(app, e, "string", "Bin2Str", rulexlib.Bin2Str(e))
	}
	//------------------------------------------------------------------------
	// Modbus
	//------------------------------------------------------------------------
	{
		addAppLib(app, e, "modbus", "F5", rulexlib.F5(e))
		addAppLib(app, e, "modbus", "F6", rulexlib.F6(e))
		addAppLib(app, e, "modbus", "F15", rulexlib.F15(e))
		addAppLib(app, e, "modbus", "F16", rulexlib.F16(e))
	}
	//------------------------------------------------------------------------
	// 注册 rhinopi GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	{
		// DO1 DO2
		addAppLib(app, e, "rhinopi", "DO1Set", rulexlib.H3DO1Set(e))
		addAppLib(app, e, "rhinopi", "DO1Get", rulexlib.H3DO1Get(e))
		addAppLib(app, e, "rhinopi", "DO2Set", rulexlib.H3DO2Set(e))
		addAppLib(app, e, "rhinopi", "DO2Get", rulexlib.H3DO2Get(e))
		// DI1 DI2 DI3
		addAppLib(app, e, "rhinopi", "DI1Get", rulexlib.H3DI1Get(e))
		addAppLib(app, e, "rhinopi", "DI2Get", rulexlib.H3DI2Get(e))
		addAppLib(app, e, "rhinopi", "DI3Get", rulexlib.H3DI3Get(e))

	}
	{
		addAppLib(app, e, "misc", "XOR", rulexlib.XOR(e))
		addAppLib(app, e, "misc", "CRC16", rulexlib.CRC16(e))
	}
	{
		// 树莓派4B
		addAppLib(app, e, "raspi4b", "GPIOGet", rulexlib.RASPI4_GPIOGet(e))
		addAppLib(app, e, "raspi4b", "GPIOSet", rulexlib.RASPI4_GPIOSet(e))
	}
	{
		// // 玩客云WS1608
		addAppLib(app, e, "ws1608", "GPIOGet", rulexlib.WKYWS1608_GPIOGet(e))
		addAppLib(app, e, "ws1608", "GPIOSet", rulexlib.WKYWS1608_GPIOSet(e))
	}
	// math 数学库
	addAppLib(app, e, "math", "TFloat", rulexlib.TruncateFloat(e))
	// Audio
	addAppLib(app, e, "audio", "PlayMusic", rulexlib.PlayMusic(e))
	// rpc
	addAppLib(app, e, "rpc", "Request", rulexlib.Request(e))
	// jq
	addAppLib(app, e, "jq", "Execute", rulexlib.JqSelect(e))
}
