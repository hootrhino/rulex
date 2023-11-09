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

# HTTP 数据客户端
## 简介
本组件的主要功能是实现HTTP请求，将数据以POST的形式发送到目标HTTP接口。
## 配置
```go
type HTTPConfig struct {
	Url     string            `json:"url" title:"URL"`
	Headers map[string]string `json:"headers" title:"HTTP Headers"`
}
```
参数含义
- url: 请求地址
- headers：HTTP请求头，是K-V键值对

## 示例
```lua
function(args)
    local err = data:ToHttp('HttpOut', data)
	print("[LUA DataToHttp] ==>", err)
	return true, args
end
```