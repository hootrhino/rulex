package appstack

import (
	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* AddLib: 根据 KV形式加载库(推荐)
*  - Global: 命名空间
*   - funcName: 函数名称
 */
func (app *Application) addLib(rx typex.RuleX, Global string, funcName string,
	f func(l *lua.LState) int) {
	rulexTb := app.vm.G.Global
	app.vm.SetGlobal(Global, rulexTb)
	mod := app.vm.SetFuncs(rulexTb, map[string]lua.LGFunction{
		funcName: f,
	})
	app.vm.Push(mod)
}

/*
*
* 加载app库函数
*
 */
func (app *Application) loadAppLib(e typex.RuleX) {
	// 消息转发
	app.addLib(e, "rulexlib", "DataToHttp", rulexlib.DataToHttp(e))
	app.addLib(e, "rulexlib", "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	app.addLib(e, "rulexlib", "JqSelect", rulexlib.JqSelect(e))
	app.addLib(e, "rulexlib", "JQ", rulexlib.JqSelect(e))
	// 日志
	app.addLib(e, "rulexlib", "log", rulexlib.Log(e))
	// 二进制操作
	app.addLib(e, "rulexlib", "MB", rulexlib.MatchBinary(e))
	app.addLib(e, "rulexlib", "B2BS", rulexlib.ByteToBitString(e))
	app.addLib(e, "rulexlib", "Bit", rulexlib.GetABitOnByte(e))
	app.addLib(e, "rulexlib", "B2I64", rulexlib.ByteToInt64(e))
	app.addLib(e, "rulexlib", "B64S2B", rulexlib.B64S2B(e))
	app.addLib(e, "rulexlib", "BS2B", rulexlib.BitStringToBytes(e))
	app.addLib(e, "rulexlib", "HToN", rulexlib.HToN(e))
	app.addLib(e, "rulexlib", "HsubToN", rulexlib.HsubToN(e))
	// 浮点数处理
	app.addLib(e, "rulexlib", "Bin2F32", rulexlib.BinToFloat32(e))
	app.addLib(e, "rulexlib", "Bin2F64", rulexlib.BinToFloat64(e))
	// URL处理
	app.addLib(e, "rulexlib", "UrlBuild", rulexlib.UrlBuild(e))
	app.addLib(e, "rulexlib", "UrlBuildQS", rulexlib.UrlBuildQS(e))
	app.addLib(e, "rulexlib", "UrlParse", rulexlib.UrlParse(e))
	app.addLib(e, "rulexlib", "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	app.addLib(e, "rulexlib", "DataToTdEngine", rulexlib.DataToTdEngine(e))
	app.addLib(e, "rulexlib", "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	app.addLib(e, "rulexlib", "Time", rulexlib.Time(e))
	app.addLib(e, "rulexlib", "TsUnix", rulexlib.TsUnix(e))
	app.addLib(e, "rulexlib", "TsUnixNano", rulexlib.TsUnixNano(e))
	app.addLib(e, "rulexlib", "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	app.addLib(e, "rulexlib", "VSet", rulexlib.StoreSet(e))
	app.addLib(e, "rulexlib", "VGet", rulexlib.StoreGet(e))
	app.addLib(e, "rulexlib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	app.addLib(e, "rulexlib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	app.addLib(e, "rulexlib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	app.addLib(e, "rulexlib", "RUUID", rulexlib.SelfRuleUUID(e, app.UUID))
	// Codec
	app.addLib(e, "rulexlib", "RPCENC", rulexlib.RPCEncode(e))
	app.addLib(e, "rulexlib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	app.addLib(e, "rulexlib", "ReadDevice", rulexlib.ReadDevice(e))
	app.addLib(e, "rulexlib", "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	app.addLib(e, "rulexlib", "ReadSource", rulexlib.ReadSource(e))
	app.addLib(e, "rulexlib", "WriteSource", rulexlib.WriteSource(e))
	// String
	app.addLib(e, "rulexlib", "T2Str", rulexlib.T2Str(e))
	app.addLib(e, "rulexlib", "Throw", rulexlib.Throw(e))
	//------------------------------------------------------------------------
	// IotHUB 库, 主要是为了适配iothub的回复消息， 注意：这个规范是w3c的
	// https://www.w3.org/TR/wot-thing-description
	//------------------------------------------------------------------------
	app.addLib(e, "iothub", "PropertySuccess", rulexlib.PropertyReplySuccess(e))
	app.addLib(e, "iothub", "PropertyFailed", rulexlib.PropertyReplyFailed(e))
	app.addLib(e, "iothub", "ActionSuccess", rulexlib.ActionReplySuccess(e))
	app.addLib(e, "iothub", "ActionFailed", rulexlib.ActionReplyFailed(e))
	//------------------------------------------------------------------------
	// 设备操作
	//------------------------------------------------------------------------
	app.addLib(e, "device", "DCACall", rulexlib.DCACall(e))
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	app.addLib(e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
	app.addLib(e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
	//------------------------------------------------------------------------
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	app.addLib(e, "eekit", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	app.addLib(e, "eekit", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))

}
