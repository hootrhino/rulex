AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'
function Main(arg)
    local deviceUUID = 'deviceUUID'
    local udpServerUUID = 'udpServerUUID'
    local Subdevices = { '01', '02' }

    while true do
        for _, sn in ipairs(Subdevices) do
            local binary, err1 = applib:ReadDevice(deviceUUID, sn)
            if err1 ~= nil then
                print("ReadDevice error:", err1)
                goto END
            end
            local jsonS, err2 = applib:Bin2Str(binary)
            if err2 ~= nil then
                print("Bin2Str error:", err1)
                goto END
            end
            print("ReadDevice => ", jsonS)
            local jsonT, err3 = applib:J2T(jsonS)
            if err3 ~= nil then
                print("J2T error:", err1)
                goto END
            end
            local udpDataJson = applib:T2J {
                model = 'A',
                sn = string.sub(jsonT['out'], 1, 6),
                state = string.sub(jsonT['out'], 9, 10)
            }
            print("UdpData => ", udpDataJson)
            local err4 = applib:DataToUdp(udpServerUUID, udpDataJson)
            print('DataToUdp success? =>', err4 == nil)
            time:Sleep(1000)
        end
        ::END::
    end
end
