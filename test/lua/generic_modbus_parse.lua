-- {
--     "d1":{
--         "tag":"d1",
--         "function":3,
--         "slaverId":1,
--         "address":0,
--         "quantity":2,
--         "value":"AiYBDA=="
--     },
--     "d2":{
--         "tag":"d2",
--         "function":3,
--         "slaverId":2,
--         "address":0,
--         "quantity":2,
--         "value":"AicBCQ=="
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
    local jt = rulexlib:J2T(data)
    for k, v in pairs(jt) do
        local ht = rulexlib:MB('>hv:16 tv:16', v['value'], false)
        print(k, "Raw value:", ht['hv'], ht['tv'])
        local humi = rulexlib:B2I64('>', rulexlib:BS2B(ht['hv']))
        local temp = rulexlib:B2I64('>', rulexlib:BS2B(ht['tv']))
        local ts = rulexlib:TsUnixNano()
        local jsont = {
            method = 'report',
            clientToken = ts,
            timestamp = ts,
            params = {
                temp = temp,
                humi = humi
            }
        }
        print(k, "Parsed value:", rulexlib:T2J(jsont))
    end
    return true, args
end}
