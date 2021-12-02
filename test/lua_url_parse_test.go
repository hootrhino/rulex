package test

import (
	"testing"

	"github.com/cjoudrey/gluaurl"
	"github.com/yuin/gopher-lua"
)

func Test_Url_parse_test(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("url", gluaurl.Loader)

	if err := L.DoString(`
			local url = require("url")
			parsed_url = url.parse("http://example.com/")
			print(parsed_url.host)
    `); err != nil {
		panic(err)
	}
}
