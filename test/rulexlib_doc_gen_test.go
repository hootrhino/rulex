package test

import (
	"testing"

	"github.com/i4de/rulex/rulexlib"
)

func Test_Gen_rulexlib_doc(t *testing.T) {
	doc := rulexlib.RulexLibDoc{
		Name:        "RULEX-标准库文档",
		Version:     "V1.0.0",
		ReleaseTime: "2022-07-30",
	}
	doc.AddFunc(rulexlib.Fun{
		NameSpace: "rulexlib",
		FunName:   "Sprintf",
		FunArgs: []rulexlib.FunArg{
			{
				Pos:         1,
				Type:        "String",
				Description: "格式化字符串",
			},
			{
				Pos:         2,
				Type:        "Any",
				Description: "格式化的值",
			},
		},
		ReturnValue: []rulexlib.ReturnValue{
			{
				Pos:         1,
				Type:        "String",
				Description: "返回字符串",
			},
			{
				Pos:         2,
				Type:        "String",
				Description: "Error信息",
			},
		},
		Description: "格式化文本",
		Example:     `local s, err = rulexlib:Sprintf("Hello%s", "World")`,
	})
	doc.BuildDoc()
}
