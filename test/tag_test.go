package test

import (
	"encoding/json"
	"rulex/core"
	"testing"
)

func Test_RenderConfig(t *testing.T) {
	//
	type __config struct {
		A    int32   `json:"a" validate:"required" title:"title--a" info:"aaaa" hidden:"true"`
		B    int64   `json:"b" validate:"required" title:"title--b" info:"bbbb" placeholder:"BBBBBBBBBBBB"`
		C    string  `json:"c" validate:"required" title:"title--c" info:"cccc" options:"cv1,cv2|cv3,cv4"`
		D    float32 `json:"d" validate:"required" title:"title--d" info:"dddd"`
		F    []int   `json:"f" validate:"required" title:"title--f" info:"ffff"`
		File string  `json:"file" validate:"required" title:"File" info:"File" file:"uploadfile"`
	}

	xcfgs, err := core.RenderConfig("MQTT", "MQTT", __config{})
	if err != nil {
		t.Fatal(err)
	} else {
		b, _ := json.MarshalIndent(xcfgs, "", " ")
		t.Log(string(b))
	}
}
