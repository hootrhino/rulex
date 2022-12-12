package test

import (
	"testing"

	"gopkg.in/ini.v1"
)

func Test_map_ini_file(t *testing.T) {
	type IP struct {
		Value []string `ini:"value,omitempty,allowshadow"`
	}
	cfg, _ := ini.ShadowLoad("./data/test_ini_1.ini")
	i := IP{}
	err1 := cfg.Section("IP").MapTo(&i)
	if err1 != nil {
		t.Fatal(err1)
	}
}
