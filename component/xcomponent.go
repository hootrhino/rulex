// Copyright (C) 2024 wwhai
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

package component

import (
	"github.com/hootrhino/rulex/typex"
)

type XComponentMetaInfo struct {
	UUID    string `json:"uuid"`    // UUID
	Name    string `json:"name"`    // 组件名
	Version string `json:"version"` // 版本
}

/*
*
* RULEX 系统组件
*
 */
type XComponent interface {
	Init(cfg map[string]any) error      // 配置
	Start(r typex.RuleX) error          // 启动
	Stop() error                        // 停止
	PluginMetaInfo() XComponentMetaInfo // 元信息
}
