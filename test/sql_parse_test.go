package test

import (
	"encoding/json"
	"rulenginex/x"

	"github.com/go-playground/assert/v2"
	"github.com/ngaut/log"

	"testing"
)

func TestJq1(t *testing.T) {

	jqExpression := `.[] | select(.id == 1)`
	inputData := []interface{}{
		map[string]interface{}{"id": 1, "name": "A1"},
		map[string]interface{}{"id": 2, "name": "A2"},
		map[string]interface{}{"id": 3, "name": "A3"},
		map[string]interface{}{"id": 4, "name": "A4"},
	}

	l1, _ := x.Select(jqExpression, &inputData)
	l2, _ := x.Select(jqExpression, &inputData)
	l3, _ := x.Select(jqExpression, &inputData)
	json1, _ := json.Marshal(l1)
	json2, _ := json.Marshal(l2)
	json3, _ := json.Marshal(l3)
	assert.Equal(t, `[{"id":1,"name":"A1"}] `, string(json1))
	assert.Equal(t, `[{"id":1,"name":"A1"}] `, string(json1))
	assert.Equal(t, `[{"id":1,"name":"A1"}] `, string(json1))
	log.Debug(string(json1))
	log.Debug(string(json2))
	log.Debug(string(json3))

}

func TestJq2(t *testing.T) {

	jqExpression1 := `.[] | select(.id == 1)|select(.temp == 10)`
	jqExpression2 := `.[] | select(.id == 3)|select(.temp > 10)`
	jqExpression3 := `.[] | select(.temp > 100.11)`
	jqExpression4 := `.[] | select(.hum == 44.5566)`
	inputData := []interface{}{
		map[string]interface{}{"id": 1, "name": "A1", "temp": 10, "hum": 20},
		map[string]interface{}{"id": 2, "name": "A2", "temp": 100.2343, "hum": 0},
		map[string]interface{}{"id": 3, "name": "A3", "temp": 0.03, "hum": 20.34},
		map[string]interface{}{"id": 4, "name": "A4", "temp": 12345676.4322454, "hum": 44.5566},
	}
	l1, _ := x.Select(jqExpression1, &inputData)
	json1, _ := json.Marshal(l1)
	log.Debug(string(json1))
	l2, _ := x.Select(jqExpression2, &inputData)
	json2, _ := json.Marshal(l2)
	log.Debug(string(json2))
	l3, _ := x.Select(jqExpression3, &inputData)
	json3, _ := json.Marshal(l3)
	log.Debug(string(json3))
	l4, _ := x.Select(jqExpression4, &inputData)
	json4, _ := json.Marshal(l4)
	log.Debug(string(json4))

}
