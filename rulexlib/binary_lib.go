package rulexlib

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

// 提取的 Key, 最长不能超过32个字母
var pattern = `[a-zA-Z0-9]{1,32}:[1-9]+`
var regExpr *regexp.Regexp = regexp.MustCompile(pattern)

/*
*
* 二进制匹匹配, [<|> K1:LEN1 K2:LEN2... ]返回一个K-V table
* 其中K是字符串, V是二进制字符串
*
 */
func MatchBinary(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		expr := state.ToString(2)
		data := state.ToString(3)
		returnMore := state.ToBool(4)
		t := lua.LTable{}
		for _, kl := range Match(expr, []byte(data), returnMore) {
			t.RawSetString(kl.K, lua.LString(kl.BS))
		}
		state.Push(&t)
		return 1
	}
}

/*
*
* 二进制串转成十六进制串
*
 */
func MatchBinaryHex(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		expr := state.ToString(2)
		data := state.ToString(3)
		returnMore := state.ToBool(4)
		t := lua.LTable{}
		for _, kl := range Match(expr, []byte(data), returnMore) {
			t.RawSetString(kl.K, lua.LString(kl.ToHexString()))
		}
		state.Push(&t)
		return 1
	}
}

func GetABitOnByte(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		if state.Get(2).Type() != lua.LTNumber {
			state.Push(nil)
			return 1
		}
		b := uint8(state.ToInt(2))
		pos := uint8(state.ToInt(3))
		if v, err := getABitOnByte(b, pos); err != nil {
			state.Push(nil)
		} else {
			state.Push(lua.LNumber(v))
		}
		return 1
	}
}

func ByteToBitString(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		data := state.ToString(2)
		state.Push(lua.LValue(lua.LString(byteToBitString([]byte(data)))))
		return 1
	}
}

//------------------------------------------------------------------------------------
// 自定义实现函数
//------------------------------------------------------------------------------------

//
// 从一个字节里面提取某 1 个位的值, 只有 0 1 两个值
//

func getABitOnByte(b byte, position uint8) (uint8, error) {
	if position > 8 {
		return 0, errors.New("position must greater than 8")
	}
	var mask byte = 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position, nil
	}
	return (b & (mask << int(position))) >> position, nil

}

//
// 字节转位串
//

func byteToBitString(b []byte) string {
	s := ""
	for _, v := range b {
		s += fmt.Sprintf("%08b", v)
	}
	return s
}

type Kl struct {
	K  string //Key
	L  uint   //Length
	BS string //BitString
}

func (k Kl) String() string {
	return fmt.Sprintf("KL@ K: %v,L: %v,BS: %v", k.K, k.L, k.BS)
}
func (k Kl) ToHexString() string {
	b, _ := bitStringToBytes(k.BS)
	return fmt.Sprintf("%X", b)
}
func (k Kl) ToByte() []byte {
	b, _ := bitStringToBytes(k.BS)
	return b
}

// Example data: 12345678
// Big-Endian [>]=12345678:  高位字节放在内存的低地址端, 低位字节放在内存的高地址端。
// Little-Endian [<]=87654321: 低位字节放在内存的低地址段, 高位字节放在内存的高地址端
//

func ByteToInt64(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		endian := state.ToString(2)
		data := state.ToString(3)
		if endian == ">" {
			v := ByteToInt([]byte(data), binary.BigEndian)
			state.Push(lua.LNumber(v))
			return 1

		}
		if endian == "<" {
			v := ByteToInt([]byte(data), binary.LittleEndian)
			state.Push(lua.LNumber(v))
			return 1

		}
		state.Push(nil)
		return 1
	}
}

// /
func ByteToInt(b []byte, order binary.ByteOrder) uint64 {
	var err error
	// uint8
	if len(b) == 1 {
		var x uint8
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			glogger.GLogger.Error(b, err)
		}
		return uint64(x)
	}
	// uint16
	if len(b) > 1 && len(b) <= 2 {
		var x uint16
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			glogger.GLogger.Error(b, err)
		}
		return uint64(x)
	}
	// uint32
	if len(b) > 2 && len(b) <= 4 {
		var x uint32
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			glogger.GLogger.Error(b, err)
		}
		return uint64(x)
	}
	// uint64
	if len(b) > 4 {
		var x uint64
		err = binary.Read(bytes.NewBuffer(b), order, &x)
		if err != nil {
			glogger.GLogger.Error(b, err)
		}
		return x
	}

	return 0
}

/*
*
* 位串转字节
*
 */

func BitStringToBytes(rx typex.RuleX) func(*lua.LState) int {
	return func(state *lua.LState) int {
		data := strings.Replace(state.ToString(2), "\"", "", -1)
		// data 可能有 '"' 字符
		b, err := bitStringToBytes(data)
		if err != nil {
			state.Push(nil)
		} else {
			state.Push(lua.LString(string(b)))
		}
		return 1
	}
}

func bitStringToBytes(s string) ([]byte, error) {
	b := make([]byte, (len(s)+(8-1))/8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '1' {
			return nil, errors.New("bitstring to bytes error, bit value out of range, only can be 0 or 1")
		}
		b[i>>3] |= (c - '0') << uint(7-i&7)
	}
	return b, nil
}

/*
*
* 位串转字节
*
 */
func AsByteSlice(b string) []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

/*
*
* 位串转十六进制
*
 */
func AsHexSlice(b string) []string {
	var out []string
	byteSlice := AsByteSlice(b)
	for _, b := range byteSlice {
		out = append(out, "0x"+hex.EncodeToString([]byte{b}))
	}
	return out
}

/*
*
* 大小端判断
*
 */
