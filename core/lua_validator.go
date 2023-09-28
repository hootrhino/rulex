// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/component/interpipeline"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

const (
	SUCCESS_KEY string = "Success"
	FAILED_KEY  string = "Failed"
	ACTIONS_KEY string = "Actions"
)

// LUA Callback : Success
func ExecuteSuccess(vm *lua.LState) (interface{}, error) {
	return interpipeline.Execute(vm, SUCCESS_KEY)
}

// LUA Callback : Failed

func ExecuteFailed(vm *lua.LState, arg lua.LValue) (interface{}, error) {
	return interpipeline.Execute(vm, FAILED_KEY, arg)
}

/*
*
* Execute Lua Callback
*
 */
func ExecuteActions(rule *typex.Rule, arg lua.LValue) (lua.LValue, error) {
	// 原始 lua 数据结构
	luaOriginTable := rule.LuaVM.GetGlobal(ACTIONS_KEY)
	if luaOriginTable != nil && luaOriginTable.Type() == lua.LTTable {
		// 断言成包含回调的 table
		funcsTable := luaOriginTable.(*lua.LTable)
		funcs := make(map[string]*lua.LFunction, funcsTable.Len())
		var err error = nil
		funcsTable.ForEach(func(idx, f lua.LValue) {
			if f.Type() == lua.LTFunction {
				funcs[idx.String()] = f.(*lua.LFunction)
			} else {
				err = errors.New(f.String() + " not a lua function")
				return
			}
		})
		if err != nil {
			return nil, err
		}
		if rule.Status != typex.RULE_STOP {
			return interpipeline.RunPipline(rule.LuaVM, funcs, arg)
		}
		// if stopped, log warning information
		glogger.GLogger.Warn("Rule has stopped:" + rule.UUID)
		return lua.LNil, nil

	}
	return nil, errors.New("'Actions' not a lua table or not exist")

}

// VerifyLuaSyntax Verify Lua Syntax
func VerifyLuaSyntax(r *typex.Rule) error {
	tempVm := lua.NewState(lua.Options{
		SkipOpenLibs:     true,
		RegistrySize:     0,
		RegistryMaxSize:  0,
		RegistryGrowStep: 0,
	})

	if err := tempVm.DoString(r.Success); err != nil {
		return err
	}
	if tempVm.GetGlobal(SUCCESS_KEY).Type() != lua.LTFunction {
		return errors.New("'Success' callback function missed")
	}

	if err := tempVm.DoString(r.Failed); err != nil {
		return err
	}
	if tempVm.GetGlobal(FAILED_KEY).Type() != lua.LTFunction {
		return errors.New("'Failed' callback function missed")
	}
	if err := tempVm.DoString(r.Actions); err != nil {
		return err
	}
	//
	// validate lua syntax
	//
	actionsTable := tempVm.GetGlobal(ACTIONS_KEY)
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
	// 释放语法验证阶段的临时虚拟机
	tempVm.Close()
	tempVm = nil
	// 交给规则脚本
	r.LuaVM.DoString(r.Success)
	r.LuaVM.DoString(r.Actions)
	r.LuaVM.DoString(r.Failed)
	return nil
}
