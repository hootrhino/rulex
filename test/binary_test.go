package test

import (
	"github.com/hootrhino/rulex/rulexlib"

	"testing"
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
