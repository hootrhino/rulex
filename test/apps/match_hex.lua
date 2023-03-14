AppNAME = 'applib:MatchHex'
AppVERSION = '0.0.1'
function Main(arg)
    -- 十六进制提取器
    local MatchHexS = applib:MatchHex("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")
    for key, value in pairs(MatchHexS) do
        print('applib:MatchHex', key, value)
    end
    return 0
end
