// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
