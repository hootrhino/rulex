---@diagnostic disable: undefined-global
-- {
--     "temp":{
--         "tag":"temp",
--         "weight":0,
--         "initValue":0,
--         "function":3,
--         "slaverId":1,
--         "address":0,
--         "quantity":1,
--         "value":"00fb"
--     }
-- }
--
-- 该款设备是一个瑞蒙德水温度探测器，其有一个Modbus寄存器(2字节)用来存储数据
-- 资料：http://www.remond.ltd/pr.jsp?_pp=0_464_13
--
Actions =
{
    function(args)
        local dataT, err = rulexlib:J2T(data)
        if (err ~= nil) then
            print('parse json error:', err)
            return true, args
        end
        for _, value in pairs(dataT) do
            local MatchHexS = rulexlib:MatchUInt("temp:[0,1]", value['value'])
            local ts = rulexlib:TsUnixNano()
            local Json = rulexlib:T2J(
                {
                    timestamp = ts,
                    temp = MatchHexS['temp'] * 0.1,
                }
            )
            rulexlib:Debug(Json)
        end
        return true, args
    end
}
