package x

import "sync"

//
// Rule is for desplay, XResource is actully worker
// XResource{
//  inEndId
//  enabled
//  status.
//}
//
type XResource interface {
	Test(inEndId string) bool      //0
	Register(inEndId string) error //1
	Start() error                  //2
	Enabled() bool
	Reload()
	Pause()
	Status() State
	Stop()
}

//
// Stream from resource and to target
//
type XTarget interface {
	Test(outEndId string) bool      //0
	Register(outEndId string) error //1
	Start() error                   //2
	Enabled() bool
	Reload()
	Pause()
	Status() State
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
	InEndId string
	Enable  bool
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
