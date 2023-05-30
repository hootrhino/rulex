package typex

import "gopkg.in/ini.v1"

//
// 插件开发步骤：
// 1 在ini配置文件中增加配置项
// 2 实现插件接口: XPlugin
// 3 LoadPlugin(sectionK string, p typex.XPlugin)
//

// 插件的服务参数
type ServiceArg struct {
	UUID string      `json:"uuid"` // 插件UUID, Rulex用来查找插件的
	Name string      `json:"name"` // 服务名, 在服务中响应识别
	Args interface{} `json:"args"` // 服务参数
}
type ServiceResult struct {
	Out interface{} `json:"out"`
}

/*
*
* 插件的元信息结构体
*   注意：插件信息这里uuid，name有些是固定写死的，比较特殊，不要轻易改变已有的，否则会导致接口失效
*        只要是已有的尽量不要改这个UUID。
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
	Service(ServiceArg) ServiceResult // 对外提供一些服务
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}
