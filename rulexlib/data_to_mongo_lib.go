package rulexlib

import (
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
func DataToMongo(e typex.RuleX, id string, data string) {
	handleDataFormat(e, id, data)
}

type MongoLib struct {
}

func NewMongoLib() typex.XLib {
	return &MongoLib{}
}
func (l *MongoLib) Name() string {
	return "DataToMongo"
}
func (l *MongoLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		DataToMongo(rx, id, data)
		return 0
	}
}
