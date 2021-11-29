---@diagnostic disable: undefined-global
-- Success
function Success()
    rulex:log("success")
end
-- Failed
function Failed(error)
    rulex:log(error)
end

---
--- 这里展示一个远程发送指令后响应的Demo
--- 假设远程指令是打开开关，然后同步状态到云端,
--- 指令体：{
---            "cmdId" : "hu008987yp7yujjm",
---            "type" : "OPEN",
---            "sn": [
---                   "SN0001",
---                   "SN0002"
---                  ]
---        }
--- 表示打开 SN0001 SN0002 两个开关
---
Actions = {
    function(data)
        local json = require("json")
        local Tb = json.decode(data)
        local CmdId = Tb["cmdId"]
        local Type = Tb["type"]
        local SN = Tb["sn"]
        if Type == "OPEN" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x00, SN}))
            if ok then
                rulex:finishCmd(CmdId)
            else
                -- 其实没必要显式调用失败，因为服务端超时后就自己直接失败了
                rulex:failedCmd(CmdId)
            end
        end
        if Type == "OFF" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x01, SN}))
            if ok then
                rulex:finishCmd(CmdId, "OutId")
            else
                -- 其实没必要显式调用失败，因为服务端超时后就自己直接失败了
                rulex:failedCmd(CmdId, "OutId")
            end
        end
        return true, data
    end
}
