---@diagnostic disable: undefined-global
-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions =
    {
    --        ┌───────────────────────────────────────────────────────────────────────────────┐
    -- data = |0X00 0X01 0X00 0X02|0X00 0X03 0X00 0X04|0X00 0X01 0X00 0X02|0X00 0X00 0X00 0X01|
    --        └───────────────────────────────────────────────────────────────────────────────┘
    function(data)
        local json = require("json")
        local tb = rulex:MatchBinary("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["a"]))
        result['b'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["b"]))
        result['c'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["c"]))
        result['d1'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["d1"]))
        rulex:DataToMqttServer('ID', json.encode(result))
        return true, data
    end
}
