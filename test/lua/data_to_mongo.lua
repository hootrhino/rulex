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
        rulexlib:DataToMongo('OUTEND_7ce89673-d3cf-43fe-82c8-0cf4c2be50a8', data)
        return true, data
    end
}
