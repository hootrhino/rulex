# 自定义协议
该特性用于自定义协议场景下使用。例如一些私有TCP、UDP等场景下。现阶段暂时支持串口，未来随着完善会持续增加TCP、UDP支持。

## 配置
一般会有至少一个自定义协议，用一个大 MAP 来表示，下面给出一个 JSON 样例:

```json
{
    "name": "GENERIC_PROTOCOL",
    "type": "GENERIC_PROTOCOL",
    "description": "GENERIC_PROTOCOL",
    "config": {
        "commonConfig": {
            "frequency": 5,
            "autoRequest": true,
            "transport": "rs485rawserial",
            "waitTime": 10,
            "timeout": 10
        },
        "uartConfig": {
            "baudRate": 9600,
            "dataBits": 8,
            "ip": "127.0.0.1",
            "parity": "N",
            "port": 502,
            "stopBits": 1,
            "uart": "COM4"
        },
        "deviceConfig": {
            "get_uuid": {
                "name": "get_uuid",
                "description": "获取UUID",
                "protocol": {
                    "in": "FFFFFF014CB2AA55",
                    "out": "FA0101CE34AA55"
                }
            }
        }
    }
}
```

## 字段：

- name: 协议的名称, 通常代表某个设备的功能，比如读数据，开关之类的
- description: 协议的一些备注信息
- transport: 传输形式，目前支持 `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
- protocol: 协议本体
    - in: 请求参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55
    - out: 返回参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55， 这个参数一般不参与业务，主要用来做个demo对比用。

实际上观察一下样例 JSON 就知道怎么配置了。

## 常用函数
- Bytes2Hexs: 字节转成16进制字符串
  ```lua
     local s, err = hex:Bytes2Hexs({1,2,3})
     -- s 是lua的字符串: 'FFFFFF014CB2AA55'
  ```
- Hexs2Bytes: 16进制字符串转成字节
  ```lua
     local b, err = hex:Hexs2Bytes('FFFFFF014CB2AA55')
     -- b 是一个table: {0 = 0, 1 = 1}
  ```