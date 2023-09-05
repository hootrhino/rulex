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

# Modbus地址扫描仪
主要用来辅助判断总线设备是否异常，以及检查设备可能的ID号。
## 配置
```json
{
    "uuid": "MODBUS_SCANNER",
    "name": "scan",
    "args": "{\"uart\":\"COM4\",\"dataBits\":8,\"parity\":\"N\",\"stopBits\":1,\"baudRate\":4800,\"timeout\":1000}"
}
```
注意：args 是JSON字符串，其本质上是个串口配置。

## 使用
选择串口配置以后，会持续向整个地址范围发送一个探针包，当收到设备的回信时，评估其存在的可能性。
## 支持服务
### 开启扫描
请求:
```json
{
    "uuid": "MODBUS_SCANNER",
    "name": "scan",
    "args": "{\"uart\":\"/dev/ttyS1\",\"dataBits\":8,\"parity\":\"N\",\"stopBits\":1,\"baudRate\":4800,\"timeout\":1000}"
}
```
返回：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {}
}
```
### 停止扫描
请求：
```json
{
    "uuid": "MODBUS_SCANNER",
    "name": "stop",
    "args": ""
}
```
返回：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {}
}
```
## 输出
日志输出topic：`plugin/ModbusScanner/MODBUS_SCANNER`。
## 注意
- 串口配置需要正确选择。
- http请求的结果仅仅作为请求成功与否的依据，而扫描的结果则需要看日志。