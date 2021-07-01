package x

import (
	"encoding/json"

	"github.com/itchyny/gojq"
	"github.com/ngaut/log"
	"github.com/yuin/gopher-lua"
)

// Loader
func LoadJqLib(e *RuleEngine, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"Select": func(l *lua.LState) int {
			jqExpression := l.ToString(1)
			data := l.ToString(2)
			log.Debug("Jq query:", jqExpression, data)
			jsonData := map[string]interface{}{}
			if json.Unmarshal([]byte(data), &jsonData) != nil {
				vm.Push(lua.LNil)
				return 1
			} else {
				selectResult, err0 := Select(jqExpression, &[]interface{}{jsonData})
				if err0 != nil {
					vm.Push(lua.LNil)
					return 1
				} else {
					jsonString, err1 := json.Marshal(selectResult)
					if err1 != nil {
						log.Error(err1)
						vm.Push(lua.LNil)
						return 1
					} else {
						vm.Push(lua.LString(string(jsonString)))
						return 1
					}
				}
			}
		},
	})
	vm.Push(mod)
	return 1
}

func VerifyJqExpression(jqExpression string) (*gojq.Query, error) {
	query, err0 := gojq.Parse(jqExpression)
	if err0 != nil {
		log.Error("VerifyJqExpression failed:", err0)
		return nil, err0
	} else {
		return query, nil
	}
}
func Select(jqExpression string, inputData *[]interface{}) (*[]interface{}, error) {
	query, err0 := VerifyJqExpression(jqExpression)
	if err0 != nil {
		return nil, err0
	} else {
		result := []interface{}{}
		iterator := query.Run(*inputData)
		for {
			v, ok := iterator.Next()
			if !ok {
				break
			}
			if err1, ok := v.(error); ok {
				return nil, err1
			}
			result = append(result, v)
		}
		return &result, nil
	}
}
