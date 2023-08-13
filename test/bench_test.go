package test

import (
	"reflect"
	"testing"
	"time"

	lua "github.com/hootrhino/gopher-lua"
)

func TestRunLuaBench(t *testing.T) {
	luaVM := lua.NewState()
	var s1 = `
	function f(a,b)
		print("f=", a + b)
	end
`
	t1 := time.Now().UnixNano()
	err1 := luaVM.DoString(s1)
	t.Log("luaVM.DoString cost time:", time.Now().UnixNano()-t1, "ns")
	if err1 != nil {
		t.Fatal(err1)
	}
	t2 := time.Now().UnixNano()
	f := luaVM.GetGlobal("f")
	t.Log(reflect.TypeOf(f).String() == "*lua.LFunction")
	if reflect.TypeOf(f).Elem().Name() == "LFunction" {
		coroutine, _ := luaVM.NewThread()
		state, err2, _ := luaVM.Resume(coroutine, f.(*lua.LFunction), lua.LNumber(1), lua.LNumber(1))
		if state == lua.ResumeError {
			t.Error(err2)
		}
	}
	t.Log("luaVM.Resume:", time.Now().UnixNano()-t2, "ns")

}
