Actions = {
    function(args)
        print('Received data from iothub:', data)
        local source = 'iothub-mqtt'
        local device = 'device-uuid'
        local dataT, err = rulexlib:J2T(data)
        if (err ~= nil) then
            print('parse json error:', err)
            return false, data
        end
        local requestId = dataT['id']
        -- Action:
        -- {
        --     "method":"action",
        --     "id":"20a4ccfd",
        --     "actionId":"actionid",
        --     "timestamp":13124323543,
        --     "data":{
        --         "a":"1",
        --         "b":"2"
        --     }
        -- }
        if dataT['method'] == 'action' then
            local actionId = dataT['actionId']
            if actionId == 'Action get_state' then
                local readData, err = rulexlib:ReadDevice(0, device)
                if (err ~= nil) then
                    print('Action get_state error:', err)
                    return false, data
                end
                print('Action get_state:', readData)
                local _, err = iothub:ActionSuccess(source, requestId, { code = 0 })
                if (err ~= nil) then
                    print('Action Reply error:', err)
                    return false, data
                end
            end
        end
        -- Property:
        -- {
        --     "method":"property",
        --     "id":"20a4ccfd",
        --     "timestamp":13124323543,
        --     "data":{
        --         "a":"1"
        --     }
        -- }
        if dataT['method'] == 'property' then
            local _, err = iothub:PropertySuccess(source, requestId, { code = 0 })
            if (err ~= nil) then
                print('Property Reply error:', err)
                return false, data
            end
            print('IotHUB property:', dataT['data'])
        end
        return true, args
    end
}
