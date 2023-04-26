package test

import (
	"fmt"
	"testing"
	"unsafe"
)

type ssd struct {
	_   int //8
	len int //8
}

func strlen(s string) int {
	return ((*ssd)(unsafe.Pointer(&s))).len
}

// go test -timeout 30s -run ^Test_string_uptr github.com/hootrhino/rulex/test -v -count=1 -run=none
func Test_string_uptr(t *testing.T) {
	s := "HELLO_WORLD"
	sp := unsafe.Pointer(&s)
	t.Log((sp))
	t.Log(len(s))

}
func Benchmark_string_len1(b *testing.B) {
	s := ""
	for i := 0; i < 10*1024; i++ {
		s += fmt.Sprintf("%v", i)
	}
	l := 0
	b.ResetTimer()
	for i := 0; i < 10*1024; i++ {
		l += 1
	}
	b.StopTimer()
	b.Log("len1", l, b.N)

}

// go.exe test -benchmem -run=^$ -bench ^Benchmark_string_len$ github.com/hootrhino/rulex/test -v -cpu=6
func Benchmark_string_len2(b *testing.B) {
	s := ""
	for i := 0; i < 10*1024; i++ {
		s += fmt.Sprintf("%v", i)
	}
	l := 0
	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		l = len(s)
	}
	b.StopTimer()
	b.Log("len2", l, b.N)

}

// go.exe test -benchmem -run=^$ -bench ^Benchmark_string_lenssd$ github.com/hootrhino/rulex/test -v -cpu=6
func Benchmark_string_len3(b *testing.B) {
	s := ""
	for i := 0; i < 10*1024; i++ {
		s += fmt.Sprintf("%v", i)
	}
	l := 0
	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		l = strlen(s)
	}
	b.StopTimer()
	b.Log("len3", l, b.N)
}
