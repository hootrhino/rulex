package typex

import (
	"errors"
	"reflect"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

//
//  Run lua as pipline
//
func RunPipline(vm *lua.LState, funcs map[string]*lua.LFunction, arg lua.LValue) (lua.LValue, error) {
	// start 1
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

//
// validate lua callback
//
func validate(values []lua.LValue, f func() (lua.LValue, error)) (lua.LValue, error) {
	// Lua call back must have 2 args!!!
	if len(values) != 2 {
		return nil, errors.New("action callback must have 2 arguments:[bool, T]")
	} else {
		return f()
	}
}

//
//
//

// Execute Lua function
func Execute(vm *lua.LState, k string, args ...lua.LValue) (interface{}, error) {
	callable := vm.GetGlobal(k)
	name := reflect.TypeOf(callable).Elem().Name()
	if name == "LFunction" {
		return callLuaFunc(vm, callable.(*lua.LFunction), args...)
	}
	if name == "LNilType" {
		return nil, errors.New("Target:" + k + " is not exists")
	}
	return nil, errors.New("Target:" + k + " is not a lua function")
}

// callLuaFunc
func callLuaFunc(vm *lua.LState, callable *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	if callable == nil {
		return nil, errors.New("Callable function is not exists")
	} else {
		coroutine, _ := vm.NewThread()
		//
		// callback return value :lValues =[bool, T]
		//
		state, err, lValues := vm.Resume(coroutine, callable, args...)
		if state == lua.ResumeError {
			return nil, errors.New("current state is not lua.ResumeOK:" + err.Error())
		}
		if state == lua.ResumeOK {
			// only need T
			return lValues[1:], nil
		}
		return nil, errors.New("current state is not lua.ResumeOK")
	}
}
