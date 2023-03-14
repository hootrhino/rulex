package test

import (
	"testing"

	"github.com/i4de/rulex/rulexlib"
)

// 性能测试
func BenchmarkRegex(b *testing.B) {
	segments := rulexlib.MatchHexLib("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")
	b.Logf("%+v", segments)
}

func TestRegex(t *testing.T) {
	segments := rulexlib.MatchHexLib("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")

	t.Logf("%+v", segments)
}
