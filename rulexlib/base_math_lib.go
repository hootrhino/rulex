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
	"math"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
)

func truncateFloat(number float64, decimalPlaces int) float64 {
	scale := math.Pow(10, float64(decimalPlaces))
	result := math.Floor(number*scale) / scale
	return result
}

/*
*
* 取小数位 applib:Float(number, decimalPlaces) -> float
*
 */
func TruncateFloat(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		number := l.ToNumber(2)
		decimalPlaces := l.ToInt(3)
		l.Push(lua.LNumber(truncateFloat(float64(number), decimalPlaces)))
		return 1
	}
}
