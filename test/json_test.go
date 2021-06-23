package test

import (
	"encoding/json"
	"testing"
)

func TestParseJson(t *testing.T) {

	s1 := make(map[string]interface{})
	s2 := make(map[string]interface{})
	b1 := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	b2 := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	t.Log(json.Unmarshal(b1, &s1))
	t.Log(json.Unmarshal(b2, &s2))
	t.Log(s1)
	t.Log(s2)
}
