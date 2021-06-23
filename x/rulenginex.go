package x

import (
	"errors"
	"reflect"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

//
func (e *RuleEngine) Start(sc StartCallback) *map[string]interface{} {
	e.ConfigMap = &map[string]interface{}{}
	(sc)()
	return e.ConfigMap
}

//
//
func (e *RuleEngine) GetConfig(k string) interface{} {
	return (*e.ConfigMap)[k]
}

func (e *RuleEngine) LoadInEnds(in *InEnd) error {
	return tryCreateInEnd(in, e)

}

//
func tryCreateInEnd(in *InEnd, e *RuleEngine) error {
	// TODO: support more type at future
	if in.Type == "MQTT" {

		return startResources(NewMqttInEndResource(in.Id), in, e)
	}
	if in.Type == "HTTP" {
		return startResources(NewHttpInEndResource(in.Id), in, e)
	}
	if in.Type == "COAP" {
		return startResources(NewCoAPInEndResource(in.Id), in, e)
	}
	return errors.New("unsupported rule type:" + in.Type)
}

//
func startResources(r Resource, in *InEnd, e *RuleEngine) error {
	SaveInEnd(in)
	return r.Start(e, func() {
		if err := r.Register(in.Id); err == nil {
			log.Info("Start InEnd:", in.Name, ",", in.Id, " Start successfully")
		} else {
			log.Fatal("Start InEnd:", in.Name, ",", in.Id, " Start failure, error:", err)
		}
	}, func(err error) {
		log.Error(err)
	})
}

// LoadOutEnds
func (e *RuleEngine) LoadOutEnds(out *OutEnd) {
}

// LoadRules
func (e *RuleEngine) LoadRules(r *Rule) error {
	if e := VerifyCallback(r); e != nil {
		return e
	} else {
		if len(r.From) > 0 {
			for _, inId := range r.From {
				if in := GetInEnd(inId); in != nil {
					(*in.Binds)[r.Id] = *r
					SaveRule(r)
					return nil
				} else {
					return errors.New("InEnd:" + inId + " is not exists")
				}
			}
		}
	}
	return errors.New("from can not be empty")

}

// Stop
func (e *RuleEngine) Stop() {
}

// LoadStdLib
func (e *RuleEngine) LoadStdLib() {

}

// RunSuccessCallback
func (e *RuleEngine) RunSuccessCallback(ruleId string) {

}

// RunFailedCallback
func (e *RuleEngine) RunFailedCallback(ruleId string) {

}

// Work
func (e *RuleEngine) Work(in *InEnd, data string) (bool, error) {
	// log.Warn("[TODO] RuleEngine Work Find rule with in id:", in.Id)
	for _, rule := range *in.Binds {
		_, err0 := rule.ExecuteActions(lua.LString(data))
		if err0 != nil {
			rule.ExecuteFailed(lua.LString(err0.Error()))
			return false, err0
		} else {
			rule.ExecuteSuccess()
			return true, nil

		}
	}
	return false, nil
}

// Verify Lua Syntax
func VerifyCallback(r *Rule) error {
	e1 := r.VM.DoString(r.Success)
	if e1 != nil {
		return e1
	}
	e2 := r.VM.DoString(r.Failed)
	if e2 != nil {
		return e1
	}
	e3 := r.VM.DoString(r.Actions)
	if e3 != nil {
		return e1
	}
	return nil
}

//
func (r *Rule) ExecuteActions(arg lua.LValue) (lua.LValue, error) {
	table := r.VM.GetGlobal("Actions")
	if table != nil && table.Type() == lua.LTTable {
		funcs := make(map[string]*lua.LFunction)
		table.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			t := reflect.TypeOf(f).Elem().Name()
			if t == "LFunction" {
				funcs[idx.String()] = f.(*lua.LFunction)
			}
		})
		return runPipline(r.VM, funcs, arg)
	}
	return nil, errors.New("not a lua table")
}

func (r *Rule) ExecuteSuccess() {
	execute(r.VM, "Success")
}

func (r *Rule) ExecuteFailed(arg lua.LValue) {
	execute(r.VM, "Failed", arg)

}

// Execute Lua function
func execute(vm *lua.LState, k string, args ...lua.LValue) (interface{}, error) {
	callable := vm.GetGlobal(k)
	name := reflect.TypeOf(callable).Elem().Name()
	if name == "LFunction" {
		return callLuaFunc(vm, callable.(*lua.LFunction), args...)
	}
	if name == "LNilType" {
		return nil, errors.New("target:" + k + " is not exists")
	}
	return nil, errors.New("target:" + k + " is n	ot a lua function")
}

// callLuaFunc
func callLuaFunc(vm *lua.LState, callable *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	if callable == nil {
		return nil, errors.New("callable function is not exists")
	} else {
		coroutine, _ := vm.NewThread()
		state, err1, lValues := vm.Resume(coroutine, callable, args...)
		if state != lua.ResumeOK {
			return nil, err1
		} else {
			return lValues, nil
		}
	}
}
