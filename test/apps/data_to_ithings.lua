AppNAME = 'DataToIthings'
AppVERSION = '0.0.1'
function Main(arg)
    local IthingsUUID = 'INe72f62072c314ec195d2f777398db019'
    while true do
        local udpData = applib:T2J {
            method = 'subDevReport',
            SubDeviceName = 'RULEX-大屏1',
            params = {
                model = 'A001',
                sn = "1000",
                state = "1"
            }
        }
        print("DataToIthings => ", udpData)
        local _, err4 = applib:WriteSource(IthingsUUID, udpData)
        print('DataToIthings success? =>', err4 == nil)
        applib:Sleep(2000)
    end
    return 0
end
