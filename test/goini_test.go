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

// go test -timeout 30s -run ^TestParse_EtcOsRelease github.com/hootrhino/rulex/test -v -count=1
func TestParse_EtcOsRelease(t *testing.T) {
	t.Log(CatOsRelease())
}
/*
*
* Linux: cat /etc/os-release
*
*/
func CatOsRelease() (map[string]string, error) {
	returnMap := map[string]string{}
	cfg, err := ini.ShadowLoad("./data/os-release.conf")
	if err != nil {
		return nil, err
	}
	DefaultSection, err := cfg.GetSection("DEFAULT")
	if err != nil {
		return nil, err
	}
	for _, Key := range DefaultSection.KeyStrings() {
		V, _ := DefaultSection.GetKey(Key)
		returnMap[Key] = V.String()
	}
	return returnMap, nil

}
