package demo_plugin

import "rulex/core"

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}
func (dm *DemoPlugin) Load() *core.XPluginEnv {
	return core.NewXPluginEnv()
}
func (dm *DemoPlugin) Init(*core.XPluginEnv) error {
	return nil

}
//
func (dm *DemoPlugin) Install(*core.XPluginEnv) (*core.XPluginMetaInfo, error) {
	return &core.XPluginMetaInfo{
		Name:     "DemoPlugin",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}, nil
}
func (dm *DemoPlugin) Start(*core.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Uninstall(*core.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Clean() {

}
