Actions = {
    function(data)
        print('Received data from iothub:', data)
        local source = 'tencentIothub'
        local device = 'YK8Device1'
        local dataT, err = rulexlib:J2T(data)
        if (err ~= nil) then
            print('Received data from iothub parse to json error:', err)
            return false, data
        end
        -- Action
        if dataT['method'] == 'action' then
            local actionId = dataT['actionId']
            if actionId == 'get_state' then
                local readData, err = rulexlib:ReadDevice(0, device)
                if (err ~= nil) then
                    print('ReadDevice data from device error:', err)
                    return false, data
                end
                print('ReadDevice data from device:', readData)
                local readDataT, err = rulexlib:J2T(readData)
                if (err ~= nil) then
                    print('Parse ReadDevice data from device error:', err)
                    return false, data
                end
                local yk08001State = readDataT['yk08-001']
                print('yk08001State:', yk08001State['value'])
                local _, err = iothub:ActionSuccess(source, dataT['id'],
                    yk08001State['value'])
                if (err ~= nil) then
                    print('ActionReply error:', err)
                    return false, data
                end
            end
        end
        -- Property
        if dataT['method'] == 'property' then
            local schemaParams = rulexlib:J2T(dataT['data'])
            print('schemaParams:', schemaParams)
            local n1, err = rulexlib:WriteDevice(device, 0, rulexlib:T2J({
                {
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
                }
            }))
            if (err ~= nil) then
                print('WriteDevice error:', err)
                local _, err = iothub:PropertyFailed(source, dataT['id'])
                if (err ~= nil) then
                    print('Reply error:', err)
                    return false, data
                end
            else
                local _, err = iothub:PropertySuccess(source, dataT['id'], {})
                if (err ~= nil) then
                    print('Reply error:', err)
                    return false, data
                end
            end
        end
        return true, data
    end
}
