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
    n, err = rulexlib:WriteDevice('YK8Device1', rulexlib:T2J({{
        ['function'] = 15,
        ['slaverId'] = 3,
        ['address'] = 0,
        ['quantity'] = 1,
        ['value'] = data
    }}))
    return false, data
end}
