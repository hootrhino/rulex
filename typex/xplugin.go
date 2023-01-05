package typex

import "gopkg.in/ini.v1"

//
// 插件开发步骤：
// 1 在ini配置文件中增加配置项
// 2 实现插件接口: XPlugin
// 3 LoadPlugin(sectionK string, p typex.XPlugin)
//

//
// External Plugin
//
type XPlugin interface {
	Init(*ini.Section) error // 参数为外部配置
	Start(RuleX) error
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}

type XPluginMetaInfo struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Homepage string `json:"homepage"`
	HelpLink string `json:"helpLink"`
	Author   string `json:"author"`
	Email    string `json:"email"`
	License  string `json:"license"`
}
