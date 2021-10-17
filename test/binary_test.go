package test

import (
	"errors"
	"fmt"
	"testing"
)

func GetABitOnByte(b byte, position uint8) (v uint8, errs error) {
	//  --------------->
	//  7 6 5 4 3 2 1 0
	// |.|.|.|.|.|.|.|.|
	//
	if position == 0 {
		return (b & 0b00000001) >> position, nil
	}
	if position == 1 {
		return (b & 0b00000010) >> position, nil
	}
	if position == 2 {
		return (b & 0b00000100) >> position, nil
	}
	if position == 3 {
		return (b & 0b00001000) >> position, nil
	}
	if position == 4 {
		return (b & 0b00010000) >> position, nil
	}
	if position == 5 {
		return (b & 0b00100000) >> position, nil
	}
	if position == 6 {
		return (b & 0b01000000) >> position, nil
	}
	if position == 7 {
		return (b & 0b10000000) >> position, nil
	}
	return 0, errors.New("Position must between (0-8)")
}
func TestBinary(t *testing.T) {
	fmt.Printf("%08b\n", (0 << 0b00000001))
	fmt.Printf("%08b\n", (1 << 0b00000001))
	fmt.Printf("%08b\n", (2 << 0b00000001))
	v := byte(0b00010001)
	//
	fmt.Println(GetABitOnByte(v, 0))
	fmt.Println(GetABitOnByte(v, 1))
	fmt.Println(GetABitOnByte(v, 2))
	fmt.Println(GetABitOnByte(v, 3))
	fmt.Println(GetABitOnByte(v, 4))
	fmt.Println(GetABitOnByte(v, 5))
	fmt.Println(GetABitOnByte(v, 6))
	fmt.Println(GetABitOnByte(v, 7))
}
