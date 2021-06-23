package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

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
		end
	`
	err1 := luaVM.DoString(s1)
	if err1 != nil {
		panic(err1)
	}
	f := luaVM.GetGlobal("f")
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction), lua.LNumber(1), lua.LNumber(1))
		if state == lua.ResumeError {
			fmt.Println(err2.Error())
		}
	}
}
func TestRunLua(t *testing.T) {

	err1 := luaVM.DoString("print('helloworld')")
	if err1 != nil {
		panic(err1)
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

func TestRunLuaBench(t *testing.T) {
	var s1 = `
	function f(a,b)
		print("f=", a + b)
	end
`
	t1 := time.Now().UnixNano()
	err1 := luaVM.DoString(s1)
	fmt.Println("luaVM.DoString cost time:", time.Now().UnixNano()-t1, "ns")
	if err1 != nil {
		panic(err1)
	}
	t2 := time.Now().UnixNano()
	f := luaVM.GetGlobal("f")
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction), lua.LNumber(1), lua.LNumber(1))
		if state == lua.ResumeError {
			fmt.Println(err2.Error())
		}
	}
	fmt.Println("luaVM.Resume:", time.Now().UnixNano()-t2, "ns")

}
