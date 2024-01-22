<!--
 Copyright (C) 2024 wwhai

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

# 硬件资源表
这里主要用来处理硬件接口初始化，比如网卡，GPIO、串口、蓝牙等。现阶段因为前期历史原因导致部分硬件初始化操作放在HttpServer里面了，不符合模块化设计。后期需要逐步迁移到硬件管理器。

## 原理
Rulex的配置数据库会有一些初始化数据，每次启动的时候，在此处初始化硬件即可。
```go
//...
NetConfig = db.query("network")
Network.Init(NetConfig)

GpioConfig = db.query("gpio")
Gpio.Init(GpioConfig)

BleConfig = db.query("ble")
Ble.Init(BleConfig)
//...
```