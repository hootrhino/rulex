---@diagnostic disable: undefined-global
-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end

--
-- {
--     "tag":"data",
--     "function":3,
--     "address":0,
--     "quantity":4,
--     "value":"0298010d"
-- }
-- -------------------------------------------------------
-- | 状态码 | 数据长度 |    湿度   |    温度   | CRC校验  |
-- -------------------------------------------------------
-- | 0x01   |  0x04   | 0x00 0x00  | 0x00 0x00 |0x00 0x00|
-- -------------------------------------------------------
Actions = {function(data)
    local table = rulexlib:J2T(data)
    local value = table['value']
    local t = rulexlib:HsubToN(value, 5, 8)
    local h = rulexlib:HsubToN(value, 0, 4)
    local t1 = rulexlib:HToN(string.sub(value, 5, 8))
    local h2 = rulexlib:HToN(string.sub(value, 0, 4))
    print('Data=> ', rulexlib:T2J({
        Device = "TH00000001",
        Ts = rulexlib:TsUnix(),
        T = t,
        H = h,
        T1 = t1,
        H2 = h2
    }))
    return true, data
end}
