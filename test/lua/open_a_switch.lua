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
--- 这里展示一个远程发送到经过ESP8266控制的多路继电器后响应的 Demo：
--- 假设远程指令是打开开关，然后同步状态到云端, 指令体：
---     {"cmdId": "00001", "cmd" :"open","sw": [1, 2] }
---     [1, 2]表示打开 1 2 两个开关
---
Actions = {
    function(data)
        local json = require("json")
        local Tb = json.decode(data)
        local CmdId = Tb["cmdId"]
        local Cmd = Tb["cmd"]
        local SW = Tb["sw"]
        --- 开
        if Cmd == "on" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x01, SW}))
            if ok then
                rulex:finishCmd(CmdId)
            else
                -- 其实没必要显式调用失败，因为服务端超时后就自己直接失败了
                rulex:failedCmd(CmdId)
            end
        end
        --- 关
        if Cmd == "off" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x00, SW}))
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
