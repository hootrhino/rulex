package typex

//
// !!! All 'RuleEngine' parameter passed by pseudo constructure function !!!
//
import (
	"sync"
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
	Start() *map[string]interface{}
	//
	PushQueue(QueueData) error
	//
	Work(in *InEnd, data string) (bool, error)
	//
	GetConfig(k string) interface{}
	//
	LoadInEnd(in *InEnd) error
	GetInEnd(id string) *InEnd
	SaveInEnd(in *InEnd)
	RemoveInEnd(id string)
	AllInEnd() map[string]*InEnd
	//
	LoadOutEnd(out *OutEnd) error
	AllOutEnd() map[string]*OutEnd
	GetOutEnd(id string) *OutEnd
	SaveOutEnd(out *OutEnd)
	RemoveOutEnd(out *OutEnd)
	//
	LoadHook(h XHook) error
	//
	LoadPlugin(p XPlugin) error
	AllPlugins() *map[string]*XPluginMetaInfo
	//
	LoadRule(r *Rule) error
	AllRule() map[string]*Rule
	RemoveRule(uuid string) error
	//
	RunLuaCallbacks(in *InEnd, data string)
	//
	RunHooks(data string)
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
	DataModels() *map[string]XDataModel
	Reload()
	Pause()
	Status() ResourceState
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
	Stop()
}

//
// External Plugin
//
type XPlugin interface {
	Load() *XPluginEnv
	Init(*XPluginEnv) error
	Install(*XPluginEnv) (*XPluginMetaInfo, error)
	Start(*XPluginEnv) error
	Uninstall(*XPluginEnv) error
	Clean()
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
	sync.Mutex
	PointId    string // Input: Resource; Output: Target
	Enable     bool
	RuleEngine RuleX
}

//
// External Plugin
//
type XPluginEnv struct {
	env *map[string]interface{}
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

//
func NewXPluginMetaInfo() *XPluginMetaInfo {
	return &XPluginMetaInfo{}
}

//
func NewXPluginEnv() *XPluginEnv {
	return &XPluginEnv{
		env: &map[string]interface{}{},
	}
}

//
func (p *XPluginEnv) Get(k string) interface{} {
	return (*(p.env))[k]
}

//
func (p *XPluginEnv) Set(k string, v interface{}) {
	(*(p.env))[k] = v
}

// GoPlugins support
// ONLY for *nix OS!!!
type XModuleInfo struct {
}
type XModule interface {
	Load() XModuleInfo
	UnLoad() error
}

//--------------------------------------------------------------
// Remote Stream
// ┌───────────────┐          ┌────────────────┐
// │   RULEX       │ ◄─────── │   SERVER       │
// └───────────────┘          └────────────────┘
// Rulex Stream use GRPC for transport layer
//--------------------------------------------------------------
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
}

type XDriver interface {
	Test() error
	Init() error
	Work() error
	State() DriverState
	Stop() error
}
