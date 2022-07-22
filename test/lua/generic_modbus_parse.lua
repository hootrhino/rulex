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
Actions = {function(data)
    local jt = rulexlib:J2T(data)
    for k, value in pairs(jt) do
        local dt = rulexlib:MB('<tb:16 vb:16', value['value'], false)
        print(dt)
        local hv = rulexlib:B2I64(1, rulexlib:BS2B(dt['tb']))
        print(hv)
        local tv = rulexlib:B2I64(1, rulexlib:BS2B(dt['vb']))
        print(tv)
        print(k, hv, tv)
    end
    return true, data
end}
