package interpipeline

import (
	"errors"
	"strconv"

	lua "github.com/hootrhino/gopher-lua"
)

// RunPipline
//
//	Run lua as pipline
func RunPipline(vm *lua.LState, funcs map[string]*lua.LFunction, arg lua.LValue) (lua.LValue, error) {
	// start 1
	acc := 1
	return pipLine(vm, acc, funcs, arg)
}

func pipLine(vm *lua.LState, acc int, funcs map[string]*lua.LFunction, arg lua.LValue) (lua.LValue, error) {
	if acc == len(funcs) {
		values, err0 := callLuaFunc(vm, funcs[strconv.Itoa(acc)], arg)
		if err0 != nil {
			return nil, err0
		}
		return validate(values, func() (lua.LValue, error) {
			result := values[1]
			return result, nil
		})

	}
	values, err0 := callLuaFunc(vm, funcs[strconv.Itoa(acc)], arg)
	if err0 != nil {
		return nil, err0
	}
	return validate(values, func() (lua.LValue, error) {
		next := values[0]
		result := values[1]
		if next.Type() == lua.LTBool {
			if next.(lua.LBool) {
				return pipLine(vm, acc+1, funcs, result)
			}
			return result, nil
		}
		return nil, errors.New("'Action' callback first argument is must be bool")

	})

}

// validate lua callback
func validate(values []lua.LValue, f func() (lua.LValue, error)) (lua.LValue, error) {
	// Lua call back must have 2 args!!!
	if len(values) != 2 {
		return nil, errors.New("'Action' callback must have 2 return value:[bool, T]")
	} else {
		return f()
	}
}

// 执行lua函数的接口, 后期可以用这个接口来实现运行 lua 微服务
func Execute(vm *lua.LState, k string, args ...lua.LValue) (interface{}, error) {
	callable := vm.GetGlobal(k)
	if callable.Type() == lua.LTFunction {
		return callLuaFunc(vm, callable.(*lua.LFunction), args...)
	}
	return nil, errors.New("target:" + k + " is not a lua function")
}

/*
*
*

	callLuaFunc

*
*/
func callLuaFunc(vm *lua.LState, callable *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	if callable == nil {
		return nil, errors.New("callable function is not exists")
	}
	err := vm.CallByParam(lua.P{
		Fn:      callable,
		NRet:    2,
		Protect: true,
	}, args...)
	if err != nil {
		return nil, err
	}
	A1 := vm.Get(-1)
	A2 := vm.Get(-2)
	//   _
	//5 |_| -1
	//4 |_| -2
	//3 |_| -3
	//2 |_| -4
	//1 |_| -5
	// C.call(Lst, arg1,arg2) -> -1
	// 返回值在栈顶
	// Pop的作用是把压入栈的参数args...删除，防止registry无限增长
	vm.Pop(2) // 这是我很早以前写的代码，
	//        // 但是今天突然发现看不懂为啥这里要Pop栈？
	// 2023 12-01 删除
	// return data, true|false
	// -1栈顶；1栈底
	// 2024-1-23: 终于明白了，这是退栈操作，防止lua的registry溢出。
	//            今晚真是个难眠之夜，通宵把这个问题解决了。

	return []lua.LValue{A2, A1}, nil

}

// #include <lua.h>
// #include <lauxlib.h>
// #include <lualib.h>

// int main() {
//     lua_State *L = luaL_newstate();
//     luaL_openlibs(L);
//     if (luaL_dofile(L, "example.lua") != LUA_OK) {
//         fprintf(stderr, "Error running Lua script: %s\n", lua_tostring(L, -1));
//         lua_close(L);
//         return 1;
//     }
//     lua_getglobal(L, "my_lua_function");
//     lua_pushinteger(L, 10);
//     lua_pushinteger(L, 20);
//     if (lua_pcall(L, 2, 2, 0) != LUA_OK) {
//         fprintf(stderr, "Error calling Lua function: %s\n", lua_tostring(L, -1));
//         lua_close(L);
//         return 1;
//     }
//     int result1 = lua_tointeger(L, -2);
//     int result2 = lua_tointeger(L, -1);
//     printf("Result 1: %d\n", result1);
//     printf("Result 2: %d\n", result2);
//     lua_close(L);
//     return 0;
// }
