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
	"testing"
)

/*
*
* CRC 计算
*
 */

func calculateCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF

	for _, b := range data {
		crc ^= uint16(b)

		for i := 0; i < 8; i++ {
			lsb := crc & 0x0001
			crc >>= 1

			if lsb == 1 {
				crc ^= 0xA001 // 0xA001 是Modbus CRC16多项式的表示
			}
		}
	}

	return crc
}
func uint16ToBytes(val uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, val)
	return bytes
}

// go test -timeout 30s -run ^Test_calculateCRC github.com/hootrhino/rulex/test -v -count=1
func Test_calculateCRC(t *testing.T) {
	// [132 10] => 84 0A
	t.Log(uint16ToBytes(calculateCRC16([]byte{01, 03, 00, 00, 00, 01})))
}
