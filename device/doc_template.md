# OPCUA 采集协议

## 介绍

简介。。。。。

## 配置

这是设备启动所需的必要配置, 以OPCUA的为例:

```json
{
    "commonConfig":{
        "endpoint":"opc.tcp://NOAH:53530/OPCUA/SimulationServer",
        "policy":"POLICY_BASIC128RSA15",
        "mode":"MODE_SIGN",
        "auth":"AUTH_ANONYMOUS",
        "username":"1",
        "password":"1",
        "timeout":10,
        "frequency":500,
        "retryTime":10
    },
    "opcuaNodes":[
        {
            "tag":"node1",
            "description":"node 1",
            "nodeId":"ns=3;i=1013",
            "dataType":"String",
            "value":""
        },
        {
            "tag":"node2",
            "description":"node 2",
            "nodeId":"ns=3;i=1001",
            "dataType":"String",
            "value":""
        }
    ]
}
```

## 字段说明

给出上面配置里面出现的字段的详细说明

### commonConfig

| 字段名    | 类型   | 必填 | 说明                                                                                                                  |
| --------- | ------ | ---- | --------------------------------------------------------------------------------------------------------------------- |
| endpoint  | string | √    | OPC UA服务端的地址，对应gopcua库中的ClientEndpoint                                                                    |
| policy    | Enum   | √    | 安全策略URL，可以是None、Basic128Rsa15、Basic256、Basic256Sha256中的任意一个，对应gopcua库中的ClientSecurityPolicyURI |
| mode      | Enum   | √    | 安全模式，可以是None、Sign、SignAndEncrypt中的任意一个，对应gopcua库中的ClientSecurityMode                            |
| auth      | Enum   | √    | 认证模式，可以是Anonymous、UserName、Certificate中的任意一个，对应gopcua库中的ClientAuthMode                          |
| username  | string | √    | 用户名，对应gopcua库中的ClientUsername                                                                                |
| password  | string | √    | 密码，对应gopcua库中的ClientPassword                                                                                  |
| timeout   | string | √    | 超时时间，对应gopcua库中的RequestTimeout (毫秒)                                                                       |
| frequency | string | √    | 采集频率(毫秒)                                                                                                        |
| retryTime | string | √    | 出错重试次数                                                                                                          |

## 设备数据读取

如果设备支持读, 给出一个读出来的数据示例

```json
{
    "d1":{
        "tag":"d1",
        "function":3,
        "slaverId":1,
        "address":0,
        "quantity":2,
        "value":"000a0b0c0d"
    }
}
```

## 设备数据写入

如果设备支持写, 给出一个写入的数据示例

```lua
{
    "function":3,
    "slaverId":2,
    "address":0,
    "quantity":2,
    "value":"..."
}
```

## 常用函数

为了更加清楚的描述接口的使用，下面给出数据解析详细示例：

1. 打印日志

    ```lua
    Actions =
    {
        function(args)
            return true, args
        end
    }

    ```

2. 推送到Mongodb

    ```lua
    Actions =
    {
        function(args)
            return true, args
        end
    }

    ```

多多益善。

## 维护

开源参与者需要给出维护作者的邮箱，方便及时处理问题。

- <xxx@xxx.com>
