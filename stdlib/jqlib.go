package stdlib

import (
	"encoding/json"
	"rulex/typex"

	"github.com/itchyny/gojq"
	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// Loader
func LoadJqLib(e typex.RuleX, vm *lua.LState) {
	vm.Push(
		vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
			"JqSelect": func(l *lua.LState) int {
				jqExpression := l.ToString(1)
				data := l.ToString(2)
				var jsonData []interface{}
				if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
					vm.Push(lua.LNil)
					log.Error(err, jsonData, data)
					return 1
				}
				selectResult, err0 := JqSelect(jqExpression, jsonData)
				if err0 != nil {
					vm.Push(lua.LNil)
					log.Error(err0)
					return 1
				}
				resultString, err1 := json.Marshal(selectResult)
				if err1 != nil {
					log.Error(err1)
					vm.Push(lua.LNil)
					return 1
				}
				vm.Push(lua.LString(resultString))
				return 1
			},
		},
		),
	)
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

/**
* In either case, you cannot use custom type values as the query input.
* The type should be []interface{} for an array and map[string]interface{} for a map (just like decoded to an interface{} using the encoding/json package).
* You can't use []int or map[string]string, for example.
* If you want to query your custom struct, marshal to JSON, unmarshal to interface{} and use it as the query input.
 */
func JqSelect(jqExpression string, inputData []interface{}) ([]interface{}, error) {
	query, err0 := VerifyJqExpression(jqExpression)
	if err0 != nil {
		return nil, err0
	}
	var result []interface{}
	iterator := query.Run(inputData)
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
	return result, nil

}
