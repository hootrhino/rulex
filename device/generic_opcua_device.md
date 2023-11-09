# OPCUA 采集协议

## 介绍

OPCUA采集协议是基于OPCUA协议的一种采集协议，主要用于采集OPCUA协议的数据，前支持的功能有：

- 读取数据-nodeId方式
- 认证模式：用户名密码认证，匿名
- 消息模式 ：支持 encrypt，sign，none
写数据功能暂未实现，后续会陆续支持。

## 配置

一般会有至少一个自定义协议，关键字段是 `opcuanodes` ，下面给出一个 JSON 样例:

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
    "opcuanodes":[
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

### commonConfig

该参数为通用配置。

| 字段名  |  类型  | 必填  | 说明  |
| ------------ | ------------ | ------------ | ------------ |
| endpoint  |  string  |  √  | OPC UA服务端的地址，对应gopcua库中的ClientEndpoint  |
| policy  |  Enum  |  √  | 安全策略URL，可以是None、Basic128Rsa15、Basic256、Basic256Sha256中的任意一个，对应gopcua库中的ClientSecurityPolicyURI  |
| mode  |  Enum  |  √  | 安全模式，可以是None、Sign、SignAndEncrypt中的任意一个，对应gopcua库中的ClientSecurityMode |
| auth  |  Enum  |  √  | 认证模式，可以是Anonymous、UserName、Certificate中的任意一个，对应gopcua库中的ClientAuthMode |
| username  |  string  |  √  | 用户名，对应gopcua库中的ClientUsername |
| password  |  string  |  √  |密码，对应gopcua库中的ClientPassword |
| timeout  |  string  |  √  | 超时时间，对应gopcua库中的RequestTimeout (毫秒)  |
| frequency  |  string  |  √  | 采集频率(毫秒)  |
| retryTime  |  string  |  √  | 出错重试次数  |

### opcuanodes

该参数为点表字段。

| 字段名  |  类型  | 必填  | 说明  |
| ------------ | ------------ | ------------ | ------------ |
| Tag  |  string  |  √  | 点表的标签，用于标识点表  |
| NodeId  |  string  |  √  | 点表的NodeId，对应gopcua库中的NodeId  |
| DataType  |  string  |  √  | 点表的数据类型，对应gopcua库中的VariantType  |
| Value  |  string  |  √  | 点表的值，对应gopcua库中的Variant  |
| Description  |  string  |  ×  | 点表的描述，用于描述点表  |

## 设备数据读取

```lua
{
 -- 读到的数据
}
```

## 设备数据写入

```lua
{
   -- 读到的数据
}
```

## 常用函数

下面给出数据解析示例：

### 打印日志

```lua
Actions =
{
    function(args)
        return true, args
    end
}

```

### 推送到Mongodb

```lua
Actions =
{
    function(args)
        return true, args
    end
}

```

## 维护

- <xxx@xxx.com>
