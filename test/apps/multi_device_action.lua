AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'
function Main(arg)
    local deviceUUID = 'A'
    -- local udpServerUUID = 'B'
    while true do
        repeat
            local binary, err1 = applib:ReadDevice(deviceUUID, 1)
            if err1 ~= nil then
                print(1, err1)
                break
            end
            local jsonS, err2 = applib:Bin2Str(binary)
            if err2 ~= nil then
                print(2, err2)
                break
            end
            print("ReadDevice => ", jsonS)
            local jsonT, err3 = applib:J2T(jsonS)
            if err3 ~= nil then
                print(3, err3)
                break
            end
            print("Out ==> ", jsonT['out'])
            local state = string.sub(jsonT['out'], 4, 5)
            print("State ==> ", state)
            if state == '0' then
                -- 74 02 25 0A 01 02 5A AA 55 无人
                applib:ReadDevice(deviceUUID, 2)
            else
                -- 74 02 25 0A 01 01 59 AA 55 有人
                applib:ReadDevice(deviceUUID, 3)
            end
            time:Sleep(1000)
        until true
    end

    return 0
end
