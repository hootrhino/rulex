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
* 字节上某个位转逻辑值
*
 */
func bitToBool(data byte, index uint8) bool {
	return getABitOnByte(data, index) == 1
}

/*
*
* 字节转布尔值 本质上是判断是否等于0 or 1
*
 */
func byteToBool(data byte) bool {
	return data == 1
}

/*
*
* 西门子PLC的一些配置
*
 */

type S1200Block struct {
	Tag     string `json:"tag"`     // 数据tag
	Address int    `json:"address"` // 地址
	Start   int    `json:"start"`   // 起始地址
	Size    int    `json:"size"`    // 数据长度
}
type S1200BlockValue struct {
	Tag     string `json:"tag"`     // 数据tag
	Address int    `json:"address"` // 地址
	Start   int    `json:"start"`   // 起始地址
	Size    int    `json:"size"`    // 数据长度
	Value   []byte `json:"value"`
}
