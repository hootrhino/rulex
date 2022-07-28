package common

import "errors"

//--------------------------------------------------------------------------------------------------
// 内部函数
//--------------------------------------------------------------------------------------------------

/*
*
* 取某个字节上的位
*
 */
func GetABitOnByte(b byte, position uint8) (v uint8) {
	if position > 8 {
		return 0
	}
	var mask byte = 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position
	}
	return (b & (mask << int(position))) >> position

}

/*
*
* 设置字节上的某个位
*
 */
func SetABitOnByte(b *byte, position uint8, value bool) (byte, error) {
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
* 字符串转字节
*
 */
func BitStringToBytes(s string) ([]byte, error) {
	if len(s)%8 != 0 {
		return nil, errors.New("length must be integer multiple of 8")
	}
	b := make([]byte, (len(s)+(8-1))/8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '1' {
			return nil, errors.New("value out of range")
		}
		b[i>>3] |= (c - '0') << uint(7-i&7)
	}
	return b, nil
}

/*
*
* 字节上某个位转逻辑值
*
 */
func BitToBool(data byte, index uint8) bool {
	return GetABitOnByte(data, index) == 1
}
func BitToUint8(data byte, index uint8) uint8 {
	if GetABitOnByte(data, index) == 1 {
		return 1
	}
	if GetABitOnByte(data, index) == 0 {
		return 0
	}
	return 0
}

/*
*
* 字节转布尔值 本质上是判断是否等于0 or 1
*
 */
func ByteToBool(data byte) bool {
	return data == 1
}
