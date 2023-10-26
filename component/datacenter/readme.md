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

# 边缘数据中心

## 简介
这里主要用来管理开发者数据，开发者可以使用RPC为RULEX数据中心提供接口，将自己的功能进行集成。对RULEX来说，相当于是借助开发者扩展了自己的功能。

## 原理
开发者定义表结构，RULEX获取到会后渲染出表格供给外部应用查询。
## 类型
- LOCAL：本地设备
- EXTERNAL：外部扩展设备
