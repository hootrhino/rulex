package typex

import "gopkg.in/ini.v1"

//
// 插件开发步骤：
// 1 在ini配置文件中增加配置项
// 2 实现插件接口: XPlugin
// 3 LoadPlugin(sectionK string, p typex.XPlugin)
//

//
// 插件的服务参数
//
type ServiceArg struct {
	UUID string      `json:"uuid"` // 插件UUID
	Name string      `json:"name"` // 服务名
	Args interface{} `json:"args"` // 参数
}

/*
*
* 插件的元信息
*
 */
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

/*
*
* 插件: 用来增强RULEX的外部功能，本色不属于RULEX
*
 */
type XPlugin interface {
	Init(*ini.Section) error // 参数为外部配置
	Start(RuleX) error
	Service(ServiceArg) error // 对外提供一些服务
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}
