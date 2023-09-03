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

# Modbus CRC 计算器
## 简介
Modbus CRC（Cyclic Redundancy Check）是一种校验和算法，用于检测数据传输中的错误。它常用于工业自动化和通信领域，特别是在 Modbus 通信协议中。

CRC 是一种循环冗余校验，通过对数据进行多项式运算来生成一个固定长度的校验值，通常是 16 位。这个校验值随数据一起传输，并在接收端再次计算，然后与接收到的校验值进行比较，以确定数据是否在传输过程中发生了错误。

Modbus CRC 使用的多项式是 0xA001（或者说是 x^16 + x^15 + x^2 + 1），这是一种经过优化的多项式，适用于 Modbus 通信。发送端计算数据的 CRC 值并将其附加到数据帧中，接收端接收数据后也会计算 CRC 值，然后将其与接收到的 CRC 值进行比较。如果两个 CRC 值不匹配，就表示数据在传输过程中出现了错误。

Modbus CRC 的主要目的是检测数据传输中的错误，以确保数据的完整性和准确性。它是一种简单但有效的校验方法，广泛用于工业控制和自动化系统中的通信。

## 支持指令
### 计算大端CRC
请求：
```json
{
    "uuid": "MODBUS_CRC_CALCULATOR",
    "name": "crc16big",
    "args": "000102030405"
}
```
返回：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "value": "840a"
    }
}
```
### 计算小端CRC
请求：
```json
{
    "uuid": "MODBUS_CRC_CALCULATOR",
    "name": "crc16little",
    "args": "000102030405"
}
```
返回：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "value": "840a"
    }
}
```