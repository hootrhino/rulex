package x

import (
	"encoding/json"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

// LoadSqlLib
func LoadSqlLib(e *RuleEngine, vm *lua.LState) int {
	mod := vm.SetFuncs(vm.G.Global, map[string]lua.LGFunction{
		"Select": func(vm *lua.LState) int {
			data := vm.ToString(1)
			sql := vm.ToString(2)
			selectResult, err0 := Select(data, sql)
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
		},
	})
	vm.Push(mod)
	return 1
}

//
func Select(data string, sql string) (*map[string]interface{}, error) {
	// log.Debug(data, sql)
	result, err0 := jsonStringToMap(data)
	if err0 != nil {
		return nil, err0
	} else {
		parseResult, err1 := SqlParse(result, sql)
		if err1 != nil {
			return nil, err1
		} else {
			return parseResult, nil
		}
	}
}

//
func jsonStringToMap(jsonString string) (*map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}
