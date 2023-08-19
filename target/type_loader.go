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

package target

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
)

var TM typex.TargetRegistry

/*
*
* 给前端返回资源类型，这里是个蹩脚的设计
* 以实现功能为主，后续某个重构版本会做的优雅点
*
 */

func LoadTt() {
	TM = core.NewTargetTypeManager()
	TM.Register(typex.HTTP_TARGET, &typex.XConfig{})
	TM.Register(typex.MONGO_SINGLE, &typex.XConfig{})
	TM.Register(typex.MQTT_TARGET, &typex.XConfig{})
	TM.Register(typex.NATS_TARGET, &typex.XConfig{})
	TM.Register(typex.TDENGINE_TARGET, &typex.XConfig{})
}
