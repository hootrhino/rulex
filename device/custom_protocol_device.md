# 自定义协议
该特性用于自定义协议场景下使用。例如一些私有TCP、UDP等场景下。现阶段暂时支持串口，未来随着完善会持续增加TCP、UDP支持。
假设一个总线上面挂了很多不一样的设备，此时要互相操作，也可以使用该特性。

## 配置
协议分静态协议和动态协议，下面是动态协议示例，一般会有至少一个自定义协议，关键字段是 `deviceConfig` ，下面给出一个 JSON 样例:

### 动态协议

```json
{
    "name":"GENERIC_PROTOCOL",
    "type":"GENERIC_PROTOCOL",
    "description":"GENERIC_PROTOCOL",
    "config":{
        "commonConfig":{
            "transport":"rs485rawserial",
            "retryTime":5,
            "frequency":100
        },
        "uartConfig":{
            "timeout":1000,
            "baudRate":9600,
            "dataBits":8,
            "parity":"N",
            "stopBits":1,
            "uart":"COM5"
        }
    }
}
```

## 字段：

- name: 协议的名称, 通常代表某个设备的功能，比如读数据，开关之类的
- type: 1-静态；2-动态, 在动态协议里面必须为2
- description: 协议的一些备注信息
- transport: 传输形式，目前支持 `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`

## 设备数据处理
```lua
-- 动态协议请求
AppNAME = 'Read'
AppVERSION = '0.0.1'
function Main(arg)
    local Id = 'DEVICE056b93901b3b4a5b9a3d69d14dc1139f'
    while true do
        local result, err = applib:CtrlDevice(Id, "010300000002C40B")
        --result {"in":"010300000002C40B","out":"010304000100022a32"}
        print("CtrlDevice result=>", result)
        applib:Sleep(60)
    end
    return 0
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