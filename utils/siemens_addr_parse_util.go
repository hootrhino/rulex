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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

/*
*
- 符号地址（Symbolic Addressing）：
使用符号名称来表示变量或输入/输出地址。这种方式更加直观和易于理解，适用于高级编程语言和工程师使用。例如，可以使用变量名"MotorSpeed"或输入名"I1"来表示对应的地址。
- 基于字节的地址（Byte-based Addressing）：
使用字节地址和位地址的组合来表示变量或输入/输出地址。字节地址表示内存中的字节偏移，而位地址表示字节中的位偏移。例如，使用地址"DB1.DBX10.3"表示数据块1中偏移为10的字节的第3位。
- 基于字的地址（SHORT-based Addressing）：
类似于基于字节的地址，但是将地址表示为字（16位）的偏移。例如，使用地址"DB1.DBD20"表示数据块1中偏移为20的字。
- 基于地址区域的地址（Address Area-based Addressing）：
将地址按照不同的区域进行划分，如输入区域（I），输出区域（Q），数据块区域（DB）等。每个区域都有特定的地址范围。例如，使用地址"I10.3"表示输入区域的第10个输入的第3位。
*
*/

// AddressInfo 包含解析后的地址信息
type AddressInfo struct {
	AddressType     string // 寄存器类型: DB I Q
	DataBlockType   string // 数据类型: BYTE SHORT INT
	DataBlockSize   int    // 数据长度
	DataBlockOrder  string // 字节序
	DataBlockNumber int    // 数据块号
	ElementNumber   int    // 元素号
	BitNumber       int    // 位号，只针对I、Q
}

func (O AddressInfo) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

/*
*
* 解析地址
*
 */
func ParseSiemensDB(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	if len(s) < 3 {
		return AddressInfo, fmt.Errorf("invalid Address format:%s", s)
	}
	if s[:2] == "I." {
		return _ParseADDR_I(s)
	}
	if s[:2] == "Q." {
		return _ParseADDR_Q(s)
	}
	if s[:2] == "IB" {
		return _ParseADDR_IB(s)
	}
	if s[:2] == "QB" {
		return _ParseADDR_QB(s)
	}
	if s[:2] == "DB" {
		return _ParseDB_DX(s)
	}
	return AddressInfo, nil

}

// 解析DB: DB4900.DBD2108
// ^DB\d+\.\w+\.\d+$
func _ParseDB_DX(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	Len := len(s)
	if Len < 3 {
		return AddressInfo, fmt.Errorf("Invalid Address length:%s", s)
	}
	// DB4900.DBD2108
	if s[:2] == "DB" {
		AddressInfo.AddressType = "DB"
		parts := strings.Split(s[2:], ".")
		// 4900.DBD2108
		if len(parts) != 2 {
			return AddressInfo, fmt.Errorf("Invalid Element Number Address format:%s", s)
		}
		DataBlockNumber, err1 := strconv.Atoi(parts[0])
		if err1 != nil {
			return AddressInfo, fmt.Errorf("Address %s Atoi failed:%s", s, err1)
		}
		AddressInfo.DataBlockNumber = DataBlockNumber
		// DBD2108
		if len(parts[1]) < 4 {
			return AddressInfo, fmt.Errorf("Invalid Element Number Address format:%s", s)
		}
		ElementNumber, err2 := strconv.Atoi(parts[1][3:]) //DBD...
		if err2 != nil {
			return AddressInfo, fmt.Errorf("Element Number %s Atoi failed:%s", s, err2)
		}
		switch parts[1][2] {

		case 'D': // DBD: 4字节
			AddressInfo.DataBlockSize = 4
		case 'W': // DBw: 2字节
			AddressInfo.DataBlockSize = 2
		case 'B': // DBB: 1字节
			AddressInfo.DataBlockSize = 1
		case 'X': // DBX: 1字节
			AddressInfo.DataBlockSize = 1
		default:
			return AddressInfo, fmt.Errorf("Invalid Element Type:%s", parts[1][2])
		}
		AddressInfo.ElementNumber = ElementNumber
	}
	return AddressInfo, nil
}

// 解析I格式
// I（输入）寄存器：
// 读取I寄存器：I + 编号，例如 I0.0、I1.1 等。
// 写入I寄存器：I + 编号，例如 I0.0、I1.1 等。

func _ParseADDR_I(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	if len(s) < 2 {
		return AddressInfo, fmt.Errorf("invalid Address format:%s", s)
	}
	// I0.0
	parts := strings.Split(s[1:], ".")
	// 0.0
	if len(parts) != 2 {
		return AddressInfo, fmt.Errorf("Invalid ElementNumber Address format:%s", s)
	}
	//
	DataBlockNumber, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		return AddressInfo, fmt.Errorf("Address %s Atoi failed:%s", s, err1)
	}
	BitNumber, err2 := strconv.Atoi(parts[0])
	if err2 != nil {
		return AddressInfo, fmt.Errorf("Address %s Atoi failed:%s", s, err2)
	}
	AddressInfo.AddressType = "I"
	AddressInfo.DataBlockType = "BYTE"
	AddressInfo.DataBlockNumber = DataBlockNumber
	AddressInfo.BitNumber = BitNumber
	return AddressInfo, nil
}

// IB0
func _ParseADDR_IB(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	AddressInfo.AddressType = "IB"
	DataBlockNumber := 0
	AddressInfo.DataBlockNumber = DataBlockNumber
	BitNumber := 0
	AddressInfo.BitNumber = BitNumber
	return AddressInfo, nil
}

// 解析Q格式
// Q（输出）寄存器：
// 读取Q寄存器：Q + 编号，例如 Q0.0、Q1.1 等。
// 写入Q寄存器：Q + 编号，例如 Q0.0、Q1.1 等。
func _ParseADDR_Q(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	if len(s) < 2 {
		return AddressInfo, fmt.Errorf("invalid Address format:%s", s)
	}
	// I0.0
	parts := strings.Split(s[1:], ".")
	// 0.0
	if len(parts) != 2 {
		return AddressInfo, fmt.Errorf("Invalid ElementNumber Address format:%s", s)
	}
	//
	DataBlockNumber, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		return AddressInfo, fmt.Errorf("Address %s Atoi failed:%s", s, err1)
	}
	BitNumber, err2 := strconv.Atoi(parts[0])
	if err2 != nil {
		return AddressInfo, fmt.Errorf("Address %s Atoi failed:%s", s, err2)
	}
	AddressInfo.AddressType = "Q"
	AddressInfo.DataBlockNumber = DataBlockNumber
	AddressInfo.DataBlockType = "BYTE"
	AddressInfo.BitNumber = BitNumber

	return AddressInfo, nil
}
func _ParseADDR_QB(s string) (AddressInfo, error) {
	AddressInfo := AddressInfo{}
	AddressInfo.AddressType = "QB"
	DataBlockNumber := 0
	AddressInfo.DataBlockNumber = DataBlockNumber
	BitNumber := 0
	AddressInfo.BitNumber = BitNumber
	return AddressInfo, nil
}
