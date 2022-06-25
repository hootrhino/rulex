package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/sidecar"
	"github.com/i4de/rulex/source"
	"github.com/i4de/rulex/statistics"
	"github.com/i4de/rulex/target"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/shirou/gopsutil/disk"
	lua "github.com/yuin/gopher-lua"
)

//
// 规则引擎
//
type RuleEngine struct {
	SideCar sidecar.SideCar    `json:"-"`
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
		SideCar: sidecar.NewSideCarManager(typex.GCTX),
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
		glogger.GLogger.Error("PushQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}
func (e *RuleEngine) PushInQueue(in *typex.InEnd, data string) error {
	qd := typex.QueueData{
		E:    e,
		I:    in,
		O:    nil,
		Data: data,
	}
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushInQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}

/*
*
* 设备数据入流引擎
*
 */
func (e *RuleEngine) PushDeviceQueue(Device *typex.Device, data string) error {
	qd := typex.QueueData{
		D:    Device,
		E:    e,
		I:    nil,
		O:    nil,
		Data: data,
	}
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushInQueue error:", err)
		statistics.IncInFailed()
	} else {
		statistics.IncIn()
	}
	return err
}
func (e *RuleEngine) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := typex.QueueData{
		E:    e,
		D:    nil,
		I:    nil,
		O:    out,
		Data: data,
	}
	err := typex.DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushOutQueue error:", err)
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
// LoadRule: 每个规则都绑定了资源(FromSource)或者设备(FromDevice)
//
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	if err := core.VerifyCallback(r); err != nil {
		return err
	}
	e.SaveRule(r)
	//--------------------------------------------------------------
	// Load LoadBuildInLuaLib
	//--------------------------------------------------------------
	LoadBuildInLuaLib(e, r)
	glogger.GLogger.Infof("Rule [%v, %v] load successfully", r.Name, r.UUID)
	// 绑定输入资源
	for _, inUUId := range r.FromSource {
		// 查找输入定义的资源是否存在
		if in := e.GetInEnd(inUUId); in != nil {
			(in.BindRules)[r.UUID] = *r
			return nil
		} else {
			return errors.New("'InEnd':" + inUUId + " is not exists when bind resource")
		}
	}
	// 绑定设备(From 1.0.1)
	for _, devUUId := range r.FromDevice {
		// 查找输入定义的资源是否存在
		if Device := e.GetDevice(devUUId); Device != nil {
			// 绑定资源和规则，建立关联关系
			(Device.BindRules)[r.UUID] = *r
		} else {
			return errors.New("'Device':" + devUUId + " is not exists when bind resource")
		}
	}
	return nil

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
			for _, r := range inEnd.BindRules {
				if rule.UUID == r.UUID {
					delete(inEnd.BindRules, ruleId)
				}
			}
			return true
		})
		// 清空Device的绑定
		Devices := e.AllDevices()
		Devices.Range(func(key, value interface{}) bool {
			Device := value.(*typex.Device)
			for _, r := range Device.BindRules {
				if rule.UUID == r.UUID {
					delete(Device.BindRules, ruleId)
				}
			}
			return true
		})
		e.Rules.Delete(ruleId)
		rule = nil
		glogger.GLogger.Infof("Rule [%v] has been deleted", ruleId)
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
	glogger.GLogger.Info("Ready to stop rulex")
	e.InEnds.Range(func(key, value interface{}) bool {
		inEnd := value.(*typex.InEnd)
		if inEnd.Source != nil {
			glogger.GLogger.Info("Stop InEnd:", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).SetState(typex.SOURCE_DOWN)
			inEnd.Source.Stop()
			if inEnd.Source.Driver() != nil {
				inEnd.Source.Driver().Stop()
			}
		}
		return true
	})
	// 停止所有外部资源
	e.OutEnds.Range(func(key, value interface{}) bool {
		outEnd := value.(*typex.OutEnd)
		if outEnd.Target != nil {
			glogger.GLogger.Info("Stop Target:", outEnd.Name, outEnd.UUID)
			outEnd.Target.Stop()
		}
		return true
	})
	// 停止所有插件
	e.Plugins.Range(func(key, value interface{}) bool {
		plugin := value.(typex.XPlugin)
		glogger.GLogger.Info("Stop plugin:", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		return true
	})
	// 停止所有设备
	e.Devices.Range(func(key, value interface{}) bool {
		Device := value.(*typex.Device)
		glogger.GLogger.Info("Stop Device:", Device.Name)
		Device.Device.Stop()
		return true
	})
	// 外挂停了
	e.AllGoods().Range(func(key, value interface{}) bool {
		goodsProcess := value.(*sidecar.GoodsProcess)
		glogger.GLogger.Info("Stop Goods Process:", goodsProcess.UUID())
		goodsProcess.Stop()
		return true
	})

	// 回收资源
	runtime.Gosched()
	runtime.GC()

	if err := glogger.GLOBAL_LOGGER.Close(); err != nil {
		glogger.GLogger.Error(err)
	}
	if err := glogger.LUA_LOGGER.Close(); err != nil {
		glogger.GLogger.Error(err)
	}
	glogger.GLogger.Info("Stop Rulex successfully")
}

//
// 核心功能: Work, 主要就是推流进队列
//
func (e *RuleEngine) WorkInEnd(in *typex.InEnd, data string) (bool, error) {
	if err := e.PushInQueue(in, data); err != nil {
		return false, err
	}
	return true, nil
}

