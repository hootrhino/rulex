package engine

import (
	"rulex/rulexlib"
	"rulex/typex"
)

func LoadBuildInLuaLib(e typex.RuleX, r *typex.Rule) {
	//
	// Load Stdlib
	//--------------------------------------------------------------
	// 消息转发
	r.LoadLib(e, rulexlib.NewHttpLib())
	r.LoadLib(e, rulexlib.NewMqttLib())
	// JQ
	r.LoadLib(e, rulexlib.NewJqLib())
	// 日志
	r.LoadLib(e, rulexlib.NewLogLib())
	// 直达数据
	r.LoadLib(e, rulexlib.NewWriteInStreamLib())
	r.LoadLib(e, rulexlib.NewWriteOutStreamLib())
	// 二进制操作
	r.LoadLib(e, rulexlib.NewMatchBinaryLib())
	r.LoadLib(e, rulexlib.NewByteToBitStringLib())
	r.LoadLib(e, rulexlib.NewGetABitOnByteLib())
	r.LoadLib(e, rulexlib.NewByteToInt64Lib())
	r.LoadLib(e, rulexlib.NewBitStringToBytesLib())
	// JSON编解码
	r.LoadLib(e, rulexlib.NewJsonEncodeLib())
	r.LoadLib(e, rulexlib.NewJsonDecodeLib())
	// URL处理
	r.LoadLib(e, rulexlib.NewUrlBuildLib())
	r.LoadLib(e, rulexlib.NewUrlBuildQSLib())
	r.LoadLib(e, rulexlib.NewUrlParseLib())
	r.LoadLib(e, rulexlib.NewUrlResolveLib())
	// 数据持久化
	r.LoadLib(e, rulexlib.NewTdEngineLib())
	r.LoadLib(e, rulexlib.NewMongoLib())
	// From 0.0.8: 使用新版本的库加载方式
	// 时间库
	r.AddLib(e, "Time", rulexlib.Time(e))
	r.AddLib(e, "TsUnix", rulexlib.TsUnix(e))
	r.AddLib(e, "TsUnixNano", rulexlib.TsUnixNano(e))
	// 缓存器库
	r.AddLib(e, "StoreGet", rulexlib.StoreGet(e))
	r.AddLib(e, "StoreGet", rulexlib.StoreGet(e))
	r.AddLib(e, "StoreDelete", rulexlib.StoreDelete(e))
}
