local source = 'tencentIothub'
local device = 'YK8Device1'
function Rule(data)
    rulexlib:log('Received data from iothub:', data)
    local dataT, err = rulexlib:J2T(data)
    if (err ~= nil) then
        rulexlib:log('Received data from iothub parse to json error:', err)
        return false, data
    end
    -- Action
    if dataT['method'] == 'action' then
        local actionId = dataT['actionId']
        if actionId == 'get_state' then
            local readData, err = rulexlib:ReadDevice(device)
            if (err ~= nil) then
                rulexlib:log('ReadDevice data from device error:', err)
                return false, data
            end
            rulexlib:log('ReadDevice data from device:', readData)
            local readDataT, err = rulexlib:J2T(readData)
            if (err ~= nil) then
                rulexlib:log('Parse ReadDevice data from device error:', err)
                return false, data
            end
            local _, err = iothub:ActionReplySuccess(source, dataT['id'], readDataT['value'])
            if (err ~= nil) then
                rulexlib:log('ActionReply error:', err)
                return false, data
            end
        end
    end
    -- Property
    if dataT['method'] == 'property' then
        local schemaParams = dataT['data']
        local n1, err = rulexlib:WriteDevice(device, rulexlib:T2J({{
            ['function'] = 15,
            ['slaverId'] = 3,
            ['address'] = 0,
            ['quantity'] = 1,
            ['value'] = rulexlib:T2Str({
                [1] = schemaParams['sw8'],
                [2] = schemaParams['sw7'],
                [3] = schemaParams['sw6'],
                [4] = schemaParams['sw5'],
                [5] = schemaParams['sw4'],
                [6] = schemaParams['sw3'],
                [7] = schemaParams['sw2'],
                [8] = schemaParams['sw1']
            })
        }}))
        if (err ~= nil) then
            rulexlib:log('WriteDevice error:', err)
            local _, err = iothub:PropertyReplyFailed(source, dataT['id'])
            if (err ~= nil) then
                rulexlib:log('Reply error:', err)
                return false, data
            end
        else
            local _, err = iothub:PropertyReplySuccess(source, dataT['id'], {})
            if (err ~= nil) then
                rulexlib:log('Reply error:', err)
                return false, data
            end
        end
    end
    return true, data
end
