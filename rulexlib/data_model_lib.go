package rulexlib

import (
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* 改变模型值
*
 */
func SetModelValue(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		name := l.ToString(3)
		value := l.ToString(4)
		setValue(rx, uuid, name, value)
		return 0
	}
}

/*
*
* 改变值
*
 */
func setValue(rx typex.RuleX, uuid, name, value string) {

	in := rx.GetInEnd(uuid)
	if in != nil {
		DataModel := in.DataModelsMap[name]
		DataModel.Value = value
		in.DataModelsMap[name] = DataModel
	}
}
