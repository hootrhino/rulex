package demo_plugin

import (
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

type DemoPlugin struct {
	uuid string
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{
		uuid: "DEMO01",
	}
}

func (dm *DemoPlugin) Init(config *ini.Section) error {
	return nil
}

func (dm *DemoPlugin) Start(typex.RuleX) error {
	return nil
}
func (dm *DemoPlugin) Stop() error {
	return nil
}

func (hh *DemoPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "DemoPlugin",
		Version:  "v0.0.1",
		Homepage: "/",
		HelpLink: "/",
		Author:   "RHILEXTeam",
		Email:    "RHILEXTeam@hootrhino.com",
		License:  "AGPL",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *DemoPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
