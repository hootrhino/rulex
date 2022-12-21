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
	r.AddLib(e, "rulexlib", "DataToHttp", rulexlib.DataToHttp(e))
	r.AddLib(e, "rulexlib", "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	r.AddLib(e, "rulexlib", "JqSelect", rulexlib.JqSelect(e))
	r.AddLib(e, "rulexlib", "JQ", rulexlib.JqSelect(e))
	// 日志
	r.AddLib(e, "rulexlib", "log", rulexlib.Log(e))
	// 二进制操作
	r.AddLib(e, "rulexlib", "MB", rulexlib.MatchBinary(e))
	r.AddLib(e, "rulexlib", "B2BS", rulexlib.ByteToBitString(e))
	r.AddLib(e, "rulexlib", "Bit", rulexlib.GetABitOnByte(e))
	r.AddLib(e, "rulexlib", "B2I64", rulexlib.ByteToInt64(e))
	r.AddLib(e, "rulexlib", "BS2B", rulexlib.BitStringToBytes(e))
	r.AddLib(e, "rulexlib", "HToN", rulexlib.HToN(e))
	r.AddLib(e, "rulexlib", "HsubToN", rulexlib.HsubToN(e))
	// 浮点数处理
	r.AddLib(e, "rulexlib", "Bin2F32", rulexlib.BinToFloat32(e))
	r.AddLib(e, "rulexlib", "Bin2F64", rulexlib.BinToFloat32(e))
	// URL处理
	r.AddLib(e, "rulexlib", "UrlBuild", rulexlib.UrlBuild(e))
	r.AddLib(e, "rulexlib", "UrlBuildQS", rulexlib.UrlBuildQS(e))
	r.AddLib(e, "rulexlib", "UrlParse", rulexlib.UrlParse(e))
	r.AddLib(e, "rulexlib", "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	r.AddLib(e, "rulexlib", "DataToTdEngine", rulexlib.DataToTdEngine(e))
	r.AddLib(e, "rulexlib", "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	r.AddLib(e, "rulexlib", "Time", rulexlib.Time(e))
	r.AddLib(e, "rulexlib", "TsUnix", rulexlib.TsUnix(e))
	r.AddLib(e, "rulexlib", "TsUnixNano", rulexlib.TsUnixNano(e))
	r.AddLib(e, "rulexlib", "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	r.AddLib(e, "rulexlib", "VSet", rulexlib.StoreSet(e))
	r.AddLib(e, "rulexlib", "VGet", rulexlib.StoreGet(e))
	r.AddLib(e, "rulexlib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	r.AddLib(e, "rulexlib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	r.AddLib(e, "rulexlib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	r.AddLib(e, "rulexlib", "RUUID", rulexlib.SelfRuleUUID(e, r.UUID))
	// Codec
	r.AddLib(e, "rulexlib", "RPCENC", rulexlib.RPCEncode(e))
	r.AddLib(e, "rulexlib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	r.AddLib(e, "rulexlib", "ReadDevice", rulexlib.ReadDevice(e))
	r.AddLib(e, "rulexlib", "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	r.AddLib(e, "rulexlib", "ReadSource", rulexlib.ReadSource(e))
	r.AddLib(e, "rulexlib", "WriteSource", rulexlib.WriteSource(e))
	// String
	r.AddLib(e, "rulexlib", "T2Str", rulexlib.T2Str(e))
	r.AddLib(e, "rulexlib", "Throw", rulexlib.Throw(e))
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
