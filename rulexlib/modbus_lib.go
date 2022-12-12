package rulexlib

import (
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

//  --------------------------------------------
// |Function | Register Type
//  --------------------------------------------
// |	1	 | Read Coil
// |	2	 | Read Discrete Input
// |	3	 | Read Holding Registers
// |	4	 | Read Input Registers
// |	5	 | Write Single Coil
// |	6	 | Write Single Holding Register
// |	15	 | Write Multiple Coils
// |	16	 | Write Multiple Holding Registers
//  --------------------------------------------
/*
*
* Modbus Function1
*
 */
func F1(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function2
*
 */
func F2(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function3
*
 */
func F3(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function4
*
 */
func F4(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}