//
// 核心功能: Work, 主要就是推流进队列
//
func (e *RuleEngine) WorkDevice(Device *typex.Device, data string) (bool, error) {
	if err := e.PushDeviceQueue(Device, data); err != nil {
		return false, err
	}
	return true, nil
}

//
// 执行lua脚本
//
func (e *RuleEngine) RunSourceCallbacks(in *typex.InEnd, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range in.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			_, err := core.ExecuteActions(&rule, lua.LString(callbackArgs))
			if err != nil {
				glogger.GLogger.Error("RunLuaCallbacks error:", err)
				_, err := core.ExecuteFailed(rule.VM, lua.LString(err.Error()))
				if err != nil {
					glogger.GLogger.Error(err)
				}
			} else {
				_, err := core.ExecuteSuccess(rule.VM)
				if err != nil {
					glogger.GLogger.Error(err)
					return
				}
			}
		}
	}
}

//
// 执行lua脚本
//
func (e *RuleEngine) RunDeviceCallbacks(Device *typex.Device, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range Device.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			_, err := core.ExecuteActions(&rule, lua.LString(callbackArgs))
			if err != nil {
				glogger.GLogger.Error("RunLuaCallbacks error:", err)
				_, err := core.ExecuteFailed(rule.VM, lua.LString(err.Error()))
				if err != nil {
					glogger.GLogger.Error(err)
				}
			} else {
				_, err := core.ExecuteSuccess(rule.VM)
				if err != nil {
					glogger.GLogger.Error(err)
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
		glogger.GLogger.Infof("Plugin is not enable:%s", p.PluginMetaInfo().Name)
		return nil
	}

	if err := p.Init(section); err != nil {
		return err
	}
	_, ok := e.Plugins.Load(p.PluginMetaInfo().Name)
	if ok {
		return errors.New("plugin already installed:" + p.PluginMetaInfo().Name)
	}

	if err := p.Start(e); err != nil {
		return err
	}

	e.Plugins.Store(p.PluginMetaInfo().Name, p)
	glogger.GLogger.Infof("Plugin start successfully:[%v]", p.PluginMetaInfo().Name)
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
		glogger.GLogger.Infof("InEnd [%v] has been deleted", id)
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
		if outEnd.Target != nil {
			outEnd.Target.Stop()
			e.OutEnds.Delete(uuid)
			outEnd = nil
		}
		glogger.GLogger.Infof("InEnd [%v] has been deleted", uuid)
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
		"system":   utils.BToMb(m.Sys),
		"alloc":    utils.BToMb(m.Alloc),
		"total":    utils.BToMb(m.TotalAlloc),
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
		glogger.GLogger.Error(err)
	}
	return string(b)
}

/*
*
* 加载外部程序
*
 */
func (e *RuleEngine) LoadGoods(goods sidecar.Goods) error {
	return e.SideCar.Fork(goods)
}

//
// 删除外部驱动
//
func (e *RuleEngine) RemoveGoods(uuid string) error {
	if e.GetGoods(uuid) != nil {
		e.SideCar.Remove(uuid)
		return nil
	}
	return fmt.Errorf("goods %v not exists", uuid)
}

//
// 所有外部驱动
//
func (e *RuleEngine) AllGoods() *sync.Map {
	return e.SideCar.AllGoods()
}

//
// 获取某个外部驱动信息
//
func (e *RuleEngine) GetGoods(uuid string) *sidecar.Goods {
	goodsProcess := e.SideCar.Get(uuid)
	goods := sidecar.Goods{
		UUID:        goodsProcess.UUID(),
		Addr:        goodsProcess.Addr(),
		Description: goodsProcess.Description(),
		Args:        goodsProcess.Args(),
	}
	return &goods
}

//
// 取一个进程
//
func (e *RuleEngine) PickUpProcess(uuid string) *sidecar.GoodsProcess {
	return e.SideCar.Get(uuid)
}

//
// 重启源
//
func (e *RuleEngine) RestartInEnd(uuid string) error {
	if value, ok := e.InEnds.Load(uuid); ok {
		o := (value.(*typex.InEnd))
		if o.State == typex.SOURCE_UP {
			o.Source.Stop()
		}
		if err := e.LoadInEnd(o); err != nil {
			glogger.GLogger.Error("InEnd load failed:", err)
			return err
		}
		return nil
	}
	return errors.New("InEnd:" + uuid + "not exists")
}

//
// 重启目标
//
func (e *RuleEngine) RestartOutEnd(uuid string) error {
	if value, ok := e.OutEnds.Load(uuid); ok {
		o := (value.(*typex.OutEnd))
		if o.State == typex.SOURCE_UP {
			o.Target.Stop()
		}
		if err := e.LoadOutEnd(o); err != nil {
			glogger.GLogger.Error("OutEnd load failed:", err)
			return err
		}
		return nil
	}
	return errors.New("OutEnd:" + uuid + "not exists")
}

//
// 重启设备
//
func (e *RuleEngine) RestartDevice(uuid string) error {
	if value, ok := e.Devices.Load(uuid); ok {
		o := (value.(*typex.Device))
		if o.State == typex.DEV_RUNNING {
			o.Device.Stop()
		}
		if err := e.LoadDevice(o); err != nil {
			glogger.GLogger.Error("Device load failed:", err)
			return err
		}
		return nil
	}
	return errors.New("Device:" + uuid + "not exists")
}
