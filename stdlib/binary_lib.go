package stdlib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"rulex/typex"
	"strconv"
	"strings"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

var regexper *regexp.Regexp
var pattern = `[a-z]+:[1-9]+`

func init() {
	regexper = regexp.MustCompile(pattern)

}

type BinaryLib struct {
	bBuffer  *bytes.Buffer
	regexper *regexp.Regexp
}

func NewBinaryLib() typex.XLib {

	return nil
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
	mask := 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position, nil
	} else {
		return (b & (1 << mask)) >> position, nil
	}
}

//
// TODO: 下一个大版本支持，至少3个月后
//
// 这里借鉴了下Erlang的二进制语法: <<A:5,B:4>> = <<"helloworld">>
// 其中A = hello B= world
//

func ByteToBitFormatString(b []byte) string {
	s := ""
	for _, v := range b {
		s += fmt.Sprintf("%08b", v)
	}
	return s
}

type Kl struct {
	K  string
	L  uint
	BS interface{}
	I  uint64
}

func (k Kl) String() string {
	return fmt.Sprintf("K: %v, L: %v, BS: %v, I: %v", k.K, k.L, k.BS, k.I)
}

func ByteToInt(b []byte, order binary.ByteOrder) uint64 {
	var err error
	// uint8
	if len(b) == 1 {
		var x uint8
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			log.Error(b, err)
		}
		return uint64(x)
	}
	// uint16
	if len(b) > 1 && len(b) <= 2 {
		var x uint16
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			log.Error(b, err)
		}
		return uint64(x)
	}
	// uint32
	if len(b) > 2 && len(b) <= 4 {
		var x uint32
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			log.Error(b, err)
		}
		return uint64(x)
	}
	// uint64
	if len(b) > 4 && len(b) <= 8 {
		var x uint64
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			log.Error(b, err)
		}
		return x
	}

	return 0
}
func BitStringToBytes(s string) ([]byte, error) {
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
func endian(endian uint8) binary.ByteOrder {
	if endian == '>' {
		return binary.BigEndian
	}
	// 小端 0x34 0x12
	if endian == '<' {
		return binary.LittleEndian
	}
	return binary.LittleEndian
}
func Match(s string, data []byte, returnMore bool) []Kl {
	cursor := 0
	result := []Kl{}
	matched, err0 := regexp.MatchString(pattern, s[1:])
	if matched {
		bfs := ByteToBitFormatString(data)
		// <a:12 b:12
		for _, v := range regexper.FindAllString(s[1:], -1) {
			kl := strings.Split(v, ":")
			k := kl[0]
			if l, err1 := strconv.Atoi(kl[1]); err1 == nil {
				if cursor+l < len(s) {
					binString := bfs[cursor : cursor+l]
					if b, err := BitStringToBytes(binString); err == nil {
						intValue := ByteToInt(b, endian(s[0]))
						result = append(result, Kl{k, uint(l), binString, intValue})
					} else {
						log.Error(err)
						result = append(result, Kl{k, uint(l), binString, 0})
					}

				} else {
					result = append(result, Kl{k, uint(l), nil, 0})
				}
				cursor += l
			}
		}
		if returnMore {
			if cursor < len(s) {
				result = append(result, Kl{"_$", uint(len(bfs) - cursor), bfs[cursor:], 0})
			}
		}

	} else {
		log.Error(err0)
	}

	return result
}
