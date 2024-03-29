## 简介
PLC S1200是西门子（Siemens）公司的一种工业可编程逻辑控制器（Programmable Logic Controller，PLC）设备系列。PLC S1200系列广泛应用于自动化控制领域，用于监控和控制各种工业过程。

PLC（可编程逻辑控制器）是一种专用的计算机控制设备，广泛应用于工业自动化领域。PLC通过输入信号的检测和输出信号的控制，实现对生产过程的监控和控制。PLC S1200系列具有高性能和可靠性，适用于中小型自动化系统。

PLC S1200系列设备通常用于控制和监视各种工业过程，例如生产线、机械设备、自动化装置、流程控制系统等。它们可以与传感器、执行器和其他设备进行通信，通过编程逻辑来处理输入信号，并输出相应的控制信号。

PLC S1200系列设备提供了丰富的输入输出接口、通信接口和编程功能，以满足各种自动化控制需求。通过编程，用户可以定义逻辑控制规则、配置输入输出映射、实现数据处理和通信功能等。

需要注意的是，PLC S1200是西门子（Siemens）公司的商标产品，更详细的信息和技术规格可以参考西门子官方文档或与其联系。
## 地址
在西门子PLC编程中，I（输入）和Q（输出）寄存器是非常常见的，它们分别用于存储输入信号和输出信号。以下是如何在西门子PLC中读写I和Q寄存器的格式：
1. I（输入）寄存器：
   - 读取I寄存器：`I` + 编号，例如 `I0.0`、`I1.1` 等。
   - 写入I寄存器：`I` + 编号，例如 `I0.0`、`I1.1` 等。
2. Q（输出）寄存器：
   - 读取Q寄存器：`Q` + 编号，例如 `Q0.0`、`Q1.1` 等。
   - 写入Q寄存器：`Q` + 编号，例如 `Q0.0`、`Q1.1` 等。
在西门子PLC中，I和Q寄存器通常以字节（8位）为单位进行操作，但也可以以字（16位）或双字（32位）为单位。例如，`I0.0` 表示输入字节0的第一个位，而 `IB0` 表示整个输入字节0（8位）。类似地，`Q0.0` 表示输出字节0的第一个位，而 `QB0` 表示整个输出字节0（8位）。
如果你需要读取或写入一个字（16位）的I或Q寄存器，你可以使用 `IW` 或 `QW` 来表示。例如，`IW0` 表示输入字0，而 `QW0` 表示输出字0。
西门子PLC的编程语言（例如Ladder Diagram、Function Block Diagram等）提供了丰富的指令来操作这些寄存器，包括读取、写入、逻辑操作、算术运算等。在实际编程中，你需要根据具体的指令集和编程环境来选择合适的指令和格式。

## 参数
```json
{
    "name": "SIEMENS_PLC",
    "type": "SIEMENS_PLC",
    "gid": "DROOT",
    "config": {
        "host": "127.0.0.1:1500",
        "rack": 0,
        "slot": 1,
        "model": "SIEMENS_PLC",
        "timeout": 1000,
        "autoRequest": true,
        "idleTimeout": 1000,
        "frequency": 1000,
        "blocks": [
            {
                "tag": "Value",
                "frequency": 1000,
                "type": "DB",
                "address": 1,
                "start": 100,
                "size": 16
            }
        ]
    },
    "description": "SIEMENS_PLC"
}
```
## 脚本示例
```lua

-- Actions
-- 采集到的数据:
-- {
--     "tag":"Value",
--     "type":"DB",
--     "frequency":0,
--     "address":1,
--     "start":100,
--     "size":16,
--     "value":"00000001000000020000000300000004"
-- }
Actions =
{
    function(args)
        local dataT, err = json:J2T(args)
        if (err ~= nil) then
            stdlib:Debug('parse json error:' .. err)
            return true, args
        end
        for key, value in pairs(dataT) do
            --data: 00000001000000020000000300000004
            local MatchHexS = hex:MatchUInt("a:[0,3];b:[4,7];c:[8,11];d:[12,15]", value['value'])
            local ts = time:Time()
            local Json = json:T2J(
                {
                    tag = key,
                    ts = ts,
                    a = MatchHexS['a'],
                    b = MatchHexS['b'],
                    c = MatchHexS['c'],
                    d = MatchHexS['d'],
                }
            )
            stdlib:Debug(Json)
        end
        return true, args
    end
}

```