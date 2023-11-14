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
	{
		// 消息转发\数据持久化\编解码
		r.AddLib(e, "data", "ToHttp", rulexlib.DataToHttp(e))
		r.AddLib(e, "data", "ToMqtt", rulexlib.DataToMqtt(e))
		r.AddLib(e, "data", "ToUdp", rulexlib.DataToUdp(e))
		r.AddLib(e, "data", "ToTcp", rulexlib.DataToTcp(e))
		r.AddLib(e, "data", "ToTdEngine", rulexlib.DataToTdEngine(e))
		r.AddLib(e, "data", "ToMongo", rulexlib.DataToMongo(e))
		r.AddLib(e, "data", "ToNats", rulexlib.DataToNats(e))
	}
	{
		r.AddLib(e, "rpc", "Request", rulexlib.Request(e))
	}
	{
		// JQ
		r.AddLib(e, "jq", "Execute", rulexlib.JqSelect(e))
	}
	{
		// 标准库
		r.AddLib(e, "stdlib", "Debug", rulexlib.DebugRule(e, r.UUID))
		r.AddLib(e, "stdlib", "Println", rulexlib.Println(e))
		r.AddLib(e, "stdlib", "Throw", rulexlib.Throw(e))
		r.AddLib(e, "_G", "Debug", rulexlib.DebugRule(e, r.UUID))
		r.AddLib(e, "_G", "Throw", rulexlib.Throw(e))
	}
	{
		// 二进制操作
		r.AddLib(e, "binary", "MB", rulexlib.MatchBinary(e))
		r.AddLib(e, "binary", "MBHex", rulexlib.MatchBinaryHex(e))
		r.AddLib(e, "binary", "B2BS", rulexlib.ByteToBitString(e))
		r.AddLib(e, "binary", "Bit", rulexlib.GetABitOnByte(e))
		r.AddLib(e, "binary", "B2I64", rulexlib.ByteToInt64(e))
		r.AddLib(e, "binary", "B64S2B", rulexlib.B64S2B(e))
		r.AddLib(e, "binary", "BS2B", rulexlib.BitStringToBytes(e))
		r.AddLib(e, "binary", "Bin2F32", rulexlib.BinToFloat32(e))
		r.AddLib(e, "binary", "Bin2F64", rulexlib.BinToFloat64(e))
	}
	{
		// URL处理
		r.AddLib(e, "url", "UrlBuild", rulexlib.UrlBuild(e))
		r.AddLib(e, "url", "UrlBuildQS", rulexlib.UrlBuildQS(e))
		r.AddLib(e, "url", "UrlParse", rulexlib.UrlParse(e))
		r.AddLib(e, "url", "UrlResolve", rulexlib.UrlResolve(e))
	}
	{
		// 时间库
		r.AddLib(e, "time", "Time", rulexlib.Time(e))
		r.AddLib(e, "time", "TimeMs", rulexlib.TimeMs(e))
		r.AddLib(e, "time", "TsUnix", rulexlib.TsUnix(e))
		r.AddLib(e, "time", "TsUnixNano", rulexlib.TsUnixNano(e))
		r.AddLib(e, "time", "NtpTime", rulexlib.NtpTime(e))
		r.AddLib(e, "time", "Sleep", rulexlib.Sleep(e))
	}

	{
		// 缓存器库
		r.AddLib(e, "kv", "VSet", rulexlib.StoreSet(e))
		r.AddLib(e, "kv", "VGet", rulexlib.StoreGet(e))
		r.AddLib(e, "kv", "VDel", rulexlib.StoreDelete(e))
	}
	{
		// JSON
		r.AddLib(e, "json", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
		r.AddLib(e, "json", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	}
	// Get Rule ID
	r.AddLib(e, "rule", "SelfId", rulexlib.SelfRuleUUID(e, r.UUID))
	{
		// Device R/W
		r.AddLib(e, "device", "ReadDevice", rulexlib.ReadDevice(e))
		r.AddLib(e, "device", "WriteDevice", rulexlib.WriteDevice(e))
		// Source R/W
		r.AddLib(e, "device", "ReadSource", rulexlib.ReadSource(e))
		r.AddLib(e, "device", "WriteSource", rulexlib.WriteSource(e))
	}
	{
		// String
		r.AddLib(e, "string", "T2Str", rulexlib.T2Str(e))
		r.AddLib(e, "string", "Bin2Str", rulexlib.Bin2Str(e))
	}
	//------------------------------------------------------------------------
	// 十六进制编码处理
	//------------------------------------------------------------------------
	{
		r.AddLib(e, "hex", "Bytes2Hexs", rulexlib.Bytes2Hexs(e))
		r.AddLib(e, "hex", "Hexs2Bytes", rulexlib.Hexs2Bytes(e))
		r.AddLib(e, "hex", "ABCD", rulexlib.ABCD(e))
		r.AddLib(e, "hex", "DCBA", rulexlib.DCBA(e))
		r.AddLib(e, "hex", "BADC", rulexlib.BADC(e))
		r.AddLib(e, "hex", "CDAB", rulexlib.CDAB(e))
		r.AddLib(e, "hex", "HexToNum", rulexlib.HToN(e))
		r.AddLib(e, "hex", "HsubToN", rulexlib.HsubToN(e))
		r.AddLib(e, "hex", "MatchHex", rulexlib.MatchHex(e))
		r.AddLib(e, "hex", "MatchUInt", rulexlib.MatchUInt(e))
	}
	{
		// DO1 DO2
		r.AddLib(e, "rhinopi", "DO1Set", rulexlib.H3DO1Set(e))
		r.AddLib(e, "rhinopi", "DO1Get", rulexlib.H3DO1Get(e))
		r.AddLib(e, "rhinopi", "DO2Set", rulexlib.H3DO2Set(e))
		r.AddLib(e, "rhinopi", "DO2Get", rulexlib.H3DO2Get(e))
		// DI1 DI2 DI3
		r.AddLib(e, "rhinopi", "DI1Get", rulexlib.H3DI1Get(e))
		r.AddLib(e, "rhinopi", "DI2Get", rulexlib.H3DI2Get(e))
		r.AddLib(e, "rhinopi", "DI3Get", rulexlib.H3DI3Get(e))

	}
	// Modbus
	{
		r.AddLib(e, "modbus", "F5", rulexlib.F5(e))
		r.AddLib(e, "modbus", "F6", rulexlib.F6(e))
		r.AddLib(e, "modbus", "F15", rulexlib.F15(e))
		r.AddLib(e, "modbus", "F16", rulexlib.F16(e))
	}
	//------------------------------------------------------------------------
	// AI BASE
	//------------------------------------------------------------------------
	r.AddLib(e, "aibase", "Infer", rulexlib.Infer(e))
	//------------------------------------------------------------------------
	// Math
	//------------------------------------------------------------------------
	r.AddLib(e, "math", "TFloat", rulexlib.TruncateFloat(e))
	// LocalDBQuery
	r.AddLib(e, "datacenter", "DBQuery", rulexlib.LocalDBQuery(e))

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
