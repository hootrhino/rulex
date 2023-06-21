package engine

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/rulexlib"
	"github.com/hootrhino/rulex/typex"
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
	r.AddLib(e, "rulexlib", "DataToUdp", rulexlib.DataToUdp(e))
	// JQ
	r.AddLib(e, "rulexlib", "JqSelect", rulexlib.JqSelect(e))
	r.AddLib(e, "rulexlib", "JQ", rulexlib.JqSelect(e))
	// 日志
	r.AddLib(e, "rulexlib", "log", rulexlib.Log(e))
	// 二进制操作
	r.AddLib(e, "rulexlib", "MB", rulexlib.MatchBinary(e))
	r.AddLib(e, "rulexlib", "MBHex", rulexlib.MatchBinaryHex(e))
	r.AddLib(e, "rulexlib", "B2BS", rulexlib.ByteToBitString(e))
	r.AddLib(e, "rulexlib", "Bit", rulexlib.GetABitOnByte(e))
	r.AddLib(e, "rulexlib", "B2I64", rulexlib.ByteToInt64(e))
	r.AddLib(e, "rulexlib", "B64S2B", rulexlib.B64S2B(e))
	r.AddLib(e, "rulexlib", "BS2B", rulexlib.BitStringToBytes(e))
	r.AddLib(e, "rulexlib", "HToN", rulexlib.HToN(e))
	r.AddLib(e, "rulexlib", "HsubToN", rulexlib.HsubToN(e))
	r.AddLib(e, "rulexlib", "MatchHex", rulexlib.MatchHex(e))
	r.AddLib(e, "rulexlib", "MatchUInt", rulexlib.MatchUInt(e))
	// 浮点数处理
	r.AddLib(e, "rulexlib", "Bin2F32", rulexlib.BinToFloat32(e))
	r.AddLib(e, "rulexlib", "Bin2F64", rulexlib.BinToFloat64(e))
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
	r.AddLib(e, "rulexlib", "Sleep", rulexlib.Sleep(e))

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
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	r.AddLib(e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
	r.AddLib(e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
	//------------------------------------------------------------------------
	// 十六进制字节序处理
	//------------------------------------------------------------------------
	r.AddLib(e, "hex", "ABCD", rulexlib.ABCD(e))
	r.AddLib(e, "hex", "DCBA", rulexlib.DCBA(e))
	r.AddLib(e, "hex", "BADC", rulexlib.BADC(e))
	r.AddLib(e, "hex", "CDAB", rulexlib.CDAB(e))
	//------------------------------------------------------------------------
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	// EEKIT
	r.AddLib(e, "eekit", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	r.AddLib(e, "eekit", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))
	// 树莓派4B
	r.AddLib(e, "raspi4b", "GPIOGet", rulexlib.RASPI4_GPIOGet(e))
	r.AddLib(e, "raspi4b", "GPIOSet", rulexlib.RASPI4_GPIOSet(e))
	// 玩客云WS1508
	r.AddLib(e, "ws1608", "GPIOGet", rulexlib.WKYWS1608_GPIOGet(e))
	r.AddLib(e, "ws1608", "GPIOSet", rulexlib.WKYWS1608_GPIOSet(e))
	//------------------------------------------------------------------------
	// AI BASE
	//------------------------------------------------------------------------
	r.AddLib(e, "aibase", "Infer", rulexlib.Infer(e))

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
