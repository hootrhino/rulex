package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rulex/core"
	"rulex/source"
	"rulex/statistics"
	"rulex/target"
	"rulex/typex"
	"rulex/utils"
	"runtime"
	"sync"
	"time"

	"github.com/ngaut/log"
	"github.com/shirou/gopsutil/disk"
	lua "github.com/yuin/gopher-lua"
)

//
//
//
type RuleEngine struct {
	Hooks   *sync.Map          `json:"hooks"`
	Rules   *sync.Map          `json:"rules"`
	Plugins *sync.Map          `json:"plugins"`
	InEnds  *sync.Map          `json:"inends"`
	OutEnds *sync.Map          `json:"outends"`
	Drivers *sync.Map          `json:"drivers"`
	Devices *sync.Map          `json:"devices"`
	Config  *typex.RulexConfig `json:"config"`
}

//
//
//
func NewRuleEngine(config typex.RulexConfig) typex.RuleX {
	return &RuleEngine{
		Plugins: &sync.Map{},
		Hooks:   &sync.Map{},
		Rules:   &sync.Map{},
		InEnds:  &sync.Map{},
		OutEnds: &sync.Map{},
		Drivers: &sync.Map{},
		Devices: &sync.Map{},
		Config:  &config,
	}
}

//
//
//
func (e *RuleEngine) Start() *typex.RulexConfig {
	typex.StartQueue(core.GlobalConfig.MaxQueueSize)
	source.LoadSt()
	target.LoadTt()
	return e.Config
}

