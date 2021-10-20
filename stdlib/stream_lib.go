package stdlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

type WriteInStreamLib struct {
}

func NewWriteInStreamLib() typex.XLib {
	return &HttpLib{}
}
func (l *WriteInStreamLib) Name() string {
	return "WriteInStream"
}
func (l *WriteInStreamLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		WriteInStream(rx, id, data)
		return 0
	}
}

//
type WriteOutStreamLib struct {
}

func NewWriteOutStreamLib() typex.XLib {
	return &HttpLib{}
}
func (l *WriteOutStreamLib) Name() string {
	return "WriteOutStream"
}
func (l *WriteOutStreamLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		WriteOutStream(rx, id, data)
		return 0
	}
}

func WriteInStream(e typex.RuleX, id string, data string) {
	e.GetInEnd(id).Resource.OnStreamApproached(data)
}
func WriteOutStream(e typex.RuleX, id string, data string) {
	e.GetOutEnd(id).Target.OnStreamApproached(data)
}
