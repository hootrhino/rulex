package test

import (
	"rulex/stdlib/parser"
	"testing"
)

func TestParser(t *testing.T) {
	parser.Parse([]byte("[A : 10]"))
}
