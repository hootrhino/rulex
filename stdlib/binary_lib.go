package stdlib

import (
	"errors"
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

type BinaryLib struct {
}

func (l *BinaryLib) LoadLib(name string, e typex.RuleX, L *lua.LState) error {
	return nil
}
func (l *BinaryLib) UnLoadLib(name string) error {
	return nil
}

//------------------------------------------------------------------------------------
// 自定义实现函数
//------------------------------------------------------------------------------------

//
// 从一个字节里面提取某 1 个位的值，只有 0 1 两个值
// 注意：这里取得是大端模式，也就是最高位在最前面，最低位在最后面
//

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
	return 0, errors.New("position must between (0-8)")
}

//
// TODO: 下一个大版本支持，至少3个月后
//
// 这里借鉴了下Erlang的二进制语法: <<A:5,B:4>> = <<"helloworld">>
// 其中A = hello B= world
//
func GetBinary(expr string, data []byte) map[string]interface{} {
	return nil
}
