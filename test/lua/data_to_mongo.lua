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
    dataToMongo("MongoDB001", data)
    print("[LUA Actions Callback]: Save to mongodb!")
    return true, data
end}
