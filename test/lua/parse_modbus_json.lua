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
        local json = require("json")
        local table = json.decode(data)
        local address = table['address']
        local value = table['value']
        local parseTable = rulexlib:MatchBinary(">high:8 low:8", value, false)
        if address == 0 then
            print("电压:", json.encode(parseTable))
        end
        if address == 1 then
            print("电流:", json.encode(parseTable))
        end
        return true, data
    end
}
