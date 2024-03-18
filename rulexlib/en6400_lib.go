// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	archsupport "github.com/hootrhino/rulex/bspsupport"
	"github.com/hootrhino/rulex/typex"
)

func EN6400_LedOn(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.EN6400_GPIO231Set(int(1))
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}
func EN6400_LedOff(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.EN6400_GPIO231Set(int(0))
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}
