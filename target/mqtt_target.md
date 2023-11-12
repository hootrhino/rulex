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

# MQTT 客户端
## 简介
本组件实现将数据写入MQTT Server。
## 配置
```go
type MqttConfig struct {
	Host     string `json:"host" validate:"required" title:"服务地址"`
	Port     int    `json:"port" validate:"required" title:"服务端口"`
	ClientId string `json:"clientId" validate:"required" title:"客户端ID"`
	Username string `json:"username" validate:"required" title:"连接账户"`
	Password string `json:"password" validate:"required" title:"连接密码"`
	PubTopic string `json:"pubTopic" title:"上报TOPIC" info:"上报TOPIC"` // 上报数据的 Topic
	SubTopic string `json:"subTopic" title:"订阅TOPIC" info:"订阅TOPIC"` // 上报数据的 Topic
}
```
字段解释
- Host: MQTT 主机
- Port: MQTT 端口
- ClientId: MQTT ClientId
- Username: MQTT Username
- Password: MQTT Password
- PubTopic: MQTT 上行Topic
- SubTopic: MQTT 下行Topic

## 示例
```lua
function(args)
    local err = data:ToMqtt('MqttOut', data)
	print("[LUA DataToMqtt] ==>", err)
	return true, args
end
```