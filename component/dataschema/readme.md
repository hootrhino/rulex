<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

# 物模型
轻量级物模型，主要用于在边缘侧显示数据，约束数据，属于IoTHub的物模型的子集。
## 使用
当设备需要注册自己的模型时，使用`dataschema.RegisterSlot`,卸载模型`dataschema.RegisterSlot`,更新值`dataschema.SetValue`.
```lua
SchemaId = db.device.SchemaId
if SchemaId!=""->
    dataschema.RegisterSlot(mdev.PointId)
    ModbusSchemaCacheValues = db.find(SELECT * FROM `m_iot_properties` WHERE schema_id="SCHEMALD74LCS3")
    for each(ModbusSchemaCacheValues) ->
        dataschema.SetValue(mdev.PointId, K, V)
```
## 数据
请求地址：`http://127.0.0.1:2580/api/v1/devices/properties?current=1&size=10&uuid=DEVICENKRZFRYW`
Lua:
```lua
local R = dataschema:Update('DEVICENKRZFRYW', json:T2J({
    a = acc1,
    b = acc2
}))
```

JSON示例：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "current": 1,
        "size": 10,
        "total": 2,
        "records": [
            {
                "label": "b",
                "name": "b",
                "description": "",
                "type": "BOOL",
                "rw": "R",
                "unit": "t",
                "value": 112
            },
            {
                "label": "a",
                "name": "a",
                "description": "AAA",
                "type": "FLOAT",
                "rw": "R",
                "unit": "ºC",
                "value": 312
            }
        ]
    }
}
```