package demo_plugin

import (
	"rulex/typex"
)

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}

func (dm *DemoPlugin) Init() error {
	return nil
}

func (dm *DemoPlugin) Start() error {
	return nil
}
func (dm *DemoPlugin) Stop() error {
	return nil
}


func (hh *DemoPlugin) XPluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "DemoPlugin",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
