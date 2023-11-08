-- 动态协议请求
AppNAME = 'Read'
AppVERSION = '0.0.1'
function Main(arg)
    local Id = 'DEVICE056b93901b3b4a5b9a3d69d14dc1139f'
    while true do
        local result, err = applib:CtrlDevice(Id, "010300000002C40B")
        --result {"in":"010300000002C40B","name":"write","out":"010304000100022a32"}
        print("CtrlDevice result=>", result)
        time:Sleep(60)
    end
    return 0
end
