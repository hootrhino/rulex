function Main(Arg)
    local RawData, err = applib:Control('UUID', '0A01020230303FFAACC', '')
    if err ~= nil then
        error(err)
        return
    end
    -- Lua 校验
    local Valid = applib:CheckCRC(RawData, 1, 2, 3, 4)
    if ~Valid then
        error(err)
        return
    end
    -- 数据解码
    print("Success:", applib:ParseABDC(RawData))
    return 0
end
