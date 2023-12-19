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

package test

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"testing"
)

// go test -timeout 30s -run ^Test_Binary_to_Float github.com/hootrhino/rulex/test -v -count=1
func Test_Binary_to_Float(t *testing.T) {
	// Hex 40 49 0E 56 = D3.141
	dBytes, _ := hex.DecodeString("40490E56")
	V1 := math.Float32frombits(binary.LittleEndian.Uint32(dBytes))
	t.Log(V1)
	V2 := math.Float32frombits(binary.BigEndian.Uint32(dBytes))
	t.Log(V2)
}
