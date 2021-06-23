package extralib

import lua "github.com/yuin/gopher-lua"

// TODO Load libs

func LoadLibs(l *lua.LState) {
	l.PreloadModule("decodelib", LoadDecodeLib)
	l.PreloadModule("encodelib", LoadEncodeLib)
	l.PreloadModule("sqllib", LoadSqlLib)
	OpenSubset(l)
}

// OpenSubset
func OpenSubset(l *lua.LState) {
	for _, pair := range []struct {
		name string
		f    lua.LGFunction
	}{
		{lua.LoadLibName, lua.OpenPackage},
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
	} {
		if err := l.CallByParam(lua.P{
			Fn:      l.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, lua.LString(pair.name)); err != nil {
			panic(err)
		}
	}
}
