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
    --        ┌───────────────────────────────────────────────┐
    -- data = |00 00 00 01|00 00 00 01|00 00 00 01|00 00 00 01|
    --        └───────────────────────────────────────────────┘
    function(data)
        local json = require("json")
        local tb = rulex:MatchBinary("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["a"]))
        result['b'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["b"]))
        result['c'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["c"]))
        result['d1'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["d1"]))
        print("rulex:MatchBinary 2:", json.encode(result))
        rulex:DataToMqttServer('OUTEND_bcb1b88c-7d83-49a3-88c0-d048e6368089', json.encode(result))
        return true, data
    end
}
