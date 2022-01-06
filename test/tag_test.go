package test

import (
	"encoding/json"
	"rulex/core"
	"testing"
)

func Test_Name(t *testing.T) {
	//
	type mqttConfig struct {
		A    int32   `json:"a" validate:"required" title:"title--a" info:"info--aaaa"`
		B    int64   `json:"b" validate:"required" title:"title--b" info:"info--bbbb"`
		C    string  `json:"c" validate:"required" title:"title--c" info:"info--cccc" enum:"cv1,cv2|cv3,cv4"`
		D    float32 `json:"d" validate:"required" title:"title--d" info:"info--dddd"`
		F    []int   `json:"f" validate:"required" title:"title--f" info:"info--ffff"`
		File string  `json:"file" validate:"required" title:"File" info:"File" file:"uploadfile"`
	}

	xcfg, err := core.RenderConfig(mqttConfig{})
	if err != nil {
		t.Fatal(err)
	} else {
		b, _ := json.MarshalIndent(xcfg, "", " ")
		t.Log(string(b))
	}
}
