--
-- App use lua syntax, goto https://hootrhino.github.io for more document
-- APPID: APP12a32a3d45df4555a50e04472a78ab4b
--
AppNAME = "helloworld"
AppVERSION = "1.0.2"
AppDESCRIPTION = "helloworld"

function ToScreen(In, Out, Data)
    local result, err1 = applib:CtrlDevice(In, "ctrlin30ms", Data)
    print("err1=>", err1)
    print("result=>", result)
    local jsonT, err3 = applib:J2T(result)
    if err3 ~= nil then
        print(3, err3)
        return
    end
    local udpData = applib:T2J{
        model = 'HX-S20',
        sn = string.sub(jsonT['out'], 1, 2),
        state = string.sub(jsonT['out'], 14, 14)
    }
    print("udpData => ", udpData)
    local err4 = applib:DataToUdp(Out, udpData)
    print('DataToUdp success? =>', err4 == nil)
end

--
-- Main
--
-- DEVICEa3a94f6c1ae84f56a92f0f4cd7f53cc4
function Main()
    local Inid = 'DEVICE7e31f1820a7f4c92977d26bcae2aae69'
    local Outid = 'OUTfde44d8992ce4fde82a3b096352b10cf'
    local cmdlist = {
        '010300010003540B', '02030000000305F8', '0303000000030429',
        '040300000003059E', '050300000003044F', '060300000003047C',
        '07030000000305AD', '0803000000030552', '0903000000030483',
        '0A030000000304B0'
    }
    while true do

        for _, cmd in ipairs(cmdlist) do
			ToScreen(Inid,Outid,cmd)
            time:Sleep(500)
        end

    end
end

