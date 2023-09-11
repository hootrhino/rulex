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

package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/aibase"
	"github.com/hootrhino/rulex/appstack"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/interqueue"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/trailer"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/shirou/gopsutil/v3/disk"
)

// 规则引擎
type RuleEngine struct {
	Hooks             *sync.Map            `json:"hooks"`
	Rules             *sync.Map            `json:"rules"`
	Plugins           *sync.Map            `json:"plugins"`
	InEnds            *sync.Map            `json:"inends"`
	OutEnds           *sync.Map            `json:"outends"`
	Drivers           *sync.Map            `json:"drivers"`
	Devices           *sync.Map            `json:"devices"`
	Config            *typex.RulexConfig   `json:"config"`
	Trailer           typex.XTrailer       `json:"-"`
	AppStack          typex.XAppStack      `json:"-"`
	AiBaseRuntime     typex.XAiRuntime     `json:"-"`
	DeviceTypeManager typex.DeviceRegistry `json:"-"`
	SourceTypeManager typex.SourceRegistry `json:"-"`
	TargetTypeManager typex.TargetRegistry `json:"-"`
	MetricStatistics  *typex.MetricStatistics
}

func NewRuleEngine(config typex.RulexConfig) typex.RuleX {
	re := &RuleEngine{
		DeviceTypeManager: core.NewDeviceTypeManager(),
		SourceTypeManager: core.NewSourceTypeManager(),
		TargetTypeManager: core.NewTargetTypeManager(),
		Plugins:           &sync.Map{},
		Hooks:             &sync.Map{},
		Rules:             &sync.Map{},
		InEnds:            &sync.Map{},
		OutEnds:           &sync.Map{},
		Drivers:           &sync.Map{},
		Devices:           &sync.Map{},
		Config:            &config,
		MetricStatistics:  typex.NewMetricStatistics(),
	}
	// trailer
	re.Trailer = trailer.NewTrailerManager(re)
	// lua appstack manager
	re.AppStack = appstack.NewAppStack(re)
	// current only support Internal ai
	re.AiBaseRuntime = aibase.NewAIRuntime(re)
	// Queue
	interqueue.InitDataCacheQueue(re, core.GlobalConfig.MaxQueueSize)
	return re
}
func (e *RuleEngine) GetMetricStatistics() *typex.MetricStatistics {
	return e.MetricStatistics
}
func (e *RuleEngine) GetAiBase() typex.XAiRuntime {
	return e.AiBaseRuntime
}
func (e *RuleEngine) Start() *typex.RulexConfig {
	e.InitDeviceTypeManager()
	e.InitSourceTypeManager()
	e.InitTargetTypeManager()
	// 内部队列
	interqueue.InitDataCacheQueue(e, core.GlobalConfig.MaxQueueSize)
	interqueue.StartDataCacheQueue()
	interqueue.InitInteractQueue(e, core.GlobalConfig.MaxQueueSize)
	interqueue.StartInteractQueue()
	// 前后交互组件
	core.InitWebDataPipe(e)
	core.InitInternalSchemaCache()
	go core.StartWebDataPipe()
	return e.Config
}

