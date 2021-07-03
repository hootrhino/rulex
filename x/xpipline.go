package x

import (
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

//
//  Run lua as pipline
//
func runPipline(vm *lua.LState, funcs map[string]*lua.LFunction, arg lua.LValue) (lua.LValue, error) {
	acc := 1
	return pipLine(vm, acc, funcs, arg)
}

//
func pipLine(vm *lua.LState, acc int, funcs map[string]*lua.LFunction, arg lua.LValue) (lua.LValue, error) {
	if acc == len(funcs) {
		values, err0 := callLuaFunc(vm, funcs[strconv.Itoa(acc)], arg)
		if err0 != nil {
			return nil, err0
		} else {
			return validate(values, func() (lua.LValue, error) {
				result := values[1]
				return result, nil
			})
		}
	} else {
		values, err0 := callLuaFunc(vm, funcs[strconv.Itoa(acc)], arg)
		if err0 != nil {
			return nil, err0
		} else {
			return validate(values, func() (lua.LValue, error) {
				next := values[0]
				result := values[1]
				if next.Type() == lua.LTBool {
					if next.(lua.LBool) {
						return pipLine(vm, acc+1, funcs, result)
					} else {
						return result, nil
					}
				} else {
					return nil, errors.New("action callback first argument is must be bool")
				}
			})
		}
	}
}

func validate(values []lua.LValue, f func() (lua.LValue, error)) (lua.LValue, error) {
	if len(values) != 2 {
		return nil, errors.New("action callback must have 2 arguments:[bool, T]")
	} else {
		return f()
	}
}
