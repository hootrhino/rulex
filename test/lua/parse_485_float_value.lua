Actions = {
    function(args)
        local dataT, err = json:J2T(args)
        if (err ~= nil) then
            Throw('parse json error:' .. err)
            return true, args
        end
        for key, value in pairs(dataT) do
            Debug(key .. " :: " .. value['value'])
            local MatchHexS = hex:MatchHex("humi:[0,3];temp:[4,7];pres:[8,11]", value['value'])
            local ts = time:Time()
            local Json = json:T2J(
                {
                    tag = key,
                    ts = ts,
                    hum = math:TFloat(binary:Bin2F32(MatchHexS['humi'])),
                    temp = math:TFloat(binary:Bin2F32(MatchHexS['temp'])),
                    pres = math:TFloat(binary:Bin2F32(MatchHexS['pres'])),
                })
            Debug(Json)
        end
        return true, args
    end
}