//
//
//
func (e *RuleEngine) PushQueue(qd typex.QueueData) error {
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		log.Error("PushQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}
func (e *RuleEngine) PushInQueue(in *typex.InEnd, data string) error {
	qd := typex.QueueData{
		E:    e,
		In:   in,
		Out:  nil,
		Data: data,
	}
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		log.Error("PushInQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}
func (e *RuleEngine) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := typex.QueueData{
		E:    e,
		In:   nil,
		Out:  out,
		Data: data,
	}
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		log.Error("PushOutQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}

//
//
//
func (e *RuleEngine) GetPlugins() *sync.Map {
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() *sync.Map {
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
func (e *RuleEngine) GetConfig() *typex.RulexConfig {
	return e.Config
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
	if out.Type == typex.NATS_TARGET {
		return startTarget(target.NewNatsTarget(e), out, e)
	}
	if out.Type == typex.HTTP_TARGET {
		return startTarget(target.NewHTTPTarget(e), out, e)
	}
	if out.Type == typex.TDENGINE_TARGET {
		return startTarget(target.NewTdEngineTarget(e), out, e)
	}
	if out.Type == typex.GRPC_CODEC_TARGET {
		return startTarget(target.NewCodecTarget(e), out, e)
	}
	return errors.New("unsupported target type:" + out.Type.String())

}

//
// Start output target
//
// Target life cycle:
//     Register -> Start -> running/restart cycle
//
func startTarget(target typex.XTarget, out *typex.OutEnd, e typex.RuleX) error {
	//
	// 先注册, 如果出问题了直接删除就行
	//
	e.SaveOutEnd(out)

	// Load config
	config := e.GetOutEnd(out.UUID).Config
	if config == nil {
		e.RemoveOutEnd(out.UUID)
		err := fmt.Errorf("target [%v] config is nil", out.Name)
		return err
	}
	if err := target.Init(out.UUID, config); err != nil {
		log.Error(err)
		e.RemoveInEnd(out.UUID)
		return err
	}
	// 然后启动资源
	ctx, cancelCTX := context.WithCancel(typex.GCTX)
	if err := target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		log.Error(err)
		e.RemoveOutEnd(out.UUID)
		return err
	}
	// Set sources to inend
	out.Target = target
	//
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	tryIfRestartTarget(target, e, out.UUID)
	go func(ctx context.Context) {

		// 5 seconds
		//
	TICKER:
		<-ticker.C
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
				goto CHECK
			}
		}
	CHECK:
		{
			if target.Details() == nil {
				return
			}
			tryIfRestartTarget(target, e, out.UUID)
			goto TICKER
		}

	}(typex.GCTX)
	log.Infof("Target [%v, %v] load successfully", out.Name, out.UUID)
	return nil
}

//
// 监测状态, 如果挂了重启
//
func tryIfRestartTarget(target typex.XTarget, e typex.RuleX, id string) {
	if target.Status() == typex.DOWN {
		target.Details().State = typex.DOWN
		log.Warnf("Target [%v, %v] down. try to restart it", target.Details().Name, target.Details().UUID)
		target.Stop()
		runtime.Gosched()
		runtime.GC()
		ctx, cancelCTX := context.WithCancel(typex.GCTX)
		target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX})
	} else {
		target.Details().State = typex.UP
	}
}

// LoadRule
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	if err := core.VerifyCallback(r); err != nil {
		return err
	}
	// 判断是否有源ID
	if len(r.From) > 0 {
		for _, inUUId := range r.From {
			// 查找输入定义的资源是否存在
			if in := e.GetInEnd(inUUId); in != nil {
				// Bind to rule, Key:RuleId, Value: Rule
				// RULE_0f8619ef-3cf2-452f-8dd7-aa1db4ecfdde {
				// ...
				// ...
				// }
				// 绑定资源和规则，建立关联关系
				(in.Binds)[r.UUID] = *r
				//--------------------------------------------------------------
				// Load LoadBuildInLuaLib
				//--------------------------------------------------------------
				LoadBuildInLuaLib(e, r)
				//--------------------------------------------------------------
				// Save to rules map
				//--------------------------------------------------------------
				e.SaveRule(r)
				log.Infof("Rule [%v, %v] load successfully", r.Name, r.UUID)
				return nil
			} else {
				return errors.New("'InEnd':" + inUUId + " is not exists")
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
		// 清空 InEnd 的 bind 资源
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
func (e *RuleEngine) AllRule() *sync.Map {
	return e.Rules
}

//
// Stop
//
func (e *RuleEngine) Stop() {
	log.Info("Ready to stop rulex")
	e.InEnds.Range(func(key, value interface{}) bool {
		inEnd := value.(*typex.InEnd)
		if inEnd.Source != nil {
			log.Info("Stop InEnd:", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).SetState(typex.DOWN)
			inEnd.Source.Stop()
			if inEnd.Source.Driver() != nil {
				inEnd.Source.Driver().SetState(typex.STOP)
				inEnd.Source.Driver().Stop()
			}
		}
		return true
	})
	// 停止所有外部资源
	e.OutEnds.Range(func(key, value interface{}) bool {
		outEnd := value.(*typex.OutEnd)
		if outEnd.Target != nil {
			log.Info("Stop Target:", outEnd.Name, outEnd.UUID)
			outEnd.Target.Stop()
		}
		return true
	})
	// 停止所有插件
	e.Plugins.Range(func(key, value interface{}) bool {
		plugin := value.(typex.XPlugin)
		log.Info("Stop plugin:", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		return true
	})
	// 停止所有外挂 设备
	e.Devices.Range(func(key, value interface{}) bool {
		Device := value.(*typex.Device)
		log.Info("Stop Device:", Device.Name)
		Device.Device.Stop()
		return true
	})
	// 回收资源
	runtime.Gosched()
	runtime.GC()
	log.Info("Stop Rulex successfully")
}

//
// 核心功能: Work
//
func (e *RuleEngine) Work(in *typex.InEnd, data string) (bool, error) {
	if err := e.PushInQueue(in, data); err != nil {
		return false, err
	}
	return true, nil
}

//
// 执行lua脚本
//
func (e *RuleEngine) RunLuaCallbacks(in *typex.InEnd, callbackArgs string) {
	for _, rule := range in.Binds {
		if rule.Status == typex.RULE_RUNNING {
			_, err := core.ExecuteActions(&rule, lua.LString(callbackArgs))
			if err != nil {
				log.Error("RunLuaCallbacks error:", err)
				_, err := core.ExecuteFailed(rule.VM, lua.LString(err.Error()))
				if err != nil {
					log.Error(err)
				}
			} else {
				_, err := core.ExecuteSuccess(rule.VM)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}
}

//
// ┌──────┐    ┌──────┐    ┌──────┐
// │ Init ├───►│ Load ├───►│ Stop │
// └──────┘    └──────┘    └──────┘
//
func (e *RuleEngine) LoadPlugin(sectionK string, p typex.XPlugin) error {
	section := utils.GetINISection(core.INIPath, sectionK)
	key, err1 := section.GetKey("enable")
	if err1 != nil {
		return err1
	}
	enable, err2 := key.Bool()
	if err2 != nil {
		return err2
	}
	if !enable {
		log.Infof("Plugin is not enable:%s", p.PluginMetaInfo().Name)
		return nil
	}

	if err := p.Init(section); err != nil {
		return err
	}
	_, ok := e.Plugins.Load(p.PluginMetaInfo().Name)
	if ok {
		return errors.New("plugin already installed:" + p.PluginMetaInfo().Name)
	}

	if err := p.Start(); err != nil {
		return err
	}

	e.Plugins.Store(p.PluginMetaInfo().Name, p)
	log.Infof("Plugin start successfully:[%v]", p.PluginMetaInfo().Name)
	return nil

}

//
// LoadHook
//
func (e *RuleEngine) LoadHook(h typex.XHook) error {
	value, _ := e.Hooks.Load(h.Name())
	if value != nil {
		return errors.New("hook have been loaded:" + h.Name())
	}
	e.Hooks.Store(h.Name(), h)
	return nil

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
	}
	return nil
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
		inEnd.Source.Stop()
		e.InEnds.Delete(id)
		inEnd = nil
		log.Infof("InEnd [%v] has been deleted", id)
	}
}

//
//
//
func (e *RuleEngine) AllInEnd() *sync.Map {
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
	if outEnd := e.GetOutEnd(uuid); outEnd != nil {
		outEnd.Target.Stop()
		e.OutEnds.Delete(uuid)
		outEnd = nil
		log.Infof("InEnd [%v] has been deleted", uuid)
	}
}

//
//
//
func (e *RuleEngine) AllOutEnd() *sync.Map {
	return e.OutEnds
}

//-----------------------------------------------------------------
// 获取运行时快照
//-----------------------------------------------------------------
func (e *RuleEngine) SnapshotDump() string {
	inends := []interface{}{}
	rules := []interface{}{}
	plugins := []interface{}{}
	outends := []interface{}{}
	devices := []interface{}{}
	e.AllInEnd().Range(func(key, value interface{}) bool {
		inends = append(inends, value)
		return true
	})
	e.AllRule().Range(func(key, value interface{}) bool {
		rules = append(rules, value)
		return true
	})
	e.AllPlugins().Range(func(key, value interface{}) bool {
		plugins = append(plugins, (value.(typex.XPlugin)).PluginMetaInfo())
		return true
	})
	e.AllOutEnd().Range(func(key, value interface{}) bool {
		outends = append(outends, value)
		return true
	})
	e.AllDevices().Range(func(key, value interface{}) bool {
		devices = append(devices, value)
		return true
	})
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	system := map[string]interface{}{
		"version":  e.Version().Version,
		"diskInfo": int(diskInfo.UsedPercent),
		"system":   bToMb(m.Sys),
		"alloc":    bToMb(m.Alloc),
		"total":    bToMb(m.TotalAlloc),
		"osArch":   runtime.GOOS + "-" + runtime.GOARCH,
	}
	data := map[string]interface{}{
		"rules":      rules,
		"plugins":    plugins,
		"inends":     inends,
		"outends":    outends,
		"devices":    devices,
		"statistics": statistics.AllStatistics(),
		"system":     system,
		"config":     core.GlobalConfig,
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}
	return string(b)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
