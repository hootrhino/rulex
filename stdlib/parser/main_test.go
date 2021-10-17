package parser

import (
	"fmt"
	"testing"
)

type BinaryLexser struct {
	input  []byte
	index  int
	buffer string
}

func (l *BinaryLexser) Lex(lval *yySymType) int {

	if len(l.input) == 0 {
		return 0
	}
	fmt.Println("input:", string(l.input))
	for _, v := range l.input {
		fmt.Println("input:", string(v))

		if v == '[' {
			return '['
		}
		if v == '[' {
			return '['
		}
		if v == ',' {
			return ','
		}
		if v == ':' {
			return ':'
		}
	}
	return 0
}
func (l *BinaryLexser) Error(s string) {
	fmt.Println("Error:", s)
}
func TestParser(t *testing.T) {

	yyNewParser().Parse(&BinaryLexser{
		input:  []byte("[A:10]"),
		
		index:  0,
		buffer: "",
	})
	yyNewParser().Parse(&BinaryLexser{
		input:  []byte("[A:10, B:10, C:20]"),
		index:  0,
		buffer: "",
	})

}
