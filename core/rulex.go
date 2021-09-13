package core

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"rulex/statistics"
	"runtime"
	"sync"
	"time"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

//
// RuleEngine
//
type RuleEngine struct {
	sync.Mutex
	Hooks     *map[string]XHook            `json:"-"`
	Rules     *map[string]*rule            `json:"-"`
	Plugins   *map[string]*XPluginMetaInfo `json:"plugins"`
	InEnds    *map[string]*inEnd           `json:"inends"`
	OutEnds   *map[string]*outEnd          `json:"outends"`
	ConfigMap *map[string]interface{}      `json:"configMap"`
}

//
//
//
func NewRuleEngine() RuleX {
	return &RuleEngine{
		Plugins:   &map[string]*XPluginMetaInfo{},
		Hooks:     &map[string]XHook{},
		Rules:     &map[string]*rule{},
		InEnds:    &map[string]*inEnd{},
		OutEnds:   &map[string]*outEnd{},
		ConfigMap: &map[string]interface{}{},
	}
}

//
//
//
func (e *RuleEngine) GetPlugins() *map[string]*XPluginMetaInfo {
	e.Lock()
	defer e.Unlock()
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() *map[string]*XPluginMetaInfo {
	e.Lock()
	defer e.Unlock()
	return e.Plugins
}

//
func (e *RuleEngine) Start() *map[string]interface{} {
	e.ConfigMap = &map[string]interface{}{}
	//
	defaultBanner :=
		`
-----------------------------------------------------------
~~~/=====\       ██████╗ ██╗   ██╗██╗     ███████╗██╗  ██╗
~~~||\\\||--->o  ██╔══██╗██║   ██║██║     ██╔════╝╚██╗██╔╝
~~~||///||--->o  ██████╔╝██║   ██║██║     █████╗   ╚███╔╝ 
~~~||///||--->o  ██╔══██╗██║   ██║██║     ██╔══╝   ██╔██╗ 
~~~||\\\||--->o  ██║  ██║╚██████╔╝███████╗███████╗██╔╝ ██╗
~~~\=====/       ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝
-----------------------------------------------------------
`
	file, err := os.Open("conf/banner.txt")
	if err != nil {
		log.Warn("No banner found, print default banner")
		log.Info(defaultBanner)
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Warn("No banner found, print default banner")
			log.Info(defaultBanner)
		} else {
			log.Info("\n", string(data))
		}
	}
	log.Info("rulex start successfully")
	file.Close()
	return e.ConfigMap
}

//
//
func (e *RuleEngine) GetConfig(k string) interface{} {
	return (*e.ConfigMap)[k]
}

func (e *RuleEngine) LoadInEnd(in *inEnd) error {
	if in.Type == "MQTT" {
		return startResources(NewMqttInEndResource(in.Id, e), in, e)
	}
	if in.Type == "HTTP" {
		return startResources(NewHttpInEndResource(in.Id, e), in, e)
	}
	if in.Type == "COAP" {
		return startResources(NewCoAPInEndResource(in.Id, e), in, e)
	}
	if in.Type == "GRPC" {
		return startResources(NewGrpcInEndResource(in.Id, e), in, e)
	}
	if in.Type == "LoraATK" {
		return startResources(NewLoraModuleResource(in.Id, e), in, e)
	}
	return errors.New("unsupported rule type:" + in.Type)
}

//
func startResources(resource XResource, in *inEnd, e *RuleEngine) error {
	// Save to rule engine first
	// 这么作主要是为了可以 预加载 进去，然后等环境恢复了以后自动复原
	e.SaveInEnd(in)
	// 首先把资源ID给注册进去，作为资源的全局索引
	if err := resource.Register(in.Id); err != nil {
		log.Error(err)
		return err
	} else {
		// 然后启动资源
		if err1 := resource.Start(); err1 != nil {
			log.Error(err1)
		}
		// Set resources to inend
		in.Resource = resource
		testResourceState(resource, e, in.Id)
		go func(ctx context.Context) {
			// 5 seconds
			ticker := time.NewTicker(time.Duration(time.Second * 5))
			defer resource.Stop()
			for {
				<-ticker.C
				// log.Debug("Test state...", resource.Details().Id)
				if resource.Status() == DOWN {
					testResourceState(resource, e, in.Id)
				}
			}
		}(context.Background())
		return nil
	}
}

