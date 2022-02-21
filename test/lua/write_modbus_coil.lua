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
    local t = {
        ["type"] = 5,
        ["params"] = {
            ["address"] = 1,
            ["quantity"] = 1,
            ["value"] = 0xFF00
        }
    }
    rulexlib:WriteInStream('INEND', rulexlib:T2J(t))
    return false, data
end}
