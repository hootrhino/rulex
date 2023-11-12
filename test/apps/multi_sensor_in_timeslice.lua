--
--                       ModbusRS-485 BUS
-- ─────────▲──────────────▲─────────────▲─────────────▲─────────►
--          │              │             │             │
--   ┌──────┴─────┐ ┌──────┴────┐ ┌──────┴─────┐ ┌─────┴─────┐
--   │     PH     │ │    O2     │ │     TEMP   │ │    SALT   │
--   └────────────┘ └───────────┘ └────────────┘ └───────────┘
--     Addr=01          Addr=02        Addr=03      Addr=04
--
-- 一根总线上接了4个传感器，协议全部是Modbus，但是寄存器定位不同，数据长短不同;
-- 本案例就是展示如何在一根总线上挂不同类型的设备，并且做数据解析
-- 如果你没有设备，可以用test\data\吓的 Modbus Slaver(v6.3) 配置 4device_in_one_bus.mbs 测试。
AppNAME = "多传感器一总线案例"
AppVERSION = "1.0.0"
AppDESCRIPTION = "多传感器一总线案例"
--
-- Main
--
local Devices = {
    -- PH传感器有1个寄存器，从0开始计算
    PH = {
        Req = "010300000001840A",
        Decode = function(Hexs)
            local Value, error = applib:MatchUInt("ph:[0,2]", Hexs)
            if error ~= nil then
                return 0
            end
            return Value
        end
    },
    -- 溶解氧传感器有4个寄存器，从0开始计算
    -- 溶解氧(2Word) 温度(2Word) 电流(2Word) 警告(1Word)
    O2 = {
        Req = "020300000004443A",
        Decode = function(Hexs)
            local Value, error = applib:MatchUInt("oxygen:[0,2];temperature:[2,4];current:[4,6];warn:[6,7]", Hexs)
            if error ~= nil then
                return 0
            end
            return Value
        end
    },
    -- 水温传感器有1个寄存器，从0开始计算
    TEMP = {
        Req = "03030000000185E8",
        Decode = function(Hexs)
            local Value, error = applib:MatchUInt("temperature:[0,2]", Hexs)
            if error ~= nil then
                return 0
            end
            return Value
        end
    },
    -- 盐度传感器有5个寄存器，从0开始计算
    -- 分别是：电导率(2Word) 电阻率(2Word) 水温(2Word) DTS(2Word) 盐度(2Word)
    SALT = {
        Req = "040300000005859C",
        Decode = function(Hexs)
            local Value, error = applib:MatchUInt("c1:[0,2];c2:[2,4];temperature:[4,6];dts:[6,8];dts:[8,10]", Hexs)
            if error ~= nil then
                return 0
            end
            return Value
        end
    },
}
function Main(arg)
    while true do
        for K, Device in pairs(Devices) do
            local response, err = applib:CtrlDevice("DEVICEaqN8P7Ngabs4Y5dVgXXzK7", Device.Req)
            if err ~= nil then
                print(err)
            else
                local parsedValue = Device.Decode(response)
                print('Parsed Value ===>', K, '=', parsedValue)
            end
        end
        time:Sleep(100)
    end
end
