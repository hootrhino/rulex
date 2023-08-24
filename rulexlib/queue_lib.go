// Copyright (C) 2023 wwhai
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
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 对外输出数据，这个函数的数据会被推到 yqueue 里面，然后yqueue会再次把数据推到其实现的各种pipe里，
* 可能是个 websocket，也可能是个 TCP 或者 UDP Server。
*
 */
func Output(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		// TODO
		return 0
	}
}
