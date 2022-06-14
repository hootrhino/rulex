package driver

import "errors"

//--------------------------------------------------------------------------------------------------
// 内部函数
//--------------------------------------------------------------------------------------------------

/*
*
* 取某个字节上的位
*
 */
func getABitOnByte(b byte, position uint8) (v uint8) {
	mask := 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position
	}
	return (b & (1 << mask)) >> position

}

/*
*
* 设置字节上的某个位
*
 */
func setABitOnByte(b *byte, position uint8, value bool) (byte, error) {
	if position > 7 {
		return 0, errors.New("下标必须是0-7, 高位在前, 低位在后")
	}
	if value {
		return *b & 0b1111_1111, nil
	}
	masks := []byte{
		0b11111110,
		0b11111101,
		0b11111011,
		0b11110111,
		0b11101111,
		0b11011111,
		0b10111111,
		0b01111111,
	}
	return *b & masks[position], nil

}

/*
*
* 字节转逻辑
*
 */
func byteToBool1(data byte, index uint8) bool {
	return getABitOnByte(data, index) == 1
}
func byteToBool2(data byte) bool {
	return data == 1
}
