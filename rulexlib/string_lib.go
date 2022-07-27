package rulexlib

import (
	"strings"

	"github.com/i4de/rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

func T2Str(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		table := l.ToTable(2)
		args := []string{}
		table.ForEach(func(l1, value lua.LValue) {
			args = append(args, value.String())
		})
		r := strings.Join(args, "")
		l.Push(lua.LString(r))
		return 1
	}
}
