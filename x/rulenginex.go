package x

import (
	"container/list"
	"context"
	"errors"
	"reflect"
	"rulenginex/statistics"
	"sync"
	"time"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

type TargetState int

var lock sync.Mutex

const (
	UP   TargetState = 1
	DOWN TargetState = 0
)

//
//
//
type RuleEngine struct {
	Hooks     *map[string]XHook            `json:"-"`
	Plugins   *map[string]*XPluginMetaInfo `json:"plugins"`
	InEnds    *map[string]*inEnd           `json:"inends"`
	OutEnds   *map[string]*outEnd          `json:"outends"`
	ConfigMap *map[string]interface{}      `json:"configMap"`
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		Plugins:   &map[string]*XPluginMetaInfo{},
		InEnds:    &map[string]*inEnd{},
		OutEnds:   &map[string]*outEnd{},
		ConfigMap: &map[string]interface{}{},
	}
}

//
func (e *RuleEngine) Start(sc func()) *map[string]interface{} {
	e.ConfigMap = &map[string]interface{}{}
	(sc)()
	return e.ConfigMap
}

//
//
func (e *RuleEngine) GetConfig(k string) interface{} {
	return (*e.ConfigMap)[k]
}

func (e *RuleEngine) LoadInEnds(in *inEnd) error {
	return tryCreateInEnd(in, e)

}

//
func tryCreateInEnd(in *inEnd, e *RuleEngine) error {
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
func startResources(r XResource, in *inEnd, e *RuleEngine) error {
	log.Info("Starting InEnd Resources:", in.Name)

	if r.Test(in.Id) {
		e.SaveInEnd(in)
		if err := r.Register(in.Id); err != nil {
			return err
		} else {
			return r.Start(e)
		}
	} else {
		return errors.New("Resources start failed:" + in.Name)
	}
}

//
//
// LoadOutEnds
func (e *RuleEngine) LoadOutEnds(out *outEnd) error {
	return tryCreateOutEnd(out, e)
}

//
func tryCreateOutEnd(out *outEnd, e *RuleEngine) error {
	if out.Type == "mongo" {
		return startTarget(NewMongoTarget(), out, e)
	}
	return errors.New("unsupported target type:" + out.Type)

}

//
// Start output target
//
// Target life cycle:
// Register -> Start -> Test
//
//
func startTarget(target XTarget, out *outEnd, e *RuleEngine) error {
	log.Info("Starting OutEnd Target:", out.Name)
	// Important!!! Must save outend first
	e.SaveOutEnd(out)
	out.Target = target
	// Register outend to target
	if err0 := target.Register(out.Id); err0 != nil {
		return err0
	} else {
		if err1 := target.Start(e); err1 != nil {
			return err1
		} else {
			// \!!!
			testTargetState(target, e, out.Id)
			//
			go func(ctx context.Context) {
				// 5 seconds
				ticker := time.NewTicker(time.Duration(time.Second * 5))
				defer target.Stop()
				for {
					<-ticker.C
					testTargetState(target, e, out.Id)
				}
			}(context.Background())
			return nil
		}
	}
}

// Test Target State
func testTargetState(target XTarget, e *RuleEngine, id string) {
	if !target.Test(id) {
		e.GetOutEnd(id).SetState(DOWN)
		log.Errorf("Target %s DOWN", id)
	} else {
		if e.GetOutEnd(id).GetState() == DOWN {
			e.GetOutEnd(id).SetState(UP)
			log.Warnf("Target %s recover to UP", id)
		}
	}
}

// LoadRules
func (e *RuleEngine) LoadRules(r *rule) error {
	if err := VerifyCallback(r); err != nil {
		return err
	} else {
		if len(r.From) > 0 {
			for _, inId := range r.From {
				if in := e.GetInEnd(inId); in != nil {
					// Bind to rule, Key:RuleId, Value: Rule
					// RULE_0f8619ef-3cf2-452f-8dd7-aa1db4ecfdde {
					// ...
					// ...
					// }
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

// RunSuccessCallback
func (e *RuleEngine) RunSuccessCallback(ruleId string) {

}

// RunFailedCallback
func (e *RuleEngine) RunFailedCallback(ruleId string) {

}

// Work
func (e *RuleEngine) Work(in *inEnd, data string) (bool, error) {
	statistics.IncIn()
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
func VerifyCallback(r *rule) error {
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
func (r *rule) ExecuteActions(arg lua.LValue) (lua.LValue, error) {
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

func (r *rule) ExecuteSuccess() (interface{}, error) {
	return execute(r.VM, "Success")
}

func (r *rule) ExecuteFailed(arg lua.LValue) (interface{}, error) {
	return execute(r.VM, "Failed", arg)

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

//
func (e *RuleEngine) LoadPlugin(p XPlugin) error {
	env := p.Load(e)
	err0 := p.Init(env)
	if err0 != nil {
		return err0
	} else {
		metaInfo, err1 := p.Install(env)
		if err1 != nil {
			return err1
		} else {
			if (*e.Plugins)[metaInfo.Name] != nil {
				return errors.New("plugin already instaled:" + metaInfo.Name)
			} else {
				(*e.Plugins)[metaInfo.Name] = metaInfo
				if err2 := p.Start(e, env); err2 != nil {
					return err2
				}
				return nil
			}

		}
	}
}

//
//
//
func (e *RuleEngine) GetInEnd(id string) *inEnd {
	lock.Lock()
	defer lock.Unlock()
	return (*e.InEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveInEnd(in *inEnd) {
	lock.Lock()
	defer lock.Unlock()
	(*e.InEnds)[in.Id] = in
}

//
//
//
func (e *RuleEngine) RemoveInEnd(in *inEnd) {
	lock.Lock()
	defer lock.Unlock()
	delete((*e.InEnds), in.Id)
}

//
//
//
func (e *RuleEngine) AllInEnd() map[string]*inEnd {
	lock.Lock()
	defer lock.Unlock()
	return (*e.InEnds)
}

//
//
//
func (e *RuleEngine) GetOutEnd(id string) *outEnd {
	lock.Lock()
	defer lock.Unlock()
	return (*e.OutEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveOutEnd(out *outEnd) {
	lock.Lock()
	defer lock.Unlock()
	(*e.OutEnds)[out.Id] = out
}

//
//
//
func (e *RuleEngine) RemoveOutEnd(out *outEnd) {
	lock.Lock()
	defer lock.Unlock()
	delete((*e.OutEnds), out.Id)
}

//
//
//
func (e *RuleEngine) AllOutEnd() map[string]*outEnd {
	lock.Lock()
	defer lock.Unlock()
	return (*e.OutEnds)
}

//
// LoadHook
//
func (e *RuleEngine) LoadHook(h XHook) error {
	if (*e.Hooks)[h.Name()] != nil {
		return errors.New("hook have been loaded")
	} else {
		(*e.Hooks)[h.Name()] = h
		return nil
	}
}

//
// RunHooks
//
func (e *RuleEngine) RunHooks(data string) {
	for _, h := range *e.Hooks {
		if err := runHook(h, data); err != nil {
			log.Error("run hook:", h.Name(), " failed, error is:", err)
		}
	}
}
func runHook(h XHook, data string) error {
	return h.Work(data)
}
