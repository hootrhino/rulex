package engine

import (
	"context"
	"errors"
	"fmt"
	"rulex/core"
	"rulex/resource"
	"rulex/statistics"
	"rulex/stdlib"
	"rulex/target"
	"rulex/typex"
	"runtime"
	"sync"
	"time"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

//
//
// RuleEngine
//
type RuleEngine struct {
	sync.Mutex
	Hooks     map[string]typex.XHook           `json:"-"`
	Rules     map[string]*typex.Rule           `json:"-"`
	Plugins   map[string]typex.XPlugin         `json:"plugins"`
	InEnds    map[string]*typex.InEnd          `json:"inends"`
	OutEnds   map[string]*typex.OutEnd         `json:"outends"`
	Drivers   map[string]typex.XExternalDriver `json:"drivers"`
	ConfigMap map[string]interface{}           `json:"configMap"`
}

//
//
//
func NewRuleEngine() typex.RuleX {
	return &RuleEngine{
		Plugins:   map[string]typex.XPlugin{},
		Hooks:     map[string]typex.XHook{},
		Rules:     map[string]*typex.Rule{},
		InEnds:    map[string]*typex.InEnd{},
		OutEnds:   map[string]*typex.OutEnd{},
		Drivers:   map[string]typex.XExternalDriver{},
		ConfigMap: map[string]interface{}{},
	}
}

//
//
//
func (e *RuleEngine) LoadDriver(typex.XExternalDriver) error {

	return nil

}

//
//
//
func (e *RuleEngine) Start() map[string]interface{} {
	e.ConfigMap = map[string]interface{}{}
	log.Info("Init XQueue, max queue size is:", core.GlobalConfig.MaxQueueSize)
	typex.DefaultDataCacheQueue = &typex.DataCacheQueue{
		Queue: make(chan typex.QueueData, core.GlobalConfig.MaxQueueSize),
	}
	go func(ctx context.Context, xQueue typex.XQueue) {
		for {
			select {
			case qd := <-xQueue.GetQueue():
				{
					//
					// 消息队列有2种用法:
					// 1 进来的数据缓存
					// 2 出去的消息缓存
					// 只需要判断 in 或者 out 是不是 nil即可
					//
					if qd.In != nil {
						//
						// 传 In 进来为了查找 Rule 去执行
						// 内存中的数据关联性 Key: In, Value: Rules...
						//
						qd.E.RunLuaCallbacks(qd.In, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.Out != nil {
						//  传 Out 为了实现数据外流
						(*qd.E.AllOutEnd()[qd.Out.Id]).Target.To(qd.Data)
					}
				}
			default:
				{
				}
			}
		}
	}(context.Background(), typex.DefaultDataCacheQueue)
	return e.ConfigMap
}

//
//
//
func (e *RuleEngine) PushQueue(qd typex.QueueData) error {
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		log.Error("ATK_LORA_01Driver error: ", err)
	}
	return err
}

//
//
//
func (e *RuleEngine) GetPlugins() map[string]typex.XPlugin {
	e.Lock()
	defer e.Unlock()
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() map[string]typex.XPlugin {
	e.Lock()
	defer e.Unlock()
	return e.Plugins
}

//
//
//
func (e *RuleEngine) Version() typex.Version {
	return typex.DefaultVersion
}

//
//
func (e *RuleEngine) GetConfig(k string) interface{} {
	return (e.ConfigMap)[k]
}

func (e *RuleEngine) LoadInEnd(in *typex.InEnd) error {
	if in.Type == typex.MQTT {
		return startResources(resource.NewMqttInEndResource(in.Id, e), in, e)
	}
	if in.Type == typex.HTTP {
		return startResources(resource.NewHttpInEndResource(in.Id, e), in, e)
	}
	if in.Type == typex.COAP {
		return startResources(resource.NewCoAPInEndResource(in.Id, e), in, e)
	}
	if in.Type == typex.GRPC {
		return startResources(resource.NewGrpcInEndResource(in.Id, e), in, e)
	}
	if in.Type == typex.UART_MODULE {
		return startResources(resource.NewUartModuleResource(in.Id, e), in, e)
	}
	if in.Type == typex.UDP {
		return startResources(resource.NewUdpInEndResource(in.Id, e), in, e)
	}
	if in.Type == typex.MODBUS_TCP_MASTER {
		return startResources(resource.NewModbusMasterResource(in.Id, e), in, e)
	}
	if in.Type == typex.SNMP_SERVER {
		return startResources(resource.NewSNMPInEndResource(in.Id, e), in, e)
	}
	return fmt.Errorf("unsupported InEnd type:%s", in.Type)
}

//
// start Resources
//
/*
* Life cycle
+------------------+       +------------------+   +---------------+        +---------------+
|     Register     |------>|   Start          |-->|     Test      |--+ --->|  Stop         |
+------------------+  |    +------------------+   +---------------+  |     +---------------+
                      |                                              |
                      |                                              |
                      +-------------------Error ---------------------+
*/
func startResources(resource typex.XResource, in *typex.InEnd, e *RuleEngine) error {
	// Save to rule engine first
	// 这么作主要是为了可以 预加载 进去, 然后等环境恢复了以后自动复原
	e.SaveInEnd(in)
	// 首先把资源ID给注册进去, 作为资源的全局索引
	if err := resource.Register(in.Id); err != nil {
		log.Error(err)
		return err
	} else {
		// 然后启动资源
		if err1 := resource.Start(); err1 != nil {
			log.Error("Resource start failed:", resource.Details().Id, ", errors:", err1)
		}
		// Set resources to inend
		in.Resource = resource
		testResourceState(resource, e, in.Id)
		testDriverState(resource, e, in.Id)
		go func(ctx context.Context) {
			// 5 seconds
			ticker := time.NewTicker(time.Duration(time.Second * 5))
			defer resource.Stop()
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{
					}
				}
				//------------------------------------
				// 驱动挂了资源也挂了，因此检查驱动状态在先
				//------------------------------------
				testResourceState(resource, e, in.Id)
				testDriverState(resource, e, in.Id)
				//------------------------------------

			}

		}(context.Background())
		log.Infof("InEnd %v %v load successfully", in.Name, in.Id)
		return nil
	}
}

//
// test ResourceState
//
func testResourceState(resource typex.XResource, e *RuleEngine, id string) {
	if resource.Status() != typex.UP {
		e.GetInEnd(id).SetState(typex.DOWN)
		//----------------------------------
		// 当资源挂了以后先给停止, 然后重启
		//----------------------------------
		log.Warnf("Resource %v %v down. try to restart it", resource.Details().Id, resource.Details().Name)
		resource.Stop()
		//----------------------------------
		// 驱动也要停了
		//----------------------------------
		if resource.Driver() != nil {
			resource.Driver().Stop()
		}
		//----------------------------------
		// 主动垃圾回收一波
		//----------------------------------
		runtime.Gosched()
		runtime.GC()
		resource.Start()
	}
}
func testDriverState(resource typex.XResource, e *RuleEngine, id string) {
	if resource.Driver() != nil {
		// println("testDriverState:", resource.Driver().State())
		if resource.Status() == typex.UP {
			if resource.Driver().State() == typex.STOP {
				log.Warn("Driver stoped:", resource.Driver().DriverDetail().Name)
				e.GetInEnd(id).SetState(typex.DOWN)
				// Start driver
				if err := resource.Driver().Init(); err != nil {
					log.Error("Driver initial error:", err)
				} else {
					log.Info("Try to start driver: ", resource.Driver().DriverDetail().Name)
					if err := resource.Driver().Work(); err != nil {
						log.Error("Driver initial error:", err)
					} else {
						log.Info("Driver start successfully:", resource.Driver().DriverDetail().Name)
					}
				}
			}
		}

	}
}

//
// LoadOutEnd
//
func (e *RuleEngine) LoadOutEnd(out *typex.OutEnd) error {
	return tryCreateOutEnd(out, e)
}

//
// CreateOutEnd
//
func tryCreateOutEnd(out *typex.OutEnd, e typex.RuleX) error {
	if out.Type == typex.MONGO_SINGLE {
		return startTarget(target.NewMongoTarget(e), out, e)
	}
	if out.Type == typex.MQTT_TARGET {
		return startTarget(target.NewMqttTarget(e), out, e)
	}
	return errors.New("unsupported target type:" + out.Type.String())

}

//
// Start output target
//
// Target life cycle:
// Register -> Start -> Test
//
func startTarget(target typex.XTarget, out *typex.OutEnd, e typex.RuleX) error {
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
			for {
				select {
				case <-ticker.C:
					if target.Status() == typex.DOWN {
						testTargetState(target, e, out.Id)
					}
				default:
					{
					}
				}
			}
		}(context.Background())
		return nil
	}
}

