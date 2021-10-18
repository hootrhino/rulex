package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/lunixbochs/struc"
)

func Test_Run(t *testing.T) {

	type Example struct {
		A int `struc:"big"`
		// B will be encoded/decoded as a 16-bit int (a "short")
		// but is stored as a native int in the struct
		B int `struc:"int16"`

		// the sizeof key links a buffer's size to any int field
		Size int `struc:"int8,little,sizeof=Str"`
		Str  string
		// you can get freaky if you want
		Str2 string `struc:"[5]int64"`
	}
	var buf bytes.Buffer
	err := struc.Pack(&buf, &Example{1, 2, 0, "test", "test2"})
	t.Log(err)
	o := &Example{}
	err = struc.Unpack(&buf, o)
	t.Log(err)
}
func Test_Run_go_binary(t *testing.T) {
	inputs := [][]byte{
		{0x01},
		{0x02},
		{0x7f},
		{0x80, 0x01},
		{0xff, 0x01},
		{0x80, 0x02},
	}
	for _, b := range inputs {
		x, n := binary.Uvarint(b)
		if n != len(b) {
			fmt.Println("Uvarint did not consume all of in")
		}
		fmt.Println(x)
	}
}
