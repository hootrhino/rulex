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

package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/component/interqueue"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* applib:DatToComponent('ComponentId', 'SchemaId', Data)

func FilterSingleDataWithSchema(map[string]interface{}, []typex.DataDefine) map[string]interface{}
	{
	    "value": 56.1，
	    "date": "2023-09-10 13:45:56"
	}
*/

func DataToUiComponent(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		ComponentId := l.ToString(2)
		SchemaId := l.ToString(3)
		Data := l.ToTable(4)
		InMap := map[string]interface{}{}
		Data.ForEach(func(l1, l2 lua.LValue) {
			InMap[l1.String()] = l2
		})
		// SIMPLE_DATE_LINE     简单一线,一轴为时间走势
		// COMPLEX_DATE_LINE    复杂多线,一轴为时间走势
		// COMPLEX_CROSS_LINE   复杂线，其中X、Y轴可以互相变换
		if schema, ok := core.SchemaGet(SchemaId); ok {
			if schema.Type == "SIMPLE_DATE_LINE" {
				result := core.FilterSingleDataWithSchema(InMap, schema.Schema)
				dataToUiComponent(ComponentId, result)
				goto END
			}
			l.Push(lua.LString("Schema not found:" + SchemaId))
			return 1
		}
	END:
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* 将数据推到YQueue

	{
	   "topic" : "/visual/tocomponent/{$$component-UUID}",
	   "componentid" :"EVINU99",
	   "data" : {
	        "value": 56.1，
	        "date": "2023-09-10 13:45:56"
	    }
	}
*/
func dataToUiComponent(ComponentId string, data map[string]interface{}) {
	interqueue.SendData(interqueue.InteractQueueData{
		Topic:       "/visual/component/" + ComponentId,
		ComponentId: ComponentId,
		Data:        data,
	})
}
