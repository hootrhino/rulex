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
        -- {"type": "OPEN", "sn": "SN0001"}
        local Type = data["type"]
        local SN = data["sn"]
        -- 给底下的设备端发指令
        stdlib:WriteInStream('9e7c45352db1', data)
        -- 给远程端回复指令
        stdlib:WriteOutStream('9e7c45352db1', data)
        return true, data
    end
}
