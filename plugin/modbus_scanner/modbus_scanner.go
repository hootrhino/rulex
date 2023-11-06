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

package modbusscanner

import (
	"context"

	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* 串口配置
*
 */
type modbusScanner struct {
	uuid    string
	busying bool
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewModbusScanner() *modbusScanner {

	return &modbusScanner{
		uuid:    "MODBUS_SCANNER",
		busying: false,
	}
}

func (ms *modbusScanner) Init(config *ini.Section) error {
	return nil
}

func (ms *modbusScanner) Start(typex.RuleX) error {
	return nil
}
func (ms *modbusScanner) Stop() error {
	return nil
}

func (hh *modbusScanner) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "Modbus Device Scanner",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
