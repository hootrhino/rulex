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
    local dataT, err = rulexlib:J2T(data)
    if dataT['method'] == 'property' then
        local params = dataT['params']
        local cmd = {
            [1] = params['sw8'],
            [2] = params['sw7'],
            [3] = params['sw6'],
            [4] = params['sw5'],
            [5] = params['sw4'],
            [6] = params['sw3'],
            [7] = params['sw2'],
            [8] = params['sw1']
        }
        local n1, err1 = rulexlib:WriteDevice('YK8Device1', rulexlib:T2J({{
            ['function'] = 15,
            ['slaverId'] = 3,
            ['address'] = 0,
            ['quantity'] = 1,
            ['value'] = rulexlib:T2Str(cmd)
        }}))
        if (err1) then
            rulexlib:Throw(err1)
        end
        local n2, err2 = rulexlib:WriteSource('tencentIothub', rulexlib:T2J({
            method = 'reply',
            clientToken = dataT['clientToken'],
            code = 1,
            status = 'OK'
        }))
        if (err2) then
            rulexlib:Throw(err2)
        end
    end
    return true, data
end}
