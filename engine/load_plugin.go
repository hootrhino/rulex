package engine

import (
	"errors"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

// ┌──────┐    ┌──────┐    ┌──────┐
// │ Init ├───►│ Load ├───►│ Stop │
// └──────┘    └──────┘    └──────┘
func (e *RuleEngine) LoadPlugin(sectionK string, p typex.XPlugin) error {
	section := utils.GetINISection(core.INIPath, sectionK)
	/*key, err1 := section.GetKey("enable")
	if err1 != nil {
		return err1
	}
	enable, err2 := key.Bool()
	if err2 != nil {
		return err2
	}
	if !enable {
		glogger.GLogger.Infof("Plugin is not enable:%s", p.PluginMetaInfo().Name)
		return nil
	}*/

	if err := p.Init(section); err != nil {
		return err
	}
	_, ok := e.Plugins.Load(p.PluginMetaInfo().UUID)
	if ok {
		return errors.New("plugin already installed:" + p.PluginMetaInfo().Name)
	}

	if err := p.Start(e); err != nil {
		return err
	}

	e.Plugins.Store(p.PluginMetaInfo().UUID, p)
	glogger.GLogger.Infof("Plugin start successfully:[%v]", p.PluginMetaInfo().Name)
	return nil

}
