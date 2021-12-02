package core

import (
	"errors"
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

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
