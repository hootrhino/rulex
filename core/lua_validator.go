package core

import (
	"errors"
	"reflect"
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

// LUA Callback : Success
func ExecuteSuccess(vm *lua.LState) (interface{}, error) {
	return typex.Execute(vm, "Success")
}

// LUA Callback : Failed

func ExecuteFailed(vm *lua.LState, arg lua.LValue) (interface{}, error) {
	return typex.Execute(vm, "Failed", arg)
}

//
func ExecuteActions(rule *typex.Rule, arg lua.LValue) (lua.LValue, error) {
	table := rule.VM.GetGlobal("Actions")
	if table != nil && table.Type() == lua.LTTable {
		funcs := make(map[string]*lua.LFunction)
		table.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			t := reflect.TypeOf(f).Elem().Name()
			if t == "LFunction" {
				funcs[idx.String()] = f.(*lua.LFunction)
			}
		})
		return typex.RunPipline(rule.VM, funcs, arg)
	} else {
		return nil, errors.New("'Actions' not a lua table or not exist")
	}
}

// VerifyCallback Verify Lua Syntax
func VerifyCallback(r *typex.Rule) error {
	vm := r.VM
	if err := vm.DoString(r.Success); err != nil {
		return err
	}
	if vm.GetGlobal("Success").Type() != lua.LTFunction {
		return errors.New("'Success' callback function missed")
	}

	if err := vm.DoString(r.Failed); err != nil {
		return err
	}
	if vm.GetGlobal("Failed").Type() != lua.LTFunction {
		return errors.New("'Failed' callback function missed")
	}
	if err := vm.DoString(r.Actions); err != nil {
		return err
	}
	//
	// validate lua syntax
	//
	actionsTable := vm.GetGlobal("Actions")
	if actionsTable != nil && actionsTable.Type() == lua.LTTable {
		valid := true
		actionsTable.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			//
			// golang function in lua is '*lua.LFunction' type
			//
			if !(f.Type() == lua.LTFunction) {
				valid = false
			}
		})
		if !valid {
			return errors.New("invalid function type")
		}
	} else {
		return errors.New("'Actions' must be a functions table")
	}
	return nil
}
