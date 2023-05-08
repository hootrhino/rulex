package test

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

// go test -timeout 30s -run ^TestOk github.com/hootrhino/rulex/test -v -count=1
func Test_JSON_contains_byte(t *testing.T) {
	type _data_ struct {
		Data  string
		Data1 []byte
	}
	d := _data_{
		Data:  string([]rune{0x0A, 0xAC, 0x0A, 0xAC}),
		Data1: []byte{0x0A, 0xAC, 0x0A, 0xAC},
	}
	m := map[string]_data_{}
	m["1"] = d
	bytes, _ := json.Marshal(m)
	t.Log("Data1[0]: ", m["1"].Data1[0])
	t.Log("Data1[1]: ", m["1"].Data1[1])
	t.Log("Data1[2]: ", m["1"].Data1[2])
	t.Log("Data1[3]: ", m["1"].Data1[3])
	t.Log(string(bytes))
	bss, err := base64.StdEncoding.DecodeString("CqwKrA==")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bss)
}
