package typex

import (
	"context"
	"sync"
)

// Global config
type Extlib struct {
	Value []string `ini:"extlibs,,allowshadow" json:"extlibs"`
}
type RulexConfig struct {
	AppName               string `ini:"app_name" json:"appName"`
	AppId                 string `ini:"app_id" json:"appId"`
	MaxQueueSize          int    `ini:"max_queue_size" json:"maxQueueSize"`
	SourceRestartInterval int    `ini:"resource_restart_interval" json:"sourceRestartInterval"`
	GomaxProcs            int    `ini:"gomax_procs" json:"gomaxProcs"`
	EnablePProf           bool   `ini:"enable_pprof" json:"enablePProf"`
	EnableConsole         bool   `ini:"enable_console" json:"enableConsole"`
	LogLevel              string `ini:"log_level" json:"logLevel"`
	LogPath               string `ini:"log_path" json:"logPath"`
	LuaLogPath            string `ini:"lua_log_path" json:"luaLogPath"`
	MaxStoreSize          int    `ini:"max_store_size" json:"maxStoreSize"`
	AppDebugMode          bool   `ini:"app_debug_mode" json:"appDebugMode"`
	Extlibs               Extlib `ini:"extlibs,,allowshadow" json:"extlibs"`
	UpdateServer          string `ini:"update_server" json:"updateServer"`
}

// RuleX interface
type RuleX interface {
	//
	// 启动规则引擎
	//
	Start() *RulexConfig

	//
	// 执行任务
	//
	WorkInEnd(*InEnd, string) (bool, error)
	WorkDevice(*Device, string) (bool, error)
	//
	// 获取配置
	//
	GetConfig() *RulexConfig
	//
	// 加载输入
	//
	LoadInEndWithCtx(in *InEnd, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 获取输入
	//
	GetInEnd(string) *InEnd
	//
	// 保存输入
	//
	SaveInEnd(*InEnd)
	//
	// 删除输入
	//
	RemoveInEnd(string)
	//
	// 所有输入列表
	//
	AllInEnd() *sync.Map
	//
	// 加载输出
	//
	LoadOutEndWithCtx(in *OutEnd, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 所有输出
	//
	AllOutEnd() *sync.Map
	//
	// 获取输出
	//
	GetOutEnd(string) *OutEnd
	//
	// 保存输出
	//
	SaveOutEnd(*OutEnd)
	//
	// 删除输出
	//
	RemoveOutEnd(string)
	//
	// 加载Hook
	//
	LoadHook(XHook) error
	//
	// 加载插件
	//
	LoadPlugin(string, XPlugin) error
	//
	// 所有插件列表
	//
	AllPlugins() *sync.Map
	//
	// 加载规则
	//
	LoadRule(*Rule) error
	//
	// 所有规则列表
	//
	AllRule() *sync.Map
	//
	// 获取规则
	//
	GetRule(id string) *Rule
	//
	// 删除规则
	//
	RemoveRule(uuid string)
	//
	// 运行 lua 回调
	//
	RunSourceCallbacks(*InEnd, string)
	RunDeviceCallbacks(*Device, string)
	//
	// 运行 hook
	//
	RunHooks(string) //TODO Hook 未来某个版本会加强,主要用来加载本地动态库
	//
	// 获取版本
	//
	Version() Version

	//
	// 停止规则引擎
	//
	Stop()
	//
	// Snapshot Dump
	//
	SnapshotDump() string
	//
	// 加载设备
	//
	LoadDeviceWithCtx(in *Device, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 获取设备
	//
	GetDevice(string) *Device
	//
	// 保存设备
	//
	SaveDevice(*Device)
	//
	//
	//
	AllDevices() *sync.Map
	//
	// 删除设备
	//
	RemoveDevice(string)
	//
	// 重启源
	//
	RestartInEnd(uuid string) error
	//
	// 重启目标
	//
	RestartOutEnd(uuid string) error
	//
	// 重启设备
	//
	RestartDevice(uuid string) error
}

// 拓扑接入点，比如 modbus 检测点等
// UUID: gyh9uo7uh7o67u
// Name: ModbusMeter001
// Alive: true
// Tag: modbus
type TopologyPoint struct {
	UUID   string `json:"uuid"`
	Parent string `json:"parent"`
	Name   string `json:"name"`
	Alive  bool   `json:"alive"`
	Tag    string `json:"tag"`
}
