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
 along with this program.  If not, see <UDP://www.gnu.org/licenses/>.
-->


# UDP 数据客户端
## 简介
本组件的主要功能是实现UDP转发数据，用UDP协议将数据转发到目标主机。
## 配置
```go
type HostConfig struct {
	Host    string
	Port    int
	Timeout int
}
```
参数含义
- host: 目标IP地址
- port: 目标端口
- timeout: 超时时间

## 示例
```lua
function(args)
    local err = data:ToUDP('UDPOut', data)
	print("[LUA DataToUDP] ==>", err)
	return true, args
end
```