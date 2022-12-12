package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yuin/gopher-lua"
)

//global
var luaVM *lua.LState

//init
func init() {
	luaVM = lua.NewState()
}

var Script1 = `
function Success()
    print("======> success")
end
-- Failed
function Failed(error)
    print("======> failed:", error)
end

-- Actions
function Actions()
    return {
        function (data)
            
        end,
        function (data)
            
        end
    }
end
-- From
function From()
    return{"id=1","id=2"}
end
`

func TestTwoSameNameFunctions(t *testing.T) {
	var s1 = `
		function f()
			print("f1")
		end
		-- Failed
		function f()
			print("f2")
		end
	`
	err1 := luaVM.DoString(s1)
	if err1 != nil {
		panic(err1)
	}
	f := luaVM.GetGlobal("f")
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction))
		if state == lua.ResumeError {
			fmt.Println(err2.Error())
		}
	}
}
func TestAgs(t *testing.T) {
	var s1 = `
		function f(a,b)
			print("f=", a + b)
            print("f=", t1f:t1f())
		end
	`
	err1 := luaVM.DoString(s1)
	if err1 != nil {
		panic(err1)
	}
	f := luaVM.GetGlobal("f")
	luaVM.SetGlobal("t1f", luaVM.G.Global)
	luaVM.SetField(luaVM.G.Global, "t1f", luaVM.NewFunction(func(state *lua.LState) int {
		return 0
	}))
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction), lua.LNumber(1), lua.LNumber(1))
		if state == lua.ResumeError {
			fmt.Println(err2.Error())
		}
	}
}

func TestCallFailed(t *testing.T) {
	err1 := luaVM.DoString(Script1)
	if err1 != nil {
		panic(err1)
	} else {
		successCallBack := luaVM.GetGlobal("Failed").(*lua.LFunction)
		// TODO
		luaVM.SetGlobal("rule1.susscess", luaVM.NewFunction(func(L *lua.LState) int {
			return 0
		}))
		if successCallBack != nil {
			coroutine, _ := luaVM.NewThread()
			state, err2, values := luaVM.Resume(coroutine, successCallBack, lua.LString("This is error"))
			if state == lua.ResumeError {
				fmt.Println(err2.Error())
			}

			for i, lv := range values {
				fmt.Printf("%v : %v\n", i, lv)
			}

			if state == lua.ResumeOK {
				fmt.Println("yield break(ok)")
			}
		} else {
			fmt.Println("Nil")
		}

	}
}
func TestCF(t *testing.T) {
	var s1 = `
		function f()
            print("f=", M:Fun(1,2,3))
		end
	`
	err1 := luaVM.DoString(s1)
	if err1 != nil {
		panic(err1)
	}
	luaVM.SetGlobal("M", luaVM.G.Global)
	luaVM.SetField(luaVM.G.Global, "Fun", luaVM.NewFunction(func(state *lua.LState) int {
		fmt.Println("------------------------------------")
		n0 := state.ToNumber(0)
		n1 := state.ToNumber(1)
		n2 := state.ToNumber(2)
		n3 := state.ToNumber(3)
		n4 := state.ToNumber(4)
		fmt.Println(n0, n1, n2, n3, n4)
		fmt.Println("------------------------------------")
		state.Push(lua.LString("ok1"))
		state.Push(lua.LString("ok2"))
		state.Push(lua.LString("ok3"))
		return 3
	}))
	f := luaVM.GetGlobal("f")
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction))
		if state == lua.ResumeError {
			fmt.Println(err2.Error())
		}
	}
}
