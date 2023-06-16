package test

import (
	"testing"

	"github.com/hootrhino/rulex/rulexlib"
)

// 性能测试
func BenchmarkRegex(b *testing.B) {
	segments := rulexlib.MatchHexLib("A0:[0,1];A1:[1,2]",
		"0117011d0127011a0110010e")
	b.Logf("%+v", segments)
}

func TestRegex(t *testing.T) {
	// 0117 011d 0127 011a 0110 010e
	segments := rulexlib.MatchHexLib("A0:[0,1];A1:[1,2]",
		"0117011d0127011a0110010e")
	for _, sg := range segments {
		t.Logf("%+v", sg.ToUInt64())
	}
	t.Logf("%+v", segments)
}
