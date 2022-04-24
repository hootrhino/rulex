package test

import (
	"rulex/rulexlib"

	"testing"

	"github.com/go-playground/assert/v2"
)

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

func TestReverseBitOrder(t *testing.T) {
	assert.Equal(t, byte(0b1011_1111), rulexlib.ReverseBits(byte(0b1111_1101)))
	assert.Equal(t, byte(0b1100_0000), rulexlib.ReverseBits(byte(0b0000_0011)))
	assert.Equal(t, byte(0b0000_0101), rulexlib.ReverseBits(byte(0b1010_0000)))
	assert.Equal(t, byte(0b1010_1010), rulexlib.ReverseBits(byte(0b0101_0101)))
	assert.Equal(t, []byte{3, 2, 1}, rulexlib.ReverseByteOrder([]byte{1, 2, 3}))
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
