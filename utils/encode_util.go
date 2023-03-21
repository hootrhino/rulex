package utils

import (
	"unicode"
)

//
// 去掉\u0000字符
//
func TrimZero(s string) string {
	str := make([]rune, 0, len(s))
	for _, v := range s {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) {
			continue
		}
		str = append(str, v)
	}
	return string(str)
}
