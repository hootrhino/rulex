package test

import (
	"os"
	"testing"
	"text/template"
)

var s = `
---@diagnostic disable: undefined-global
-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions = {function(args)
    local t = rulexlib:J2T(data)
    local V0 = rulexlib:MB(">{{.a}}:16 {{.b}}:16 {{.c}}:16 {{.d}}:16 {{.e}}:16", t['value'], false)
    local a = rulexlib:T2J(V0['{{.a}}'])
    local b = rulexlib:T2J(V0['{{.b}}'])
    local c = rulexlib:T2J(V0['{{.c}}'])
    local d = rulexlib:T2J(V0['{{.d}}'])
    local e = rulexlib:T2J(V0['{{.e}}'])
    print('{{.a}} ==> ', {{.a}}, ' ->', rulexlib:B2I64('>', rulexlib:BS2B(a)))
    print('{{.b}} ==> ', {{.b}}, ' ->', rulexlib:B2I64('>', rulexlib:BS2B(b)))
    print('{{.c}} ==> ', {{.c}}, ' ->', rulexlib:B2I64('>', rulexlib:BS2B(c)))
    print('{{.d}} ==> ', {{.d}}, ' ->', rulexlib:B2I64('>', rulexlib:BS2B(d)))
    print('{{.e}} ==> ', {{.e}}, ' ->', rulexlib:B2I64('>', rulexlib:BS2B(e)))
    return true, args
end}

`

func Test_gen_template(*testing.T) {
	t := template.New("test")
	t = template.Must(t.Parse(s))

	t.Execute(os.Stdout, map[string]string{
		"a": "va",
		"b": "vb",
		"c": "vc",
		"d": "vd",
		"e": "ve",
	})
}
