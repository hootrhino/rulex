package test

import (
	"testing"
	"time"
)

//
type InEnd struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//
	Config map[string]interface{} `json:"config"`
}

func Test_Map_delete(t *testing.T) {
	m := map[string]InEnd{
		"M0":  {UUID: "5"},
		"M1":  {UUID: "1"},
		"M2":  {UUID: "2"},
		"M3":  {UUID: "3"},
		"M4":  {UUID: "4"},
		"M5":  {UUID: "5"},
		"M6":  {UUID: "5"},
		"M7":  {UUID: "5"},
		"M8":  {UUID: "5"},
		"M9":  {UUID: "5"},
		"M10": {UUID: "5"},
	}
	start := time.Now()
	delete(m, "M1")
	elapsed := time.Since(start)
	t.Log("Cost：", elapsed)
}
