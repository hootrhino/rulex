-- data = [
--     {"tag":"add1", "id": "001", "value": 0x0001},
--     {"tag":"add2", "id": "002", "value": 0x0002},
-- ]
function ParseData(data)
    for index, child in ipairs(data) do
        rulexlib:DataToMqtt('OUTEND', "topic/001", child.value)
    end
end
