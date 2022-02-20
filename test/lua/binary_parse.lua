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
        local tb = rulexlib:MB("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = rulexlib:ByteToInt64(1, rulexlib:BitStringToBytes(tb["a"]))
        result['b'] = rulexlib:ByteToInt64(1, rulexlib:BitStringToBytes(tb["b"]))
        result['c'] = rulexlib:ByteToInt64(1, rulexlib:BitStringToBytes(tb["c"]))
        result['d1'] = rulexlib:ByteToInt64(1, rulexlib:BitStringToBytes(tb["d1"]))
        print("rulexlib:MB 2:", json.encode(result))
        rulexlib:DataToMqttServer('OUTEND_bcb1b88c-7d83-49a3-88c0-d048e6368089', json.encode(result))
        return true, data
    end
}
