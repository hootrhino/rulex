---@diagnostic disable: undefined-global
--- 
Actions = {
    function(data)
        print("[LUA Actions Callback]: Mqtt payload:", data)
        if Result["temp"] ~= nil then
            dataToMongo("MongoDB001", Result)
            print("[LUA Actions Callback]: Save to mongodb!")
        end
        --if Result["hum"] ~= nil then
            -- TODO more target support
            -- dataToKafka("Kafka001", Result)
            -- print("[LUA Actions Callback]: Save to Kafka!")
        --end
        --if Result["hum"] ~= nil then
            -- TODO more target support
            -- dataToMysql("Mysql001", Result)
            -- print("[LUA Actions Callback]: Save to Mysql!")
        --end
        return true, data
    end
}
