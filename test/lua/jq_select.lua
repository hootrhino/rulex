---@diagnostic disable: undefined-global

-- Success
function Success()
    print("======> success")
end
-- Failed
function Failed(error)
    print("======> failed:", error)
end

-- Actions
Actions = {
    function(data)
        print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", rulex:JqSelect(".[] | select(.temp > 50000000)", data))
        return true, data
    end,
    function(data)
        print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", rulex:JqSelect(".[] | select(.hum < 20)", data))
        return true, data
    end,
    function(data)
        print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", rulex:JqSelect(".[] | select(.co2 > 50)", data))
        return true, data
    end,
    function(data)
        print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", rulex:JqSelect(".[] | select(.lex > 50)", data))
        return true, data
    end
}
