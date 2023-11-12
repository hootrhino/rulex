-- data = [
--     {"tag":"add1", "id": "001", "value": 0x0001},
--     {"tag":"add2", "id": "002", "value": 0x0002},
-- ]

function ParseData(data)
    -- data: {"in":"AA0011...","out":"AABBCDD..."}
    local DataT, err = rulexlib:J2T(data)
    if err ~= nil then
        return true, args
    end
    -- Do your business
    rulexlib:log(DataT['in'])
    rulexlib:log(DataT['out'])
end
