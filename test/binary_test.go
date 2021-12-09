package test

import (
	"encoding/binary"
	"fmt"
	"rulex/rulexlib"
	"strconv"
	"testing"
)

func TestBinary(t *testing.T) {
	fmt.Printf("%08b\n", (0 << 0b00000001))
	fmt.Printf("%08b\n", (1 << 0b00000001))
	fmt.Printf("%08b\n", (2 << 0b00000001))
	v := byte(0b00010001)
	//
	fmt.Println(rulexlib.GetABitOnByte(v, 0))
	fmt.Println(rulexlib.GetABitOnByte(v, 1))
	fmt.Println(rulexlib.GetABitOnByte(v, 2))
	fmt.Println(rulexlib.GetABitOnByte(v, 3))
	fmt.Println(rulexlib.GetABitOnByte(v, 4))
	fmt.Println(rulexlib.GetABitOnByte(v, 5))
	fmt.Println(rulexlib.GetABitOnByte(v, 6))
	fmt.Println(rulexlib.GetABitOnByte(v, 7))
}

func TestBinaryMatch(t *testing.T) {
	// 00001111 10101010 00001010
	t.Log("rulexlib.Match ==> <a:1 b:3 b:5", "000011111010101000001010")
	kls := rulexlib.Match("<a:1 b:3 b:5", []byte{0b00001111, 0b10101010, 0b00001010}, true)
	var len uint = 0
	for _, v := range kls {
		len += v.L
	}
	t.Log("Len:", len)
	t.Log(kls)
}

/*
*
* 大端模式
*
 */
func TestBinaryMatch_big(t *testing.T) {
	//aab: 01100001 01100001 01100010
	kls := rulexlib.Match(">a:8 b:8 c:8", []byte("aab"), true)
	var len uint = 0
	for _, v := range kls {
		len += v.L
		t.Log("字段:", v.K, " 二进制串:", v.BS)
	}
	t.Log("Len:", len)
	t.Log(kls)
}

/*
*
* 小端模式
*
 */
func TestBinaryMatch_little(t *testing.T) {
	//baa:  01100010 01100001 01100001
	kls := rulexlib.Match("<k1:8 k2:8 k3:8", []byte("aab"), true)
	var len uint = 0
	for _, v := range kls {
		len += v.L
		t.Log("字段:", v.K, " 二进制串:", v.BS)
	}
	t.Log("Len:", len)
	t.Log(kls)
}
func TestByteToBitFormatString(t *testing.T) {
	//
	// 假设Modbus的线圈有8个，状态如下：
	// 0 1 1 0 0 0 0 1
	// 原始数据是1个字节
	originData := []byte{0b_0110_0001, 0b_0110_0001}
	t.Log("originData:", rulexlib.ByteToBitString(originData))
	// 到了网关后被转成字符串
	formatData := string(originData)
	t.Log("formatData:", formatData)
	// 二进制语法匹配的时候，再次把字符串转成字节
	t.Logf("formatData -> originData:%08b\n", []byte(formatData))

}
func TestBinaryParseInt(t *testing.T) {
	if i, err := strconv.ParseInt("00001111", 2, 64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i)
	}
	b := []byte{0b00011000, 0b00011000} //6168

	fmt.Println(rulexlib.ByteToInt(b, binary.LittleEndian))
}

func TestBitStringToBytes(t *testing.T) {
	t.Log(rulexlib.BitStringToBytes("000000010000000110000000"))
}
