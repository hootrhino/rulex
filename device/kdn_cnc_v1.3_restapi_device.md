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

# 凯帝恩CNC

KND 部分 CNC 系统上运行了 REST API 服务器，用于向第三方开放部分数据接口。服务运行于 HTTP 标准端口（即 80 端口），当前最新版本为 v1.3， 建议使用 v1.3， 请求的基地址为/api/v1.3，兼容 v1.2, 若使用 V1.2,请求基地址为/api/v1.2（以下接口文档中若无特殊说明，均支持v1.2 和 v1.3）。
例如 CNC 的 IP 地址为 `192.168.1.101`，那么访问 CNC 的/status 接口应使用地址：v1.3 版本 `http://192.168.1.101/api/v1.3/status`；v1.2 版本 `http://192.168.1.101/api/v1.2/status`所有接口只接收 Content-Type 是 `application/json`类型的 HTTP 数据，除部分文件相关的接口外，大部分接口返回的也是 `application/json` 类型的 HTTP 数据。所有数据的编码必须是 UTF-8以下说明中，如未加特殊说明，均表示 HTTP 方法为 GET.

# 接口
## 获取当前状态
- GET: `http://192.168.1.101/api/v1.3/status`
```json
{
          "id":  15583,
          "type":  "K2000TCi_",
          "manufacturer":  "KND",
          "manufacture-time":  "20150710 133910",
          "cnc-type":  "M",
          "cnc-name":  "我的铣床",
          "soft-version":  "K2000TCi_A04_V4.2.00b",
          "fpga-version":  "2111_140606",
          "ladder-version":  "ktc_std_3.0_20160331",
          "nc-axes":  [
                    "X",
                    "Z",
                    "Y"
          ],
          "nc-relative-axes":  [
                    "X",
                    "Z",
                    "Y"
          ],
          "axes":  [
                    "X",
                    "Z",
                    "Y"
          ],
          "relative-axes":  [
                    "X",
                    "Z",
                    "Y"
          ]
}
```

## 错误信息
```json
{
 "error": -1, // 错误码
 "error-message": "未知错误" // 错误的消息
}
```