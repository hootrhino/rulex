package engine

import (
	"time"

	"github.com/i4de/rulex/rulexlib"
	"github.com/i4de/rulex/typex"
)

/*
*
* 每次新增加函数以后记得在这里手动增加个，麻烦一点但是很简单.
*
 */
func BuildInLuaLibDoc() {
	currentTime := time.Now()
	var rulexlibDoc rulexlib.RulexLibDoc = rulexlib.RulexLibDoc{
		Name:        "RULEX-标准库文档",
		Version:     typex.DefaultVersion.Version,
		ReleaseTime: currentTime.Format("2006-01-02"),
	}
	// 消息转发
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "DataToHttp",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local result, err = rulexlib:DataToHttp("UUID", "HelloWorld")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "DataToMqtt",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "MQTT地址",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到MQTT服务器",
		Example:     `local result, err = rulexlib:DataToMqtt("UUID", "HelloWorld")`,
	})
	// JQ 这两个是同一个函数 为了兼容老版本的名字 留下第一个
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "JQ",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "JQ表达式",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "JSON格式的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "筛选后的数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "使用JQ筛选数据",
		Example:     `local result, err = rulexlib:JQ("[].name", "{"name":"value"}")`,
	})
	// 日志
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "log",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "格式化字符串",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "日志文本",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "打印日志到文本中",
		Example:     `local result, err = rulexlib:log("%s", "HelloWorld")`,
	})
	// 二进制操作
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "MB",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "二进制语法",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "字节流",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "二进制按位数匹配,'<' 表示小端，'>' 表示大端, 键之间要用空格",
		Example:     `local result, err = rulexlib:MB(">a:8 b:8", "0000000100000010")`,
	})
	// 二进制字符串转字节
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "B2BS",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "将字节转成诸如'0101010101010101'这样的字符串",
		Example:     `local result, err = rulexlib:B2BS("A")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "Bit",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "string",
				Description: "要处理的字节",
			},
			{
				Pos:         2,
				Type:        "Int",
				Description: "位置",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "提取一个字节的某个位, Pos最大不能大于8",
		Example:     `local result, err = rulexlib:Bit("A", 1)`,
	})
	// 字节转Int64
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "B2I64",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "字节流",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "Int64",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "字节转Int64",
		Example:     `local result, err = rulexlib:B2I64("....")`,
	})
	// 二进制字符串转字节
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "BS2B",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "将诸如'0101010101010101'这样的字符串转成真实的二进制值",
		Example:     `local result, err = rulexlib:BS2B("0101010101010101")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "HToN",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "十六进制字符串",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "十六进制格式的字符串转成数字",
		Example:     `local result, err = rulexlib:HToN("0F")`,
	})
	// 取某个Hex字符串的子串转换成数字
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "HsubToN",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "十六进制字符串",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "起点位置",
			},
			{
				Pos:         3,
				Type:        "String",
				Description: "结束位置",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "取某个Hex字符串的子串转换成数字",
		Example:     `local result, err = rulexlib:HsubToN("0a0b0c", 0, 2)`,
	})
	//
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "UrlBuildQS",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "Table",
				Description: "URL查询参数列表,必须是K-V列表",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "LuaTable",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "URL请求构建器",
		Example:     `local result, err = rulexlib:UrlBuildQS({name = 'rulex', 'age' = 0})`,
	})

	// 数据持久化
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "DataToTdEngine",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "TdEngine目标ID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到Tdengine",
		Example:     `local err = rulexlib:DataToTdEngine("UUID", "HelloWorld")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "DataToMongo",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local err = rulexlib:DataToMongo("UUID", "HelloWorld")`,
	})
	// 时间库
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "Time",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         0,
				Type:        "-",
				Description: "无参数",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "获取时间",
		Example:     `local result, err = rulexlib:Time()`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "TsUnix",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         0,
				Type:        "-",
				Description: "无参数",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "获取时间",
		Example:     `local result, err = rulexlib:TsUnix()`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "TsUnixNano",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         0,
				Type:        "-",
				Description: "无参数",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "获取时间",
		Example:     `local result, err = rulexlib:TsUnixNano()`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "NtpTime",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         0,
				Type:        "-",
				Description: "无参数",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "获取时间",
		Example:     `local result, err = rulexlib:NtpTime()`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "T2J",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "Table",
				Description: "LUA Table",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "Lua Table转成JSON",
		Example:     `local result, err = rulexlib:T2J({k = 'v'})`,
	})
	// JSON转成Lua Table
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "J2T",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "LUA Table",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "Table",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "JSON转成Lua Table",
		Example:     `local result, err = rulexlib:J2T("{"k" = "v"}")`,
	})
	// Get Rule ID
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "RUUID",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         0,
				Type:        "-",
				Description: "无参数",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "获取当前RULE的UUID",
		Example:     `local result, err = rulexlib:RUUID()`,
	})
	// Codec
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "RPCENC",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "RPC编码",
		Example:     `local result, err = rulexlib:RPCENC("UUID", "HelloWorld")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "RPCDEC",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "RPC解码",
		Example:     `local result, err = rulexlib:RPCDEC("UUID", "HelloWorld")`,
	})
	// Device R/W
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "ReadDevice",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local result, err = rulexlib:ReadDevice("UUID", "HelloWorld")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "WriteDevice",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local result, err = rulexlib:WriteDevice("UUID", "HelloWorld")`,
	})
	// Source R/W
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "ReadSource",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local result, err = rulexlib:ReadSource("UUID", "HelloWorld")`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "WriteSource",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "UUID",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "需要处理的数据",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "该函数将数据推送到HTTP接口",
		Example:     `local result, err = rulexlib:WriteSource("UUID", "HelloWorld")`,
	})
	// String
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "T2Str",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "Table",
				Description: "要转的Table",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回数据",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: " Table 转成 String",
		Example:     `local result, err = rulexlib:T2Str({1,2,3})`,
	})
	rulexlibDoc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "Throw",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "抛出异常",
		Example:     `local err = rulexlib:Throw("Stack Overflow")`,
	})
	rulexlibDoc.BuildDoc()
}
