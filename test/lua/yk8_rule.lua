---@diagnostic disable: undefined-global
-- Success
function Success()
    --
end
-- Failed
function Failed(error)
    print(error)
end
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
