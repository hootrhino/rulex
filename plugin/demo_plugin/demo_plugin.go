package demo_plugin

import (
	"rulex/typex"
)

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}
func (dm *DemoPlugin) Load() error {
	return nil
}
func (dm *DemoPlugin) Init() error {
	return nil

}

//
func (dm *DemoPlugin) Install() error {
	return nil
}
func (dm *DemoPlugin) Start() error {
	return nil
}
func (dm *DemoPlugin) Uninstall() error {
	return nil
}
func (dm *DemoPlugin) Clean() {

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
func (hh *DemoPlugin) XPluginEnv() typex.XPluginEnv {

	return typex.XPluginEnv{}
}
