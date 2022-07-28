package engine

import (
	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
)

/*
*
* 加载标准库
*
 */
func LoadBuildInLuaLib(e typex.RuleX, r *typex.Rule) {
	// 消息转发
	r.AddLib(e, "DataToHttp", rulexlib.DataToHttp(e))
	r.AddLib(e, "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	r.AddLib(e, "JqSelect", rulexlib.JqSelect(e))
	r.AddLib(e, "JQ", rulexlib.JqSelect(e))
	// 日志
	r.AddLib(e, "log", rulexlib.Log(e))
	// 二进制操作
	r.AddLib(e, "MB", rulexlib.MatchBinary(e))
	r.AddLib(e, "B2BS", rulexlib.ByteToBitString(e))
	r.AddLib(e, "Bit", rulexlib.GetABitOnByte(e))
	r.AddLib(e, "B2I64", rulexlib.ByteToInt64(e))
	r.AddLib(e, "BS2B", rulexlib.BitStringToBytes(e))
	r.AddLib(e, "HToN", rulexlib.HToN(e))
	r.AddLib(e, "HsubToN", rulexlib.HsubToN(e))
	// URL处理
	r.AddLib(e, "UrlBuild", rulexlib.UrlBuild(e))
	r.AddLib(e, "UrlBuildQS", rulexlib.UrlBuildQS(e))
	r.AddLib(e, "UrlParse", rulexlib.UrlParse(e))
	r.AddLib(e, "UrlResolve", rulexlib.UrlResolve(e))
	// 数据持久化
	r.AddLib(e, "DataToTdEngine", rulexlib.DataToTdEngine(e))
	r.AddLib(e, "DataToMongo", rulexlib.DataToMongo(e))
	// 时间库
	r.AddLib(e, "Time", rulexlib.Time(e))
	r.AddLib(e, "TsUnix", rulexlib.TsUnix(e))
	r.AddLib(e, "TsUnixNano", rulexlib.TsUnixNano(e))
	r.AddLib(e, "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	r.AddLib(e, "VSet", rulexlib.StoreSet(e))
	r.AddLib(e, "VGet", rulexlib.StoreGet(e))
	r.AddLib(e, "VDel", rulexlib.StoreDelete(e))
	// JSON
	r.AddLib(e, "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	r.AddLib(e, "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	r.AddLib(e, "RUUID", rulexlib.SelfRuleUUID(e, r.UUID))
	// Codec
	r.AddLib(e, "RPCENC", rulexlib.RPCEncode(e))
	r.AddLib(e, "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	r.AddLib(e, "ReadDevice", rulexlib.ReadDevice(e))
	r.AddLib(e, "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	r.AddLib(e, "ReadSource", rulexlib.ReadSource(e))
	r.AddLib(e, "WriteSource", rulexlib.WriteSource(e))
	// String
	r.AddLib(e, "T2Str", rulexlib.T2Str(e))
	r.AddLib(e, "Throw", rulexlib.Throw(e))
}
