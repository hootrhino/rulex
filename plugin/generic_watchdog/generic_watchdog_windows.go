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

package genericwatchdog

import (
	"fmt"
	"runtime"

	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* 软件看门狗
*
 */
type genericWatchDog struct {
	uuid string
}

func NewGenericWatchDog() *genericWatchDog {
	return &genericWatchDog{
		uuid: "SOFT_WATCHDOG",
	}
}

func (dog *genericWatchDog) Init(config *ini.Section) error {
	return fmt.Errorf("OS support Wdog:%s", runtime.GOOS)
}

func (dog *genericWatchDog) Start(typex.RuleX) error {

	return nil
}
func (dog *genericWatchDog) Stop() error {
	return nil
}

func (hh *genericWatchDog) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "Soft WatchDog",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
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
func (cs *genericWatchDog) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
