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
        local err1 = data:ToMqtt('$UUID', data)
        if err1 ~= nil then
            -- DO YOUR FAILED HANDLE
        end
        return true, data
    end
}
