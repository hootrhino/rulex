package core

import (
	"fmt"
	"strings"
)

/*
*
* 自动生成解码器
*
 */
type Field struct {
	Name string
	Type string
	Len  uint8
}

var actions = `
Actions =
{
function(data)
%v
	return true, data
end
}
`

type GenLuaConfig struct {
	Big    bool    `json:"big"`
	More   bool    `json:"more"`
	Fields []Field `json:"fields"`
}

func GenCode(fields []Field, big bool, more bool) string {
	expr := __b(big)
	for _, field := range fields {
		expr += fmt.Sprintf("%v:%v ", field.Name, field.Len)
	}
	lua := fmt.Sprintf("\tlocal table = rulexlib:MatchBinary('%v', data, false)\n", strings.TrimSuffix(expr, " "))
	for _, field := range fields {
		lua += fmt.Sprintf("\ttable['%v'] = rulexlib:BTo%v(1, rulexlib:BSToB(tb['%v']))\n", field.Name, field.Type, field.Name)
	}
	return fmt.Sprintf(actions, strings.TrimSuffix(lua, "\n"))
}
func __b(b bool) string {
	if b {
		return ">"
	}
	return "<"
}
