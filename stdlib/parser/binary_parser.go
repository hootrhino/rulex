package parser

import (
	"fmt"
)

type BinaryParser struct {
	input  []byte
	index  int
	buffer string
	lval   yySymType
	stack  [yyInitialStackSize]yySymType
	char   int
}

func (l *BinaryParser) Lex(lval *yySymType) int {

	if len(l.input) == 0 {
		return 0
	}
	fmt.Println("input expr =>", string(l.input))
	for _, v := range l.input {

		if v == '[' {
			l.index++

		}
		if v == ']' {
			l.index++

		}
		if v == ',' {

		}
		if v == ':' {
			fmt.Println("KV===>", (l.lval.field), ":", (l.lval.length))

		}
		if string(v) != "" {

		}
		if int(v) > 0 {

		}

	}
	return 0
}
func (l *BinaryParser) Error(s string) {
	fmt.Println("Error:", s)
}
func Parse(b []byte) {
	yyNewParser().Parse(&BinaryParser{
		input:  b,
		index:  0,
		buffer: "",
	})
}
