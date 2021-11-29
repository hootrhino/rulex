package rulexlib

import (
	"rulex/typex"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

type LogLib struct {
}

func NewLogLib() typex.XLib {
	return &LogLib{}
}
func (l *LogLib) Name() string {
	return "log"
}
func (l *LogLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		content := l.ToString(2)
		log.Info("[CALLBACK]" + content)
		return 0
	}
}
