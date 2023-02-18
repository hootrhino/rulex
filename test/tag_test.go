package test

import (
	"encoding/binary"
	"testing"

	"github.com/go-playground/assert/v2"


)

func Test_binary_to_int(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x04,
		0x04, 0x03, 0x02, 0x01,
	}
	Address := binary.BigEndian.Uint32(data[0:4])
	Start := binary.BigEndian.Uint32(data[4:8])
	Size := binary.BigEndian.Uint32(data[8:12])
	assert.Equal(t, uint32(1), (Address))
	assert.Equal(t, uint32(1), (Start))
	assert.Equal(t, uint32(4), (Size))
}