//
// test ResourceState
//
func testResourceState(resource XResource, e *RuleEngine, id string) {
	if resource.Status() == UP {
		e.GetInEnd(id).SetState(UP)
	} else {
		e.GetInEnd(id).SetState(DOWN)
		// 当资源挂了以后先给停止，然后重启
		log.Warnf("resource %v down. try to restart it", resource.Details().Id)
		resource.Stop()
		runtime.Gosched()
		runtime.GC()
		resource.Start()
	}
}

//
//
// LoadOutEnd
func (e *RuleEngine) LoadOutEnd(out *outEnd) error {
	return tryCreateOutEnd(out, e)
}

//
//
//
func tryCreateOutEnd(out *outEnd, e RuleX) error {
	if out.Type == "mongo" {
		return startTarget(NewMongoTarget(e), out, e)
	}
	return errors.New("unsupported target type:" + out.Type)

}

//
// Start output target
//
// Target life cycle:
// Register -> Start -> Test
//
func startTarget(target XTarget, out *outEnd, e RuleX) error {
	// Important!!! Must save outend first
	e.SaveOutEnd(out)
	out.Target = target
	// Register outend to target
	if err0 := target.Register(out.Id); err0 != nil {
		return err0
	} else {
		if err1 := target.Start(); err1 != nil {
			log.Error(err1)
		}
		testTargetState(target, e, out.Id)
		//
		go func(ctx context.Context) {
			// 5 seconds
			ticker := time.NewTicker(time.Duration(time.Second * 5))
			defer target.Stop()
			for {
				<-ticker.C
				if target.Status() == DOWN {
					testTargetState(target, e, out.Id)
				}
			}
		}(context.Background())
		return nil
	}
}

// Test Target State
func testTargetState(target XTarget, e RuleX, id string) {
	if target.Status() == UP {
		e.GetOutEnd(id).State = UP
	} else {
		e.GetOutEnd(id).State = DOWN
		// 当资源挂了以后先给停止，然后重启
		log.Warnf("Target %v down. try to restart it", target.Details())
		target.Stop()
		runtime.Gosched()
		runtime.GC()
		target.Start()
	}
}

// LoadRule
func (e *RuleEngine) LoadRule(r *rule) error {
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
					e.SaveRule(r)
					return nil
				} else {
					return errors.New("InEnd:" + inId + " is not exists")
				}
			}
		}
	}
	return errors.New("from can not be empty")

}

//
// Remove a rule
//
func (e *RuleEngine) GetRule(id string) *rule {
	e.Lock()
	defer e.Unlock()
	return (*e.Rules)[id]
}

//
//
//
func (e *RuleEngine) SaveRule(r *rule) {
	e.Lock()
	defer e.Unlock()
	(*e.Rules)[r.Id] = r

}

//
// RemoveRule and inend--rule bindings
//
func (e *RuleEngine) RemoveRule(ruleId string) error {
	e.Lock()
	defer e.Unlock()
	if rule := e.GetRule(ruleId); rule != nil {
		for _, inEnd := range *e.InEnds {
			for _, rule := range *inEnd.Binds {
				if rule.Id == ruleId {
					delete(*inEnd.Binds, ruleId)
				}
			}
		}
		delete((*e.Rules), ruleId)
		return nil
	} else {
		return errors.New("rule:" + ruleId + " not exists")
	}
}

//
//
//
func (e *RuleEngine) AllRule() map[string]*rule {
	e.Lock()
	defer e.Unlock()
	return (*e.Rules)
}

//
// Stop
//
func (e *RuleEngine) Stop() {
	log.Debug("Stop rulex")
	for _, inEnd := range *e.InEnds {
		if inEnd.Resource != nil {
			inEnd.Resource.Stop()
		}
	}

	for _, outEnds := range *e.OutEnds {
		if outEnds.Target != nil {
			outEnds.Target.Stop()
		}
	}
	runtime.Gosched()
	runtime.GC()
}

