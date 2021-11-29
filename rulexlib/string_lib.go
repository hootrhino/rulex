package rulexlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

type StringLib struct {
}

func NewStringLib() typex.XLib {
	return &StringLib{}

}
func (l *StringLib) Name() string {
	return "String"
}
func (l *StringLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		return 0
	}
}
