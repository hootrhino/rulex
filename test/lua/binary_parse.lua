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
    function(args)
        local json = require("json")
        local tb = rulexlib:MB("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = rulexlib:B2I64(1, rulexlib:BS2B(tb["a"]))
        result['b'] = rulexlib:B2I64(1, rulexlib:BS2B(tb["b"]))
        result['c'] = rulexlib:B2I64(1, rulexlib:BS2B(tb["c"]))
        result['d1'] = rulexlib:B2I64(1, rulexlib:BS2B(tb["d1"]))
        print("rulexlib:MB 2:", json.encode(result))
        data:ToMqtt('OUTEND', json.encode(result))
        return true, args
    end
}

