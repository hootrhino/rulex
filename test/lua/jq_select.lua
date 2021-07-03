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
        local V = Select(".[] | select(.temp > 50)", data)
        if V ~=nil then
            dataToKafka("kafka001", V)
        end
        return true, data
    end,
    function(data)
        local V = Select(".[] | select(.hum < 30)", data)
        if V ~=nil then
            dataToMongo("mongo001", V)
        end
        return true, data
    end,
    function(data)
        local V = Select(".[] | select(.lex > 500)", data)
        if V ~=nil then
            dataToMysql("mysql001", V)
        end
        return true, data
    end,
    function(data)
        local V = Select(".[] | select(.co2 > 30) | select(.co2 < 50)", data)
        if V ~=nil then
            data["co2"] = 0
        end
        return true, data
    end
}
