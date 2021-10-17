package stdlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

type StringLib struct {
}

func (l *StringLib) LoadLib(name string, e typex.RuleX, L *lua.LState) error {
	return nil
}
func (l *StringLib) UnLoadLib(name string) error {
	return nil
}
