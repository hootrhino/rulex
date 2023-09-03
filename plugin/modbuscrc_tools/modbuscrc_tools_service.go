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

package modbuscrctools

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/hootrhino/rulex/typex"
)

func (mc *modbusCRCCalculator) Service(arg typex.ServiceArg) typex.ServiceResult {
	// 大端
	if arg.Name == "crc16big" {
		switch s := arg.Args.(type) {
		case string:
			{
				bytes, err := hex.DecodeString(s)
				if err != nil {
					return typex.ServiceResult{Out: []map[string]interface{}{
						{"error": err.Error()},
					}}

				}
				return typex.ServiceResult{Out: []map[string]interface{}{
					{"value": fmt.Sprintf("%02x", uint16ToBytes(calculateCRC16(bytes)))},
				}}
			}
		}
	}
	// 小端
	if arg.Name == "crc16little" {
		switch s := arg.Args.(type) {
		case string:
			{
				bytes, err := hex.DecodeString(stringReverse(s))
				if err != nil {
					return typex.ServiceResult{Out: []map[string]interface{}{
						{"error": err.Error()},
					}}

				}
				return typex.ServiceResult{Out: []map[string]interface{}{
					{"value": fmt.Sprintf("%02x", uint16ToBytes(calculateCRC16(bytes)))},
				}}
			}
		}
	}

	return typex.ServiceResult{Out: []map[string]interface{}{
		{"error": "Unsupported operate:" + arg.Name},
	}}
}

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
				crc ^= 0xA001
			}
		}
	}
	return crc
}
func uint16ToBytes(val uint16) [2]byte {
	bytes := [2]byte{}
	binary.LittleEndian.PutUint16(bytes[:], val)
	return bytes
}
func stringReverse(in string) string {
	var bytes []byte = []byte(in)
	for i := 0; i < len(in)/2; i++ {
		tmp := bytes[len(in)-i-1]
		bytes[len(in)-i-1] = bytes[i]
		bytes[i] = tmp
	}
	return string(bytes)
}
