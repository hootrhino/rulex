package appstack

import (
	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* AddLib: 根据 KV形式加载库(推荐)F
*  - Global: 命名空间
*   - funcName: 函数名称
 */
func AddAppLib(app *typex.Application,
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
* 加载app库函数
*
 */
func LoadAppLib(app *typex.Application, e typex.RuleX) {
	// 消息转发
	AddAppLib(app, e, "applib", "DataToHttp", rulexlib.DataToHttp(e))
	AddAppLib(app, e, "applib", "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	AddAppLib(app, e, "applib", "JqSelect", rulexlib.JqSelect(e))
	AddAppLib(app, e, "applib", "JQ", rulexlib.JqSelect(e))
	// 日志
	AddAppLib(app, e, "applib", "log", rulexlib.Log(e))
	// 二进制操作
	AddAppLib(app, e, "applib", "MB", rulexlib.MatchBinary(e))
	AddAppLib(app, e, "applib", "B2BS", rulexlib.ByteToBitString(e))
	AddAppLib(app, e, "applib", "Bit", rulexlib.GetABitOnByte(e))
	AddAppLib(app, e, "applib", "B2I64", rulexlib.ByteToInt64(e))
	AddAppLib(app, e, "applib", "B64S2B", rulexlib.B64S2B(e))
	AddAppLib(app, e, "applib", "BS2B", rulexlib.BitStringToBytes(e))
	AddAppLib(app, e, "applib", "HToN", rulexlib.HToN(e))
	AddAppLib(app, e, "applib", "HsubToN", rulexlib.HsubToN(e))
	// 浮点数处理
	AddAppLib(app, e, "applib", "Bin2F32", rulexlib.BinToFloat32(e))
	AddAppLib(app, e, "applib", "Bin2F64", rulexlib.BinToFloat64(e))
	// URL处理
	AddAppLib(app, e, "applib", "UrlBuild", rulexlib.UrlBuild(e))
	AddAppLib(app, e, "applib", "UrlBuildQS", rulexlib.UrlBuildQS(e))
	AddAppLib(app, e, "applib", "UrlParse", rulexlib.UrlParse(e))
	AddAppLib(app, e, "applib", "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	AddAppLib(app, e, "applib", "DataToTdEngine", rulexlib.DataToTdEngine(e))
	AddAppLib(app, e, "applib", "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	AddAppLib(app, e, "applib", "Time", rulexlib.Time(e))
	AddAppLib(app, e, "applib", "TsUnix", rulexlib.TsUnix(e))
	AddAppLib(app, e, "applib", "TsUnixNano", rulexlib.TsUnixNano(e))
	AddAppLib(app, e, "applib", "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	AddAppLib(app, e, "applib", "VSet", rulexlib.StoreSet(e))
	AddAppLib(app, e, "applib", "VGet", rulexlib.StoreGet(e))
	AddAppLib(app, e, "applib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	AddAppLib(app, e, "applib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	AddAppLib(app, e, "applib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	AddAppLib(app, e, "applib", "RUUID", rulexlib.SelfRuleUUID(e, app.UUID))
	// Codec
	AddAppLib(app, e, "applib", "RPCENC", rulexlib.RPCEncode(e))
	AddAppLib(app, e, "applib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	AddAppLib(app, e, "applib", "ReadDevice", rulexlib.ReadDevice(e))
	AddAppLib(app, e, "applib", "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	AddAppLib(app, e, "applib", "ReadSource", rulexlib.ReadSource(e))
	AddAppLib(app, e, "applib", "WriteSource", rulexlib.WriteSource(e))
	// String
	AddAppLib(app, e, "applib", "T2Str", rulexlib.T2Str(e))
	AddAppLib(app, e, "applib", "Throw", rulexlib.Throw(e))
	//------------------------------------------------------------------------
	// IotHUB 库, 主要是为了适配iothub的回复消息， 注意：这个规范是w3c的
	// https://www.w3.org/TR/wot-thing-description
	//------------------------------------------------------------------------
	AddAppLib(app, e, "iothub", "PropertySuccess", rulexlib.PropertyReplySuccess(e))
	AddAppLib(app, e, "iothub", "PropertyFailed", rulexlib.PropertyReplyFailed(e))
	AddAppLib(app, e, "iothub", "ActionSuccess", rulexlib.ActionReplySuccess(e))
	AddAppLib(app, e, "iothub", "ActionFailed", rulexlib.ActionReplyFailed(e))
	//------------------------------------------------------------------------
	// 设备操作
	//------------------------------------------------------------------------
	AddAppLib(app, e, "device", "DCACall", rulexlib.DCACall(e))
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	AddAppLib(app, e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
	AddAppLib(app, e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
	//------------------------------------------------------------------------
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	AddAppLib(app, e, "eekit", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	AddAppLib(app, e, "eekit", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))

}
