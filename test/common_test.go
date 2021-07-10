package test

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	from := "12345,"
	for _, v := range strings.Split(from, ",") {
		t.Log(v == "")
	}
}
