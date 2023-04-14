# 自定义协议
该特性用于自定义协议场景下使用。例如一些私有TCP、UDP等场景下。现阶段暂时支持串口，未来随着完善会持续增加TCP、UDP支持。
假设一个总线上面挂了很多不一样的设备，此时要互相操作，也可以使用该特性。

## 配置
一般会有至少一个自定义协议，关键字段是 `deviceConfig` ，下面给出一个 JSON 样例:

```json
{
    "name": "GENERIC_PROTOCOL",
    "type": "GENERIC_PROTOCOL",
    "description": "GENERIC_PROTOCOL",
    "config": {
        "commonConfig": {
            "transport": "rs485rawserial",
            "retryTime": 5,
            "waitTime": 60
        },
        "uartConfig": {
            "timeout": 1000,
            "baudRate": 9600,
            "dataBits": 8,
            "parity": "N",
            "stopBits": 1,
            "uart": "COM15"
        },
        "deviceConfig": {
            "read": {
                "autoRequest": false,
                "autoRequestGap": 60,
                "bufferSize": 9,
                "checkAlgorithm": "NONECHECK",
                "onCheckError": "IGNORE",
                "checksumBegin": 0,
                "checksumEnd": 6,
                "checksumValuePos": 6,
                "description": "read",
                "name": "read",
                "rw": 1,
                "protocol": {
                    "in": "D400070A0001D8AA55",
                    "out": "D400070A0001D8AA55"
                }
            },
            "write": {
                "autoRequest": false,
                "autoRequestGap": 60,
                "bufferSize": 9,
                "checkAlgorithm": "NONECHECK",
                "onCheckError": "IGNORE",
                "checksumBegin": 0,
                "checksumEnd": 6,
                "checksumValuePos": 6,
                "description": "write",
                "name": "write",
                "rw": 2,
                "protocol": {
                    "in": "D400070A0001D8AA55",
                    "out": "D400070A0001D8AA55"
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
- algorithm：校验算法，支持 `XOR`,`CRC16`,`NONECHECK`, 默认为不认证，即`NONECHECK`
- rw: 读写权限，有3个值: 1:RO 2:WO 3:RW
- bufferSize: 期望读取字节数，这个比较重要，最好是经过精确计算
- checksumBegin: 从哪个位置开始校验
- checksumEnd: 从哪里结束校验
- checksumValuePos: 校验对比值的位置
- protocol: 协议本体
    - in: 请求参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55
    - out: 返回参数, 用大写十六进制表示法，否则会解析失败, 例如：FFFFFF014CB2AA55， 这个参数一般不参与业务，主要用来返回取值。

实际上观察一下样例 JSON 就知道怎么配置了。

## 设备数据读取
```lua
-- get_uuid 就是配置的读指令
local binary1, err1 = applib:ReadDevice("ID12345", "get_uuid")
if err1 ~= nil then
    print(err1)
end
```
## 设备数据写入
```lua
-- ctrl1 就是配置的写指令
-- 参数：id, 指令，参数
local binary1, err1 = applib:WriteDevice("ID12345", "ctrl1", "args")
if err1 ~= nil then
    print(err1)
end
```
注意:**数据在总线形式下并发读写是有独占机制，这里加了锁来处理**
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
- applib:MatchHex 提取十六进制
  ```lua
     -- 第一个参数为提取表达式
     -- 格式为: "name1:[start, end];name2:[start, end]···"
      AppNAME = 'applib:MatchHex'
      AppVERSION = '0.0.1'
      function Main(arg)
          -- 十六进制提取器
          local MatchHexS = applib:MatchHex("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")
          for key, value in pairs(MatchHexS) do
              print('applib:MatchHex', key, value)
          end
          return 0
      end
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
- hex:ABCD 十六进制字节序按照 ABCD 顺序调整
  ```lua
		local V, err = rulexlib:ABCD("AABBCCDDEEFF")
    -- V = FFEEDDCCBBAA
  ```

- hex:DCBA 十六进制字节序按照 DCBA 顺序调整
  ```lua
		local V, err = rulexlib:DCBA("FFEEDDCCBBAA")
    -- V = AABBCCDDEEFF
  ```

- hex:BADC 十六进制字节序按照 BADC 顺序调整
  ```lua
		local V, err = rulexlib:BADC("CDAB12EF")
    -- V = ABCDEF12
  ```

- hex:CDAB 十六进制字节序按照 CDAB 顺序调整
  ```lua
		local V, err = rulexlib:CDAB("ABCDEF12")
    -- V = CDAB12EF
  ```