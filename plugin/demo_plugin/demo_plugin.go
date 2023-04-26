package demo_plugin

import (
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
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
		Name:     "DemoPlugin",
		Version:  "0.0.1",
		Homepage: "www.github.com/hootrhino/rulex",
		HelpLink: "www.github.com/hootrhino/rulex",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *DemoPlugin) Service(arg typex.ServiceArg) error {
	return nil
}
