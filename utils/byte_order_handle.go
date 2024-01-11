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

package utils

import (
	"encoding/hex"
	"fmt"
	"math"
)

/*
*
*解析 Modbus 的值 有符号,
注意：如果想解析值，必须不能超过4字节，目前常见的数一般都是4字节，也许后期会有8字节，但是目前暂时不支持
*
*/
func ParseSignedValue(DataBlockType string, DataBlockOrder string,
	Weight float32, byteSlice [4]byte) string {
	switch DataBlockType {
	case "RAW":
		{
			return hex.EncodeToString(byteSlice[:])
		}
	case "BYTE":
		{
			return fmt.Sprintf("%d", byteSlice[0])
		}
	case "SHORT":
		{
			// AB: 1234
			// BA: 3412
			if DataBlockOrder == "AB" {
				uint16Value := int16(byteSlice[3]) | int16(byteSlice[2])<<8
				return fmt.Sprintf("%d", uint16Value*int16(Weight))

			}
			if DataBlockOrder == "BA" {
				uint16Value := int16(byteSlice[2]) | int16(byteSlice[3])<<8
				return fmt.Sprintf("%d", uint16Value*int16(Weight))
			}

		}
	case "INT":
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := int32(byteSlice[0]) | int32(byteSlice[1])<<8 |
				int32(byteSlice[2])<<16 | int32(byteSlice[3])<<24
			return fmt.Sprintf("%d", intValue*int32(Weight))

		}
		if DataBlockOrder == "CDAB" {
			slice := [4]byte{}
			slice[0], slice[1] = byteSlice[2], byteSlice[3]
			slice[2], slice[3] = byteSlice[0], byteSlice[1]
			intValue := int32(slice[0]) | int32(slice[1])<<8 |
				int32(slice[2])<<16 | int32(slice[3])<<24
			return fmt.Sprintf("%d", intValue*int32(Weight))
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := int32(byteSlice[3]) | int32(byteSlice[2])<<8 |
				int32(byteSlice[1])<<16 | int32(byteSlice[0])<<24
			return fmt.Sprintf("%d", intValue*int32(Weight))
		}
	case "FLOAT": // 3.14159:DCBA -> 40490FDC
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := int32(byteSlice[0]) | int32(byteSlice[1])<<8 |
				int32(byteSlice[2])<<16 | int32(byteSlice[3])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue*Weight)
		}
		if DataBlockOrder == "CDAB" {
			intValue := int32(byteSlice[2]) | int32(byteSlice[3])<<8 |
				int32(byteSlice[0])<<16 | int32(byteSlice[1])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue*Weight)
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := int32(byteSlice[3]) | int32(byteSlice[2])<<8 |
				int32(byteSlice[1])<<16 | int32(byteSlice[0])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue*Weight)
		}
	default:
		return ParseUSignedValue(DataBlockType, DataBlockOrder, Weight, byteSlice)
	}
	return ""
}

/*
*
*解析西门子的值 无符号
*
 */
func ParseUSignedValue(DataBlockType string, DataBlockOrder string,
	Weight float32, byteSlice [4]byte) string {
	switch DataBlockType {
	case "USHORT":
		{
			// AB: 1234
			// BA: 3412
			if DataBlockOrder == "AB" {
				uint16Value := uint16(byteSlice[3]) | uint16(byteSlice[2])<<8
				return fmt.Sprintf("%d", uint16Value*uint16(Weight))

			}
			if DataBlockOrder == "BA" {
				uint16Value := uint16(byteSlice[2]) | uint16(byteSlice[3])<<8
				return fmt.Sprintf("%d", uint16Value*uint16(Weight))
			}

		}
	case "UINT":
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := uint32(byteSlice[0]) | uint32(byteSlice[1])<<8 |
				uint32(byteSlice[2])<<16 | uint32(byteSlice[3])<<24
			return fmt.Sprintf("%d", intValue*uint32(Weight))

		}
		if DataBlockOrder == "CDAB" {
			slice := [4]byte{}
			slice[0], slice[1] = byteSlice[2], byteSlice[3]
			slice[2], slice[3] = byteSlice[0], byteSlice[1]
			intValue := uint32(slice[0]) | uint32(slice[1])<<8 |
				uint32(slice[2])<<16 | uint32(slice[3])<<24
			return fmt.Sprintf("%d", intValue*uint32(Weight))
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := uint32(byteSlice[3]) | uint32(byteSlice[2])<<8 |
				uint32(byteSlice[1])<<16 | uint32(byteSlice[0])<<24
			return fmt.Sprintf("%d", intValue*uint32(Weight))
		}
	case "UFLOAT": // 3.14159:DCBA -> 40490FDC
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := int32(byteSlice[0]) | int32(byteSlice[1])<<8 |
				int32(byteSlice[2])<<16 | int32(byteSlice[3])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue*Weight)
		}
		if DataBlockOrder == "CDAB" {
			intValue := int32(byteSlice[2]) | int32(byteSlice[3])<<8 |
				int32(byteSlice[0])<<16 | int32(byteSlice[1])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", math.Abs(float64(floatValue*Weight)))
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := int32(byteSlice[3]) | int32(byteSlice[2])<<8 |
				int32(byteSlice[1])<<16 | int32(byteSlice[0])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", math.Abs(float64(floatValue*Weight)))
		}
	}
	return ""
}

/*
*
* 默认字节序
*
 */
func GetDefaultDataOrder(Type, Order string) string {
	if Order == "" {
		switch Type {
		case "INT", "UINT", "FLOAT", "UFLOAT":
			return "DCBA"
		case "BYTE", "I", "Q":
			return "A"
		case "SHORT", "USHORT":
			return "BA"
		case "LONG", "ULONG":
			return "HGFEDCBA"
		}
	}
	return Order
}

/*
*
* 处理空指针初始值
*
 */
func HandleZeroValue[V int16 | int32 | int64 | float32 | float64](v *V) *V {
	if v == nil {
		return new(V)
	}
	return v
}
