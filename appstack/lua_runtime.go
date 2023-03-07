package appstack

import (
	"fmt"

	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
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
	// 检查名称
	AppNAME := tempVm.GetGlobal("AppNAME")
	if AppNAME == nil {
		return fmt.Errorf("'AppNAME' field not exists")
	}
	if AppNAME.Type() != lua.LTString {
		return fmt.Errorf("'AppNAME' must be string")
	}
	// 检查类型
	AppVERSION := tempVm.GetGlobal("AppVERSION")
	if AppVERSION == nil {
		return fmt.Errorf("'AppVERSION' field not exists")
	}
	if AppVERSION.Type() != lua.LTString {
		return fmt.Errorf("'AppVERSION' must be string")
	}
	// 检查描述信息
	AppDESCRIPTION := tempVm.GetGlobal("AppDESCRIPTION")
	if AppDESCRIPTION == nil {
		if AppDESCRIPTION.Type() != lua.LTString {
			return fmt.Errorf("'AppDESCRIPTION' must be string")
		}
	}

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
func addAppLib(app *typex.Application,
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
	addAppLib(app, e, "applib", "DataToHttp", rulexlib.DataToHttp(e))
	addAppLib(app, e, "applib", "DataToMqtt", rulexlib.DataToMqtt(e))
	// JQ
	addAppLib(app, e, "applib", "JqSelect", rulexlib.JqSelect(e))
	addAppLib(app, e, "applib", "JQ", rulexlib.JqSelect(e))
	// 日志
	addAppLib(app, e, "applib", "log", rulexlib.Log(e))
	// 二进制操作
	addAppLib(app, e, "applib", "MB", rulexlib.MatchBinary(e))
	addAppLib(app, e, "applib", "B2BS", rulexlib.ByteToBitString(e))
	addAppLib(app, e, "applib", "Bit", rulexlib.GetABitOnByte(e))
	addAppLib(app, e, "applib", "B2I64", rulexlib.ByteToInt64(e))
	addAppLib(app, e, "applib", "B64S2B", rulexlib.B64S2B(e))
	addAppLib(app, e, "applib", "BS2B", rulexlib.BitStringToBytes(e))
	addAppLib(app, e, "applib", "HToN", rulexlib.HToN(e))
	addAppLib(app, e, "applib", "HsubToN", rulexlib.HsubToN(e))
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
	addAppLib(app, e, "applib", "TsUnix", rulexlib.TsUnix(e))
	addAppLib(app, e, "applib", "TsUnixNano", rulexlib.TsUnixNano(e))
	addAppLib(app, e, "applib", "NtpTime", rulexlib.NtpTime(e))
	// 缓存器库
	addAppLib(app, e, "applib", "VSet", rulexlib.StoreSet(e))
	addAppLib(app, e, "applib", "VGet", rulexlib.StoreGet(e))
	addAppLib(app, e, "applib", "VDel", rulexlib.StoreDelete(e))
	// JSON
	addAppLib(app, e, "applib", "T2J", rulexlib.JSONE(e)) // Lua Table -> JSON
	addAppLib(app, e, "applib", "J2T", rulexlib.JSOND(e)) // JSON -> Lua Table
	// Get Rule ID
	addAppLib(app, e, "applib", "RUUID", rulexlib.SelfRuleUUID(e, app.UUID))
	// Codec
	addAppLib(app, e, "applib", "RPCENC", rulexlib.RPCEncode(e))
	addAppLib(app, e, "applib", "RPCDEC", rulexlib.RPCDecode(e))
	// Device R/W
	addAppLib(app, e, "applib", "ReadDevice", rulexlib.ReadDevice(e))
	addAppLib(app, e, "applib", "WriteDevice", rulexlib.WriteDevice(e))
	// Source R/W
	addAppLib(app, e, "applib", "ReadSource", rulexlib.ReadSource(e))
	addAppLib(app, e, "applib", "WriteSource", rulexlib.WriteSource(e))
	// String
	addAppLib(app, e, "applib", "T2Str", rulexlib.T2Str(e))
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
	// 注册GPIO操作函数到LUA运行时
	//------------------------------------------------------------------------
	addAppLib(app, e, "eekit", "GPIOGet", rulexlib.EEKIT_GPIOGet(e))
	addAppLib(app, e, "eekit", "GPIOSet", rulexlib.EEKIT_GPIOSet(e))

}
