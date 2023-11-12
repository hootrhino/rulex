-- {
--     "d1":{
--         "tag":"d1",
--         "function":3,
--         "slaverId":1,
--         "address":0,
--         "quantity":2,
--         "value":""
--     }
-- }
---@diagnostic disable: undefined-global
-- Success
function Success()
    -- rulexlib:log("success")
end
-- Failed
function Failed(error)
    rulexlib:log("Error:", error)
end

-- Actions
Actions = {function(args)
    for _, v in pairs(rulexlib:J2T(data)) do
        local ts = rulexlib:TsUnixNano()
        local jsont = {
            method = 'report',
            requestId = ts,
            timestamp = ts,
            params = v['value']
        }
        data:ToMqtt('mqttOutEnd', rulexlib:T2J(jsont))
    end
    return true, args
end}
