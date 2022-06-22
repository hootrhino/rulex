package test

import (
	"encoding/json"
	"testing"

	"github.com/i4de/rulex/core"
)

func Test_gen_lua_code(t *testing.T) {
	c := core.GenLuaConfig{
		Big:  true,
		More: false,
		Fields: []core.Field{
			{
				Name: "a",
				Type: "Int",
				Len:  1,
			},
			{
				Name: "b",
				Type: "String",
				Len:  2,
			},
			{
				Name: "c",
				Type: "Float",
				Len:  9,
			},
		},
	}
	b, _ := json.MarshalIndent(c, "", "  ")
	t.Log(string(b))
	t.Log(core.GenCode(c.Fields, c.Big, c.More))
}
