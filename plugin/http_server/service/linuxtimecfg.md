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

`timedatectl` 是用于管理系统时间和日期设置的命令行工具。它通常在Linux系统上用于配置和查看系统时间、日期、时区以及NTP时间同步设置。以下是一个简单的 `timedatectl` 教程，帮助你入门：

**1. 查看系统当前时间和日期：**

要查看系统当前的时间和日期，只需运行以下命令：

```bash
timedatectl
```

这将显示当前的时间、日期、时区以及NTP时间同步状态。

**2. 查看系统时区：**

要查看系统当前的时区设置，运行以下命令：

```bash
timedatectl show --property=Timezone --value
```

这将显示当前的时区名称（例如："Asia/Shanghai"）。

**3. 设置系统时区：**

如果需要更改系统的时区，可以运行以下命令：

```bash
sudo timedatectl set-timezone Your_Timezone
```

将 `Your_Timezone` 替换为所需的时区名称。例如，要将时区设置为"Asia/Shanghai"，可以运行：

```bash
sudo timedatectl set-timezone Asia/Shanghai
```

**4. 启用/禁用NTP时间同步：**

要启用NTP时间同步，可以运行以下命令：

```bash
sudo timedatectl set-ntp true
```

要禁用NTP时间同步，可以运行：

```bash
sudo timedatectl set-ntp false
```

启用NTP时间同步会自动将系统时间与NTP服务器同步。

**5. 手动设置系统时间：**

你还可以手动设置系统的时间和日期。例如，要将系统时间设置为2023年9月7日下午3点30分，可以运行以下命令：

```bash
sudo timedatectl set-time "2023-09-07 15:30:00"
```

**6. 检查时钟是否被校准：**

要检查系统时钟是否被校准，可以运行以下命令：

```bash
timedatectl show --property=NTPSynchronized --value
```

如果值为 "yes"，表示时钟已被校准。

这些是一些基本的 `timedatectl` 命令和用法。你可以通过运行 `man timedatectl` 命令来查看更多详细的信息和选项，或者查阅相关文档以获取更多信息。