package parser

import (
	"fmt"
)

type BinaryParser struct {
	input  []byte
	index  int
	buffer string
	lValue yySymType
}

func (l *BinaryParser) Lex(lval *yySymType) int {

	if len(l.input) == 0 {
		return 0
	}
	fmt.Println("input expr =>", string(l.input))

	for _, v := range l.input {
		fmt.Println("Current is:", string(v))
		if v == '<' {
			return lb
		}
		if v == '>' {
			return rb
		}
		if v == ':' {
			return as
		}
		// if v == ':' {
		fmt.Println("KV ===>", (l.lValue.pair), ":", (l.lValue.l))
		// }
	}
	return 0
}
func (l *BinaryParser) Error(s string) {
	fmt.Println("Error:", s, "error @:", l.index, string(l.input[l.index]))
}
func Parse(b []byte) {
	yyNewParser().Parse(&BinaryParser{
		input:  b,
		index:  0,
		buffer: "",
	})
}
