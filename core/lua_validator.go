package core

import (
	"errors"
	"reflect"
	"rulex/typex"

	lua "github.com/yuin/gopher-lua"
)

// Verify Lua Syntax
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
	// validate Syntax
	actionsTable := vm.GetGlobal("Actions")
	if actionsTable != nil && actionsTable.Type() == lua.LTTable {
		valid := false
		actionsTable.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			valid = (reflect.TypeOf(f).Elem().Name() == "LFunction")
		})
		if !valid {
			return errors.New("Invalid function type")
		}
	} else {
		return errors.New("'Actions' must be a functions table")
	}
	return nil
}