// Work
func (e *RuleEngine) Work(in *inEnd, data string) (bool, error) {
	statistics.IncIn()
	//
	// Run Lua
	//
	e.runLuaCallbacks(in, data)
	//
	// Run Hook
	//
	e.runHooks(data)
	return false, nil
}
func (e *RuleEngine) runLuaCallbacks(in *inEnd, data string) {
	for _, rule := range *in.Binds {
		_, err := rule.ExecuteActions(lua.LString(data))
		if err != nil {
			log.Error(err)
			rule.ExecuteFailed(lua.LString(err.Error()))
		} else {
			rule.ExecuteSuccess()
		}
	}
}

// Verify Lua Syntax
func VerifyCallback(r *rule) error {
	vm := r.VM
	e1 := vm.DoString(r.Success)
	if e1 != nil {
		return e1
	}
	if vm.GetGlobal("Success").Type() != lua.LTFunction {
		return errors.New("success not submit")
	}
	e2 := vm.DoString(r.Failed)
	if e2 != nil {
		return e2
	}
	if vm.GetGlobal("Failed").Type() != lua.LTFunction {
		return errors.New("failed not submit")
	}
	e3 := vm.DoString(r.Actions)
	if e3 != nil {
		return e3
	}
	// validate Syntax
	actionsTable := vm.GetGlobal("Actions")
	if actionsTable != nil && actionsTable.Type() == lua.LTTable {
		valid := false
		actionsTable.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			valid = (reflect.TypeOf(f).Elem().Name() == "LFunction")
		})
		if !valid {
			return errors.New("invalid function type")
		}
	} else {
		return errors.New("actions must be a functions table")
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
	} else {
		return nil, errors.New("actions not a lua table or not exist")
	}
}

// LUA Callback : Success
func (r *rule) ExecuteSuccess() (interface{}, error) {
	return execute(r.VM, "Success")
}

// LUA Callback : Failed

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
	return nil, errors.New("target:" + k + " is not a lua function")
}

// callLuaFunc
func callLuaFunc(vm *lua.LState, callable *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	if callable == nil {
		return nil, errors.New("callable function is not exists")
	} else {
		coroutine, _ := vm.NewThread()
		state, err, lValues := vm.Resume(coroutine, callable, args...)
		if state != lua.ResumeOK {
			return nil, err
		} else {
			return lValues, nil
		}
	}
}

//
func (e *RuleEngine) LoadPlugin(p XPlugin) error {
	env := p.Load()
	err0 := p.Init(env)
	if err0 != nil {
		return err0
	} else {
		metaInfo, err1 := p.Install(env)
		if err1 != nil {
			return err1
		} else {
			if (*e.Plugins)[metaInfo.Name] != nil {
				return errors.New("plugin already installed:" + metaInfo.Name)
			} else {
				(*e.Plugins)[metaInfo.Name] = metaInfo
				if err2 := p.Start(env); err2 != nil {
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
	e.Lock()
	defer e.Unlock()
	return (*e.InEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveInEnd(in *inEnd) {
	e.Lock()
	defer e.Unlock()
	(*e.InEnds)[in.Id] = in
}

//
//
//
func (e *RuleEngine) RemoveInEnd(id string) {
	e.Lock()
	defer e.Unlock()
	delete((*e.InEnds), id)
}

//
//
//
func (e *RuleEngine) AllInEnd() map[string]*inEnd {
	e.Lock()
	defer e.Unlock()
	return (*e.InEnds)
}

//
//
//
func (e *RuleEngine) GetOutEnd(id string) *outEnd {
	e.Lock()
	defer e.Unlock()
	return (*e.OutEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveOutEnd(out *outEnd) {
	e.Lock()
	defer e.Unlock()
	(*e.OutEnds)[out.Id] = out
}

//
//
//
func (e *RuleEngine) RemoveOutEnd(out *outEnd) {
	e.Lock()
	defer e.Unlock()
	delete((*e.OutEnds), out.Id)
}

//
//
//
func (e *RuleEngine) AllOutEnd() map[string]*outEnd {
	e.Lock()
	defer e.Unlock()
	return (*e.OutEnds)
}

//
// LoadHook
//F
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
func (e *RuleEngine) runHooks(data string) {
	for _, h := range *e.Hooks {
		if err := runHook(h, data); err != nil {
			h.Error(err)
		}
	}
}
func runHook(h XHook, data string) error {
	return h.Work(data)
}
