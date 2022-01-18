package rulexlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
func DataToTdEngineServer(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

type TdEngineLib struct {
}

func NewTdEngineLib() typex.XLib {
	return &TdEngineLib{}
}
func (l *TdEngineLib) Name() string {
	return "DataToTdEngineServer"
}
func (l *TdEngineLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		//
		// SQL: INSERT INTO meter VALUES (NOW, %v, %v....);
		//
		data := l.ToString(3) // Data must arrays [1,2,3,4....]
		DataToTdEngineServer(rx, id, data)
		return 0
	}
}
