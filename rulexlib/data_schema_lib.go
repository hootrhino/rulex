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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rulexlib

import (
	"encoding/json"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/component/dataschema"
	"github.com/hootrhino/rulex/typex"
)

/*
*
  - 更新物模型的值
    第一个参数是设备ID；后面是Json，映射到数据模型

*
*/
func DataSchemaValueUpdate(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		deviceId := l.ToString(2)
		SchemaValue := l.ToString(3)
		IoTPropertySlotMap := dataschema.GetSlot(deviceId)
		if IoTPropertySlotMap == nil {
			l.Push(lua.LString("Iot Schema Slot Not Exists"))
			return 1
		}
		LuaArgsJson := map[string]any{}
		if err := json.Unmarshal([]byte(SchemaValue),
			&LuaArgsJson); err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		for K, V := range LuaArgsJson {
			IoTPropertyValue := IoTPropertySlotMap[K]
			IoTPropertyValue.Value = V
			IoTPropertyValue.LastFetchTime = uint64(time.Now().UnixMilli())
			dataschema.SetValue(deviceId, K, IoTPropertyValue)
		}
		l.Push(lua.LNil)
		return 1
	}
}
