package typex

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

//
// 驱动的数据模型
//
type XDataModel struct {
	Type      ModelType
	Name      string
	MaxLength int
	MinLength int
}

//
// RuleX interface
//
type RuleX interface {
	Start() *sync.Map
	//
	PushQueue(QueueData) error
	//
	Work(*InEnd, string) (bool, error)
	//
	GetConfig(k string) interface{}
	//
	LoadInEnd(*InEnd) error
	GetInEnd(string) *InEnd
	SaveInEnd(*InEnd)
	RemoveInEnd(string)
	AllInEnd() *sync.Map
	//
	LoadOutEnd(*OutEnd) error
	AllOutEnd() *sync.Map
	GetOutEnd(string) *OutEnd
	SaveOutEnd(*OutEnd)
	RemoveOutEnd(string)
	//
	LoadHook(XHook) error
	//
	LoadPlugin(XPlugin) error
	AllPlugins() *sync.Map
	//
	LoadRule(*Rule) error
	AllRule() *sync.Map
	GetRule(id string) *Rule
	RemoveRule(uuid string)
	//
	RunLuaCallbacks(*InEnd, string)
	//
	RunHooks(string)
	//
	//
	Version() Version
	//
	Stop()
}

//
// XResource: 终端资源，比如实际上的 MQTT 客户端
//
type XResource interface {
	Details() *InEnd
	Test(inEndId string) bool      //0
	Register(inEndId string) error //1
	Start() error                  //2
	Enabled() bool
	DataModels() []XDataModel
	Reload()
	Pause()
	Status() ResourceState
	OnStreamApproached(data string) error
	Driver() XExternalDriver
	Stop()
}

//
// Stream from resource and to target
//
type XTarget interface {
	Details() *OutEnd
	Test(outEndId string) bool      //0
	Register(outEndId string) error //1
	Start() error                   //2
	Enabled() bool
	Reload()
	Pause()
	Status() ResourceState
	To(data interface{}) error
	OnStreamApproached(data string) error
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
// 外挂驱动，比如串口，PLC等，驱动可以挂在输入或者输出资源上。
// 典型案例：
// 1. MODBUS TCP模式 ,数据输入后转JSON输出到串口屏幕上
// 2. MODBUS TCP模式外挂了很多继电器,来自云端的 PLC 控制指令先到网关，然后网关决定推送到哪个外挂
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
