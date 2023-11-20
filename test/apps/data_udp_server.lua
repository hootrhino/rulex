AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'
function Main(arg)
    for i = 1, 10, 1 do
        local data = { name = 'Demo', sn = 'A123456', state = '00' }
        local err = data:ToUdp('UdpServer', applib:T2J(data))
        applib:log('DataToUdp success? =>', err == nil)
        time:Sleep(100)
    end
    return 0
end
