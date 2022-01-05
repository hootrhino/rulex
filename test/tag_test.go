package test

import (
	"encoding/json"
	"rulex/core"
	"testing"
)

func Test_Name(t *testing.T) {
	//
	type mqttConfig struct {
		A int32   `json:"a" validate:"required" title:"a" info:"aaaa"`
		B int64   `json:"b" validate:"required" title:"b" info:"bbbb"`
		C string  `json:"c" validate:"required" title:"c" info:"cccc"`
		D float32 `json:"d" validate:"required" title:"d" info:"dddd"`
		F []int   `json:"f" validate:"required" title:"f" info:"ffff"`
	}

	xcfg, err := core.RenderConfig(mqttConfig{})
	if err != nil {
		t.Fatal(err)
	} else {
		b, _ := json.Marshal(xcfg)
		t.Log(string(b))
	}
}
