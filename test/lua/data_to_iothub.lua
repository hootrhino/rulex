---@diagnostic disable: undefined-global
--- 数据透传到Iothub, 其中数据需要按照物模型来
-- Success
function Success()
    -- rulexlib:log("success")
end

-- Failed
function Failed(error)
    -- rulexlib:log("Error:", error)
end

-------------------------------------------------------------------------------
-- {
--     "id": "1",
--     "method":"thing.event.property.post",
--     "params": {
--          "ais_data":"!AIVDO,1,1,,A,13Mum40000000G?vMdHG6Hi2>220S0@0,0*67"
--     }
-- }
-------------------------------------------------------------------------------
-- {
--     "type":"VDM",
--     "gwid":"HR0001",
--     "message_id":19,
--     "user_id":413825345,
--     "name":"YUXINHUO16626",
--     "sog":3.2,
--     "longitude":114.347,
--     "latitude":30.62909,
--     "cog":226.3,
--     "true_heading":511,
--     "timestamp":35
-- }
-- Actions
Actions = { function(args)
    local error1, JsonT = json:J2T(args)
    local t =
    {
        id = string:MakeUid(),
        method = "thing.event.property.post",
        params = {
            ais_data = JsonT['ais_data']
        }
    }
    local jsons = json:T2J(t)
    local error = data:ToMqtt('OUTQAQXBVCU', jsons)
    if error ~= nil then
        stdlib:Throw(error)
    end
    return true, args
end }
