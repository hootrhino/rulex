package test

import (
	"rulex/utils"
	"testing"
)

func Test_cgo_getpid(t *testing.T) {
	t.Log(utils.GetPid())
}
