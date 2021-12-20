package typex

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

//
// RuleX interface
//
type RuleX interface {
	//
	// 启动规则引擎
	//
	Start() *sync.Map
	//
	// 消息推到队列
	//
	PushQueue(QueueData) error
	//
	// 执行任务
	//
	Work(*InEnd, string) (bool, error)
	//
	// 获取配置
	//
	GetConfig(k string) interface{}
	//
	// 加载输入
	//
	LoadInEnd(*InEnd) error
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
	LoadOutEnd(*OutEnd) error
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
	LoadPlugin(XPlugin) error
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
	RunLuaCallbacks(*InEnd, string)
	//
	// 运行 hook
	//
	RunHooks(string)
	//
	// 获取版本
	//
	Version() Version

	//
	// 停止规则引擎
	//
	Stop()
}

//
// 拓扑接入点，比如 modbus 检测点等
// UUID: gyh9uo7uh7o67uijh
// Name: ModbusMeter001
// Alive: true
// Tag: modbus
//
type TopologyPoint struct {
	UUID   string `json:"uuid"`
	Parent string `json:"parent"`
	Name   string `json:"name"`
	Alive  bool   `json:"alive"`
	Tag    string `json:"tag"`
}

//
// XResource: 终端资源, 比如实际上的 MQTT 客户端
//
type XResource interface {
	//
	// 测试资源是否可用
	//
	Test(intEndId string) bool

	//
	// 注册InEndID到资源
	//
	Register(intEndId string) error
	//
	// 启动资源
	//
	Start() error
	//
	// 资源是否被启用
	//
	Enabled() bool
	//
	// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
	//
	DataModels() []XDataModel
	//
	// 获取前端表单定义
	//
	Configs() []XConfig
	//
	// 重载: 比如可以在重启的时候把某些数据保存起来
	//
	Reload()
	//
	// 挂起资源, 用来做暂停资源使用
	//
	Pause()
	//
	// 获取资源状态
	//
	Status() ResourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *InEnd
	//
	// 不经过规则引擎处理的直达数据接口
	//
	OnStreamApproached(data string) error
	//
	// 驱动接口, 通常用来和硬件交互
	//
	Driver() XExternalDriver
	//
	//
	//
	Topology() []TopologyPoint
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}

//
// Stream from resource and to target
//
type XTarget interface {
	//
	// 测试资源是否可用
	//
	Test(outEndId string) bool

	//
	// 注册InEndID到资源
	//
	Register(outEndId string) error
	//
	// 启动资源
	//
	Start() error
	//
	// 资源是否被启用
	//
	Enabled() bool
	//
	// 重载: 比如可以在重启的时候把某些数据保存起来
	//
	Reload()
	//
	// 挂起资源, 用来做暂停资源使用
	//
	Pause()
	//
	// 获取资源状态
	//
	Status() ResourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *OutEnd
	//
	// 数据出口
	//
	To(data interface{}) error
	//
	// 不经过规则引擎处理的直达数据
	//
	OnStreamApproached(data string) error
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}

//
// External Plugin
//
type XPlugin interface {
	Init() error
	Start() error
	Stop() error
	XPluginMetaInfo() XPluginMetaInfo
}

//
// XHook for enhancement rulex with golang
//
type XHook interface {
	Work(data string) error
	Error(error)
	Name() string
}

//
// XStatus for resource status
//
type XStatus struct {
	PointId    string // Input: Resource; Output: Target
	Enable     bool
	RuleEngine RuleX
}

//
type XPluginMetaInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Homepage string `json:"homepage"`
	HelpLink string `json:"helpLink"`
	Author   string `json:"author"`
	Email    string `json:"email"`
	License  string `json:"license"`
}

// GoPlugins support
// ONLY for *nix OS!!!
type XModuleInfo struct {
}
type XModule interface {
	Load() XModuleInfo
	UnLoad() error
}

//------------------------------------------------
// 					Remote Stream
//------------------------------------------------
// ┌───────────────┐          ┌────────────────┐
// │   RULEX       │ <─────── │   SERVER       │
// │   RULEX       │  ───────>│   SERVER       │
// └───────────────┘          └────────────────┘
//------------------------------------------------
type XStream interface {
	Start() error
	OnStreamApproached(data string) error
	State() XStatus
	Close()
}

//
// User's Protocol
//
type XProtocol interface {
	// 是否增加用户解码器接口？
	// 这里就需要用lua来实现这么一套工具
}

//
// 外挂驱动, 比如串口, PLC等, 驱动可以挂在输入或者输出资源上。
// 典型案例:
// 1. MODBUS TCP模式 ,数据输入后转JSON输出到串口屏幕上
// 2. MODBUS TCP模式外挂了很多继电器,来自云端的 PLC 控制指令先到网关, 然后网关决定推送到哪个外挂
//
type DriverDetail struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type XExternalDriver interface {
	Test() error
	Init() error
	Work() error
	State() DriverState
	SetState(DriverState)
	//---------------------------------------------------
	// 读写接口是给LUA标准库用的, 驱动只管实现读写逻辑即可
	//---------------------------------------------------
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	//---------------------------------------------------
	DriverDetail() *DriverDetail
	Stop() error
}

//
// XLib: 库函数接口
//
type XLib interface {
	Name() string
	LibFun(RuleX) func(*lua.LState) int
}