func Endian(endian byte) binary.ByteOrder {
	// < 小端 0x34 0x12
	if endian == '>' {
		return binary.BigEndian
	}
	// > 大端 0x12 0x34
	if endian == '<' {
		return binary.LittleEndian
	}
	return binary.LittleEndian
}

// 字节空位补位0
func append0Prefix(n int) string {
	if (n % 8) == 7 {
		return "0"
	}
	if (n % 8) == 6 {
		return "00"
	}
	if (n % 8) == 5 {
		return "000"
	}
	if (n % 8) == 4 {
		return "0000"
	}
	if (n % 8) == 3 {
		return "00000"
	}
	if (n % 8) == 2 {
		return "000000"
	}
	if (n % 8) == 1 {
		return "0000000"
	}
	return ""
}

//--------------------------------------------------------------
// rulexlib:Match
//--------------------------------------------------------------

func Match(expr string, data []byte, returnMore bool) []Kl {
	cursor := 0
	result := []Kl{}
	// glogger.GLogger.Debug(pattern, expr[1:])

	matched, err := regexp.MatchString(pattern, expr[1:])
	if matched {
		endian := expr[0]
		// < 小端
		if Endian(endian) == binary.LittleEndian {
			bfs := byteToBitString(data)
			// <a:12 b:12
			reverseBfs := ReverseString(bfs)
			buildResult(returnMore, cursor, reverseBfs, expr[1:], &result)
		}
		// > 大端
		if Endian(endian) == binary.BigEndian {
			bfs := byteToBitString(data)
			// <a:12 b:12
			buildResult(returnMore, cursor, bfs, expr[1:], &result)
		}

	} else {
		glogger.GLogger.Error(matched, err)
	}
	return result
}

func append0(cursor int, bfs string, returnMore bool, result *[]Kl) {
	if returnMore {
		if cursor < len(bfs) {
			// 是否最后补上剩下的字节
			*result = append(*result, Kl{"_", uint(len(bfs) - cursor), append0Prefix(len(bfs[cursor:])) + bfs[cursor:]})
		}
	}
}
func buildResult(returnMore bool, cursor int, bfs string, expression string, result *[]Kl) {
	for _, v := range regExpr.FindAllString(expression, -1) {
		k_l := strings.Split(v, ":")
		k := k_l[0]
		if l, err1 := strconv.Atoi(k_l[1]); err1 == nil {
			if cursor+l <= len(bfs) {
				binString := bfs[cursor : cursor+l]
				*result = append(*result, Kl{
					K:  k,
					L:  uint(l),
					BS: append0Prefix(len(binString)) + binString,
				})
			} else {
				*result = append(*result, Kl{k, uint(l), ""})
			}
			cursor += l
		}
	}
	append0(cursor, bfs, returnMore, result)
}

/*
*
* 逆转位顺序
*
 */
func ReverseBits(b byte) byte {
	result := byte(0)
	for i := 7; i >= 0; i-- {
		result |= ((b << i) >> 7) << i
	}
	return result
}

/*
*
* 逆转字节顺序
*
 */
func ReverseByteOrder(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

/*
*
* 逆转字符串
*
 */
// Reverse reverses the input while respecting UTF8 encoding and combined characters
func ReverseString(text string) string {
	textRunes := []rune(text)
	textRunesLength := len(textRunes)
	if textRunesLength <= 1 {
		return text
	}

	i, j := 0, 0
	for i < textRunesLength && j < textRunesLength {
		j = i + 1
		for j < textRunesLength && isMark(textRunes[j]) {
			j++
		}

		if isMark(textRunes[j-1]) {
			// Reverses Combined Characters
			reverse(textRunes[i:j], j-i)
		}

		i = j
	}

	// Reverses the entire array
	reverse(textRunes, textRunesLength)

	return string(textRunes)
}

/*
*
* 字节逆序
*
 */
func reverse(runes []rune, length int) {
	for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
}

// isMark determines whether the rune is a marker
func isMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

/*
*
* Hex to number
*
 */
func HToN(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		s := l.ToString(2)
		if iv, err := HexToNumber(s); err != nil {
			l.Push(lua.LNil)
		} else {
			l.Push(lua.LNumber(iv))
		}
		return 1
	}
}

/*
*
* 取某个Hex字符串的子串转换成数字
*
 */
func HsubToN(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		s := l.ToString(2)
		start := l.ToInt(3)
		offset := l.ToInt(4)
		if iv, err := HexToNumber(s[start:offset]); err != nil {
			l.Push(lua.LNil)
		} else {
			l.Push(lua.LNumber(iv))
		}
		return 1
	}
}

/*
*
* 十六进制字符串转数字
*
 */
func HexToNumber(s string) (int64, error) {
	iv, err := strconv.ParseInt(s, 16, len(s)*8)
	if err != nil {
		return 0, err
	}
	return iv, nil
}

/*
*
* 二进制转浮点数，参考资料：
* https://blog.51cto.com/u_12512821/2363818
* https://www.ruanyifeng.com/blog/2010/06/ieee_floating-point_representation.html
*
 */

func BinToFloat32(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		bin := l.ToString(2)
		bits := binary.BigEndian.Uint32([]byte(bin))
		l.Push(lua.LNumber(math.Float32frombits(bits)))
		return 1
	}
}
func BinToFloat64(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		bin := l.ToString(2)
		bits := binary.BigEndian.Uint64([]byte(bin))
		l.Push(lua.LNumber(math.Float64frombits(bits)))
		return 1
	}
}

/*
*
* Base64 to byte: 用来处理golang的JSON转换byte[]问题
*
 */
func B64S2B(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		b64s := l.ToString(2)
		bss, err := base64.StdEncoding.DecodeString(b64s)
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LString(bss))
			l.Push(lua.LNil)
		}
		return 2
	}
}
