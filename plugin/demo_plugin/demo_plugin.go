package demo_plugin

import (
	"rulex/typex"
)

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}
func (dm *DemoPlugin) Load() *typex.XPluginEnv {
	return typex.NewXPluginEnv()
}
func (dm *DemoPlugin) Init(*typex.XPluginEnv) error {
	return nil

}

//
func (dm *DemoPlugin) Install(*typex.XPluginEnv) (*typex.XPluginMetaInfo, error) {
	return &typex.XPluginMetaInfo{
		Name:     "DemoPlugin",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}, nil
}
func (dm *DemoPlugin) Start(*typex.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Uninstall(*typex.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Clean() {

}
