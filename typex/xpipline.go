package typex

import (
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

// RunPipline
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
					return nil, errors.New("'Action' callback first argument is must be bool")
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
		return nil, errors.New("'Action' callback must have 2 arguments:[bool, T]")
	} else {
		return f()
	}
}

//
//
// 执行lua函数的接口, 后期可以用这个接口来实现运行 lua 微服务
//
func Execute(vm *lua.LState, k string, args ...lua.LValue) (interface{}, error) {
	callable := vm.GetGlobal(k)
	if callable.Type() == lua.LTFunction {
		return callLuaFunc(vm, callable.(*lua.LFunction), args...)
	}
	return nil, errors.New("target:" + k + " is not a lua function")
}

// callLuaFunc
func callLuaFunc(vm *lua.LState, callable *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	if callable == nil {
		return nil, errors.New("callable function is not exists")
	} else {
		coroutine, _ := vm.NewThread()
		//
		// callback return value :lValues =[bool, T]
		//
		state, err, lValues := vm.Resume(coroutine, callable, args...)
		if state == lua.ResumeOK {
			//
			// only need T
			//
			return lValues[:2], nil
		}
		return nil, errors.New("lua run error, message is: " + err.Error())
	}
}
