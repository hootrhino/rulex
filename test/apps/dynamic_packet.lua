-- 动态协议请求
AppNAME = 'DataToIthings'
AppVERSION = '0.0.1'
function Main(arg)
    local Id = 'DEVICEd78ad724852f4cb9a29c6dd6bf6c2f99'
    while true do
        local result, err = applib:CtrlDevice(Id, "write", "010300000002C40B")
        --result {"in":"010300000002C40B","name":"write","out":"010304000100022a32"}
        print("CtrlDevice result=>", result)
        print("XOR=>", misc:XOR("0", 0))
        applib:Sleep(2000)
    end
    return 0
end
