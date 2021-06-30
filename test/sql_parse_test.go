package test

import (
	"encoding/json"
	"rulenginex/x"

	"github.com/go-playground/assert/v2"
	"github.com/ngaut/log"

	"testing"
)

func TestJq(t *testing.T) {

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