// Test Target State
func testTargetState(target typex.XTarget, e typex.RuleX, id string) {
	if target.Status() == typex.UP {
		e.GetOutEnd(id).State = typex.UP
	} else {
		e.GetOutEnd(id).State = typex.DOWN
		// 当资源挂了以后先给停止, 然后重启
		log.Warnf("Target %v down. try to restart it", target.Details())
		target.Stop()
		runtime.Gosched()
		runtime.GC()
		target.Start()
	}
}

// LoadRule
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	if err := core.VerifyCallback(r); err != nil {
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
					//
					// Load Stdlib
					//--------------------------------------------------------------
					r.LoadLib(e, stdlib.NewBinaryLib())
					r.LoadLib(e, stdlib.NewMongoLib())
					r.LoadLib(e, stdlib.NewHttpLib())
					r.LoadLib(e, stdlib.NewMqttLib())
					r.LoadLib(e, stdlib.NewJqLib())
					r.LoadLib(e, stdlib.NewWriteInStreamLib())
					r.LoadLib(e, stdlib.NewWriteOutStreamLib())
					r.LoadLib(e, stdlib.NewByteToBitStringLib())
					r.LoadLib(e, stdlib.NewGetABitOnByteLib())
					//--------------------------------------------------------------
					// Save to rules map
					//
					e.SaveRule(r)
					log.Infof("Rule %v %v load successfully", r.Name, r.Id)
					return nil
				} else {
					return errors.New("'InEnd':" + inId + " is not exists")
				}
			}
		}
	}
	return errors.New("'From' can not be empty")

}

