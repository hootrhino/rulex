---@diagnostic disable: undefined-global
-- Success
function Success()
    --
end
-- Failed
function Failed(error)
    print(error)
end
-- '{"type":5, "params":[{"address":1,"quantity":1,"values":1}]}'
-- Actions
Actions = {function(data)
    local t = {}
    local e = {}
    e["address"] = 1
    e["quantity"] = 1
    e["values"] = 1
    t["type"] = 5
    t["params"] = e
    rulexlib:WriteInStream('INEND_85211405-7285-446b-a120-ede14767a89f', rulexlib:JsonEncode(t))
    return false, data
end}
