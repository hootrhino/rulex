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
    function(args)
        print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", rulexlib:JQ(".[] | select(.temp > 50000000)", data))
        return true, args
    end,
    function(args)
        print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", rulexlib:JQ(".[] | select(.hum < 20)", data))
        return true, args
    end,
    function(args)
        print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", rulexlib:JQ(".[] | select(.co2 > 50)", data))
        return true, args
    end,
    function(args)
        print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", rulexlib:JQ(".[] | select(.lex > 50)", data))
        return true, args
    end
}
