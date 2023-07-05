-- 应用名称 "ws1608:GPIO Test"

function Main(arg)
    while true do
        ws1608:GPIOSet("red", 1)
        applib:Sleep(2000)
        ws1608:GPIOSet("red", 0)
        applib:Sleep(2000)
    end
end
