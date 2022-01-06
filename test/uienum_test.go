package test

import (
	"rulex/core"
	"testing"
)

func Test_split_string(t *testing.T) {
	str1 := "k1,v1|k2,v2"
	data, err := core.RenderSelect(str1)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}
