package utils

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
function(args)
%v
	return true, args
end
}
`

type GenLuaConfig struct {
	Big    bool    `json:"big"`    // 大小端
	More   bool    `json:"more"`   // 需要剩下的字节？
	Fields []Field `json:"fields"` // 字段列表
}

func GenCode(fields []Field, big bool, more bool) string {
	expr := __b(big)
	for _, field := range fields {
		expr += fmt.Sprintf("%v:%v ", field.Name, field.Len)
	}
	lua := fmt.Sprintf("\tlocal table = rulexlib:MB('%v', data, false)\n", strings.TrimSuffix(expr, " "))
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
