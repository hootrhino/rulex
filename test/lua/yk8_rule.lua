---@diagnostic disable: undefined-global
-- Success
function Success()
    --
end
-- Failed
function Failed(error)
    print(error)
end
-- 属性下发以后，格式如下:
-- data = {
--     "method": "property",
--     "requestId": "20a4ccfd",
--     "timestamp": 0123456,
--     "params": {
--         "sw1": "1"
--         "sw2": "1"
--         "sw3": "1"
--         "sw4": "1"
--         "sw5": "1"
--         "sw6": "1"
--         "sw7": "1"
--         "sw8": "1"
--     }
-- }
-- Actions
Actions = {function(data)

    rulexlib:WriteDevice('mqttOutEnd-iothub', rulexlib:T2J({{
        ['function'] = 15,
        ['slaverId'] = 1,
        ['address'] = 0,
        ['quantity'] = 1,
        ['value'] = '00011000'
    }, {
        ['function'] = 15,
        ['slaverId'] = 1,
        ['address'] = 0,
        ['quantity'] = 1,
        ['value'] = '01111000'
    }}))
    return true, data
end}