//
// Remove a rule
//
func (e *RuleEngine) GetRule(id string) *typex.Rule {
	e.Lock()
	defer e.Unlock()
	return (e.Rules)[id]
}

//
//
//
func (e *RuleEngine) SaveRule(r *typex.Rule) {
	e.Lock()
	defer e.Unlock()
	(e.Rules)[r.Id] = r

}

//
// RemoveRule and inend--rule bindings
//
func (e *RuleEngine) RemoveRule(ruleId string) error {
	e.Lock()
	defer e.Unlock()
	if rule := e.GetRule(ruleId); rule != nil {
		for _, inEnd := range e.InEnds {
			for _, rule := range *inEnd.Binds {
				if rule.Id == ruleId {
					delete(*inEnd.Binds, ruleId)
				}
			}
		}
		delete((e.Rules), ruleId)
		return nil
	} else {
		return errors.New("'Rule':" + ruleId + " not exists")
	}
}

//
//
//
func (e *RuleEngine) AllRule() map[string]*typex.Rule {
	e.Lock()
	defer e.Unlock()
	return (e.Rules)
}

//
// Stop
//
func (e *RuleEngine) Stop() {
	log.Info("Stopping Rulex......")
	context.Background().Done()
	for _, inEnd := range e.InEnds {
		if inEnd.Resource != nil {
			log.Info("Stop InEnd:", inEnd.Name, inEnd.Id)
			inEnd.Resource.Stop()
			if inEnd.Resource.Driver() != nil {
				inEnd.Resource.Driver().Stop()
			}
		}
	}

	for _, outEnd := range e.OutEnds {
		if outEnd.Target != nil {
			log.Info("Stop Target:", outEnd.Name, outEnd.Id)
			outEnd.Target.Stop()
		}
	}

	for _, plugin := range e.Plugins {
		plugin.Uninstall()
		plugin.Clean()

	}
	runtime.Gosched()
	runtime.GC()
	log.Info("Stop Rulex successfully")
}