func (e *RuleEngine) GetPlugins() *sync.Map {
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() *sync.Map {
	return e.Plugins
}

func (e *RuleEngine) Version() typex.Version {
	return typex.DefaultVersion
}

func (e *RuleEngine) GetConfig() *typex.RulexConfig {
	return e.Config
}

// Stop
func (e *RuleEngine) Stop() {
	glogger.GLogger.Info("[*] Ready to stop rulex")
	e.InEnds.Range(func(key, value interface{}) bool {
		inEnd := value.(*typex.InEnd)
		if inEnd.Source != nil {
			glogger.GLogger.Info("Stop InEnd:", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).State = typex.SOURCE_STOP
			inEnd.Source.Stop()
			if inEnd.Source.Driver() != nil {
				inEnd.Source.Driver().Stop()
			}
		}
		glogger.GLogger.Info("Stop InEnd:", inEnd.Name, inEnd.UUID, " Successfully")
		return true
	})
	// 停止所有外部资源
	e.OutEnds.Range(func(key, value interface{}) bool {
		outEnd := value.(*typex.OutEnd)
		if outEnd.Target != nil {
			glogger.GLogger.Info("Stop NewTarget:", outEnd.Name, outEnd.UUID)
			e.GetOutEnd(outEnd.UUID).State = typex.SOURCE_STOP
			outEnd.Target.Stop()
			glogger.GLogger.Info("Stop NewTarget:", outEnd.Name, outEnd.UUID, " Successfully")
		}
		return true
	})
	// 停止所有插件
	e.Plugins.Range(func(key, value interface{}) bool {
		plugin := value.(typex.XPlugin)
		glogger.GLogger.Info("Stop plugin:", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		glogger.GLogger.Info("Stop plugin:", plugin.PluginMetaInfo().Name, " Successfully")
		return true
	})
	// 停止所有设备
	e.Devices.Range(func(key, value interface{}) bool {
		Device := value.(*typex.Device)
		glogger.GLogger.Info("Stop Device:", Device.Name)
		e.GetDevice(Device.UUID).State = typex.DEV_STOP
		Device.Device.Stop()
		glogger.GLogger.Info("Stop Device:", Device.Name, " Successfully")
		return true
	})
	// 外挂停了
	e.Trailer.Stop()
	// 所有的APP停了
	e.AppStack.Stop()
	glogger.GLogger.Info("[√] Stop Rulex successfully")
	if err := glogger.Close(); err != nil {
		fmt.Println("Close logger error: ", err)
	}
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkInEnd(in *typex.InEnd, data string) (bool, error) {
	if err := interqueue.DefaultDataCacheQueue.PushInQueue(in, data); err != nil {
		return false, err
	}
	return true, nil
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkDevice(Device *typex.Device, data string) (bool, error) {
	if err := interqueue.DefaultDataCacheQueue.PushDeviceQueue(Device, data); err != nil {
		return false, err
	}
	return true, nil
}

/*
*
* 执行针对资源端的规则脚本
*
 */
func (e *RuleEngine) RunSourceCallbacks(in *typex.InEnd, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range in.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			if rule.Type == "lua" {
				_, err := core.ExecuteActions(&rule, lua.LString(callbackArgs))
				if err != nil {
					glogger.GLogger.Error("RunLuaCallbacks error:", err)
					_, err := core.ExecuteFailed(rule.LuaVM, lua.LString(err.Error()))
					if err != nil {
						glogger.GLogger.Error(err)
					}
				} else {
					_, err := core.ExecuteSuccess(rule.LuaVM)
					if err != nil {
						glogger.GLogger.Error(err)
						return // lua 是规则链，有短路原则，中途出错会中断
					}
				}
			}
		}
	}
}

/*
*
* 执行针对设备端的规则脚本
*
 */
func (e *RuleEngine) RunDeviceCallbacks(Device *typex.Device, callbackArgs string) {
	for _, rule := range Device.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			if rule.Type == "lua" {
				_, err := core.ExecuteActions(&rule, lua.LString(callbackArgs))
				if err != nil {
					glogger.GLogger.Error("RunLuaCallbacks error:", err)
					_, err := core.ExecuteFailed(rule.LuaVM, lua.LString(err.Error()))
					if err != nil {
						glogger.GLogger.Error(err)
					}
				} else {
					_, err := core.ExecuteSuccess(rule.LuaVM)
					if err != nil {
						glogger.GLogger.Error(err)
						return
					}
				}
			}

		}
	}
}

// LoadHook
func (e *RuleEngine) LoadHook(h typex.XHook) error {
	value, _ := e.Hooks.Load(h.Name())
	if value != nil {
		return errors.New("hook have been loaded:" + h.Name())
	}
	e.Hooks.Store(h.Name(), h)
	return nil

}

// RunHooks
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

func (e *RuleEngine) GetInEnd(uuid string) *typex.InEnd {
	v, ok := (e.InEnds).Load(uuid)
	if ok {
		return v.(*typex.InEnd)
	}
	return nil
}

func (e *RuleEngine) SaveInEnd(in *typex.InEnd) {
	e.InEnds.Store(in.UUID, in)
}

func (e *RuleEngine) RemoveInEnd(id string) {
	if inEnd := e.GetInEnd(id); inEnd != nil {
		inEnd.Source.Stop()
		e.InEnds.Delete(id)
		inEnd = nil
		glogger.GLogger.Infof("InEnd [%v] has been deleted", id)
	}
}

func (e *RuleEngine) AllInEnd() *sync.Map {
	return e.InEnds
}

func (e *RuleEngine) GetOutEnd(id string) *typex.OutEnd {
	v, ok := e.OutEnds.Load(id)
	if ok {
		return v.(*typex.OutEnd)
	} else {
		return nil
	}

}

func (e *RuleEngine) SaveOutEnd(out *typex.OutEnd) {
	e.OutEnds.Store(out.UUID, out)

}

func (e *RuleEngine) RemoveOutEnd(uuid string) {
	if outEnd := e.GetOutEnd(uuid); outEnd != nil {
		if outEnd.Target != nil {
			outEnd.Target.Stop()
			e.OutEnds.Delete(uuid)
			outEnd = nil
		}
		glogger.GLogger.Infof("OutEnd [%v] has been deleted", uuid)
	}
}

func (e *RuleEngine) AllOutEnd() *sync.Map {
	return e.OutEnds
}

