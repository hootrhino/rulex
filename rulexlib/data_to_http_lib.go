package rulexlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

type HttpLib struct {
}

func NewHttpLib() typex.XLib {
	return &HttpLib{}
}

//
//
//
func DataToHttpServer(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

func (l *HttpLib) Name() string {
	return "DataToHttpServer"
}
func (l *HttpLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		DataToHttpServer(rx, id, data)
		return 0
	}
}
