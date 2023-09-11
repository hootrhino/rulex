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

package core

import "github.com/hootrhino/rulex/typex"

/*
*
* 主要用来处理模型和数据
单线条：

	{
	    "value": 56.1，
	    "date": "2023-09-10 13:45:56"
	}

多线条：

	{
	    "category" :"temp",
	    "value": 56.1，
	    "date": "2023-09-10 13:45:56"
	}

*
*/

/*
*
* 单线条, data长度为1，其中有一个轴固定为‘date’
^
|                       *****
|                     **     **
| **                 **         **                        *
|  **               **            **                     **
|   **            **               **                  **
|    **         **                 **                **
|      **     **                    **            **
|        ******                        **       **
|                                        ********
+------------------------------------------------------------->
*/
func FilterSingleDataWithSchema(data map[string]interface{},
	dataDefine []typex.DataDefine) map[string]interface{} {
	if len(dataDefine) != 1 {
		return nil
	}
	Define := dataDefine[0]
	result := map[string]interface{}{}
	result["date"] = data["date"]
	result[Define.Name] = data[Define.Name]
	return result
}
