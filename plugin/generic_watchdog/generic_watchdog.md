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

# 通用看门狗

## 简介
用于监视应用程序或系统的运行状态，并在检测到问题时采取预定的操作。这有助于确保系统或应用程序在异常情况下能够自动恢复或采取适当的措施，以防止系统崩溃或无响应。

## 注意
该插件只针对有看门狗的硬件使用，如果没有看门狗，该插件无效。一般来说存在`/dev/watchdog`设备的硬件都应该支持。