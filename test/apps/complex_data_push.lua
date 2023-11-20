AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'
function Main(arg)
    local deviceUUID = 'A'
    local udpServerUUID = 'B'

    while true do
        repeat
            local binary, err1 = applib:ReadDevice(deviceUUID, "p1")
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
            local udpData = applib:T2J {
                model = 'A',
                sn = string.sub(jsonT['out'], 1, 6),
                state = string.sub(jsonT['out'], 9, 10)
            }
            print("udpData => ", udpData)
            local err4 = data:ToUdp(udpServerUUID, udpData)
            print('DataToUdp success? =>', err4 == nil)
            time:Sleep(1000)
        until true
    end

    return 0
end
