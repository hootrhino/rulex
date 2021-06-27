---@diagnostic disable: undefined-global
-- From
function From()
    return {"INEND1"}
end
-- Success
function Success()
    print("======> success")
end
-- Failed
function Failed(error)
    print("======> failed:", error)
end

-- Actions
Actions = {function(data)
    print("[LUA Actions Callback]: Mqtt payload:", data)
    Result = Select(data, "select temp,hum from INPUT_DATA where temp > '100' and hum < '24'")
    print("[LUA Actions Callback]Result ===>", Result)
    return true, data
end}
