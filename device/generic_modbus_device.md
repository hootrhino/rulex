## 简介

通用 Modbus 资源，可以用来实现常见的 modbus 协议寄存器读写等功能。

## 配置

```json
{
    "name": "GENERIC_MODBUS",
    "type": "GENERIC_MODBUS",
    "actionScript": "",
    "description": "GENERIC_MODBUS",
    "config": {
        "frequency": 100,
        "autoRequest": true,
        "mode": "RTU",
        "timeout": 5,
        "config": {
            "baudRate": 9600,
            "dataBits": 8,
            "ip": "127.0.0.1",
            "parity": "N",
            "port": 502,
            "stopBits": 1,
            "uart": "COM5"
        },
        "registers": [
            {
                "slaverId": 1,
                "function": 3,
                "address": 3,
                "quantity": 3,
                "tag": "d1"
            },
            {
                "slaverId": 2,
                "function": 3,
                "address": 3,
                "quantity": 3,
                "tag": "d2"
            }
        ]
    }
}
```

## 数据

### 读

#### 1. CMD

无指令，直接读取出来的数据是一个数组:

```go
{
    "d1":{
        "tag":"d1",
        "weight":0,
        "initValue":0,
        "function":3,
        "slaverId":1,
        "address":0,
        "quantity":2,
        "value":""
    },
    "d2":{
        "tag":"d2",
        "weight":0,
        "initValue":0,
        "function":3,
        "slaverId":1,
        "address":0,
        "quantity":2,
        "value":""
    }
}
```

#### 2. Args

无参数

#### 3. 数据样例

```json
[
    {
        "tag" :"t1",
        "function" :3,
        "slaverId" :1,
        "address" :0,
        "quantity" :1,
        "value" : 100
    }
]
```

### 写

#### 1. CMD

无指令

#### 2. Args

```go
[
    {
        tag      string
        function int
        slaverId byte
        address  uint16
        quantity uint16
        value    []byte
    }
]
```

#### 3. 数据样例

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

## 案例

```lua
-- Success
function Success()
    print("Success")
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions = {
    function(data)
    local json = [
            {
                "tag" :"t1",
                "function" :3,
                "slaverId" :1,
                "address" :0,
                "quantity" :1,
                "value" : 100
            }
        ]
        local _, err = rulexlib:WriteDevice(device, rulexlib:T2J(json))
        if (err ~= nil) then
            log('WriteDevice open err: ', err)
            return false, data
        end
        return true, data
    end
}

```

## 说明

仅仅是个通用 modbus 处理器。

## 社区

### 维护者

- wwhai
