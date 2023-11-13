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

# TCP 客户端
主要用来透传数据。
## 配置
直接透传模式：
```json
{
    "type": "TCP_TRANSPORT",
    "name": "TCP_TRANSPORT Server",
    "description": "TCP_TRANSPORT Server",
    "config": {
        "commonConfig": {
            "dataMode": "RAW_STRING",
            "allowPing": true,
            "pingPacket": "PING"
        },
        "hostConfig": {
            "host": "127.0.0.1",
            "port": 6005,
            "timeout": 3000
        }
    }
}
```

十六进制透传模式：
```json
{
    "type": "TCP_TRANSPORT",
    "name": "TCP_TRANSPORT Server",
    "description": "TCP_TRANSPORT Server",
    "config": {
        "commonConfig": {
            "dataMode": "HEX_STRING",
            "allowPing": true,
            "pingPacket": "PING"
        },
        "hostConfig": {
            "host": "127.0.0.1",
            "port": 6005,
            "timeout": 3000
        }
    }
}
```
## 脚本示例
```lua
function Main(arg)
	stdlib:Debug("Hello World:" .. time:Time())
	while true do
		local err1 = data:ToTcp('OUTD8HC47', "Hello World:" .. time:Time())
		if err1 ~= nil then
			stdlib:Debug(err)
		end
		time:Sleep(1000)
	end
	return 0
end
```