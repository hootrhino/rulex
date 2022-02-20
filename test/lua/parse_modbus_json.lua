---@diagnostic disable: undefined-global
-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end

--
-- data = {"function":3,"address":0,"quantity":1," value":1}
-- data = {"function":3,"address":1,"quantity":1," value":1}
-- -----------------------
-- |电压      |电流
-- -----------------------
-- |0x00 0x00 | 0x00 0x00
-- -----------------------
Actions = {
    function(data)
        rulexlib:log("data:", data)
        local json = require("json")
        local table = json.decode(data)
        local value = table['value']
        local parseTable = rulexlib:MB(">u:16 v:16", value, false)
        rulexlib:log("UA:", json.encode(parseTable))
        return true, data
    end
}
