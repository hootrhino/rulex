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
	// // 检查名称
	// AppNAME := tempVm.GetGlobal("AppNAME")
	// if AppNAME == nil {
	// 	return fmt.Errorf("'AppNAME' field not exists")
	// }
	// if AppNAME.Type() != lua.LTString {
	// 	return fmt.Errorf("'AppNAME' must be string")
	// }
	// // 检查类型
	// AppVERSION := tempVm.GetGlobal("AppVERSION")
	// if AppVERSION == nil {
	// 	return fmt.Errorf("'AppVERSION' field not exists")
	// }
	// if AppVERSION.Type() != lua.LTString {
	// 	return fmt.Errorf("'AppVERSION' must be string")
	// }
	// // 检查描述信息
	// AppDESCRIPTION := tempVm.GetGlobal("AppDESCRIPTION")
	// if AppDESCRIPTION == nil {
	// 	if AppDESCRIPTION.Type() != lua.LTString {
	// 		return fmt.Errorf("'AppDESCRIPTION' must be string")
	// 	}
	// }

	// 检查函数入口
	AppMain := tempVm.GetGlobal("Main")
	if AppMain == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMain.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 释放语法验证阶段的临时虚拟机
	tempVm.Close()
	tempVm = nil
	return nil
}

/*
*
* AddLib: 根据 KV形式加载库(推荐)F
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
	addAppLib(app, e, "applib", "DataToHttp", rulexlib.DataToHttp(e))
	addAppLib(app, e, "applib", "DataToMqtt", rulexlib.DataToMqtt(e))

	addAppLib(app, e, "applib", "DataToUdp", rulexlib.DataToUdp(e))
	// JQ
	addAppLib(app, e, "applib", "JqSelect", rulexlib.JqSelect(e))
	addAppLib(app, e, "applib", "JQ", rulexlib.JqSelect(e))
	addAppLib(app, e, "applib", "Debug", rulexlib.DebugAPP(e, app.UUID))
	// 二进制操作
	addAppLib(app, e, "applib", "MB", rulexlib.MatchBinary(e))
	addAppLib(app, e, "applib", "MBHex", rulexlib.MatchBinaryHex(e))
	addAppLib(app, e, "applib", "B2BS", rulexlib.ByteToBitString(e))
	addAppLib(app, e, "applib", "Bit", rulexlib.GetABitOnByte(e))
	addAppLib(app, e, "applib", "B2I64", rulexlib.ByteToInt64(e))
	addAppLib(app, e, "applib", "B64S2B", rulexlib.B64S2B(e))
	addAppLib(app, e, "applib", "BS2B", rulexlib.BitStringToBytes(e))
	addAppLib(app, e, "applib", "HToN", rulexlib.HToN(e))
	addAppLib(app, e, "applib", "HsubToN", rulexlib.HsubToN(e))
	addAppLib(app, e, "applib", "MatchHex", rulexlib.MatchHex(e))
	addAppLib(app, e, "applib", "MatchUInt", rulexlib.MatchUInt(e))
	// 浮点数处理
	addAppLib(app, e, "applib", "Bin2F32", rulexlib.BinToFloat32(e))
	addAppLib(app, e, "applib", "Bin2F64", rulexlib.BinToFloat64(e))
	// URL处理
	addAppLib(app, e, "applib", "UrlBuild", rulexlib.UrlBuild(e))
	addAppLib(app, e, "applib", "UrlBuildQS", rulexlib.UrlBuildQS(e))
	addAppLib(app, e, "applib", "UrlParse", rulexlib.UrlParse(e))
	addAppLib(app, e, "applib", "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	addAppLib(app, e, "applib", "DataToTdEngine", rulexlib.DataToTdEngine(e))
	addAppLib(app, e, "applib", "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	addAppLib(app, e, "applib", "Time", rulexlib.Time(e))
	addAppLib(app, e, "applib", "TimeMs", rulexlib.TimeMs(e))
	addAppLib(app, e, "applib", "TsUnix", rulexlib.TsUnix(e))
	addAppLib(app, e, "applib", "TsUnixNano", rulexlib.TsUnixNano(e))
	addAppLib(app, e, "applib", "NtpTime", rulexlib.NtpTime(e))
	addAppLib(app, e, "applib", "Sleep", rulexlib.Sleep(e))
	// 缓存器库
	addAppLib(app, e, "applib", "VSet", rulexlib.StoreSet(e))
	addAppLib(app, e, "applib", "VGet", rulexlib.StoreGet(e))
	addAppLib(app, e, "applib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	addAppLib(app, e, "applib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	addAppLib(app, e, "applib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table

	// Codec
	addAppLib(app, e, "applib", "RPCENC", rulexlib.RPCEncode(e))
	addAppLib(app, e, "applib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	addAppLib(app, e, "applib", "ReadDevice", rulexlib.ReadDevice(e))
	addAppLib(app, e, "applib", "WriteDevice", rulexlib.WriteDevice(e))
	// Ctrl Device: request --> response
	addAppLib(app, e, "applib", "CtrlDevice", rulexlib.CtrlDevice(e))
	// Source R/W
	addAppLib(app, e, "applib", "ReadSource", rulexlib.ReadSource(e))
	addAppLib(app, e, "applib", "WriteSource", rulexlib.WriteSource(e))
	// String
	addAppLib(app, e, "applib", "T2Str", rulexlib.T2Str(e))
	addAppLib(app, e, "applib", "Bin2Str", rulexlib.Bin2Str(e))
	addAppLib(app, e, "applib", "Throw", rulexlib.Throw(e))
	//------------------------------------------------------------------------
	// IotHUB 库, 主要是为了适配iothub的回复消息， 注意：这个规范是w3c的
	// https://www.w3.org/TR/wot-thing-description
	//------------------------------------------------------------------------
	addAppLib(app, e, "iothub", "PropertySuccess", rulexlib.PropertyReplySuccess(e))
	addAppLib(app, e, "iothub", "PropertyFailed", rulexlib.PropertyReplyFailed(e))
	addAppLib(app, e, "iothub", "ActionSuccess", rulexlib.ActionReplySuccess(e))
	addAppLib(app, e, "iothub", "ActionFailed", rulexlib.ActionReplyFailed(e))
	//------------------------------------------------------------------------
	// 设备操作
	//------------------------------------------------------------------------
	addAppLib(app, e, "device", "DCACall", rulexlib.DCACall(e))
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	addAppLib(app, e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
	addAppLib(app, e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
	//------------------------------------------------------------------------
	// 十六进制字节序处理
	//------------------------------------------------------------------------
	addAppLib(app, e, "hex", "ABCD", rulexlib.ABCD(e))
	addAppLib(app, e, "hex", "DCBA", rulexlib.DCBA(e))
	addAppLib(app, e, "hex", "BADC", rulexlib.BADC(e))
	addAppLib(app, e, "hex", "CDAB", rulexlib.CDAB(e))
	//------------------------------------------------------------------------
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	// EEKIT H3
	addAppLib(app, e, "eekith3", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	addAppLib(app, e, "eekith3", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))
	// DO1 DO2
	addAppLib(app, e, "eekith3", "H3DO1Set", rulexlib.H3DO1Set(e))
	addAppLib(app, e, "eekith3", "H3DO1Get", rulexlib.H3DO1Get(e))
	// DI1 DI2
	addAppLib(app, e, "eekith3", "H3DI1Get", rulexlib.H3DI1Get(e))
	addAppLib(app, e, "eekith3", "H3DI2Get", rulexlib.H3DI2Get(e))
	addAppLib(app, e, "eekith3", "H3DI3Get", rulexlib.H3DI3Get(e))
	//
	addAppLib(app, e, "eekith3", "H3DO2Set", rulexlib.H3DO2Set(e))
	addAppLib(app, e, "eekith3", "H3DO2Get", rulexlib.H3DO2Get(e))
	// 树莓派4B
	addAppLib(app, e, "raspi4b", "GPIOGet", rulexlib.RASPI4_GPIOGet(e))
	addAppLib(app, e, "raspi4b", "GPIOSet", rulexlib.RASPI4_GPIOSet(e))
	// 玩客云WS1608
	addAppLib(app, e, "ws1608", "GPIOGet", rulexlib.WKYWS1608_GPIOGet(e))
	addAppLib(app, e, "ws1608", "GPIOSet", rulexlib.WKYWS1608_GPIOSet(e))
	//------------------------------------------------------------------------
	// 校验数据
	//------------------------------------------------------------------------
	addAppLib(app, e, "misc", "XOR", rulexlib.XOR(e))
	addAppLib(app, e, "misc", "CRC16", rulexlib.CRC16(e))
	//------------------------------------------------------------------------
	// yqueue
	//------------------------------------------------------------------------
	addAppLib(app, e, "pipe", "Output", rulexlib.Output(e))
	//------------------------------------------------------------------------
	// Audio
	//------------------------------------------------------------------------
	addAppLib(app, e, "audio", "PlayMusic", rulexlib.PlayMusic(e))
	//------------------------------------------------------------------------
	// Math
	//------------------------------------------------------------------------
	addAppLib(app, e, "math", "TFloat", rulexlib.TruncateFloat(e))
	//------------------------------------------------------------------------
	// UI 数据写入
	//------------------------------------------------------------------------
	addAppLib(app, e, "ui", "LoadData", rulexlib.DataToUiComponent(e))
	//------------------------------------------------------------------------
	// Modbus
	//------------------------------------------------------------------------
	addAppLib(app, e, "modbus", "F5", rulexlib.F5(e))
	addAppLib(app, e, "modbus", "F6", rulexlib.F6(e))
	addAppLib(app, e, "modbus", "F15", rulexlib.F15(e))
	addAppLib(app, e, "modbus", "F16", rulexlib.F16(e))

}
