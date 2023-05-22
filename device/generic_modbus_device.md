# Modbus 通用采集器

## 简介

通用 Modbus 资源，可以用来实现常见的 modbus 协议寄存器读写等功能。

## 配置

```json
{
    "name": "485温湿度传感器",
    "type": "GENERIC_MODBUS",
    "config": {
        "commonConfig": {
            "frequency": 2000,
            "retryTime": 5,
            "autoRequest": true,
            "mode": "RTU",
            "timeout": 20
        },
        "rtuConfig": {
            "timeout": 30,
            "baudRate": 4800,
            "dataBits": 8,
            "parity": "N",
            "stopBits": 1,
            "uart": "COM5"
        },
        "registers": [
            {
                "weight": 1,
                "initValue": 0,
                "slaverId": 2,
                "address": 0,
                "quantity": 2,
                "tag": "t2",
                "function": 3,
                "value": ""
            },
            {
                "weight": 1,
                "initValue": 0,
                "slaverId": 1,
                "address": 0,
                "quantity": 2,
                "tag": "t1",
                "function": 3,
                "value": ""
            }
        ]
    }
}
```

## 数据样例

```json
{
    "d1":{
        "tag":"d1",
        "function":3,
        "slaverId":1,
        "address":0,
        "quantity":6,
        "value":"0117011d0127011a0110010e"
    }
}
```

- value: 十六进制字符串

## 常用函数

为了更加清楚的描述接口的使用，下面给出数据解析详细示例：

1. 打印日志

    ```lua
    Actions =
    {
        function(data)
            print(data)
            return true, data
        end
    }

    ```

2. 推送到Mongodb

    ```lua
    Actions =
    {
        function(data)
            local dataTable = rulexlib:J2T(data)
            local Value = dataTable['value']
            rulexlib:DataToMongo('uuid', Value)
            return true, data
        end
    }

    ```

多多益善。

## 维护

开源参与者需要给出维护作者的邮箱，方便及时处理问题。

- <xxx@xxx.com>
