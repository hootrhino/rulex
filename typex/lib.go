package typex

import lua "github.com/yuin/gopher-lua"

//
// XLib: 库函数接口
//
type XLib interface {
	Name() string
	LibFun(RuleX) func(*lua.LState) int
}
