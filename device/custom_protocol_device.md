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
            "transport": "rs485rawserial",
            "waitTime": 10
        },
        "uartConfig": {
            "timeout": 10,
            "baudRate": 9600,
            "dataBits": 8,
            "parity": "N",
            "stopBits": 1,
            "uart": "COM3"
        },
        "deviceConfig": {
            "1": {
                "autoRequest": true,
                "autoRequestGap": 5000,
                "bufferSize": 7,
                "checksum": "CRC16",
                "checksumBegin": 1,
                "checksumEnd": 2,
                "checksumValuePos": 7,
                "description": "获取UUID",
                "name": "get_uuid",
                "protocol": {
                    "in": "FFFFFF014CB2AA55",
                    "out": "FFFFFF014CB2AA55"
                },
                "rw": 1,
                "timeout": 1000
            }
        }
    }
}
```

## 字段：

- name: 协议的名称, 通常代表某个设备的功能，比如读数据，开关之类的
- description: 协议的一些备注信息
- transport: 传输形式，目前支持 `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
- algorithm：校验算法，支持 `XOR`,`CRC16`,`NONECHECK`, 默认为不认证，即`NONECHECK`
- protocol: 协议本体
    - in: 请求参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55
    - out: 返回参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55， 这个参数一般不参与业务，主要用来做个demo对比用。

实际上观察一下样例 JSON 就知道怎么配置了。

## 常用函数
下面列出一些常用的函数:

- hex:Bytes2Hexs: 字节转成16进制字符串
  ```lua
     local s, err = hex:Bytes2Hexs({1,2,3})
     -- s 是lua的字符串: 'FFFFFF014CB2AA55'
  ```
- hex:Hexs2Bytes: 16进制字符串转成字节
  ```lua
     local b, err = hex:Hexs2Bytes('FFFFFF014CB2AA55')
     -- b 是一个table: {0 = 0, 1 = 1}
  ```
- eekit:GPIOSet: 控制GPIO
  ```lua
     local err = eekit:GPIOSet(6, 1)
  ```
- eekit:GPIOGet 16进制字符串转成字节
  ```lua
     local value, err = eekit:GPIOGet(6)
     -- value 的值为 0 或者 1
  ```
- applib:GPIOGet 提取十六进制
  ```lua
     -- 第一个参数为提取表达式
     -- 格式为: "name1:[start, end];name2:[start, end]···"
     local MatchHexS = rulexlib:MatchHex("name1:[1,3];name2:[4,5]", "FFFFFF014CB2AA55")
  ```
- applib:MB 二进制匹匹配, 返回值为二进制的字符串表示法
  ```lua
     -- 第一个参数为提取表达式
     -- 格式为: [<|> K1:LEN1 K2:LEN2... ], 返回一个K-V table
		local V6 = rulexlib:T2J(rulexlib:MB("<a:5 b:3 c:1", "aab", false))

  ```
- applib:MBHex 二进制匹匹配, 返回值为二进制的十六进制表示法
  ```lua
     -- 第一个参数为提取表达式
     -- 格式为: [<|> K1:LEN1 K2:LEN2... ], 返回一个K-V table
		local V6 = rulexlib:T2J(rulexlib:MBHex("<a:5 b:3 c:1", "aab", false))

  ```