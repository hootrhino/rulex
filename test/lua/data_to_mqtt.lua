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
        stdlib:DataToMqttServer('OUTEND_58fa7728-b82e-4124-8380-f56eec2731b9', data)
        return true, data
    end
}
