package rulexlib

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 十六进制字符串转byte数组
*
 */
func Hexs2Bytes(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		hexs := l.ToString(2)
		s, e := hex.DecodeString(hexs)
		if e != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(e.Error()))
		} else {
			table := lua.LTable{}
			for _, v := range s {
				table.Append(lua.LNumber(v))
			}
			l.Push(&table)
			l.Push(lua.LNil)
		}
		return 2
	}
}

/*
*
* byte数组转十六进制字符串
*
 */
func Bytes2Hexs(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		bytes := l.ToString(2)
		l.Push(lua.LString(hex.EncodeToString([]byte(bytes))))
		l.Push(lua.LNil)
		return 2
	}
}

/*
*---------------------------------------------------------------------------
* 十六进制字符串匹配: MatchHex("FFFFFF014CB2AA55", "age:[1,232];sex:[4,5]")
*---------------------------------------------------------------------------
 */
func MatchHex(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		exprS := l.ToString(2)
		hexS := l.ToString(3)
		mhs := MatchHexLib(exprS, hexS)
		ntb := lua.LTable{}
		for _, v := range mhs {
			ntb.RawSetString(v.Name, lua.LString(v.ToHexString()))
		}
		l.Push(&ntb)
		return 1
	}
}

/*
*
* 匹配十六进制转成整数
*  MatchHex("FFFFFF014CB2AA55", "age:[1,1];sex:[4,5]")
*
 */
func MatchUInt(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		exprS := l.ToString(2)
		hexS := l.ToString(3)
		mhs := MatchHexLib(exprS, hexS)
		ntb := lua.LTable{}
		for _, v := range mhs {
			size := len(v.Value)

			if size == 2 {
				ntb.RawSetString(v.Name, lua.LNumber(v.ToUint16()))
			}
			if size == 4 {
				ntb.RawSetString(v.Name, lua.LNumber(v.ToUint32()))
			}
			if size == 8 {
				ntb.RawSetString(v.Name, lua.LNumber(v.ToUInt64()))
			}
			if size > 8 {
				ntb.RawSetString(v.Name, lua.LNumber(-0xFFFFFFFF))
			}
		}
		l.Push(&ntb)
		return 1
	}
}

type HexSegment struct {
	Name  string
	Value []byte
}

func (sgm HexSegment) ToHexString() string {
	return fmt.Sprintf("%X", sgm.Value)
}

func (sgm HexSegment) ToUint16() uint16 {
	value := binary.BigEndian.Uint16(sgm.Value)
	return value
}
func (sgm HexSegment) ToUint32() uint32 {
	value := binary.BigEndian.Uint32(sgm.Value)
	return value
}
func (sgm HexSegment) ToUInt64() uint64 {
	value := binary.BigEndian.Uint64(sgm.Value)
	return value
}

// 全局正则表达式编译器, 这是已经验证过的正则表达式，所以一定编译成功，故不检查error
var regexMatcher, _ = regexp.Compile(`[a-zA-Z0-9]+:\[[0-9]+,[0-9]+\]`)

/*
*
* 匹配十六进制字符
*
 */
func MatchHexLib(regExpr, hexStr string) []HexSegment {
	match := regexMatcher.FindAllString(regExpr, -1)
	if len(match) == 0 {
		return nil
	}
	var segments []HexSegment
	for _, item := range match {
		splits := strings.Split(item, ":")
		if len(splits) != 2 {
			return nil
		}

		name := splits[0]
		start, end := extIndex(splits[1])
		subHex := extHex(hexStr, start, end)
		value, _ := hex.DecodeString(subHex)

		segments = append(segments, HexSegment{name, value})
	}
	return segments
}

func extIndex(str string) (start, end int) {
	indexStr := strings.TrimSuffix(strings.TrimPrefix(str, "["), "]")
	split := strings.Split(indexStr, ",")
	if len(split) != 2 {
		return -1, -1
	}
	start, err := strconv.Atoi(split[0])
	if err != nil {
		return -1, -1
	}
	end, err2 := strconv.Atoi(split[1])
	if err2 != nil {
		return -1, -1
	}
	return start, end
}

func extHex(hexStr string, start, end int) string {
	if start < 0 || end < 0 || start > end || end*2 > len(hexStr) {
		return ""
	}
	return hexStr[start*2 : (end+1)*2]
}
