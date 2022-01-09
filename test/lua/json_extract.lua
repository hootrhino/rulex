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
Actions = {
    -- data= [
    --          {"name":"modbus", "value": [0,1,1,1,0,1,0,1,1]},
    --          {"name":"uart", "value": [0,1]}
    --       ]
    function(data)
        local jsonData = json.decode(data)
        for _, value in ipairs(jsonData) do
            local name = value["name"]
            if name =="modbus" then
                rulexlib:DataToMqttServer('f56eec2731b9', {jsonData[0], jsonData[2], jsonData[5]})
                return true, data
            end
        end
        return false, data
    end
}