// Work
func (e *RuleEngine) Work(in *typex.InEnd, data string) (bool, error) {
	statistics.IncIn()
	err := e.PushQueue(typex.QueueData{
		In:   in,
		Out:  nil,
		E:    e,
		Data: data,
	})
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (e *RuleEngine) RunLuaCallbacks(in *typex.InEnd, data string) {
	for _, rule := range *in.Binds {
		_, err := rule.ExecuteActions(lua.LString(data))
		if err != nil {
			log.Error("RunLuaCallbacks error:", err)
			rule.ExecuteFailed(lua.LString(err.Error()))
		} else {
			rule.ExecuteSuccess()
		}
	}
}

//
func (e *RuleEngine) LoadPlugin(p typex.XPlugin) error {
	err0 := p.Init()
	if err0 != nil {
		return err0
	}
	err1 := p.Install()
	if err1 != nil {
		return err1
	}
	if (e.Plugins)[p.XPluginMetaInfo().Name] != nil {
		return errors.New("plugin already installed:" + p.XPluginMetaInfo().Name)
	}
	(e.Plugins)[p.XPluginMetaInfo().Name] = p
	if err2 := p.Start(); err2 != nil {
		return err2
	}
	return nil

}

//
// LoadHook
//
func (e *RuleEngine) LoadHook(h typex.XHook) error {
	if (e.Hooks)[h.Name()] != nil {
		return errors.New("hook have been loaded:" + h.Name())
	} else {
		(e.Hooks)[h.Name()] = h
		return nil
	}
}

//
// RunHooks
//
func (e *RuleEngine) RunHooks(data string) {
	for _, h := range e.Hooks {
		if err := runHook(h, data); err != nil {
			h.Error(err)
		}
	}
}
func runHook(h typex.XHook, data string) error {
	return h.Work(data)
}

//
//
//
func (e *RuleEngine) GetInEnd(id string) *typex.InEnd {
	e.Lock()
	defer e.Unlock()
	return (e.InEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveInEnd(in *typex.InEnd) {
	e.Lock()
	defer e.Unlock()
	(e.InEnds)[in.Id] = in
}

//
//
//
func (e *RuleEngine) RemoveInEnd(id string) {
	e.Lock()
	defer e.Unlock()
	delete((e.InEnds), id)
}

//
//
//
func (e *RuleEngine) AllInEnd() map[string]*typex.InEnd {
	e.Lock()
	defer e.Unlock()
	return (e.InEnds)
}

//
//
//
func (e *RuleEngine) GetOutEnd(id string) *typex.OutEnd {
	e.Lock()
	defer e.Unlock()
	return (e.OutEnds)[id]
}

//
//
//
func (e *RuleEngine) SaveOutEnd(out *typex.OutEnd) {
	e.Lock()
	defer e.Unlock()
	(e.OutEnds)[out.Id] = out
}

//
//
//
func (e *RuleEngine) RemoveOutEnd(out *typex.OutEnd) {
	e.Lock()
	defer e.Unlock()
	delete((e.OutEnds), out.Id)
}

//
//
//
func (e *RuleEngine) AllOutEnd() map[string]*typex.OutEnd {
	e.Lock()
	defer e.Unlock()
	return (e.OutEnds)
}
