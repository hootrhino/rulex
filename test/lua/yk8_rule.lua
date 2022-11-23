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
-- data =
-- {
--     "method": "property",
--     "requestId": "20a4ccfd",
--     "timestamp": 123456,
--     "params": {
--         "sw1": "1",
--         "sw2": "1",
--         "sw3": "1",
--         "sw4": "1",
--         "sw5": "1",
--         "sw6": "1",
--         "sw7": "1",
--         "sw8": "1"
--     }
-- }
-- Actions
Actions = { function(data)
    local dataT, err = rulexlib:J2T(data)
    -- 兼容多种平台
    if dataT['method'] == 'control' or dataT['method'] == 'property' then
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
        local n, err = rulexlib:WriteDevice('YK8Device1', 0, rulexlib:T2J({ {
            ['function'] = 15,
            ['slaverId'] = 3,
            ['address'] = 0,
            ['quantity'] = 1,
            ['value'] = rulexlib:T2Str(cmd)
        } }))
        if (err) then
            rulexlib:Throw(err)
        end
        -- read data publish to mqtt
        local n, err = rulexlib:WriteSource('tencentIothub', rulexlib:T2J({
            method = 'control_reply',
            clientToken = dataT['clientToken'],
            code = 0,
            status = 'OK'
        }))
        if (err) then
            rulexlib:Throw(err)
        end
    end
    return true, data
end }
