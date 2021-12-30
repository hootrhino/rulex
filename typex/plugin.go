package typex

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
// External Plugin
//
type XPlugin interface {
	Init() error
	Start() error
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}
