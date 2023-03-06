package engine

import (
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
)

/*
*
* 加载标准库, 为什么是每个LUA脚本都加载一次？主要是为了隔离，互不影响
*
 */
func LoadBuildInLuaLib(e typex.RuleX, r *typex.Rule) {
	// 消息转发
	r.AddLib(e, "applib", "DataToHttp", rulexlib.DataToHttp(e))
	r.AddLib(e, "applib", "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	r.AddLib(e, "applib", "JqSelect", rulexlib.JqSelect(e))
	r.AddLib(e, "applib", "JQ", rulexlib.JqSelect(e))
	// 日志
	r.AddLib(e, "applib", "log", rulexlib.Log(e))
	// 二进制操作
	r.AddLib(e, "applib", "MB", rulexlib.MatchBinary(e))
	r.AddLib(e, "applib", "B2BS", rulexlib.ByteToBitString(e))
	r.AddLib(e, "applib", "Bit", rulexlib.GetABitOnByte(e))
	r.AddLib(e, "applib", "B2I64", rulexlib.ByteToInt64(e))
	r.AddLib(e, "applib", "B64S2B", rulexlib.B64S2B(e))
	r.AddLib(e, "applib", "BS2B", rulexlib.BitStringToBytes(e))
	r.AddLib(e, "applib", "HToN", rulexlib.HToN(e))
	r.AddLib(e, "applib", "HsubToN", rulexlib.HsubToN(e))
	// 浮点数处理
	r.AddLib(e, "applib", "Bin2F32", rulexlib.BinToFloat32(e))
	r.AddLib(e, "applib", "Bin2F64", rulexlib.BinToFloat64(e))
	// URL处理
	r.AddLib(e, "applib", "UrlBuild", rulexlib.UrlBuild(e))
	r.AddLib(e, "applib", "UrlBuildQS", rulexlib.UrlBuildQS(e))
	r.AddLib(e, "applib", "UrlParse", rulexlib.UrlParse(e))
	r.AddLib(e, "applib", "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	r.AddLib(e, "applib", "DataToTdEngine", rulexlib.DataToTdEngine(e))
	r.AddLib(e, "applib", "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	r.AddLib(e, "applib", "Time", rulexlib.Time(e))
	r.AddLib(e, "applib", "TsUnix", rulexlib.TsUnix(e))
	r.AddLib(e, "applib", "TsUnixNano", rulexlib.TsUnixNano(e))
	r.AddLib(e, "applib", "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	r.AddLib(e, "applib", "VSet", rulexlib.StoreSet(e))
	r.AddLib(e, "applib", "VGet", rulexlib.StoreGet(e))
	r.AddLib(e, "applib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	r.AddLib(e, "applib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	r.AddLib(e, "applib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	r.AddLib(e, "applib", "RUUID", rulexlib.SelfRuleUUID(e, r.UUID))
	// Codec
	r.AddLib(e, "applib", "RPCENC", rulexlib.RPCEncode(e))
	r.AddLib(e, "applib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	r.AddLib(e, "applib", "ReadDevice", rulexlib.ReadDevice(e))
	r.AddLib(e, "applib", "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	r.AddLib(e, "applib", "ReadSource", rulexlib.ReadSource(e))
	r.AddLib(e, "applib", "WriteSource", rulexlib.WriteSource(e))
	// String
	r.AddLib(e, "applib", "T2Str", rulexlib.T2Str(e))
	r.AddLib(e, "applib", "Throw", rulexlib.Throw(e))
	//------------------------------------------------------------------------
	// IotHUB 库, 主要是为了适配iothub的回复消息， 注意：这个规范是w3c的
	// https://www.w3.org/TR/wot-thing-description
	//------------------------------------------------------------------------
	r.AddLib(e, "iothub", "PropertySuccess", rulexlib.PropertyReplySuccess(e))
	r.AddLib(e, "iothub", "PropertyFailed", rulexlib.PropertyReplyFailed(e))
	r.AddLib(e, "iothub", "ActionSuccess", rulexlib.ActionReplySuccess(e))
	r.AddLib(e, "iothub", "ActionFailed", rulexlib.ActionReplyFailed(e))
	//------------------------------------------------------------------------
	// 设备操作
	//------------------------------------------------------------------------
	r.AddLib(e, "device", "DCACall", rulexlib.DCACall(e))
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	r.AddLib(e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
	r.AddLib(e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
	//------------------------------------------------------------------------
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	r.AddLib(e, "eekit", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	r.AddLib(e, "eekit", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))

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
