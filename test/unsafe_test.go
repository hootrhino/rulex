package test

import (
	"fmt"
	"testing"
)

func TestUnsafe(t *testing.T) {

	ss := ""
	b := []byte{0b11110000, 0b00001111}
	for _, v := range b {
		ss += fmt.Sprintf("%b", v)
	}
	t.Log(ss)
}
