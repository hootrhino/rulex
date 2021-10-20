package stdlib

import (
	"encoding/json"
	"rulex/typex"

	"github.com/itchyny/gojq"
	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

type JqLib struct {
}

func NewJqLib() typex.XLib {
	return &HttpLib{}
}
func (l *JqLib) Name() string {
	return "JqSelect"
}
func (l *JqLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		// LUA Args: Jq, Data ->
		// stack:  ------------
		//         |   Nil(0)  |
		//         ------------
		//         |   Jq Exp  |
		//         ------------
		//         |   Data    |
		//         ------------
		// Doc: https://github.com/lichuang/Lua-Source-Internal/blob/master/doc/ch03-Lua%E8%99%9A%E6%8B%9F%E6%9C%BA%E6%A0%88%E7%BB%93%E6%9E%84%E5%8F%8A%E7%9B%B8%E5%85%B3%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84.md
		jqExpression := stateStack.ToString(2)
		data := stateStack.ToString(3)
		var jsonData []interface{}
		if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
			stateStack.Push(lua.LNil)
			log.Error("Internal Error: ", err, ", InputData:", string(data))
		}
		selectResult, err0 := JqSelect(jqExpression, jsonData)
		if err0 != nil {
			stateStack.Push(lua.LNil)
			log.Error("JqSelect Error:", err0)
		}
		resultString, err1 := json.Marshal(selectResult)
		if err1 != nil {
			stateStack.Push(lua.LNil)
			log.Error("Json Marshal 'selectResult' error:", err1)
		}

		if string(resultString) == "[null]" {
			stateStack.Push(lua.LNil)
		} else {
			stateStack.Push(lua.LString(resultString))
		}
		return 1
	}
}

func VerifyJqExpression(jqExpression string) (*gojq.Query, error) {
	if query, err := gojq.Parse(jqExpression); err != nil {
		log.Error("VerifyJqExpression failed:", jqExpression, ", error:", err)
		return nil, err
	} else {
		return query, nil
	}
}

// JqSelect
/**
* In either case, you cannot use custom type values as the query input.
* The type should be []interface{} for an array and map[string]interface{} for a map (just like decoded to an interface{} using the encoding/json package).
* You can't use []int or map[string]string, for example.
* If you want to query your custom struct, marshal to JSON, unmarshal to interface{} and use it as the query input.
 */
func JqSelect(jqExpression string, inputData interface{}) ([]interface{}, error) {
	/**
	Input Data Json:
			[
				{  // },
				{  // }
			]
	*/
	query, err0 := VerifyJqExpression(jqExpression)
	if err0 != nil {
		return nil, err0
	}
	result := []interface{}{}
	iterator := query.Run(inputData)
	for {

		v, ok := iterator.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		//
		// iterator will return nil value, but we needn't nil.
		//
		if v != nil {
			result = append(result, v)
		}
	}
	return result, nil

}
