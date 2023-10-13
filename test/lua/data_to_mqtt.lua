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
        print('Data ======> ', data)
        data:ToMqtt('OUTEND', data)
        return true, data
    end
}
