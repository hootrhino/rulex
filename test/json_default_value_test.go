package test

import (
	"encoding/json"
	"testing"
)

type _data1 struct {
	Tag     string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	EndFlag string `json:"endFlag" validate:"required" title:"采集频率"`
}

func Test_json_default_value(t *testing.T) {
	jsonS := `{
		"tag"    : "2016",
		"endFlag": "\n"
	}
	`
	data := _data1{}
	err := json.Unmarshal([]byte(jsonS), &data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data.Tag)
	t.Log([]byte(data.EndFlag)[0] == '\n')

}
