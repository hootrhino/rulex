package core

import (
	"errors"

	"github.com/antonmedv/expr"
	lua "github.com/i4de/gopher-lua"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

const (
	SUCCESS_KEY string = "Success"
	FAILED_KEY  string = "Failed"
	ACTIONS_KEY string = "Actions"
)

// LUA Callback : Success
func ExecuteSuccess(vm *lua.LState) (interface{}, error) {
	return typex.Execute(vm, SUCCESS_KEY)
}

// LUA Callback : Failed

func ExecuteFailed(vm *lua.LState, arg lua.LValue) (interface{}, error) {
	return typex.Execute(vm, FAILED_KEY, arg)
}

/*
*
* Execute Expression
* https://expr.medv.io/docs/Getting-Started
*
 */
func ExecuteExpression(rule *typex.Rule, env map[string]interface{}) (interface{}, error) {
	return expr.Run(rule.ExprVM, env)
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
			return typex.RunPipline(rule.LuaVM, funcs, arg)
		}
		// if stopped, log warning information
		glogger.GLogger.Warn("Rule has stopped:" + rule.UUID)
		return lua.LNil, nil

	} else {
		return nil, errors.New("'Actions' not a lua table or not exist")
	}
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

/*
*
* 验证expr表达式的语法
*
 */
func VerifyExprSyntax(r *typex.Rule) error {
	env := map[string]interface{}{}
	_, err := expr.Compile(r.Expression, expr.Env(env))
	if err != nil {
		return err
	}
	return nil
}
