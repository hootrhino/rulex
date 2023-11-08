AppNAME = "test_demo"
AppVERSION = "1.0.0"
AppDESCRIPTION = "从设备里面读取数据然后推送到MQTT"

function Main(arg)
    for i = 1, 10, 1 do
        local data, err1 = applib:ReadDevice("uuid", 0, "192.168.1.1:502")
        if err1 ~= nil then
            applib:log(err1)
            return 0
        end
        local err2 = applib:DataToMqtt('UUID', applib:T2J({
            temp = i,
            humi = 13.45
        }))
        if err2 ~= nil then
            applib:log(err1)
            return 0
        end
        time:Sleep(1000)
    end
    return 0
end
