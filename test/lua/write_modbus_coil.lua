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
    rulexlib:WriteInStream('INEND_77c16142-f849-48c3-b150-34aed2d0d9ae', rulexlib:JsonEncode(t))
    return false, data
end}