// -----------------------------------------------------------------
// 获取运行时快照
// -----------------------------------------------------------------
func (e *RuleEngine) SnapshotDump() string {
	inends := []interface{}{}
	rules := []interface{}{}
	plugins := []interface{}{}
	outends := []interface{}{}
	devices := []interface{}{}
	drivers := []interface{}{}
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
		Device := value.(*typex.Device)
		if Device.Device.Driver() != nil {
			devices = append(devices, Device.Device.Driver())
		}
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
		"drivers":    drivers,
		"statistics": e.MetricStatistics,
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
func (e *RuleEngine) LoadGoods(goods typex.Goods) error {
	return e.Trailer.Fork(goods)
}

// 删除外部驱动
func (e *RuleEngine) RemoveGoods(uuid string) error {
	if e.GetGoods(uuid) != nil {
		e.Trailer.Remove(uuid)
		return nil
	}
	return fmt.Errorf("goods %v not exists", uuid)
}

// 所有外部驱动
func (e *RuleEngine) AllGoods() *sync.Map {
	return e.Trailer.AllGoods()
}

// 获取某个外部驱动信息
func (e *RuleEngine) GetGoods(uuid string) *typex.Goods {
	goodsProcess := e.Trailer.Get(uuid)
	goods := typex.Goods{
		UUID:        goodsProcess.Uuid,
		Addr:        goodsProcess.Addr,
		Description: goodsProcess.Description,
		Args:        goodsProcess.Args,
	}
	return &goods
}

// 取一个进程
func (e *RuleEngine) PickUpProcess(uuid string) *typex.GoodsProcess {
	return e.Trailer.Get(uuid)
}

// 重启源
func (e *RuleEngine) RestartInEnd(uuid string) error {
	if _, ok := e.InEnds.Load(uuid); ok {
		return nil
	}
	return errors.New("InEnd:" + uuid + "not exists")
}

// 重启目标
func (e *RuleEngine) RestartOutEnd(uuid string) error {
	if _, ok := e.OutEnds.Load(uuid); ok {
		return nil
	}
	return errors.New("OutEnd:" + uuid + "not exists")
}

// 重启设备
func (e *RuleEngine) RestartDevice(uuid string) error {
	if _, ok := e.Devices.Load(uuid); ok {
		return nil
	}
	return errors.New("Device:" + uuid + "not exists")
}

/*
*
* 初始化设备管理器
*
 */
func (e *RuleEngine) InitDeviceTypeManager() error {
	e.DeviceTypeManager.Register(typex.TSS200V02,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewTS200Sensor,
		},
	)
	e.DeviceTypeManager.Register(typex.YK08_RELAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewYK8Controller,
		},
	)
	e.DeviceTypeManager.Register(typex.RTU485_THER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewRtu485Ther,
		},
	)
	e.DeviceTypeManager.Register(typex.S1200PLC,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewS1200plc,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_MODBUS,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericModbusDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_UART,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericUartDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_SNMP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericSnmpDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.USER_G776,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewUsrG776DTU,
		},
	)
	e.DeviceTypeManager.Register(typex.ICMP_SENDER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewIcmpSender,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_PROTOCOL,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewCustomProtocolDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_OPCUA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericOpcuaDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_CAMERA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewVideoCamera,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_AIS,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewAISDeviceMaster,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_BACNET_IP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericBacnetIpDevice,
		},
	)
	return nil
}

/*
*
* 初始化输入资源管理器
*
 */
func (e *RuleEngine) InitSourceTypeManager() error {
	e.SourceTypeManager.Register(typex.MQTT,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewMqttInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.HTTP,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewHttpInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.COAP,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewCoAPInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.GRPC,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewGrpcInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.NATS_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewNatsSource,
		},
	)
	e.SourceTypeManager.Register(typex.RULEX_UDP,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewUdpInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.GENERIC_IOT_HUB,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewIoTHubSource,
		},
	)
	return nil
}

/*
*
* 初始化输出资源管理器
*
 */
func (e *RuleEngine) InitTargetTypeManager() error {
	e.TargetTypeManager.Register(typex.MONGO_SINGLE,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMongoTarget,
		},
	)
	e.TargetTypeManager.Register(typex.MQTT_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMqttTarget,
		},
	)
	e.TargetTypeManager.Register(typex.NATS_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewNatsTarget,
		},
	)
	e.TargetTypeManager.Register(typex.HTTP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewHTTPTarget,
		},
	)
	e.TargetTypeManager.Register(typex.TDENGINE_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewTdEngineTarget,
		},
	)
	e.TargetTypeManager.Register(typex.GRPC_CODEC_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewCodecTarget,
		},
	)
	e.TargetTypeManager.Register(typex.UDP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewUdpTarget,
		},
	)
	e.TargetTypeManager.Register(typex.SQLITE_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewSqliteTarget,
		},
	)
	e.TargetTypeManager.Register(typex.USER_G776_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewUserG776,
		},
	)
	return nil
}
