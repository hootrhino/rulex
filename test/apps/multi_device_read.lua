AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'
AppDESCRIPTION = '多个设备的场景'
function Main(arg)
    local deviceUUID = 'A'
    while true do
        ::continue::
        local binary1, err1 = applib:ReadDevice(deviceUUID, "devie1")
        if err1 ~= nil then
            print(1, err1)
            goto continue
        end
        print(binary1)
        local binary2, err1 = applib:ReadDevice(deviceUUID, "devie2")
        if err1 ~= nil then
            print(1, err1)
            goto continue
        end
        print(binary2)
        time:Sleep(1000)
    end
    return 0
end
