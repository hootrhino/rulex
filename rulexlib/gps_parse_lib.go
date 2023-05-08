package rulexlib

import (
	"github.com/adrianmo/go-nmea"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* 解析GPS数据："$GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70"
*
 */
func ParseGPS(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		sentence := l.ToString(2)
		ss, err := nmea.Parse(sentence)
		if err != nil {
			glogger.GLogger.Fatal(err)
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
			return 2
		}
		l.Push(lua.LString(ss.String()))
		l.Push(lua.LNil)
		return 2
	}
}
