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

package softwatchdog

import (
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* 软件看门狗
*
 */
type softWatchDog struct {
	uuid string
}

func NewSoftWatchDog() *softWatchDog {
	return &softWatchDog{
		uuid: "SOFT_WATCHDOG",
	}
}

func (dog *softWatchDog) Init(config *ini.Section) error {
	return nil
}

func (dog *softWatchDog) Start(typex.RuleX) error {
	return nil
}
func (dog *softWatchDog) Stop() error {
	return nil
}

func (hh *softWatchDog) PluginMetaInfo() typex.XPluginMetaInfo {
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
func (cs *softWatchDog) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
