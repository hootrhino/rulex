package test

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	parse "github.com/hootrhino/gopher-lua/parse"
)

// global
var luaVM *lua.LState

// init
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

//

func Test_loop_close(t *testing.T) {
	var s1 = `
		function f()
			while true do
			print("Hello World")
			end
		end
		f()
	`
	var luaVM = lua.NewState()
	go func() {
		err1 := luaVM.DoString(s1)
		if err1 != nil {
			panic(err1)
		}
	}()
	time.Sleep(2 * time.Second)
	luaVM.Close()
	time.Sleep(3 * time.Second)
}

// CompileLua reads the passed lua file from disk and compiles it.
func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

// Example shows how to share the compiled byte code from a lua script between multiple VMs.
func TestCompileLua(t *testing.T) {
	codeToShare, _ := CompileLua("lua/_exit.lua")
	t.Log(codeToShare.Code)
	// a := lua.NewState()
	// b := lua.NewState()
	// c := lua.NewState()
	// DoCompiledFile(a, codeToShare)
	// DoCompiledFile(b, codeToShare)
	// DoCompiledFile(c, codeToShare)
}

func Test_Stack_order(t *testing.T) {
	var s1 = `
	    A=1
		B=2
		function __f1()
		end
		function __f2()
		end
		function __f3()
		end
	`
	var luaVM = lua.NewState()

	err1 := luaVM.DoString(s1)
	if err1 != nil {
		panic(err1)
	}
	luaVM.G.Global.ForEach(func(l1, l2 lua.LValue) {

		if l2.Type() == lua.LTFunction {
			if l1.String()[:3] == "__f" {
				fc := l2.(*lua.LFunction)
				t.Log(fc.Proto)
			}

		}
	})
	luaVM.Close()
}
