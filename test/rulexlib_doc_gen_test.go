package test

import (
	"testing"

	"github.com/i4de/rulex/engine"
)

func Test_Gen_rulexlib_doc(t *testing.T) {
	engine.BuildInLuaLibDoc()
}
