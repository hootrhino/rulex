## 简介
PLC S1200是西门子（Siemens）公司的一种工业可编程逻辑控制器（Programmable Logic Controller，PLC）设备系列。PLC S1200系列广泛应用于自动化控制领域，用于监控和控制各种工业过程。

PLC（可编程逻辑控制器）是一种专用的计算机控制设备，广泛应用于工业自动化领域。PLC通过输入信号的检测和输出信号的控制，实现对生产过程的监控和控制。PLC S1200系列具有高性能和可靠性，适用于中小型自动化系统。

PLC S1200系列设备通常用于控制和监视各种工业过程，例如生产线、机械设备、自动化装置、流程控制系统等。它们可以与传感器、执行器和其他设备进行通信，通过编程逻辑来处理输入信号，并输出相应的控制信号。

PLC S1200系列设备提供了丰富的输入输出接口、通信接口和编程功能，以满足各种自动化控制需求。通过编程，用户可以定义逻辑控制规则、配置输入输出映射、实现数据处理和通信功能等。

需要注意的是，PLC S1200是西门子（Siemens）公司的商标产品，更详细的信息和技术规格可以参考西门子官方文档或与其联系。

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