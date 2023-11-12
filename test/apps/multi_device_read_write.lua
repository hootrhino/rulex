AppNAME = 'MultiDeviceTest'
AppVERSION = '0.0.1'
AppDESCRIPTION = 'MultiDeviceTest'
function Main(arg)
    local deviceUUID = 'DEVICE19d716e11ddb43b4a2a080f547fff6c1'
    while true do
        ::continue::
        local binary1, err1 = applib:ReadDevice(deviceUUID, "0011AABBCCFF00")
        if err1 ~= nil then
            print(1, err1)
            goto continue
        end
        local d1_state = string.sub(binary1, 8, 9)
        if d1_state == "1" then
            applib:WriteDevice(deviceUUID, "l1", "0011AABBCCFF00")
        end
        if d1_state == "0" then
            applib:WriteDevice(deviceUUID, "l2", "0011AABBCCFF11")
        end
        time:Sleep(1000)
    end
end
