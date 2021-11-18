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
	Hooks     sync.Map `json:"-"`
	Rules     sync.Map `json:"-"`
	Plugins   sync.Map `json:"plugins"`
	InEnds    sync.Map `json:"inends"`
	OutEnds   sync.Map `json:"outends"`
	Drivers   sync.Map `json:"drivers"`
	ConfigMap sync.Map `json:"configMap"`
}

//
//
//
func NewRuleEngine() typex.RuleX {
	return &RuleEngine{
		Plugins:   sync.Map{},
		Hooks:     sync.Map{},
		Rules:     sync.Map{},
		InEnds:    sync.Map{},
		OutEnds:   sync.Map{},
		Drivers:   sync.Map{},
		ConfigMap: sync.Map{},
	}
}

//
//
//
func (e *RuleEngine) Start() sync.Map {
	e.ConfigMap = sync.Map{}
	log.Info("Init XQueue, max queue size is:", core.GlobalConfig.MaxQueueSize)
	typex.DefaultDataCacheQueue = &typex.DataCacheQueue{
		Queue: make(chan typex.QueueData, core.GlobalConfig.MaxQueueSize),
	}
	go func(ctx context.Context, xQueue typex.XQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			case qd := <-xQueue.GetQueue():
				{
					//
					// 消息队列有2种用法:
					// 1 进来的数据缓存
					// 2 出去的消息缓存
					// 只需要判断 in 或者 out 是不是 nil即可
					//
					if qd.In != nil {
						qd.E.RunLuaCallbacks(qd.In, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.Out != nil {
						outEnds := qd.E.AllOutEnd()
						v, ok := outEnds.Load(qd.Out.UUID)
						if ok {
							v.(*typex.OutEnd).Target.To(qd.Data)
						}
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
func (e *RuleEngine) GetPlugins() sync.Map {
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() sync.Map {
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
	v, ok := e.ConfigMap.Load(k)
	if ok {
		return v
	} else {
		return map[string]interface{}{}
	}
}

func (e *RuleEngine) LoadInEnd(in *typex.InEnd) error {
	if in.Type == typex.MQTT {
		return startResources(resource.NewMqttInEndResource(in.UUID, e), in, e)
	}
	if in.Type == typex.HTTP {
		return startResources(resource.NewHttpInEndResource(in.UUID, e), in, e)
	}
	if in.Type == typex.COAP {
		return startResources(resource.NewCoAPInEndResource(in.UUID, e), in, e)
	}
	if in.Type == typex.GRPC {
		return startResources(resource.NewGrpcInEndResource(in.UUID, e), in, e)
	}
	if in.Type == typex.UART_MODULE {
		return startResources(resource.NewUartModuleResource(in.UUID, e), in, e)
	}
	if in.Type == typex.UDP {
		return startResources(resource.NewUdpInEndResource(in.UUID, e), in, e)
	}
	if in.Type == typex.MODBUS_TCP_MASTER {
		return startResources(resource.NewModbusMasterResource(in.UUID, e), in, e)
	}
	if in.Type == typex.SNMP_SERVER {
		return startResources(resource.NewSNMPInEndResource(in.UUID, e), in, e)
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
	//
	// 先注册，如果出问题了直接删除就行
	//
	e.SaveInEnd(in)
	if !resource.Test(in.UUID) {
		e.RemoveInEnd(in.UUID)
		return errors.New("resource test error")
	}
	// 首先把资源ID给注册进去, 作为资源的全局索引
	if err := resource.Register(in.UUID); err != nil {
		log.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	// Set resources to inend
	in.Resource = resource
	// 然后启动资源
	if err := startResource(resource, e, in.UUID); err != nil {
		log.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	go func(ctx context.Context) {
		// 5 seconds
		ticker := time.NewTicker(time.Duration(time.Second * 5))
		for {
			//
			<-ticker.C
			{
				//
				// 通过HTTP删除资源的时候，会把数据清了，只要检测到资源没了，这里也退出
				//
				if resource.Details() == nil {
					return
				}
				//------------------------------------
				// 驱动挂了资源也挂了，因此检查驱动状态在先
				//------------------------------------
				tryIfRestartResource(resource, e, in.UUID)
				// checkDriverState(resource, e, in.UUID)
				//------------------------------------
			}
		}

	}(context.Background())
	log.Infof("InEnd [%v, %v] load successfully", in.Name, in.UUID)
	return nil
}

/*
*
* 检查是否需要重新拉起资源
* 这里也有优化点：不能手动控制内存回收可能会产生垃圾
*
 */
func checkDriverState(resource typex.XResource, e *RuleEngine, id string) {
	if resource.Driver() != nil {
		// 只有资源启动状态才拉起驱动
		if resource.Status() == typex.UP {
			// 必须资源启动, 驱动才有重启意义
			if resource.Driver().State() == typex.STOP {
				log.Warn("Driver stopped:", resource.Driver().DriverDetail().Name)
				// 只需要把资源给拉闸, 就会触发重启
				resource.Stop()
			}
		}
	}

}

//
// test ResourceState
//
func tryIfRestartResource(resource typex.XResource, e *RuleEngine, id string) {
	checkDriverState(resource, e, id)
	if resource.Status() == typex.DOWN {
		resource.Details().SetState(typex.DOWN)
		//----------------------------------
		// 当资源挂了以后先给停止, 然后重启
		//----------------------------------
		log.Warnf("Resource %v %v down. try to restart it", resource.Details().UUID, resource.Details().Name)
		resource.Stop()
		//----------------------------------
		// 主动垃圾回收一波
		//----------------------------------
		runtime.Gosched()
		runtime.GC() // GC 比较慢, 但是是良性卡顿, 问题不大
		startResource(resource, e, id)
	} else {
		resource.Details().SetState(typex.UP)
	}
}

//
//
//
func startResource(resource typex.XResource, e *RuleEngine, id string) error {
	if err := resource.Start(); err != nil {
		log.Error("Resource start error:", err)
		if resource.Status() == typex.UP {
			resource.Stop()
		}
		if resource.Driver() != nil {
			if resource.Driver().State() == typex.RUNNING {
				resource.Driver().Stop()
			}
		}
		return err
	} else {
		//----------------------------------
		// 驱动也要停了
		//----------------------------------
		if resource.Driver() != nil {
			if resource.Driver().State() == typex.RUNNING {
				resource.Driver().Stop()
			}
			// Start driver
			if err := resource.Driver().Init(); err != nil {
				log.Error("Driver initial error:", err)
				return errors.New("Driver initial error:" + err.Error())
			} else {
				log.Infof("Try to start driver: [%v]", resource.Driver().DriverDetail().Name)
				if err := resource.Driver().Work(); err != nil {
					log.Error("Driver work error:", err)
					return errors.New("Driver work error:" + err.Error())
				} else {
					log.Infof("Driver start successfully: [%v]", resource.Driver().DriverDetail().Name)
					return nil
				}
			}
		}
		return nil
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
	//
	// 先注册，如果出问题了直接删除就行
	//
	e.SaveOutEnd(out)
	// 首先把资源ID给注册进去, 作为资源的全局索引
	if err := target.Register(out.UUID); err != nil {
		log.Error(err)
		e.RemoveOutEnd(out.UUID)
		return err
	}
	// 然后启动资源
	if err := target.Start(); err != nil {
		log.Error(err)
		e.RemoveOutEnd(out.UUID)
		return err
	}
	// Set resources to inend
	out.Target = target
	go func(ctx context.Context) {
		// 5 seconds
		ticker := time.NewTicker(time.Duration(time.Second * 5))
		for {
			//
			<-ticker.C
			{
				if target.Details() == nil {
					return
				}
				tryIfRestartTarget(target, e, out.UUID)

			}
		}

	}(context.Background())
	log.Infof("Target [%v, %v] load successfully", out.Name, out.UUID)
	return nil
}

//
// 监测状态，如果挂了重启
//
func tryIfRestartTarget(target typex.XTarget, e typex.RuleX, id string) {
	if target.Status() == typex.DOWN {
		target.Details().SetState(typex.DOWN)
		log.Warnf("Target [%v, %v] down. try to restart it", target.Details().Name, target.Details().UUID)
		target.Stop()
		runtime.Gosched()
		runtime.GC()
		target.Start()
	} else {
		target.Details().SetState(typex.UP)
	}
}

// LoadRule
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	if err := core.VerifyCallback(r); err != nil {
		return err
	} else {
		if len(r.From) > 0 {
			for _, inUUId := range r.From {
				if in := e.GetInEnd(inUUId); in != nil {
					// Bind to rule, Key:RuleId, Value: Rule
					// RULE_0f8619ef-3cf2-452f-8dd7-aa1db4ecfdde {
					// ...
					// ...
					// }
					(in.Binds)[r.UUID] = *r
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
					r.LoadLib(e, stdlib.NewByteToInt64Lib())
					r.LoadLib(e, stdlib.NewBitStringToBytesLib())
					//--------------------------------------------------------------
					// Save to rules map
					//
					e.SaveRule(r)
					log.Infof("Rule [%v, %v] load successfully", r.Name, r.UUID)
					return nil
				} else {
					return errors.New("'InEnd':" + inUUId + " is not exists")
				}
			}
		}
	}
	return errors.New("'From' can not be empty")

}

//
// GetRule a rule
//
func (e *RuleEngine) GetRule(id string) *typex.Rule {
	v, ok := (e.Rules).Load(id)
	if ok {
		return v.(*typex.Rule)
	} else {
		return nil
	}
}

//
//
//
func (e *RuleEngine) SaveRule(r *typex.Rule) {
	e.Rules.Store(r.UUID, r)
}

//
// RemoveRule and inend--rule bindings
//
func (e *RuleEngine) RemoveRule(ruleId string) {
	if rule := e.GetRule(ruleId); rule != nil {

		// 清空 InEnd 的 bind
		inEnds := e.AllInEnd()
		inEnds.Range(func(key, value interface{}) bool {
			inEnd := value.(*typex.InEnd)
			for _, r := range inEnd.Binds {
				if rule.UUID == r.UUID {
					delete(inEnd.Binds, ruleId)
				}
			}
			return true
		})
		e.Rules.Delete(ruleId)
		rule = nil
		log.Infof("Rule [%v] has been deleted", ruleId)
	}
}

//
//
//
func (e *RuleEngine) AllRule() sync.Map {
	return e.Rules
}

//
// Stop
//
func (e *RuleEngine) Stop() {
	log.Info("Ready to stop rulex")
	e.InEnds.Range(func(key, value interface{}) bool {
		inEnd := value.(*typex.InEnd)
		if inEnd.Resource != nil {
			log.Info("Stop InEnd:", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).SetState(typex.DOWN)
			inEnd.Resource.Stop()
			if inEnd.Resource.Driver() != nil {
				inEnd.Resource.Driver().SetState(typex.STOP)
				inEnd.Resource.Driver().Stop()
			}
		}
		return true
	})
	e.OutEnds.Range(func(key, value interface{}) bool {
		outEnd := value.(*typex.OutEnd)
		if outEnd.Target != nil {
			log.Info("Stop Target:", outEnd.Name, outEnd.UUID)
			outEnd.Target.Stop()
		}
		return true
	})

	e.Plugins.Range(func(key, value interface{}) bool {
		plugin := value.(typex.XPlugin)
		log.Info("Stop plugin:", plugin.XPluginMetaInfo().Name)
		plugin.Stop()
		return true
	})

	context.Background().Done()
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

//
// 执行lua脚本
//
func (e *RuleEngine) RunLuaCallbacks(in *typex.InEnd, data string) {
	for _, rule := range in.Binds {
		if rule.Status == typex.RULE_RUNNING {
			_, err := rule.ExecuteActions(lua.LString(data))
			if err != nil {
				log.Error("RunLuaCallbacks error:", err)
				rule.ExecuteFailed(lua.LString(err.Error()))
			} else {
				rule.ExecuteSuccess()
			}
		}
	}
}

//
func (e *RuleEngine) LoadPlugin(p typex.XPlugin) error {
	if err := p.Init(); err != nil {
		return err
	}

	_, ok := e.Plugins.Load(p.XPluginMetaInfo().Name)
	if ok {
		return errors.New("plugin already installed:" + p.XPluginMetaInfo().Name)
	}

	if err := p.Start(); err != nil {
		return err
	}

	e.Plugins.Store(p.XPluginMetaInfo().Name, p)
	log.Infof("Plugin start successfully:[%v]", p.XPluginMetaInfo().Name)
	return nil

}

//
// LoadHook
//
func (e *RuleEngine) LoadHook(h typex.XHook) error {
	value, _ := e.Hooks.Load(h.Name())
	if value != nil {
		return errors.New("hook have been loaded:" + h.Name())
	} else {
		e.Hooks.Store(h.Name(), h)
		return nil
	}
}

//
// RunHooks
//
func (e *RuleEngine) RunHooks(data string) {
	e.Hooks.Range(func(key, value interface{}) bool {
		if err := runHook(value.(typex.XHook), data); err != nil {
			value.(typex.XHook).Error(err)
		}
		return true
	})
}
func runHook(h typex.XHook, data string) error {
	return h.Work(data)
}

//
//
//
func (e *RuleEngine) GetInEnd(uuid string) *typex.InEnd {
	v, ok := (e.InEnds).Load(uuid)
	if ok {
		return v.(*typex.InEnd)
	} else {
		return nil
	}
}

//
//
//
func (e *RuleEngine) SaveInEnd(in *typex.InEnd) {
	e.InEnds.Store(in.UUID, in)
}

//
//
//
func (e *RuleEngine) RemoveInEnd(id string) {
	if inEnd := e.GetInEnd(id); inEnd != nil {
		e.InEnds.Delete(id)
		inEnd = nil
		log.Infof("InEnd [%v] has been deleted", id)
	}
}

//
//
//
func (e *RuleEngine) AllInEnd() sync.Map {
	return e.InEnds
}

//
//
//
func (e *RuleEngine) GetOutEnd(id string) *typex.OutEnd {
	v, ok := e.OutEnds.Load(id)
	if ok {
		return v.(*typex.OutEnd)
	} else {
		return nil
	}

}

//
//
//
func (e *RuleEngine) SaveOutEnd(out *typex.OutEnd) {
	e.OutEnds.Store(out.UUID, out)

}

//
//
//
func (e *RuleEngine) RemoveOutEnd(uuid string) {
	if inEnd := e.GetOutEnd(uuid); inEnd != nil {
		e.OutEnds.Delete(uuid)
		inEnd = nil
		log.Infof("InEnd [%v] has been deleted", uuid)
	}
}

//
//
//
func (e *RuleEngine) AllOutEnd() sync.Map {
	return e.OutEnds
}
