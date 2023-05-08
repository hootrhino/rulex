package typex

import lua "github.com/hootrhino/gopher-lua"

// XLib: 库函数接口; TODO: V0.1.2废弃
type XLib interface {
	Name() string
	LibFun(RuleX) func(*lua.LState) int
}
