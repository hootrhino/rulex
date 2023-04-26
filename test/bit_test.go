package test

import (
	"testing"

	"github.com/hootrhino/rulex/rulexlib"

	"github.com/go-playground/assert/v2"
)

func TestReverseBitOrder(t *testing.T) {
	assert.Equal(t, byte(0b1011_1111), rulexlib.ReverseBits(byte(0b1111_1101)))
	assert.Equal(t, byte(0b1100_0000), rulexlib.ReverseBits(byte(0b0000_0011)))
	assert.Equal(t, byte(0b0000_0101), rulexlib.ReverseBits(byte(0b1010_0000)))
	assert.Equal(t, byte(0b1010_1010), rulexlib.ReverseBits(byte(0b0101_0101)))
	assert.Equal(t, []byte{3, 2, 1}, rulexlib.ReverseByteOrder([]byte{1, 2, 3}))
}
