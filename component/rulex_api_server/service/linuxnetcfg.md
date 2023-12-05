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

`nmcli` 是 NetworkManager 的命令行工具，用于配置和管理网络连接。它允许你通过终端进行各种网络设置，包括配置无线网络、以太网连接、VPN、连接到无线热点等。以下是一个简单的`nmcli`教程，帮助你入门：

**1. 列出所有网络连接：**

要列出系统上所有的网络连接，可以使用以下命令：

```bash
nmcli connection show
```

这将显示所有配置的网络连接以及它们的状态。

**2. 查看连接详细信息：**

要查看特定连接的详细信息，可以运行以下命令，将`connection_name`替换为你要查看的连接名称：

```bash
nmcli connection show connection_name
```

这将显示该连接的详细信息，包括连接类型、IP地址、网关等。

**3. 创建新的网络连接：**

要创建新的网络连接，你可以运行以下命令，并按照提示输入连接的相关信息：

```bash
nmcli connection add type connection_type ifname interface_name
```

- `connection_type`：连接的类型，例如，`ethernet`（以太网连接）或 `wifi`（无线连接）。
- `interface_name`：网络接口的名称，例如，`eth0`（以太网接口）或 `wlan0`（无线接口）。

然后，根据你的连接类型和需求，输入其他相关信息，如IP地址、网关、DNS等。

**4. 修改现有的网络连接：**

要修改现有的网络连接，可以使用以下命令：

```bash
nmcli connection modify connection_name setting_name new_value
```

- `connection_name`：要修改的连接的名称。
- `setting_name`：要修改的设置名称，例如，`ipv4.method`（IPv4 配置方法）或 `wifi.ssid`（无线SSID）。
- `new_value`：新的设置值。

通过此命令，你可以修改各种连接设置，例如，启用或禁用自动连接、更改IP配置等。

**5. 连接到WiFi网络：**

要连接到一个WiFi网络，可以使用以下命令：

```bash
nmcli device wifi connect SSID password PASSWORD
```

- `SSID`：WiFi网络的名称。
- `PASSWORD`：WiFi网络的密码（如果需要）。

**6. 断开连接：**

要断开当前活动的连接，可以使用以下命令：

```bash
nmcli connection down connection_name
```

- `connection_name`：要断开的连接的名称。

**7. 启动连接：**

要启动已经配置但当前处于断开状态的连接，可以使用以下命令：

```bash
nmcli connection up connection_name
```

- `connection_name`：要启动的连接的名称。

这是一个基本的`nmcli`教程，帮助你入门。`nmcli`还有很多其他功能和选项，你可以通过运行 `nmcli --help` 来查看更多命令和用法。如果需要更详细的配置，请查阅相关文档或使用 `man nmcli` 命令来访问手册页。