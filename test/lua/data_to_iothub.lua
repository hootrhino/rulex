---@diagnostic disable: undefined-global
--- 数据透传到Iothub, 其中数据需要按照物模型来
-- Success
function Success()
    -- rulexlib:log("success")
end
-- Failed
function Failed(error)
    rulexlib:log("Error:", error)
end

-- Actions
Actions = {function(args)
    local ts = rulexlib:TsUnixNano()
    local t = {
        method = 'report',
        clientToken = ts,
        timestamp = ts,
        params = { -- 这个数据可以从data参数提取出来
            temp = 30.5,
            humi = 60.5
        }
    }
    local jsons = rulexlib:T2J(t)
    data:ToIoTHub(jsons)
    return true, args
end}
